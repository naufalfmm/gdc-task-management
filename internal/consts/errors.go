package consts

import "errors"

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrInvalidToken  = errors.New("invalid token")

	ErrDefaultError = errors.New("there is something error")
)

var (
	ErrorCodeBadRequest int = 400000

	ErrorCodeWrongPassword int = 401001
	ErrorCodeInvalidToken  int = 401002

	ErrorCodeInternalError int = 500000
	ErrorCodePanicRecover  int = 500999
)
