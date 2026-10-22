package consts

import (
	"errors"
	"fmt"
)

type errNotFound struct {
	domain string
	data   any
	err    error
}

func (e errNotFound) Error() string {
	return fmt.Errorf("%s %v %s", e.domain, e.data, e.err.Error()).Error()
}

func (e errNotFound) Unwrap() error {
	return e.err
}

func ErrNotFound(domain string, data any) error {
	return errNotFound{
		domain: domain,
		data:   data,
		err:    ErrNotFoundCmp,
	}
}

var (
	ErrWrongPassword error = errors.New("wrong password")
	ErrInvalidToken  error = errors.New("invalid token")
	ErrInvalidParam        = func(key string) error {
		return fmt.Errorf("invalid param %s", key)
	}

	ErrInvalidTeam = errors.New("invalid team")

	ErrNotFoundCmp = errors.New("not found")

	ErrDifferentTeam         error = errors.New("user and assignee don't have the same team")
	ErrMissingIdempotencyKey error = errors.New("idempotency key is missing")

	ErrDefaultError error = errors.New("there is something error")
)

var (
	ErrorCodeBadRequest            int = 400000
	ErrorCodeDifferentTeam         int = 400001
	ErrorCodeMissingIdempotencyKey int = 400002

	ErrorCodeWrongPassword int = 401001
	ErrorCodeInvalidToken  int = 401002

	ErrorCodeInvalidTeam int = 403001

	ErrorCodeNotFound int = 404000

	ErrorCodeInternalError int = 500000
	ErrorCodePanicRecover  int = 500999
)
