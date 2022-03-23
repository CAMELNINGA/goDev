package http

import (
	"net/http"
)

func (a *adapter) wrap(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			a.logger.WithFields(generateFields(r)).WithError(err).Error("Error handling request")
		}
	}
}
func getHello(w http.ResponseWriter, _ *http.Request) error {
	if _, err := w.Write([]byte("Hello!")); err != nil {
		return jError(w, err)
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
