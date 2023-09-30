package main

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	byteYaml, err := os.ReadFile(path.Join(wd, "config.yml"))
	if err != nil {
		panic(err)
	}
	stringYaml := os.ExpandEnv(string(byteYaml))

	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(strings.NewReader(stringYaml))

	config.Registry.SetViper(v)
	config.Registry.Register(map[string]any{
		"database": &database.Config{},
	})
}

func addPolicies(enforcer *casbin.SyncedCachedEnforcer) {
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

	err := enforcer.GetAdapter().(*gormadapter.Adapter).Transaction(enforcer, func(e casbin.IEnforcer) error {
		result, err := e.RemovePolicies(policies)
		if err != nil {
			log.Println(err)
			return err
		}
		if !result {
			log.Println("policies is not in [casbin_rules] table PROBABLY")
		}

		result, err = enforcer.AddPolicies(policies)
		if err != nil {
			log.Println(err)
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

func addGroupingPolicies(db *gorm.DB, enforcer *casbin.SyncedCachedEnforcer) {
	user := &model.User{}
	tx := db.Where("email = ?", "admin@admin.com").
		First(user)
	if err := tx.Error; err != nil {
		panic(err)
	}

	groupingPolicies := [][]string{
		{strconv.Itoa(int(user.Id)), "admin"},
	}

	err := enforcer.GetAdapter().(*gormadapter.Adapter).Transaction(enforcer, func(e casbin.IEnforcer) error {
		result, err := enforcer.RemoveGroupingPolicies(groupingPolicies)
		if err != nil {
			log.Println(err)
			return err
		}
		if !result {
			log.Println("grouping policies is not in [casbin_rules] table PROBABLY")
		}

		result, err = enforcer.AddGroupingPolicies(groupingPolicies)
		if err != nil {
			log.Println(err)
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
	db := provider.ProvideDB()
	enforcer := provider.ProvideCasbin(db)
	log.SetFlags(log.LstdFlags | log.Llongfile)

	addPolicies(enforcer)
	addGroupingPolicies(db, enforcer)
}
