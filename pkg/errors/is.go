package errors

func IsNotFound(err error) bool {
	var target ErrNotFound
	return As(err, &target)
}

func IsInvalidParameter(err error) bool {
	var target ErrInvalidParameter
	return As(err, &target)
}

func IsUnauthorized(err error) bool {
	var target ErrUnauthorized
	return As(err, &target)
}

func IsUnavailable(err error) bool {
	var target ErrUnavailable
	return As(err, &target)
}

func IsForbidden(err error) bool {
	var target ErrForbidden
	return As(err, &target)
}

func IsUnprocessableEntity(err error) bool {
	var target ErrUnprocessableEntity
	return As(err, &target)
}

func IsSystem(err error) bool {
	var target ErrSystem
	return As(err, &target)
}

func IsUnknown(err error) bool {
	var target ErrUnknown
	return As(err, &target)
}
