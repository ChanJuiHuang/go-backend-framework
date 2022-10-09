# go-backend-framework

## Download go-backend-framework installer and install framework
```bash
VERSION=<version-id> # ex: VERSION=1.0.0
curl -L https://github.com/ChanJuiHuang/go-backend-framework/archive/refs/tags/v$VERSION.tar.gz | tar -zxv --strip-components=1 go-backend-framework-$VERSION/gbf-installer.sh
./gbf-installer.sh

# Compile go-backend-framework tool kit
cd <project-name> && make kit
./bin/kit -h

# If you enable [database] and [casbin], you could run below:
# 1. [make mysql-migration args=up]. Please set db name, password and so on to the [.env] file.
# 2. [./bin/kit -run-seeder]
# 3. [./bin/kit -import-casbin-policies]
# 4. [./bin/kit -generate-jwt-key]
# All above them to set the initial data.

# If you enable [database], [redis] and [casbin], you could execute [make test] to run sample test.

# Enable database, redis and casbin can try all sample code!

# Check the server is running
go run cmd/server/main.go
curl http://127.0.0.1:8080/api/ping
```

## The sample code
The code in module, controller, scheduler and database are sample codes.
If you don't need them, you can remove it.

The permission folder in the kit folder is based on the user module.
If you want to customize your user module, please remove it.

