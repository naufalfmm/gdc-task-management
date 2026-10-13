package cmd

import (
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

func rollback(d *sql.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		if err := checkConnection(c.Context, d); err != nil {
			return err
		}

		ver := c.Int64("version")

		if err := goose.SetDialect("postgres"); err != nil {
			return err
		}

		if ver > 0 {
			if err := goose.DownTo(d, "migrations/sql", ver); err != nil {
				return err
			}
		} else {
			if err := goose.Down(d, "migrations/sql"); err != nil {
				return err
			}
		}

		return nil
	}
}

func Rollback(d *sql.DB) *cli.Command {
	return &cli.Command{
		Name:    "rollback",
		Usage:   "rollback --version <version>",
		Aliases: []string{"r"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "version",
				Aliases:  []string{"v"},
				Required: false,
			},
		},
		Action: rollback(d),
	}
}
