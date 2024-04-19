package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"

	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
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

func resetPolicies() {
	policies := []gormadapter.CasbinRule{
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
		err := tx.Table("casbin_rules").
			Where("ptype = ?", "p").
			Delete(&struct{}{}).
			Error
		if err != nil {
			return err
		}

		if err := tx.Table("casbin_rules").Create(policies).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func resetGroupingPolicies() {
	user := &model.User{}
	database := service.Registry.Get("database").(*gorm.DB)
	db := database.Where("email = ?", "admin@admin.com").
		First(user)
	if err := db.Error; err != nil {
		panic(err)
	}

	groupingPolicies := []gormadapter.CasbinRule{
		{Ptype: "g", V0: fmt.Sprintf("%d", user.Id), V1: "admin"},
	}
	err := database.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("casbin_rules").
			Where("ptype = ?", "g").
			Delete(&struct{}{}).
			Error
		if err != nil {
			return err
		}

		if err := tx.Table("casbin_rules").Create(groupingPolicies).Error; err != nil {
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
				Name:  "reset-policies",
				Usage: "reset policies",
				Action: func(cCtx *cli.Context) error {
					resetPolicies()
					return nil
				},
			},
			{
				Name:  "reset-group-policies",
				Usage: "reset group policies",
				Action: func(cCtx *cli.Context) error {
					resetGroupingPolicies()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
