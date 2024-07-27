package httputil

import (
	"net/http"

	"dxkite.cn/meownest/pkg/errors"
)

func statusFromError(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	switch {
	case errors.IsNotFound(err):
		return http.StatusNotFound
	case errors.IsInvalidParameter(err):
		return http.StatusBadRequest
	case errors.IsUnauthorized(err):
		return http.StatusUnauthorized
	case errors.IsUnavailable(err):
		return http.StatusServiceUnavailable
	case errors.IsForbidden(err):
		return http.StatusForbidden
	case errors.IsUnprocessableEntity(err):
		return http.StatusUnprocessableEntity
	case errors.IsSystem(err) || errors.IsUnknown(err):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func nameFromError(err error) string {
	if err == nil {
		return "UnknownError"
	}

	switch {
	case errors.IsNotFound(err):
		return "NotFound"
	case errors.IsInvalidParameter(err):
		return "InvalidParameter"
	case errors.IsUnauthorized(err):
		return "Unauthorized"
	case errors.IsUnavailable(err):
		return "Unavailable"
	case errors.IsForbidden(err):
		return "Forbidden"
	case errors.IsSystem(err):
		return "SystemError"
	case errors.IsUnknown(err):
		return "UnknownError"
	default:
		return "UnknownError"
	}
}
