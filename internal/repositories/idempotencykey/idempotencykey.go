package idempotencykey

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/naufalfmm/gdc-task-management/internal/domains/idempotencykeys"
	"github.com/naufalfmm/gdc-task-management/internal/models/daos"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
)

type Repository interface {
	CreateIdempotencyKey(ctx context.Context, idem daos.IdempotencyKey) (daos.IdempotencyKey, error)
	GetValidIdempotencyKey(ctx context.Context, key string) (daos.IdempotencyKey, error)
}

type repository struct {
	querier idempotencykeys.Querier
}

func NewRepository(querier idempotencykeys.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func (r *repository) CreateIdempotencyKey(ctx context.Context, idem daos.IdempotencyKey) (daos.IdempotencyKey, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = idempotencykeys.New(tx)
	}

	ik, err := q.CreateIdempotencyKey(ctx, idempotencykeys.CreateIdempotencyKeyParams{
		Key:    idem.Key,
		Type:   idem.Type,
		TypeID: idem.TypeID,
		ExpiresAt: pgtype.Timestamp{
			Time:  idem.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return daos.IdempotencyKey{}, err
	}

	return daos.IdempotencyKey{
		Key:       ik.Key,
		Type:      ik.Type,
		TypeID:    ik.TypeID,
		ExpiresAt: ik.ExpiresAt.Time,
		CreatedAt: ik.CreatedAt.Time,
	}, nil
}

func (r *repository) GetValidIdempotencyKey(ctx context.Context, key string) (daos.IdempotencyKey, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = idempotencykeys.New(tx)
	}

	ik, err := q.GetValidIdempotencyKey(ctx, key)
	if err != nil {
		return daos.IdempotencyKey{}, err
	}

	return daos.IdempotencyKey{
		Key:       ik.Key,
		Type:      ik.Type,
		TypeID:    ik.TypeID,
		ExpiresAt: ik.ExpiresAt.Time,
		CreatedAt: ik.CreatedAt.Time,
	}, nil
}
