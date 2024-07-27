package httputil

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"dxkite.cn/meownest/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

var Validator = validator.New()
var initValidator = &sync.Once{}

func init() {
	initValidator.Do(func() {
		Validator.SetTagName("binding")
		Validator.RegisterTagNameFunc(func(field reflect.StructField) string {
			if name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]; name != "" && name != "-" {
				return name
			}
			if name := strings.SplitN(field.Tag.Get("form"), ",", 2)[0]; name != "" && name != "-" {
				return name
			}
			return ""
		})
	})
}

func ReadJSON(ctx context.Context, r *http.Request, out interface{}) error {
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

	return nil
}

func ReadQuery(ctx context.Context, r *http.Request, out interface{}) error {
	if err := Bind("query", simplifyMap(r.URL.Query()), out); err != nil {
		return errors.InvalidParameter(errors.Wrap(err, "read query error"))
	}
	return nil
}

func ReadForm(ctx context.Context, r *http.Request, out interface{}) error {
	if err := r.ParseForm(); err != nil {
		return errors.InvalidParameter(errors.Wrap(err, "read form error"))
	}
	if err := Bind("form", simplifyMap(r.Form), out); err != nil {
		return errors.InvalidParameter(errors.Wrap(err, "read query error"))
	}
	return nil
}

func Validate(ctx context.Context, input interface{}) error {
	if err := Validator.Struct(input); err != nil {
		return errors.InvalidParameter(errors.Wrap(err, "validate error"))
	}
	return nil
}

func simplifyMap(values map[string][]string) map[string]string {
	newMap := map[string]string{}
	for k, v := range values {
		newMap[k] = v[len(v)-1]
	}
	return newMap
}

func Bind(name string, values map[string]string, out interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: name, Result: out})
	if err != nil {
		return errors.System(err)
	}
	if err := dec.Decode(values); err != nil {
		return errors.InvalidParameter(err)
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
