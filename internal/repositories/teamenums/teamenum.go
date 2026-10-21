package teamenums

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/naufalfmm/gdc-task-management/internal/domains/teamenums"
	"github.com/naufalfmm/gdc-task-management/internal/models/daos"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
)

type Repository interface {
	GetTeamEnum(ctx context.Context, teamID int64, enum string) (daos.TeamEnum, error)
}

type repository struct {
	querier teamenums.Querier
}

func NewRepository(querier teamenums.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func (r repository) GetTeamEnum(ctx context.Context, teamID int64, enum string) (daos.TeamEnum, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = teamenums.New(tx)
	}

	te, err := q.GetTeamEnum(ctx, teamenums.GetTeamEnumParams{
		TeamID: pgtype.Int8{
			Valid: true,
			Int64: teamID,
		},
		Enum: enum,
	})
	if err != nil {
		return daos.TeamEnum{}, err
	}

	teamEnum := daos.TeamEnum{
		ID:        te.ID,
		TeamID:    te.TeamID.Int64,
		Enum:      te.Enum,
		Data:      te.Data,
		CreatedAt: te.CreatedAt.Time,
	}

	if te.DefaultDatum.Valid {
		teamEnum.DefaultDatum = &te.DefaultDatum.String
	}

	if te.UpdatedAt.Valid {
		teamEnum.UpdatedAt = &te.UpdatedAt.Time
	}

	return teamEnum, nil
}
