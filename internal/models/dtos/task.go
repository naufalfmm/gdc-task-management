package dtos

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/consts"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/errors"
)

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`

	IdempotencyKey string   `json:"-"`
	UserAuth       UserAuth `json:"-"`
}

func (r *CreateTaskRequest) FromGinContext(gc *gin.Context) error {
	if err := gc.ShouldBindJSON(r); err != nil {
		return errors.GenerateError(*r, err)
	}

	idemKey := gc.GetHeader(consts.IdempotencyKeyHeader)
	if idemKey == "" {
		return consts.ErrMissingIdempotencyKey
	}

	r.IdempotencyKey = idemKey

	userHeader, _ := gc.Get(consts.XUserHeader)
	r.UserAuth = userHeader.(UserAuth)

	return nil
}

type CreateTaskResponse struct {
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type AssignTaskRequest struct {
	AssigneeID int64 `json:"assignee_id" validate:"required"`

	ID       int64    `json:"-"`
	UserAuth UserAuth `json:"-"`
}

func (r *AssignTaskRequest) FromGinContext(gc *gin.Context) error {
	if err := gc.ShouldBindJSON(r); err != nil {
		return errors.GenerateError(*r, err)
	}

	paramID := gc.Param("id")
	if paramID == "" {
		return consts.ErrInvalidParam("id")
	}

	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		return consts.ErrInvalidParam("id")
	}

	r.ID = id

	userHeader, _ := gc.Get(consts.XUserHeader)
	r.UserAuth = userHeader.(UserAuth)

	return nil
}

type AssignTaskResponse struct {
	Code       string `json:"code"`
	UserID     int64  `json:"user_id"`
	AssigneeID int64  `json:"assignee_id"`
}

func (r AssignTaskResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(r)
}

type GetTaskDetailRequest struct {
	Code     string   `json:"-"`
	UserAuth UserAuth `json:"-"`
}

func (r *GetTaskDetailRequest) FromGinContext(gc *gin.Context) error {
	r.Code = gc.Param("code")
	if r.Code == "" {
		return consts.ErrInvalidParam("code")
	}

	userHeader, _ := gc.Get(consts.XUserHeader)
	r.UserAuth = userHeader.(UserAuth)

	return nil
}

type GetTaskDetailResponse struct {
	Code        string     `json:"code"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	UserID      int64      `json:"user_id"`
	TeamID      int64      `json:"team_id"`
	AssigneeID  *int64     `json:"assignee_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
