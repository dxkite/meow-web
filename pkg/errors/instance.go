package errors

type customError interface {
	Unwrap() error
	customError()
}

type errNotFound struct{ error }

func (errNotFound) NotFound()     {}
func (*errNotFound) customError() {}
func (e errNotFound) Unwrap() error {
	return e.error
}

func NotFound(err error) error {
	if err == nil || IsNotFound(err) {
		return err
	}
	return errNotFound{err}
}

type errInvalidParameter struct{ error }

func (errInvalidParameter) InvalidParameter() {}
func (errInvalidParameter) customError()      {}
func (e errInvalidParameter) Unwrap() error {
	return e.error
}

func InvalidParameter(err error) error {
	if err == nil || IsInvalidParameter(err) {
		return err
	}
	return errInvalidParameter{err}
}

type errUnauthorized struct{ error }

func (errUnauthorized) Unauthorized() {}
func (errUnauthorized) customError()  {}
func (e errUnauthorized) Unwrap() error {
	return e.error
}

func Unauthorized(err error) error {
	if err == nil || IsUnauthorized(err) {
		return err
	}
	return errUnauthorized{err}
}

type errUnavailable struct{ error }

func (errUnavailable) Unavailable() {}
func (errUnavailable) customError() {}
func (e errUnavailable) Unwrap() error {
	return e.error
}

func Unavailable(err error) error {
	if err == nil || IsUnavailable(err) {
		return err
	}
	return errUnavailable{err}
}

type errForbidden struct{ error }

func (errForbidden) Forbidden()   {}
func (errForbidden) customError() {}
func (e errForbidden) Unwrap() error {
	return e.error
}

func Forbidden(err error) error {
	if err == nil || IsForbidden(err) {
		return err
	}
	return errForbidden{err}
}

type errUnprocessableEntity struct{ error }

func (errUnprocessableEntity) UnprocessableEntity() {}
func (errUnprocessableEntity) customError()         {}
func (e errUnprocessableEntity) Unwrap() error {
	return e.error
}

func UnprocessableEntity(err error) error {
	if err == nil || IsUnprocessableEntity(err) {
		return err
	}
	return errUnprocessableEntity{err}
}

type errSystem struct{ error }

func (errSystem) System()      {}
func (errSystem) customError() {}
func (e errSystem) Unwrap() error {
	return e.error
}

func System(err error) error {
	if err == nil || IsSystem(err) {
		return err
	}
	return errSystem{err}
}

type errUnknown struct{ error }

func (errUnknown) Unknown()     {}
func (errUnknown) customError() {}
func (e errUnknown) Unwrap() error {
	return e.error
}

func Unknown(err error) error {
	if err == nil || IsUnknown(err) {
		return err
	}
	return errUnknown{err}
}
