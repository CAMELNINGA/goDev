package http

import (
	"Yaratam/internal/domain"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"time"
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

func (a *adapter) UploadMultipartFile(file io.ReadCloser, username string, unit string, fileName string) (string, error) {
	var timeHTTPClient = &http.Client{Timeout: 30 * time.Second}
	apiLink := a.config.UploadURL + "/upload"

	defer file.Close()

	body, writer := io.Pipe()

	req, err := http.NewRequest(http.MethodPut, apiLink, body)
	if err != nil {
		return "", err
	}

	mwriter := multipart.NewWriter(writer)
	req.Header.Add("Content-Type", mwriter.FormDataContentType())
	req.Header.Set("token", "0")

	errchan := make(chan error)

	go func() {
		defer close(errchan)
		defer writer.Close()
		defer mwriter.Close()

		w, err := mwriter.CreateFormFile("upload", fileName)
		if err != nil {
			errchan <- err
			return
		}

		if written, err := io.Copy(w, file); err != nil {
			errchan <- fmt.Errorf("error copying %s (%d bytes written): %v", fileName, written, err)
			return
		}

		if err := mwriter.Close(); err != nil {
			errchan <- err
			return
		}
	}()

	resp, err := timeHTTPClient.Do(req)
	merr := <-errchan

	if err != nil || merr != nil {
		a.logger.Error("http and multipart error", zap.Error(err), zap.Error(merr))
		return "", fmt.Errorf("http error: %v, multipart error: %v", err, merr)
	}
	if resp.StatusCode != 200 {
		a.logger.Error("multipart error, status code", zap.Int("StatusCode", resp.StatusCode))
		return "", domain.ErrNoOKAPI
	}
	var res struct {
		Link string `json:"path,omitempty"`
		Err  bool   `json:"ok,omitempty"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if res.Link == "" {
		a.logger.Error("multipart error, got empty link", res.Err)
		return "", domain.ErrNoOKAPI
	}

	return res.Link, nil
}
