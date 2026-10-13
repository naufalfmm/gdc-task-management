package main

import (
	"database/sql"
	"os"

	"github.com/naufalfmm/gdc-task-management/config"
	"github.com/naufalfmm/gdc-task-management/migrations/cmd"
	"github.com/urfave/cli/v2"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	d, err := sql.Open("pgx", conf.GetDSN())
	if err != nil {
		panic(err)
	}
	defer d.Close()

	app := cli.NewApp()
	app.Commands = []*cli.Command{
		cmd.Migrate(d),
		cmd.Rollback(d),
		cmd.Create(d),
	}

	if err = app.Run(os.Args); err != nil {
		panic(err)
	}
}
