package main

import (
	"flag"
	"os"

	_ "github.com/joho/godotenv/autoload"

	_ "github.com/ChanJuiHuang/go-backend-framework/cmd/kit/internal/env"

	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/cmd/kit/http"
	"github.com/ChanJuiHuang/go-backend-framework/cmd/kit/key"
	"github.com/ChanJuiHuang/go-backend-framework/cmd/kit/permission"
	"github.com/ChanJuiHuang/go-backend-framework/database/seeder"
)

func main() {
	var envPath string

	generateJwtKeyFlag := flag.Bool("generate-jwt-key", false, "generate jwt private and public key")
	flag.StringVar(&envPath, "env-path", ".env", "dot env file path")
	runSeederFlag := flag.Bool("run-seeder", false, "run database seeder")
	updateRootUserPasswordFlag := flag.Bool("update-root-user-password", false, "update root user password")
	generateRootAccessTokenFlag := flag.Bool("generate-root-access-token", false, "generate root access token")
	routeListFlag := flag.Bool("route-list", false, "route list")
	importCasbinPoliciesFlag := flag.Bool("import-casbin-policies", false, "import casbin policies")
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if *generateJwtKeyFlag {
		key.GenerateJwtKey(envPath)
	}

	if *runSeederFlag {
		seeder.Run()
	}

	if *updateRootUserPasswordFlag {
		permission.UpdateRootUserPassword()
	}

	if *generateRootAccessTokenFlag {
		permission.GenerateRootAccessToken()
	}

	if *routeListFlag {
		http.RouteList()
	}

	if *importCasbinPoliciesFlag {
		provider.App.ImportCasbinPolicies()
	}
}
