package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/infrastructures/rest/controllers/auth"
	"github.com/naufalfmm/gdc-task-management/internal/infrastructures/rest/controllers/tasks"
	"github.com/naufalfmm/gdc-task-management/internal/middlewares"
)

type Routes struct {
	authController auth.Controller
	taskController tasks.Controller

	middlewares middlewares.Middlewares
}

func NewRoutes(authController auth.Controller, taskController tasks.Controller, middl middlewares.Middlewares) Routes {
	return Routes{
		authController: authController,
		taskController: taskController,
		middlewares:    middl,
	}
}

func (r Routes) Register(ge *gin.Engine) {
	ge.RedirectFixedPath = true
	ge.Use(r.middlewares.PanicRecover(), r.middlewares.ImplementCors())

	auth := ge.Group("/auth")
	auth.POST("/register", r.authController.Register)
	auth.POST("/login", r.authController.Login)

	task := ge.Group("/tasks")
	task.POST("", r.middlewares.VerifyToken(), r.taskController.CreateTask)
	task.GET("/:id", r.middlewares.VerifyToken(), r.taskController.GetTaskDetail)
	task.POST("/:id/assign", r.middlewares.VerifyToken(), r.taskController.AssignTask)
}
