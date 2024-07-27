package httputil

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"

	"dxkite.cn/meownest/pkg/errors"
)

func ReadJSON(r *http.Request, out interface{}) error {
	err := checkContentType(r, "application/json")
	if err != nil {
		return err
	}
	if r.Body == nil || r.ContentLength == 0 {
		return nil
	}

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(out)
	defer r.Body.Close()
	if err != nil {
		if err == io.EOF {
			return errors.InvalidParameter(errors.New("invalid JSON: got EOF while reading request body"))
		}
		return errors.InvalidParameter(errors.Wrap(err, "invalid JSON"))
	}

	if dec.More() {
		return errors.InvalidParameter(errors.New("unexpected content after JSON"))
	}

	if err := Validator.Struct(out); err != nil {
		return errors.InvalidParameter(errors.Wrap(err, "validate error"))
	}

	return nil
}

func checkContentType(r *http.Request, mimeType string) error {
	ct := r.Header.Get("Content-Type")
	if ct == "" && (r.Body == nil || r.ContentLength == 0) {
		return nil
	}
	return matchesContentType(ct, mimeType)
}

func matchesContentType(contentType, expectedType string) error {
	mimetype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return errors.InvalidParameter(errors.Wrapf(err, "malformed Content-Type header (%s)", contentType))
	}
	if mimetype != expectedType {
		return errors.InvalidParameter(errors.Errorf("unsupported Content-Type header (%s): must be '%s'", contentType, expectedType))
	}
	return nil
}
