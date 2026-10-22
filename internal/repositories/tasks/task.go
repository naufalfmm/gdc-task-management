package tasks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/naufalfmm/gdc-task-management/internal/domains/tasks"
	"github.com/naufalfmm/gdc-task-management/internal/models/daos"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
)

type Repository interface {
	CreateTask(ctx context.Context, task daos.Task) (daos.Task, error)
	GetTaskByID(ctx context.Context, id int64) (daos.Task, error)
	GetTaskByCode(ctx context.Context, code string) (daos.Task, error)
	AssignTaskByID(ctx context.Context, id int64, assigneeID int64) (daos.Task, error)
	GetLatestActiveTeamTask(ctx context.Context, teamID int64) (daos.Task, error)
	UpdateTask(ctx context.Context, updatedTask daos.Task) (daos.Task, error)
	DeleteTask(ctx context.Context, deletedTask daos.Task) (daos.Task, error)
}

type repository struct {
	querier tasks.Querier
}

func NewRepository(querier tasks.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func newDaoTask(t tasks.Task) daos.Task {
	task := daos.Task{
		ID:          t.ID,
		Code:        t.Code,
		Title:       t.Title,
		Description: t.Description.String,
		Status:      t.Status,
		UserID:      t.UserID.Int64,
		CreatedAt:   t.CreatedAt.Time,
	}

	if t.AssigneeID.Valid {
		task.AssigneeID = &t.AssigneeID.Int64
	}

	if t.UpdatedAt.Valid {
		task.UpdatedAt = &t.UpdatedAt.Time
	}

	if t.DeletedAt.Valid {
		task.DeletedAt = &t.DeletedAt.Time
	}

	return task
}

func (r *repository) CreateTask(ctx context.Context, task daos.Task) (daos.Task, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = tasks.New(tx)
	}

	createTaskParams := tasks.CreateTaskParams{
		Title: task.Title,
		Description: pgtype.Text{
			String: task.Description,
			Valid:  true,
		},
		Status: task.Status,
		UserID: pgtype.Int8{
			Int64: task.UserID,
			Valid: true,
		},
		AssigneeID: pgtype.Int8{
			Valid: false,
		},
	}

	if task.AssigneeID != nil {
		createTaskParams.AssigneeID = pgtype.Int8{
			Int64: *task.AssigneeID,
			Valid: true,
		}
	}

	t, err := q.CreateTask(ctx, createTaskParams)
	if err != nil {
		return daos.Task{}, nil
	}

	return newDaoTask(t), nil
}

func (r *repository) GetTaskByID(ctx context.Context, id int64) (daos.Task, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = tasks.New(tx)
	}

	t, err := q.GetTaskByID(ctx, id)
	if err != nil {
		return daos.Task{}, err
	}

	return newDaoTask(t), nil
}

func (r *repository) GetTaskByCode(ctx context.Context, code string) (daos.Task, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = tasks.New(tx)
	}

	t, err := q.GetTaskByCode(ctx, code)
	if err != nil {
		return daos.Task{}, err
	}

	return newDaoTask(t), nil
}

func (r *repository) AssignTaskByID(ctx context.Context, id int64, assigneeID int64) (daos.Task, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = tasks.New(tx)
	}

	t, err := q.AssignTask(ctx, tasks.AssignTaskParams{
		ID: id,
		AssigneeID: pgtype.Int8{
			Int64: assigneeID,
			Valid: true,
		},
	})
	if err != nil {
		return daos.Task{}, err
	}

	return newDaoTask(t), nil
}

func (r *repository) GetLatestActiveTeamTask(ctx context.Context, teamID int64) (daos.Task, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = tasks.New(tx)
	}

	t, err := q.GetLatestActiveTeamTask(ctx, pgtype.Int8{Int64: teamID, Valid: true})
	if err != nil {
		return daos.Task{}, err
	}

	return newDaoTask(t), nil
}

func (r *repository) UpdateTask(ctx context.Context, updatedTask daos.Task) (daos.Task, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = tasks.New(tx)
	}

	t, err := q.UpdateTask(ctx, tasks.UpdateTaskParams{
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      updatedTask.Status,

		ID:   updatedTask.ID,
		Code: updatedTask.Code,
	})
	if err != nil {
		return daos.Task{}, err
	}

	return newDaoTask(t), nil
}

func (r *repository) DeleteTask(ctx context.Context, deletedTask daos.Task) (daos.Task, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = tasks.New(tx)
	}

	t, err := q.DeleteTask(ctx, tasks.DeleteTaskParams{
		ID:   deletedTask.ID,
		Code: deletedTask.Code,
	})
	if err != nil {
		return daos.Task{}, err
	}

	return newDaoTask(t), nil
}
