package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/authtoken/jwt"
)

type Middlewares interface {
	PanicRecover() gin.HandlerFunc
	ImplementCors() gin.HandlerFunc
	VerifyToken() gin.HandlerFunc
}

type middlewares struct {
	jwt jwt.JWT
}

func NewMiddleware(jwt jwt.JWT) Middlewares {
	return &middlewares{
		jwt: jwt,
	}
}
