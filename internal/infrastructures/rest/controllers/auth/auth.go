package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/consts"
	"github.com/naufalfmm/gdc-task-management/internal/models/dtos"
	"github.com/naufalfmm/gdc-task-management/internal/usecases/auth"
)

type Controller struct {
	authUsecase auth.Usecase
}

func NewAuthController(authUsecase auth.Usecase) Controller {
	return Controller{
		authUsecase: authUsecase,
	}
}

func (c Controller) Register(gc *gin.Context) {
	var req dtos.RegisterRequest
	if err := req.FromGinContext(gc); err != nil {
		gc.JSON(http.StatusBadRequest, dtos.Response[any]{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
			Error:   &err,
		})

		return
	}

	resp, err := c.authUsecase.RegisterUser(gc.Request.Context(), req)
	if err != nil {
		c.buildErrorRegister(gc, err)
		return
	}

	gc.JSON(http.StatusOK, dtos.Response[dtos.AuthResponse]{
		Code: "SUCCESS",
		Data: resp,
	})
}

func (c Controller) buildErrorRegister(gc *gin.Context, err error) {
	statusCode := http.StatusInternalServerError

	gc.JSON(statusCode, dtos.Response[error]{
		Code:    "ANY_ERROR",
		Message: err.Error(),
		Error:   &err,
	})
}

func (c Controller) Login(gc *gin.Context) {
	var req dtos.LoginRequest
	if err := req.FromGinContext(gc); err != nil {
		gc.JSON(http.StatusBadRequest, dtos.Response[any]{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
			Error:   &err,
		})

		return
	}

	resp, err := c.authUsecase.Login(gc.Request.Context(), req)
	if err != nil {
		c.buildErrorLogin(gc, err)
		return
	}

	gc.JSON(http.StatusOK, dtos.Response[dtos.AuthResponse]{
		Code: "SUCCESS",
		Data: resp,
	})
}

func (c Controller) buildErrorLogin(gc *gin.Context, err error) {
	statusCode := http.StatusInternalServerError

	if errors.Is(err, consts.ErrWrongPassword) {
		gc.JSON(http.StatusUnauthorized, dtos.Response[error]{
			Code:    "WRONG_PASSWORD",
			Message: err.Error(),
			Error:   &err,
		})
		return
	}

	gc.JSON(statusCode, dtos.Response[error]{
		Code:    "ANY_ERROR",
		Message: err.Error(),
		Error:   &err,
	})
}
