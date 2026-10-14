package dtos

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/authtoken"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/validator"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	TeamName string `json:"team_name" validate:"required"`
}

func (r *RegisterRequest) ValidatePassword() error {
	return nil
}

func (r *RegisterRequest) FromGinContext(gc *gin.Context) error {
	if err := gc.ShouldBindJSON(r); err != nil {
		return validator.NewFromValidationErrors(*r, err)
	}

	return nil
}

type UserAuth struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginClaims struct {
	*jwt.RegisteredClaims
	UserAuth UserAuth `json:"user"`
}

func (l *LoginClaims) Valid() error {
	return nil
}

func (l *LoginClaims) SetExp(exp time.Duration) authtoken.AuthTokenClaim {
	if l.RegisteredClaims == nil {
		l.RegisteredClaims = &jwt.RegisteredClaims{}
	}
	l.ExpiresAt = jwt.NewNumericDate(time.Now().Add(exp))

	return l
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	Email    string        `json:"email"`
	Name     string        `json:"name"`
	TeamName string        `json:"team_name"`
	Token    TokenResponse `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequest) FromGinContext(gc *gin.Context) error {
	if err := gc.ShouldBindJSON(r); err != nil {
		return validator.NewFromValidationErrors(*r, err)
	}

	return nil
}
