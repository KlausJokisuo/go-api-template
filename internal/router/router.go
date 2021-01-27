package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/http"
	"testapi/internal/users"
	"time"
)

func Get(dbClient *pgxpool.Pool) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/users", users.Get(users.NewUserRepository(dbClient)))

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.WithFields(log.Fields{
			"method": method,
			"route":  route,
		}).Info()
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		return r, err
	}

	return r, nil
}
