package main

import (
	"strconv"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user/model"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func init() {
	globalConfig := newGlobalConfig()
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	registerConfig(*globalConfig)

	registerProvider()
}

func addPolicies() {
	policies := [][]string{
		{"admin", "/api/admin/policy", "POST"},
		{"admin", "/api/admin/policy", "DELETE"},
		{"admin", "/api/admin/policy/subject", "GET"},
		{"admin", "/api/admin/policy/subject/:subject", "GET"},
		{"admin", "/api/admin/policy/subject", "DELETE"},
		{"admin", "/api/admin/policy/reload", "POST"},
		{"admin", "/api/admin/grouping-policy", "POST"},
		{"admin", "/api/admin/grouping-policy/:userId", "GET"},
		{"admin", "/api/admin/grouping-policy", "DELETE"},
	}
	enforcer := provider.Registry.Casbin()
	logger := provider.Registry.Logger()

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
	db := provider.Registry.DB().Where("email = ?", "admin@admin.com").
		First(user)
	if err := db.Error; err != nil {
		panic(err)
	}

	groupingPolicies := [][]string{
		{strconv.Itoa(int(user.Id)), "admin"},
	}

	enforcer := provider.Registry.Casbin()
	logger := provider.Registry.Logger()

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
