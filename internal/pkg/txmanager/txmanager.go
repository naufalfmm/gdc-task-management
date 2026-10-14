package txmanager

import "context"

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxManager interface {
	Begin(ctx context.Context) (Tx, context.Context, error)
}

type txKey struct{}

func InjectTx[T any](ctx context.Context, tx T) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func ExtractTx[T any](ctx context.Context) (T, bool) {
	txCtx := ctx.Value(txKey{})

	tx, ok := txCtx.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return tx, ok
}
