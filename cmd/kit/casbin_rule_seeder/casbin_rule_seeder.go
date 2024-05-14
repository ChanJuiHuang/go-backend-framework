package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"

	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/permission"
	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewDefaultConfig,
		&registrar.RegisterExecutor,
	)
}

func resetPermissions() {
	casbinRules := []gormadapter.CasbinRule{
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy", V2: "POST"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy", V2: "DELETE"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/subject", V2: "GET"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/subject/:subject", V2: "GET"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/subject/:subject/user", V2: "GET"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/subject", V2: "DELETE"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/reload", V2: "POST"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/grouping-policy", V2: "POST"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/user/:userId/grouping-policy", V2: "GET"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/grouping-policy", V2: "DELETE"},
	}
	database := service.Registry.Get("database").(*gorm.DB)
	err := database.Transaction(func(tx *gorm.DB) error {
		if err := permission.DeleteCasbinRule(tx, "ptype = ?", "p"); err != nil {
			return err
		}

		if err := permission.CreateCasbinRule(tx, casbinRules); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func resetRoles() {
	user := &model.User{}
	database := service.Registry.Get("database").(*gorm.DB)
	db := database.Where("email = ?", "admin@admin.com").
		First(user)
	if err := db.Error; err != nil {
		panic(err)
	}

	casbinRoles := []gormadapter.CasbinRule{
		{Ptype: "g", V0: fmt.Sprintf("%d", user.Id), V1: "admin"},
	}
	err := database.Transaction(func(tx *gorm.DB) error {
		if err := permission.DeleteCasbinRule(tx, "ptype = ?", "g"); err != nil {
			return err
		}

		if err := permission.CreateCasbinRule(tx, casbinRoles); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "reset-permission",
				Usage: "reset permission",
				Action: func(cCtx *cli.Context) error {
					resetPermissions()
					return nil
				},
			},
			{
				Name:  "reset-role",
				Usage: "reset role",
				Action: func(cCtx *cli.Context) error {
					resetRoles()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