## The local development server
Use [air](https://github.com/cosmtrek/air) to run development server.
You can follow .air.toml.example to set air configuration, or customize your configuration file.

## Makefile
The makefile contains **build command**, **test command** and **database migration command**.

## Using packages
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
* [envconfig](https://github.com/kelseyhightower/envconfig)
* [wire](https://github.com/google/wire)
* [argon2](https://pkg.go.dev/golang.org/x/crypto@v0.5.0/argon2)
* [oauth](https://pkg.go.dev/golang.org/x/oauth2@v0.4.0)
* [rate-limit](https://pkg.go.dev/golang.org/x/time@v0.3.0/rate)
* [json-iterator](https://github.com/json-iterator/go)
* [testify](https://github.com/stretchr/testify)
* [swagger](https://github.com/swaggo/swag)

## Folder structure
```
go-backend-framework
├── LICENSE
├── Makefile
├── README.md
├── app
│   ├── config
│   │   ├── app.go
│   │   ├── casbin.go
│   │   ├── database.go
│   │   ├── http.go
│   │   ├── log.go
│   │   └── redis.go
│   ├── http
│   │   ├── controller
│   │   │   └── user
│   │   │       ├── email_login.go
│   │   │       ├── oauth_login.go
│   │   │       ├── token_refresh.go
│   │   │       ├── user_logout.go
│   │   │       ├── user_register.go
│   │   │       ├── user_search.go
│   │   │       └── user_update.go
│   │   ├── middleware
│   │   │   ├── authenticate_middleware.go
│   │   │   ├── authorize_middleware.go
│   │   │   ├── global_middleware.go
│   │   │   ├── rate_limit_middleware.go
│   │   │   ├── record_access_log_middleware.go
│   │   │   ├── recover_server_middleware.go
│   │   │   └── verify_csrf_token_middleware.go
│   │   ├── response
│   │   │   ├── error.go
│   │   │   └── error_response.go
│   │   ├── route
│   │   │   ├── api_route.go
│   │   │   ├── scheduler_route.go
│   │   │   └── swgger_route.go
│   │   ├── scheduler
│   │   │   └── user
│   │   │       └── refresh_token_record_delete.go
│   │   └── server.go
│   ├── module
│   │   └── user
│   │       ├── model
│   │       │   ├── email_user.go
│   │       │   ├── google_user.go
│   │       │   ├── refresh_token_record.go
│   │       │   └── user.go
│   │       ├── oauth_service.go
│   │       └── user_service.go
│   └── provider
│       ├── application.go
│       ├── application_provider.go
│       ├── casbin_provider.go
│       ├── database_provider.go
│       ├── internal
│       │   └── init
│       │       └── init.go
│       ├── json_provider.go
│       ├── log_provider.go
│       ├── modifier_provider.go
│       ├── redis_provider.go
│       ├── validator_provider.go
│       ├── wire.go
│       └── wire_gen.go
├── bin
│   ├── kit
│   └── server
├── casbin_model.conf
├── casbin_model.conf.example
├── cmd
│   ├── kit
│   │   ├── http
│   │   │   └── route_list.go
│   │   ├── key
│   │   │   └── generate_jwt_key.go
│   │   ├── main.go
│   │   └── permission
│   │       ├── generate_root_access_token.go
│   │       └── update_root_user_password.go
│   └── server
│       └── main.go
├── database
│   ├── migration
│   │   ├── 20230110113645_create_users_table.sql
│   │   ├── 20230110113926_create_refresh_token_records_table.sql
│   │   ├── 20230110114857_create_email_users_table.sql
│   │   └── 20230110124046_create_google_users_table.sql
│   ├── seeder
│   │   ├── seeder.go
│   │   └── user_seeder.go
│   └── testing-migration
│       ├── 20230110113645_create_users_table.sql
│       ├── 20230110113926_create_refresh_token_records_table.sql
│       ├── 20230110114857_create_email_users_table.sql
│       └── 20230110124046_create_google_users_table.sql
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── gbf-installer.sh
├── go.mod
├── go.sum
├── storage
│   └── log
│       ├── access.log
│       └── app.log
├── test
│   ├── example
│   │   └── example_test.go
│   ├── internal
│   │   ├── env
│   │   │   └── env.go
│   │   └── test
│   │       └── test.go
│   └── user
│       ├── email_login_test.go
│       ├── oauth_login_test.go
│       ├── refresh_token_record_delete_test.go
│       ├── token_refresh_test.go
│       ├── user_logout_test.go
│       ├── user_register_test.go
│       ├── user_search_test.go
│       └── user_update_test.go
└── util
    ├── argon2.go
    ├── database.go
    ├── error.go
    ├── google_oauth.go
    ├── jwt.go
    ├── mocked_oauth.go
    ├── oauth.go
    ├── paginator.go
    ├── random.go
    ├── rate_limiter.go
    ├── stacktrace.go
    └── time.go
```

### The app folder
The app folder includes core code in the application.
You can add the custom folder in the app folder.

---
### The provider folder
The provider folder provider the basic service for developer.
For example, database, redis, log... and so on.
All providers use [wire](https://github.com/google/wire) to do dependency injection.

---
### The config folder
The config folder is the configuration for services from the provider folder.

---
### The module folder
The module folder contains all code of business logic that can called service.
You can design the sub-folder in module folder. It does not always follow raw folder structure.

---
### The http folder
The http folder contains all http relation functions.

#### The controller folder
The controller folder use services from module to handle http request.
You can design the sub-folder in controller folder. It does not always follow raw folder structure.

#### The middleware folder
The middleware folder contains middlewares for http request.

#### The response folder
The response folder contains uniform response for http response.
Currently, it has error response.

#### The route folder
The route folder contains http routes.

#### The scheduler folder
The scheduler folder contains http api for cronjob.

---
### The bin folder
The bin folder contains binary file that made from the go build command.

---
### The cmd folder
The cmd folder contain the code for command.

---
### The database folder
The database folder contains migrations, testing migrations and seeder.

#### The migration folder
The migration folder contain database migration that made from the [goose](https://github.com/pressly/goose) command.
The makefile integrate the migration command. You can see it in makefile.

#### The testing-migration folder
The testing-migration folder contain database migration for testing.
The testing migration use sqlite in memory database to testing.

#### The seeder folder
The seeder folder contains database seeder.

---
### The docs folder
The docs folder is generated by [swagger](https://github.com/swaggo/swag) that contains api documents.

---
### The storage folder
The storage folder contains app log and access log.

---
### The test folder
The test folder contains testing code that written by [testify](https://github.com/stretchr/testify).
The test cases in test folder always test the public method, because public method includes private method.
If you want to test the private method, you could design your testing folder or test case.

All test data is stored in the sqlite, so you have install sqlite and write testing migration to test. If you don't want to use sqlite, you have to change the **test/internal/env/env.go**.

---
### The util folder
The util folder contains general usage code for daily develop.
Maybe that will become independent package in the future.
