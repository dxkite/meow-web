package httputil

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"{{ .PackageName }}/pkg/errors"
)

func Result(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if data == nil {
		w.Header().Set("Content-Length", "0")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body, _ := json.Marshal(data)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

type HttpError struct {
	Error        string   `json:"error"`
	ErrorDetails []string `json:"error_details"`
}

func Error(ctx context.Context, w http.ResponseWriter, err error) {
	statusCode := statusFromError(err)
	name := nameFromError(err)
	details := []string{}
	for _, e := range errors.Unwrap(err) {
		details = append(details, e.Error())
	}
	Result(ctx, w, statusCode, &HttpError{
		Error:        name,
		ErrorDetails: details,
	})
}
