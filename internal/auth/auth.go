package auth

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	"testapi/pkg/config"
)

//https://github.com/go-chi/jwtauth
func Get(cfg *config.Config) (*jwtauth.JWTAuth, error) {
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("unable to start server: %v", "JWT secret empty")
	}

	return jwtauth.New("HS256", cfg.JWTSecret, nil), nil
}

//Tee tarkistus ettei environment variablet ole nulleja (ei ole asetettu .env tiedostoon tai johonkin muualle)
