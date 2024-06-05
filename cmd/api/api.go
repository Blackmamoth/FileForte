package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Blackmamoth/vault/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.Compress(5, "gzip"))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Route not Found for [%s] %s", r.Method, r.URL.Path)))
	})

	config.Logger.INFO(fmt.Sprintf("Application running on port: %s", s.addr))
	return http.ListenAndServe(fmt.Sprintf(":%s", s.addr), router)
}
