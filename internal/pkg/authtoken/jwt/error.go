package jwt

import "errors"

var (
	ErrUnclaimedToken = errors.New("unclaimed token")
	ErrExpiredToken   = errors.New("expired token")
)
