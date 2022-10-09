include .env

.PHONY: all run clean test
SHELL:=/bin/bash
SERVER_BIN:=server
SERVER_FILE:=${PWD}/cmd/server/main.go

KIT_BIN:=kit
KIT_FILE:=${PWD}/cmd/kit/*.go

GOFILES:=$(shell find . -type f -name "*.go")
TAGS:="jsoniter"

OBJECTS:=kit

all:${OBJECTS}
	go build -o ${PWD}/bin/${SERVER_BIN} -v -tags ${TAGS} -ldflags "-s -w" ${SERVER_FILE}

run:
	go run -race -tags ${TAGS} ${SERVER_FILE}

debug-server:
	go build -o ${PWD}/bin/${SERVER_BIN}-debug -v -race -tags ${TAGS} -gcflags="-dwarflocationlists=true" ${SERVER_FILE}

kit:
	go build -o ${PWD}/bin/${KIT_BIN} -v -race -tags ${TAGS} -ldflags "-s -w" ${KIT_FILE}

clean:
	rm -rf ${PWD}/bin/*

mysql-migration:
	goose -dir database/migration -allow-missing -s mysql "${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?parseTime=true&loc=UTC" ${args}

pgsql-migration:
	goose -dir database/migration -allow-missing -s postgres "user=${DB_USERNAME} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} dbname=${DB_DATABASE} sslmode=disable" ${args}

sqlite-migration:
	goose -dir database/migration -allow-missing -s sqlite3 ${DB_DATABASE} ${args}

test:
	$(eval args?=./test/...)
	go test ${args}

benchmark:
	$(eval args?=./test/...)
	go test -bench=. -run=none -benchmem ${args}
