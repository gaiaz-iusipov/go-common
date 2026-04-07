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
	if target, ok := errors.AsType[wrappedError](err); ok {
		return target.Code
	}
	return http.StatusInternalServerError
}

var _ error = (*wrappedError)(nil)

type wrappedError struct {
	Err  error
	Code int
}

func (e wrappedError) Error() string {
	return e.Err.Error()
}
