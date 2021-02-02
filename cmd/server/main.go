package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	auth2 "testapi/internal/auth"
	"testapi/internal/router"
	"testapi/pkg/config"
	"testapi/pkg/db"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

//https://www.reddit.com/r/golang/comments/i3vb9z/switching_to_pgx_when_to_use_connection_pool_and/
func run() error {

	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("unable to read environment variables: %w", err)
	}

	database, err := db.Get(cfg.DbDSN)
	if err != nil {
		return fmt.Errorf("unable to connect database: %w", err)
	}

	tokenAuth, err := auth2.Get(cfg)
	if err != nil {
		return fmt.Errorf("unable to get token auth: %w", err)
	}

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:¨

	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)

	r, err := router.Get(database.Client, tokenAuth)
	if err != nil {
		return fmt.Errorf("unable to get routes: %w", err)
	}

	if err := http.ListenAndServe(fmt.Sprint(":", cfg.ServerPort), r); err != nil {
		return fmt.Errorf("unable to start server: %w", err)
	}

	return nil
}
