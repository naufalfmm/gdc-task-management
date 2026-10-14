package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/infrastructures/rest/controllers/auth"
	"github.com/naufalfmm/gdc-task-management/internal/middlewares"
)

type Routes struct {
	authController auth.Controller

	middlewares middlewares.Middlewares
}

func NewRoutes(authController auth.Controller, middl middlewares.Middlewares) Routes {
	return Routes{
		authController: authController,
		middlewares:    middl,
	}
}

func (r Routes) Register(ge *gin.Engine) {
	ge.RedirectFixedPath = true
	ge.Use(r.middlewares.PanicRecover(), r.middlewares.ImplementCors())

	auth := ge.Group("/auth")
	auth.POST("/register", r.authController.Register)
}
