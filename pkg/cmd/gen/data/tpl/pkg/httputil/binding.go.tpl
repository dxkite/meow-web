package httputil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"{{ .PackageName }}/pkg/errors"
	"{{ .PackageName }}/pkg/httputil/form"
	"github.com/go-playground/validator/v10"
)

const (
	MIMEJSON              = "application/json"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

var MultipartFormMemory int64 = 32 << 20

var Validator = validator.New()
var initValidator = &sync.Once{}

func init() {
	initValidator.Do(func() {
		// 使用自定义校验
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
	err := checkContentType(r, MIMEJSON)
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

func ReadQuery(ctx context.Context, r *http.Request, out interface{}) error {
	return form.MappingForm(form.NewForm(r.URL.Query()), out)
}

func ReadForm(ctx context.Context, r *http.Request, out interface{}) error {
	contentType := r.Header.Get("Content-Type")
	mineType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return errors.InvalidParameter(errors.Wrapf(err, "malformed Content-Type header (%s)", contentType))
	}
	var val form.Value
	switch mineType {
	case MIMEMultipartPOSTForm:
		if err := r.ParseMultipartForm(MultipartFormMemory); err != nil {
			return errors.InvalidParameter(errors.Wrapf(err, "malformed body"))
		}
		val = form.NewMultipartForm(r.MultipartForm)
	case MIMEPOSTForm:
		if err := r.ParseForm(); err != nil {
			return errors.InvalidParameter(errors.Wrapf(err, "malformed body"))
		}
		val = form.NewForm(r.Form)
	default:
		return errors.InvalidParameter(errors.Wrapf(err, "unsupported Content-Type header (%s): must be '%s' or '%s'", contentType, MIMEMultipartPOSTForm, MIMEPOSTForm))
	}
	return form.MappingForm(val, out)
}

func ReadRequest(ctx context.Context, r *http.Request, out interface{}) error {
	if r.ContentLength == 0 {
		return ReadQuery(ctx, r, out)
	}
	contentType := r.Header.Get("Content-Type")
	mineType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return errors.InvalidParameter(errors.Wrapf(err, "malformed Content-Type header (%s)", contentType))
	}
	switch mineType {
	case MIMEJSON:
		return ReadJSON(ctx, r, out)
	case MIMEMultipartPOSTForm, MIMEPOSTForm:
		return ReadForm(ctx, r, out)
	}
	return errors.InvalidParameter(fmt.Errorf("invalid Content-Type: %s", contentType))
}

func Validate(ctx context.Context, input interface{}) error {
	if err := Validator.Struct(input); err != nil {
		return errors.InvalidParameter(errors.Wrap(err, "validate error"))
	}
	return nil
}
