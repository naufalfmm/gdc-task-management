package auth

import (
	"errors"
	"net/http"
	"time"

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
		gc.JSON(http.StatusBadRequest, dtos.Response[dtos.ErrorResponse]{
			Status:    false,
			Code:      consts.ErrorCodeBadRequest,
			Message:   err.Error(),
			Timestamp: time.Now(),
			Data: dtos.ErrorResponse{
				Error: err,
			},
		})

		return
	}

	resp, err := c.authUsecase.RegisterUser(gc.Request.Context(), req)
	if err != nil {
		c.buildErrorRegister(gc, err)
		return
	}

	gc.JSON(http.StatusOK, dtos.Response[dtos.AuthResponse]{
		Status:    true,
		Message:   "Success",
		Timestamp: time.Now(),
		Data:      resp,
	})
}

func (c Controller) buildErrorRegister(gc *gin.Context, err error) {
	statusCode := http.StatusInternalServerError

	gc.JSON(statusCode, dtos.Response[dtos.ErrorResponse]{
		Status:    false,
		Code:      consts.ErrorCodeInternalError,
		Message:   consts.ErrDefaultError.Error(),
		Timestamp: time.Now(),
		Data: dtos.ErrorResponse{
			Error: consts.ErrDefaultError,
		},
	})
}

func (c Controller) Login(gc *gin.Context) {
	var req dtos.LoginRequest
	if err := req.FromGinContext(gc); err != nil {
		gc.JSON(http.StatusBadRequest, dtos.Response[dtos.ErrorResponse]{
			Status:    false,
			Code:      consts.ErrorCodeBadRequest,
			Message:   err.Error(),
			Timestamp: time.Now(),
			Data: dtos.ErrorResponse{
				Error: err,
			},
		})

		return
	}

	resp, err := c.authUsecase.Login(gc.Request.Context(), req)
	if err != nil {
		c.buildErrorLogin(gc, err)
		return
	}

	gc.JSON(http.StatusOK, dtos.Response[dtos.AuthResponse]{
		Status:    true,
		Message:   "Success",
		Timestamp: time.Now(),
		Data:      resp,
	})
}

func (c Controller) buildErrorLogin(gc *gin.Context, err error) {
	statusCode := http.StatusInternalServerError

	if errors.Is(err, consts.ErrWrongPassword) {
		gc.JSON(statusCode, dtos.Response[dtos.ErrorResponse]{
			Status:    false,
			Code:      consts.ErrorCodeWrongPassword,
			Message:   err.Error(),
			Timestamp: time.Now(),
			Data: dtos.ErrorResponse{
				Error: err,
			},
		})
		return
	}

	gc.JSON(statusCode, dtos.Response[dtos.ErrorResponse]{
		Status:    false,
		Code:      consts.ErrorCodeInternalError,
		Message:   consts.ErrDefaultError.Error(),
		Timestamp: time.Now(),
		Data: dtos.ErrorResponse{
			Error: consts.ErrDefaultError,
		},
	})
}
