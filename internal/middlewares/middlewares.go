package middlewares

import "github.com/gin-gonic/gin"

type Middlewares interface {
	PanicRecover() gin.HandlerFunc
	ImplementCors() gin.HandlerFunc
}

type middlewares struct {
}

func NewMiddleware() Middlewares {
	return &middlewares{}
}
