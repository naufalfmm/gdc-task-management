package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
)

type pgxTxManager struct {
	pool *pgxpool.Pool
}

func NewPgxTxManager(pool *pgxpool.Pool) txmanager.TxManager {
	return &pgxTxManager{pool: pool}
}

type pgxTx struct {
	tx pgx.Tx
}

func (w *pgxTx) Commit(ctx context.Context) error {
	return w.tx.Commit(ctx)
}

func (w *pgxTx) Rollback(ctx context.Context) error {
	return w.tx.Rollback(ctx)
}

func (m *pgxTxManager) Begin(ctx context.Context) (txmanager.Tx, context.Context, error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin tx: %w", err)
	}

	// Bungkus pgx.Tx ke dalam interface domain.Tx
	txWrapper := &pgxTx{tx: tx}

	// Inject tx asli ke dalam context untuk dibaca repository
	ctx = txmanager.InjectTx(ctx, tx)

	return txWrapper, ctx, nil
}
