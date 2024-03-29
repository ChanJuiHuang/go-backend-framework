package main

import (
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewDefaultConfig,
		&registrar.RegisterExecutor,
	)
}

func addPolicies() {
	policies := [][]string{
		{"admin", "/api/admin/policy", "POST"},
		{"admin", "/api/admin/policy", "DELETE"},
		{"admin", "/api/admin/policy/subject", "GET"},
		{"admin", "/api/admin/policy/subject/:subject", "GET"},
		{"admin", "/api/admin/policy/subject/:subject/user", "GET"},
		{"admin", "/api/admin/policy/subject", "DELETE"},
		{"admin", "/api/admin/policy/reload", "POST"},
		{"admin", "/api/admin/grouping-policy", "POST"},
		{"admin", "/api/admin/user/:userId/grouping-policy", "GET"},
		{"admin", "/api/admin/grouping-policy", "DELETE"},
	}
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	logger := service.Registry.Get("logger").(*zap.Logger)

	err := enforcer.GetAdapter().(*gormadapter.Adapter).Transaction(enforcer, func(e casbin.IEnforcer) error {
		result, err := e.RemovePolicies(policies)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if !result {
			logger.Warn("policies are not in [casbin_rules] table PROBABLY")
		}

		result, err = enforcer.AddPolicies(policies)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if !result {
			panic("add policies failed")
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func addGroupingPolicies() {
	user := &model.User{}
	database := service.Registry.Get("database").(*gorm.DB)
	db := database.Where("email = ?", "admin@admin.com").
		First(user)
	if err := db.Error; err != nil {
		panic(err)
	}

	groupingPolicies := [][]string{
		{strconv.Itoa(int(user.Id)), "admin"},
	}

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	logger := service.Registry.Get("logger").(*zap.Logger)

	err := enforcer.GetAdapter().(*gormadapter.Adapter).Transaction(enforcer, func(e casbin.IEnforcer) error {
		result, err := e.RemoveGroupingPolicies(groupingPolicies)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if !result {
			logger.Warn("grouping policies are not in [casbin_rules] table PROBABLY")
		}

		result, err = e.AddGroupingPolicies(groupingPolicies)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if !result {
			panic("add grouping policies failed")
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	addPolicies()
	addGroupingPolicies()
}
