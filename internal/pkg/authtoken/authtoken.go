package authtoken

import (
	"time"
)

type Data[T any] interface {
	CreatedBy() T
}

type AuthTokenClaim interface {
	Valid() error
	SetExp(exp time.Duration) AuthTokenClaim
}

type AuthToken interface {
	Generate(claim AuthTokenClaim) (accessToken string, refreshToken string, err error)
	Decode(t string, claim AuthTokenClaim) (AuthTokenClaim, error)
}
