package auth

import (
	"errors"
	"github.com/go-chi/jwtauth"
	"testapi/pkg/config"
)

func Get(cfg *config.Config) (*jwtauth.JWTAuth, error) {
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT secret empty")
	}
	return jwtauth.New("HS256", []byte(cfg.JWTSecret), nil), nil
}
