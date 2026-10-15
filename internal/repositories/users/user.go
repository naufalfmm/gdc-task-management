package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/naufalfmm/gdc-task-management/internal/domains/users"
	"github.com/naufalfmm/gdc-task-management/internal/models/daos"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
)

type Repository interface {
	InsertUser(ctx context.Context, user daos.User) (daos.User, error)
	GetUserByEmail(ctx context.Context, email string) (daos.User, error)
	GetUserByID(ctx context.Context, id int64) (daos.User, error)
}

type repository struct {
	querier users.Querier
}

func NewRepository(querier users.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func (r *repository) InsertUser(ctx context.Context, user daos.User) (daos.User, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = users.New(tx)
	}

	u, err := q.InsertUser(ctx, users.InsertUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
		TeamID:       pgtype.Int8{Int64: user.TeamID, Valid: true},
	})
	if err != nil {
		return daos.User{}, err
	}

	return daos.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		TeamID:       u.TeamID.Int64,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (daos.User, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = users.New(tx)
	}

	u, err := q.GetUserByEmail(ctx, email)
	if err != nil {
		return daos.User{}, err
	}

	return daos.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		TeamID:       u.TeamID.Int64,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}, nil
}

func (r *repository) GetUserByID(ctx context.Context, id int64) (daos.User, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = users.New(tx)
	}

	u, err := q.GetUserByID(ctx, id)
	if err != nil {
		return daos.User{}, err
	}

	return daos.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		TeamID:       u.TeamID.Int64,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}, nil
}
