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
	AssignTaskByID(ctx context.Context, id int64, assigneeID int64) (daos.Task, error)
}

type repository struct {
	querier tasks.Querier
}

func NewRepository(querier tasks.Querier) Repository {
	return &repository{
		querier: querier,
	}
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

	createdTask := daos.Task{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description.String,
		Status:      t.Status,
		UserID:      t.UserID.Int64,
		CreatedAt:   t.CreatedAt.Time,
	}

	if t.AssigneeID.Valid {
		createdTask.AssigneeID = &t.AssigneeID.Int64
	}

	if t.UpdatedAt.Valid {
		createdTask.UpdatedAt = &t.UpdatedAt.Time
	}

	return createdTask, nil
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

	task := daos.Task{
		ID:          t.ID,
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

	return task, nil
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

	task := daos.Task{
		ID:          t.ID,
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

	return task, nil
}
