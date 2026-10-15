package consts

import (
	"errors"
	"fmt"
)

var (
	ErrWrongPassword error = errors.New("wrong password")
	ErrInvalidToken  error = errors.New("invalid token")
	ErrInvalidParam        = func(key string) error {
		return fmt.Errorf("invalid param %s", key)
	}

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

	ErrorCodeInternalError int = 500000
	ErrorCodePanicRecover  int = 500999
)
