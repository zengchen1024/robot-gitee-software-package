package repository

type ErrorResourceNotFound struct {
	error
}

func NewErrorResourceNotFound(err error) ErrorResourceNotFound {
	return ErrorResourceNotFound{err}
}

func IsErrorResourceNotFound(err error) bool {
	_, ok := err.(ErrorResourceNotFound)

	return ok
}
