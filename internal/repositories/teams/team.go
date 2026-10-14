package teams

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/naufalfmm/gdc-task-management/internal/domains/teams"
	"github.com/naufalfmm/gdc-task-management/internal/models/daos"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
)

type Repository interface {
	GetTeamByName(ctx context.Context, name string) (daos.Team, error)
	InsertTeam(ctx context.Context, team daos.Team) (daos.Team, error)
}

type repository struct {
	querier teams.Querier
}

func NewRepository(querier teams.Querier) Repository {
	return &repository{
		querier: querier,
	}
}

func (r *repository) GetTeamByName(ctx context.Context, name string) (daos.Team, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = teams.New(tx)
	}

	t, err := q.GetTeamByName(ctx, name)
	if err != nil {
		return daos.Team{}, err
	}

	return daos.Team{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt.Time,
		UpdatedAt: t.UpdatedAt.Time,
	}, nil
}

func (r *repository) InsertTeam(ctx context.Context, team daos.Team) (daos.Team, error) {
	q := r.querier

	tx, ok := txmanager.ExtractTx[pgx.Tx](ctx)
	if ok {
		q = teams.New(tx)
	}

	t, err := q.InsertTeam(ctx, team.Name)
	if err != nil {
		return daos.Team{}, err
	}

	return daos.Team{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt.Time,
		UpdatedAt: t.UpdatedAt.Time,
	}, nil
}
