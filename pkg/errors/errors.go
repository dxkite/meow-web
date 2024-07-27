package errors

type ErrNotFound interface {
	NotFound()
}

type ErrInvalidParameter interface {
	InvalidParameter()
}

type ErrUnprocessableEntity interface {
	UnprocessableEntity()
}

type ErrUnauthorized interface {
	Unauthorized()
}

type ErrUnavailable interface {
	Unavailable()
}

type ErrForbidden interface {
	Forbidden()
}

type ErrSystem interface {
	System()
}

type ErrUnknown interface {
	Unknown()
}
