# go-backend-framework

## Quick Start
### Download go-backend-framework installer and install framework
```bash
VERSION=<version-id> # ex: VERSION=2.0.0
curl -L https://github.com/ChanJuiHuang/go-backend-framework/archive/refs/tags/v$VERSION.tar.gz | tar -zxv --strip-components=1 go-backend-framework-$VERSION/gbf-installer.sh
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
* [go-redis](https://github.com/go-redis/redis)
* [zap](https://github.com/uber-go/zap)
* [lumberjack](https://github.com/natefinch/lumberjack)
* [godotenv](https://github.com/joho/godotenv)
* [viper](https://github.com/spf13/viper)
* [pkg/errors](https://github.com/pkg/errors)
* [casbin](https://github.com/casbin/casbin)
* [structs](https://github.com/fatih/structs)
* [argon2](https://pkg.go.dev/golang.org/x/crypto@v0.5.0/argon2)
* [rate-limit](https://pkg.go.dev/golang.org/x/time@v0.3.0/rate)
* [testify](https://github.com/stretchr/testify)
* [swag](https://github.com/swaggo/swag)

## Folder structure
```
go-backend-framework
├── bin
│   ├── app
│   ├── jwt
│   ├── policy_seeder
│   ├── route
│   └── seeder
├── cmd
│   ├── app
│   │   ├── init.go
│   │   └── main.go
│   ├── kit
│   │   ├── jwt
│   │   │   └── jwt.go
│   │   ├── policy_seeder
│   │   │   ├── init.go
│   │   │   └── policy_seeder.go
│   │   ├── route
│   │   │   └── route.go
│   │   └── seeder
│   │       ├── init.go
│   │       └── seeder.go
│   └── template
│       ├── init.go
│       └── template.go
├── config.production.yml
├── config.testing.yml
├── config.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── gbf-installer.sh
├── go.mod
├── go.sum
├── internal
│   ├── global
│   │   └── config.go
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
│   │   ├── provider
│   │   │   ├── provider.go
│   │   │   └── registry.go
│   │   └── user
│   │       ├── model
│   │       │   └── user.go
│   │       └── user.go
│   └── test
│       ├── admin.go
│       ├── http.go
│       ├── migration.go
│       ├── test.go
│       └── user.go
├── LICENSE
├── Makefile
├── pkg
│   ├── app
│   │   └── app.go
│   ├── argon2
│   │   ├── argon2.go
│   │   └── argon2_test.go
│   ├── authentication
│   │   ├── authenticator.go
│   │   ├── authenticator_test.go
│   │   └── config.go
│   ├── config
│   │   └── registry.go
│   ├── database
│   │   ├── config.go
│   │   └── database.go
│   ├── logger
│   │   ├── config.go
│   │   └── logger.go
│   ├── pagination
│   │   └── pagination.go
│   ├── random
│   │   ├── random.go
│   │   └── random_test.go
│   ├── redis
│   │   ├── config.go
│   │   └── redis.go
│   └── stacktrace
│       └── stacktrace.go
├── README.md
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

2. config
    * It is a config registry. It must be import by any package to register its config.
    * Config registry uses viper and config.yml to integrate and defined the default config.
    * If config.yml has environment variables, you will use [dot env] to integrate.

---
### The internal folder
The internal folder only use for this project. It always design for business logic.

### The internal/pkg folder
Design for reused business logic.

#### The internal/pkg/provider folder
It is a provider registry. It uses the single pattern to register the instance of struct when the app starts.

If you want to add your instance of struct, you must modify the provider registry.

---
### The internal/global folder
Currently, It only has a config for this project.

Global means that the config can be use for any where.

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
