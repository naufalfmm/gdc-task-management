package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/naufalfmm/gdc-task-management/config"
	"github.com/naufalfmm/gdc-task-management/internal/consts"
	"github.com/naufalfmm/gdc-task-management/internal/models/daos"
	"github.com/naufalfmm/gdc-task-management/internal/models/dtos"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/notification"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
	"github.com/naufalfmm/gdc-task-management/internal/repositories/idempotencykeys"
	"github.com/naufalfmm/gdc-task-management/internal/repositories/tasks"
	"github.com/naufalfmm/gdc-task-management/internal/repositories/teams"
	"github.com/naufalfmm/gdc-task-management/internal/repositories/users"
)

type Usecase interface {
	CreateTask(ctx context.Context, r dtos.CreateTaskRequest) (dtos.CreateTaskResponse, error)
	AssignTask(ctx context.Context, r dtos.AssignTaskRequest) (dtos.AssignTaskResponse, error)
	GetTaskDetail(ctx context.Context, r dtos.GetTaskDetailRequest) (dtos.GetTaskDetailResponse, error)
}

type usecase struct {
	txManager    txmanager.TxManager
	notification notification.Notification

	taskRepositories           tasks.Repository
	idempotencyKeyRepositories idempotencykeys.Repository
	userRepositories           users.Repository
	teamRepositories           teams.Repository

	conf config.Config
}

func NewTaskUsecase(txManager txmanager.TxManager, notification notification.Notification, taskRepositories tasks.Repository, idempotencyKeyRepositories idempotencykeys.Repository, userRepositories users.Repository, teamRepositories teams.Repository, conf config.Config) Usecase {
	return &usecase{
		txManager:    txManager,
		notification: notification,

		taskRepositories:           taskRepositories,
		idempotencyKeyRepositories: idempotencyKeyRepositories,
		userRepositories:           userRepositories,
		teamRepositories:           teamRepositories,

		conf: conf,
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
		Code:        task.Code,
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

	creator, err := u.userRepositories.GetUserByID(ctx, r.UserAuth.ID)
	if err != nil {
		return daos.Task{}, err
	}

	taskCode, err := u.generateTeamCode(ctx, creator.TeamID)
	if err != nil {
		return daos.Task{}, err
	}

	task, err := u.taskRepositories.CreateTask(ctx, daos.Task{
		Title:       r.Title,
		Code:        taskCode,
		Description: r.Description,
		Status:      daos.TaskStatusToDo,
		UserID:      r.UserAuth.ID,
		TeamID:      creator.TeamID,
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

func (u usecase) generateTeamCode(ctx context.Context, teamID int64) (string, error) {
	latestTask, err := u.taskRepositories.GetLatestActiveTeamTask(ctx, teamID)
	if err != nil {
		return "", err
	}

	team, err := u.teamRepositories.GetTeamByID(ctx, teamID)
	if err != nil {
		return "", err
	}

	_, latestTaskNumber := latestTask.SplitCode()

	return fmt.Sprintf("%s-%d", team.Code, latestTaskNumber+1), nil
}

func (u usecase) AssignTask(ctx context.Context, r dtos.AssignTaskRequest) (dtos.AssignTaskResponse, error) {
	tx, ctx, err := u.txManager.Begin(ctx)
	if err != nil {
		return dtos.AssignTaskResponse{}, err
	}

	defer tx.Rollback(ctx)

	assignee, err := u.userRepositories.GetUserByID(ctx, r.AssigneeID)
	if err != nil {
		return dtos.AssignTaskResponse{}, err
	}

	user, err := u.userRepositories.GetUserByID(ctx, r.UserAuth.ID)
	if err != nil {
		return dtos.AssignTaskResponse{}, err
	}

	if user.TeamID != assignee.TeamID {
		return dtos.AssignTaskResponse{}, consts.ErrDifferentTeam
	}

	t, err := u.taskRepositories.AssignTaskByID(ctx, r.ID, r.AssigneeID)
	if err != nil {
		return dtos.AssignTaskResponse{}, err
	}

	resp := dtos.AssignTaskResponse{
		Code:       t.Code,
		UserID:     t.UserID,
		AssigneeID: *t.AssigneeID,
	}

	// send notification
	if err := u.notification.Send(ctx, "Assign Task", resp); err != nil {
		return dtos.AssignTaskResponse{}, err
	}

	tx.Commit(ctx)

	return resp, nil
}

func (u usecase) GetTaskDetail(ctx context.Context, r dtos.GetTaskDetailRequest) (dtos.GetTaskDetailResponse, error) {
	task, err := u.taskRepositories.GetTaskByCode(ctx, r.Code)
	if err != nil {
		return dtos.GetTaskDetailResponse{}, err
	}

	user, err := u.userRepositories.GetUserByID(ctx, r.UserAuth.ID)
	if err != nil {
		return dtos.GetTaskDetailResponse{}, err
	}

	if task.TeamID != user.TeamID {
		return dtos.GetTaskDetailResponse{}, consts.ErrInvalidTeam
	}

	return dtos.GetTaskDetailResponse{
		Code:        task.Code,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UserID:      task.UserID,
		TeamID:      task.TeamID,
		AssigneeID:  task.AssigneeID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}
