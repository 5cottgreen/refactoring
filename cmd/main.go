package main

import (
	"log"
	"net/http"
	"time"

	"github.com/5cottgreen/refactoring/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	handler "github.com/5cottgreen/refactoring/internal"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("'%s' faild to get configuration", err.Error())
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", handler.CheckHealth)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/fetch", handler.FetchUsers)
				r.Post("/create", handler.CreateUser)

				r.Route("/get", func(r chi.Router) {
					r.Get("/{id}", handler.GetUser)
				})

				r.Route("/update", func(r chi.Router) {
					r.Patch("/{id}", handler.UpdateUser)
				})

				r.Route("/delete", func(r chi.Router) {
					r.Delete("/{id}", handler.DeleteUser)
				})
			})
		})
	})

	addr := conf.Server.Host + ":" + conf.Server.Host

	http.ListenAndServe(addr, r)
}
