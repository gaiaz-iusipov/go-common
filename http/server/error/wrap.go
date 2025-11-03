package httpservererror

import (
	"errors"
	"net/http"
)

func New(text string, code int) error {
	return Wrap(errors.New(text), code)
}

func Wrap(err error, code int) error {
	return wrappedError{
		Err:  err,
		Code: code,
	}
}

func Unwrap(err error) int {
	var target wrappedError
	if errors.As(err, &target) {
		return target.Code
	}
	return http.StatusInternalServerError
}

type wrappedError struct {
	Err  error
	Code int
}

func (e wrappedError) Error() string {
	return e.Err.Error()
}
