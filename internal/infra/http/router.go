package http

import (
	chiprometheus "github.com/T-M-A/chi-prometheus"
	"github.com/go-chi/chi/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func (a *adapter) newRouter() (http.Handler, error) {
	r := chi.NewRouter()

	// Set default middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Use(chiprometheus.NewMiddleware("goDev"))
	r.Route("/api", func(r chi.Router) {
		r.Route("/test", func(r chi.Router) {
			r.Get("/hello", a.wrap(getHello))
		})

	})

	return r, nil
}
