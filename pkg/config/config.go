package config

import (
	"flag"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/subosito/gotenv"
	"os"
)

type Config struct {
	DbDSN      string
	ServerPort string
	JWTSecret  string
}

func (c Config) ValidateConfig() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DbDSN, validation.Required),
		validation.Field(&c.ServerPort, validation.Required),
		validation.Field(&c.JWTSecret, validation.Required),
	)
}

func Get() (*Config, error) {
	_ = gotenv.Load()
	conf := Config{}

	flag.StringVar(&conf.DbDSN, "db-dsn", os.Getenv("DB-DSN"), "Database DSN")
	flag.StringVar(&conf.ServerPort, "port", os.Getenv("PORT"), "Server Port")
	flag.StringVar(&conf.JWTSecret, "jwt-secret", os.Getenv("JWT-SECRET"), "JWT Secret")
	flag.Parse()

	if err := conf.ValidateConfig(); err != nil {
		return nil, err
	}
	return &conf, nil
}
