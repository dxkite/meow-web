package httputil

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

func Response(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if data == nil {
		w.Header().Set("Content-Length", "0")
		return
	}
	body, _ := json.Marshal(data)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

func ResponseError(ctx context.Context, w http.ResponseWriter, err error) {
	statusCode := statusFromError(err)
	name := nameFromError(err)
	Response(ctx, w, statusCode, &HttpErrorDetail{
		Code:    name,
		Message: err.Error(),
	})
}
