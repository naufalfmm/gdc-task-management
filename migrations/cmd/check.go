package cmd

import (
	"context"
	"database/sql"
)

func checkConnection(ctx context.Context, d *sql.DB) error {
	if err := d.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
