package jwt

import (
	"errors"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/authtoken"
)

const (
	DEFAULT_ACCESS_TOKEN_EXPIRATION  = 24 * time.Hour
	DEFAULT_REFRESH_TOKEN_EXPIRATION = 7 * 24 * time.Hour
	DEFAULT_ALGORITHM                = "HS256"
)

type JWT struct {
	privateKey string
	alg        string
	accexp     time.Duration
	refexp     time.Duration

	publicKey string
}

func NewJWT(privateKey, publicKey, alg string, accexp, refexp time.Duration) authtoken.AuthToken {
	if accexp <= 0 {
		accexp = DEFAULT_ACCESS_TOKEN_EXPIRATION
	}

	if refexp <= 0 {
		refexp = DEFAULT_REFRESH_TOKEN_EXPIRATION
	}

	if alg == "" {
		alg = DEFAULT_ALGORITHM
	}

	return &JWT{
		privateKey: privateKey,
		alg:        alg,
		accexp:     accexp,
		refexp:     refexp,
		publicKey:  publicKey,
	}
}

func (j *JWT) generateToken(claim authtoken.AuthTokenClaim, exp time.Duration, privateKey string) (string, error) {
	claim = claim.SetExp(exp)
	jwtClaims, ok := claim.(gojwt.Claims)
	if !ok {
		return "", errors.New("claim does not implement jwt.Claims interface")
	}

	newToken := gojwt.New(gojwt.GetSigningMethod(j.alg))
	newToken.Claims = jwtClaims

	return newToken.SignedString([]byte(privateKey))
}

func (j *JWT) Generate(claim authtoken.AuthTokenClaim) (string, string, error) {
	signedAccessToken, err := j.generateToken(claim, j.accexp, j.privateKey)
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err := j.generateToken(claim, j.refexp, j.privateKey)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func (j *JWT) Decode(t string, claim authtoken.AuthTokenClaim) (authtoken.AuthTokenClaim, error) {
	jwtClaims, ok := claim.(gojwt.Claims)
	if !ok {
		return nil, errors.New("claim does not implement gojwt.Claims")
	}

	tokenClaims, err := gojwt.ParseWithClaims(t, jwtClaims, func(token *gojwt.Token) (interface{}, error) {
		return []byte(j.publicKey), nil
	})
	if err != nil {
		if errors.Is(err, gojwt.ErrTokenExpired) {
			return nil, ErrExpiredToken // Return error spesifik expired jika ada
		}
		return nil, ErrUnclaimedToken
	}

	if tokenClaims == nil || tokenClaims.Claims == nil {
		return nil, ErrUnclaimedToken
	}

	tokClaims := tokenClaims.Claims.(authtoken.AuthTokenClaim)
	return tokClaims, nil
}
