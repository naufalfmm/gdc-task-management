package dtos

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/errors"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`

	IdempotencyKey string   `json:"-"`
	UserAuth       UserAuth `json:"-"`
}

func (r *CreateTaskRequest) FromGinContext(gc *gin.Context) error {
	if err := gc.ShouldBindJSON(r); err != nil {
		return errors.GenerateError(*r, err)
	}

	return nil
}

type CreateTaskResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
