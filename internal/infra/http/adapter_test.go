package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSayHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/test/hello", nil)
	w := httptest.NewRecorder()
	getHello(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "Hello!" {
		t.Errorf("expected ABC got %v", string(data))
	}
}
