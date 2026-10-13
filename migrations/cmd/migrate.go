package cmd

import (
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

func migrate(d *sql.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		if err := checkConnection(c.Context, d); err != nil {
			return err
		}

		if err := goose.SetDialect("postgres"); err != nil {
			return err
		}

		if err := goose.Up(d, "migrations/sql"); err != nil {
			return err
		}

		return nil
	}
}

func Migrate(d *sql.DB) *cli.Command {
	return &cli.Command{
		Name:    "migrate",
		Usage:   "migrate",
		Aliases: []string{"m"},
		Action:  migrate(d),
	}
}
