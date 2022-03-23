package http

import (
	"fmt"
	"net/http"
)

func (a *adapter) wrap(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			a.logger.WithFields(generateFields(r)).WithError(err).Error("Error handling request")
		}
	}
}
func (a *adapter) sayHello(w http.ResponseWriter, r *http.Request) error {
	for k, values := range r.Header {
		fmt.Print(k, ": ")
		for _, v := range values {
			fmt.Print(v, ", ")
		}
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello!"))
	return err
}
