package tasks

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/consts"
	"github.com/naufalfmm/gdc-task-management/internal/models/dtos"
	"github.com/naufalfmm/gdc-task-management/internal/usecases/tasks"
)

type Controller struct {
	taskUsecase tasks.Usecase
}

func NewTaskController(taskUsecase tasks.Usecase) Controller {
	return Controller{
		taskUsecase: taskUsecase,
	}
}

func (c Controller) CreateTask(gc *gin.Context) {
	var req dtos.CreateTaskRequest
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

	resp, err := c.taskUsecase.CreateTask(gc.Request.Context(), req)
	if err != nil {
		c.buildErrorCreateTask(gc, err)
		return
	}

	gc.JSON(http.StatusOK, dtos.Response[dtos.CreateTaskResponse]{
		Status:    true,
		Message:   "Success",
		Timestamp: time.Now(),
		Data:      resp,
	})
}

func (c Controller) buildErrorCreateTask(gc *gin.Context, err error) {
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

func (c Controller) AssignTask(gc *gin.Context) {
	var req dtos.AssignTaskRequest
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

	resp, err := c.taskUsecase.AssignTask(gc.Request.Context(), req)
	if err != nil {
		c.buildErrorForAssignTask(gc, err)
		return
	}

	gc.JSON(http.StatusOK, dtos.Response[dtos.AssignTaskResponse]{
		Status:    true,
		Message:   "Success",
		Timestamp: time.Now(),
		Data:      resp,
	})
}

func (c Controller) buildErrorForAssignTask(gc *gin.Context, err error) {
	statusCode := http.StatusInternalServerError

	if errors.Is(err, consts.ErrDifferentTeam) {
		gc.JSON(http.StatusBadRequest, dtos.Response[dtos.ErrorResponse]{
			Status:    false,
			Code:      consts.ErrorCodeDifferentTeam,
			Message:   err.Error(),
			Timestamp: time.Now(),
			Data: dtos.ErrorResponse{
				Error: err,
			},
		})
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

func (c Controller) GetTaskDetail(gc *gin.Context) {
	var req dtos.GetTaskDetailRequest
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

	resp, err := c.taskUsecase.GetTaskByID(gc.Request.Context(), req)
	if err != nil {
		c.buildErrorForAssignTask(gc, err)
		return
	}

	gc.JSON(http.StatusOK, dtos.Response[dtos.GetTaskDetailResponse]{
		Status:    true,
		Message:   "Success",
		Timestamp: time.Now(),
		Data:      resp,
	})
}

func (c Controller) buildErrorForGetTaskDetail(gc *gin.Context, err error) {
	statusCode := http.StatusInternalServerError

	if errors.Is(err, consts.ErrInvalidTeam) {
		gc.JSON(http.StatusForbidden, dtos.Response[dtos.ErrorResponse]{
			Status:    false,
			Code:      consts.ErrorCodeInvalidTeam,
			Message:   err.Error(),
			Timestamp: time.Now(),
			Data: dtos.ErrorResponse{
				Error: err,
			},
		})
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
