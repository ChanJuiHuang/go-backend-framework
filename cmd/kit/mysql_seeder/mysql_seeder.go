package main

import (
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"

	"github.com/ChanJuiHuang/go-backend-framework/internal/migration/mysql/seeder"
	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewDefaultConfig,
		&registrar.RegisterExecutor,
	)
}

func main() {
	seederExecutor := seeder.NewSeederExecutor()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "show",
				Usage: "show all seeders",
				Action: func(cCtx *cli.Context) error {
					seederExecutor.ShowSeeders()
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Run the seeders. EX: database_seeder run seeder1,seeder2 (run specific seeders). database_seeder run (run all seeders).",
				Action: func(cCtx *cli.Context) error {
					seederExecutor.Run(strings.Split(cCtx.Args().First(), ","))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
