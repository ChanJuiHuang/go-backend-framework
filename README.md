# go-backend-framework

## Quick Start
### Download go-backend-framework installer and install framework
You can see the whole versions in [release page](https://github.com/chan-jui-huang/go-backend-framework/releases).
```bash
VERSION=<version-id> # ex: VERSION=2.X.X
curl -L https://github.com/chan-jui-huang/go-backend-framework/archive/refs/tags/v$VERSION.tar.gz | tar -zxv --strip-components=1 go-backend-framework-$VERSION/gbf-installer.sh
./gbf-installer.sh

# Compile go-backend-framework tool kit
cd <project-name> && make
```

### Run the app
#### Before running the app, please install below
1. mysql
2. sqlite (to run test case)

### Set .env
1. fill the db environment variable in [.env] file
2. make mysql-migration args=up
3. execute the [./bin/database_seeder]
4. execute the [./bin/policy_seeder]
5. execute the [./bin/jwt]
6. execute the [./bin/jwt -env=.env.testing]

### Run test case
1. make test

### Run the app
1. go run cmd/app/*
2. curl http://127.0.0.1:8080/api/ping

## Introduction
The folder structure follows the [project-layout](https://github.com/golang-standards/project-layout).

The code in the internal folder that is designed for project, so you can change it by yourself.

The pkg folder inside the internal folder is that it is reused package only for THIS project.

The pkg folder outside the internal folder is the reused package. In the future, I plan to make it to standalone package.

The permission management use the [casbin](https://github.com/casbin/casbin). Please read document to understand.

## The local development server
Use [air](https://github.com/cosmtrek/air) to run development server.
You can follow .air.toml.example to set air configuration, or customize your configuration file.

## Makefile
The makefile contains **build command**, **test command** and **database migration command** ...etc.

## Used packages
* [gin](https://github.com/gin-gonic/gin)
* [go-playground-validator](https://github.com/go-playground/validator)
* [go-playground-form](https://github.com/go-playground/form)
* [go-playground-mold](https://github.com/go-playground/mold)
* [jwt-go](https://github.com/golang-jwt/jwt)
* [gorm](https://github.com/go-gorm/gorm)
* [goose](https://github.com/pressly/goose)
* [go-redis](https://github.com/redis/go-redis)
* [zap](https://github.com/uber-go/zap)
* [lumberjack](https://github.com/natefinch/lumberjack)
* [godotenv](https://github.com/joho/godotenv)
* [viper](https://github.com/spf13/viper)
* [pkg/errors](https://github.com/pkg/errors)
* [casbin](https://github.com/casbin/casbin)
* [mapstructure](https://github.com/mitchellh/mapstructure)
* [structs](https://github.com/fatih/structs)
* [argon2](https://pkg.go.dev/golang.org/x/crypto@v0.5.0/argon2)
* [rate-limit](https://pkg.go.dev/golang.org/x/time@v0.3.0/rate)
* [testify](https://github.com/stretchr/testify)
* [swag](https://github.com/swaggo/swag)

## Folder structure
```
go-backend-framework
├── LICENSE
├── Makefile
├── README.md
├── bin
├── cmd
│   ├── app
│   │   └── main.go
│   ├── kit
│   │   ├── database_seeder
│   │   │   └── database_seeder.go
│   │   ├── http_route
│   │   │   └── http_route.go
│   │   ├── jwt
│   │   │   └── jwt.go
│   │   └── policy_seeder
│   │       └── policy_seeder.go
│   └── template
│       └── template.go
├── config.production.yml
├── config.testing.yml
├── config.yml
├── deployment
│   └── docker
│       └── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── gbf-installer.sh
├── go.mod
├── go.sum
├── internal
│   ├── http
│   │   ├── controller
│   │   │   ├── admin
│   │   │   │   ├── admin_create_grouping_policy.go
│   │   │   │   ├── admin_create_grouping_policy_test.go
│   │   │   │   ├── admin_create_policy.go
│   │   │   │   ├── admin_create_policy_test.go
│   │   │   │   ├── admin_delete_grouping_policy.go
│   │   │   │   ├── admin_delete_grouping_policy_test.go
│   │   │   │   ├── admin_delete_policy.go
│   │   │   │   ├── admin_delete_policy_subject.go
│   │   │   │   ├── admin_delete_policy_subject_test.go
│   │   │   │   ├── admin_delete_policy_test.go
│   │   │   │   ├── admin_get_grouping_policy.go
│   │   │   │   ├── admin_get_grouping_policy_test.go
│   │   │   │   ├── admin_get_policy_subject.go
│   │   │   │   ├── admin_get_policy_subject_test.go
│   │   │   │   ├── admin_reload_policy.go
│   │   │   │   ├── admin_reload_policy_test.go
│   │   │   │   ├── admin_search_policy_subject.go
│   │   │   │   └── admin_search_policy_subject_test.go
│   │   │   └── user
│   │   │       ├── user_login.go
│   │   │       ├── user_login_test.go
│   │   │       ├── user_me.go
│   │   │       ├── user_me_test.go
│   │   │       ├── user_register.go
│   │   │       ├── user_register_test.go
│   │   │       ├── user_update.go
│   │   │       ├── user_update_password.go
│   │   │       ├── user_update_password_test.go
│   │   │       └── user_update_test.go
│   │   ├── middleware
│   │   │   ├── access_log_middleware.go
│   │   │   ├── authentication_middleware.go
│   │   │   ├── authorization_middleware.go
│   │   │   ├── global_middleware.go
│   │   │   ├── rate_limit_middleware.go
│   │   │   ├── recover_middleware.go
│   │   │   └── verify_csrf_token_middleware.go
│   │   ├── response
│   │   │   ├── error_message.go
│   │   │   ├── error_message_test.go
│   │   │   ├── error_response.go
│   │   │   └── response.go
│   │   ├── route
│   │   │   ├── admin
│   │   │   │   └── api_route.go
│   │   │   ├── api_route.go
│   │   │   ├── swagger_route.go
│   │   │   └── user
│   │   │       └── api_route.go
│   │   └── server.go
│   ├── migration
│   │   ├── 20230818113729_create_users_table.sql
│   │   ├── seeder
│   │   │   ├── seeder.go
│   │   │   └── user_seeder.go
│   │   └── test
│   │       └── 20230818113729_create_users_table.sql
│   ├── pkg
│   │   └── user
│   │       ├── model
│   │       │   └── user.go
│   │       └── user.go
│   ├── registrar
│   │   ├── authentication_registrar.go
│   │   ├── casbin_registrar.go
│   │   ├── database_registrar.go
│   │   ├── logger_registrar.go
│   │   ├── mapstructure_decoder_registrar.go
│   │   ├── redis_registrar.go
│   │   ├── register_executor.go
│   │   ├── registrar_test.go
│   │   └── simple_register_executor.go
│   ├── scheduler
│   │   ├── job
│   │   │   └── example_job.go
│   │   └── scheduler.go
│   └── test
│       ├── admin.go
│       ├── http.go
│       ├── migration.go
│       ├── test.go
│       └── user.go
└── storage
    └── log
        ├── access.log
        └── app.log
```

## Folder and package introduction
Introduction that it will introduce the **IMPORTANT** folders and packages.

### The pkg folder
The app folder includes core code in the application.
You can add the custom folder in the app folder.

#### package:
1. app
    * Define the lifecycle for your app.
    * The lifecycle has three stage that are **starting**, **started** and **terminated**.
    * You can use callback function at any stage.

2. booter
    * The booter package define how to boot the your app.
    * Booter has config registry and service registry to store the config and instance of the struct by singleton pattern.
    * To store config and instance. You can follow the internal/registrar package. It has the example to show how to do it.

---
### The internal folder
The internal folder only use for this project. It always design for business logic.

### The internal/pkg folder
Design for reused business logic.

---
### The internal/migration folder
It contains database migration that make from the [goose](https://github.com/pressly/goose).

The makefile integrates the migration commands. You can use it to add new migration and run any migration command.

#### The internal/migration/test folder
It contains database migrations for testing.
The testing migration use sqlite in memory database to testing.

#### The internal/migration/seeder folder
It contains database seeders.

---
### The internal/http folder
It contains all http relation functions.

#### The internal/http/controller folder
It use internal/pkg, third party package and so on to construct the http endpoint.

#### The internal/http/middleware folder
It contains middlewares for http request and response.

#### The internal/http/response folder
It defines response format for the http response.

#### The internal/http/route folder
It defines http routes.

---
### The internal/test folder
It has testing initializer and useful testing function that customizes for the go backend framework.

If you use the code base of the go backend framework, I suggest you to use [testify](https://github.com/stretchr/testify) to initialize your test case before you run the test case.

All testing data are stored in the sqlite, so you have install sqlite and write testing migration to test. If you don't want to use sqlite, you have to change the **config.testing.yml**.

Please see the testing example in the internal/http/controller!

---
### The cmd folder
It has all command from the project.

#### The cmd/app folder
It is the main command in this project. If you execute it, the http server will be started.

#### The cmd/kit/* folder
It has the useful commands to help you develop the app.

For example, create jwt, database seeder executer, http route list and predefined permissions.

#### The cmd/template folder
It is the command template. If you want to use config registry and provider registry, you can use this template.

---
### The storage folder
It use for static file. Currently, it only has app log and access log.

---
### The bin folder
It contains binary file from the cmd folder.

---
### The docs folder
The docs folder is generated by [swag](https://github.com/swaggo/swag). It is swagger api document.
