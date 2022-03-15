package http

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (a *adapter) newRouter() (http.Handler, error) {
	r := chi.NewRouter()
	return r, nil
}
