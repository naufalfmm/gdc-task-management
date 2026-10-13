package cmd

import (
	"database/sql"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

func create(d *sql.DB) cli.ActionFunc {
	return func(c *cli.Context) error {
		name := strings.ToLower(c.String("name"))
		name = strings.ReplaceAll(name, " ", "_")

		if err := goose.SetDialect("postgres"); err != nil {
			return err
		}

		if err := goose.Create(d, "migrations/sql", name, "sql"); err != nil {
			return err
		}

		return nil
	}
}

func Create(d *sql.DB) *cli.Command {
	return &cli.Command{
		Name:    "create",
		Usage:   "create --name <name>",
		Aliases: []string{"c"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Required: true,
			},
		},
		Action: create(d),
	}
}
