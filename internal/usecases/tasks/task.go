package tasks

import (
	"context"
	"time"

	"github.com/naufalfmm/gdc-task-management/config"
	"github.com/naufalfmm/gdc-task-management/internal/models/daos"
	"github.com/naufalfmm/gdc-task-management/internal/models/dtos"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
	"github.com/naufalfmm/gdc-task-management/internal/repositories/idempotencykey"
	"github.com/naufalfmm/gdc-task-management/internal/repositories/tasks"
)

type Usecase interface {
	CreateTask(ctx context.Context, r dtos.CreateTaskRequest) (dtos.CreateTaskResponse, error)
}

type usecase struct {
	txManager                  txmanager.TxManager
	taskRepositories           tasks.Repository
	idempotencyKeyRepositories idempotencykey.Repository

	conf config.Config
}

func NewTaskUsecase(txManager txmanager.TxManager, taskRepositories tasks.Repository, idempotencyKeyRepositories idempotencykey.Repository, conf config.Config) Usecase {
	return &usecase{
		txManager:                  txManager,
		taskRepositories:           taskRepositories,
		idempotencyKeyRepositories: idempotencyKeyRepositories,
		conf:                       conf,
	}
}

func (u usecase) CreateTask(ctx context.Context, r dtos.CreateTaskRequest) (dtos.CreateTaskResponse, error) {
	tx, ctx, err := u.txManager.Begin(ctx)
	if err != nil {
		return dtos.CreateTaskResponse{}, err
	}

	defer tx.Rollback(ctx)

	task, err := u.getOrCreateTask(ctx, r)
	if err != nil {
		return dtos.CreateTaskResponse{}, err
	}

	tx.Commit(ctx)

	return dtos.CreateTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
	}, nil
}

func (u usecase) getOrCreateTask(ctx context.Context, r dtos.CreateTaskRequest) (daos.Task, error) {
	idemKey, err := u.idempotencyKeyRepositories.GetValidIdempotencyKey(ctx, r.IdempotencyKey)
	if err == nil && idemKey.Type == daos.IdempotencyKeyTypeTask {
		task, err := u.taskRepositories.GetTaskByID(ctx, idemKey.TypeID)
		if err != nil {
			return daos.Task{}, err
		}

		return task, nil
	}

	task, err := u.taskRepositories.CreateTask(ctx, daos.Task{
		Title:       r.Title,
		Description: r.Description,
		Status:      daos.TaskStatusToDo,
		UserID:      r.UserAuth.ID,
	})
	if err != nil {
		return daos.Task{}, err
	}

	if _, err := u.idempotencyKeyRepositories.CreateIdempotencyKey(ctx, daos.IdempotencyKey{
		Key:       r.IdempotencyKey,
		Type:      daos.IdempotencyKeyTypeTask,
		TypeID:    task.ID,
		ExpiresAt: time.Now().Add(u.conf.IdempotencyKeyExpiration),
	}); err != nil {
		return daos.Task{}, err
	}

	return task, nil
}
