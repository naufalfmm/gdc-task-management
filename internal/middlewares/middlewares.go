package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/authtoken"
)

type Middlewares interface {
	PanicRecover() gin.HandlerFunc
	ImplementCors() gin.HandlerFunc
	VerifyToken() gin.HandlerFunc
}

type middlewares struct {
	authToken authtoken.AuthToken
}

func NewMiddleware(authToken authtoken.AuthToken) Middlewares {
	return &middlewares{
		authToken: authToken,
	}
}
