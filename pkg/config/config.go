package config

import (
	"flag"
	"github.com/subosito/gotenv"
	"os"
)

type Config struct {
	DbDSN      string
	ServerPort string
	JWTSecret  string
}

func Get() (*Config, error) {
	err := gotenv.Load()
	if err != nil {
		return nil, err
	}

	conf := &Config{}

	flag.StringVar(&conf.DbDSN, "db_dsn", os.Getenv("DB_DSN"), "Database DSN")
	flag.StringVar(&conf.ServerPort, "port", os.Getenv("PORT"), "Server Port")
	flag.StringVar(&conf.JWTSecret, "jwt_secret", os.Getenv("JWT_SECRET"), "JWT Secret")
	flag.Parse()
	return conf, nil
}
