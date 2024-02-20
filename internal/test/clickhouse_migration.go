package test

import (
	"database/sql"
	"fmt"
	"path"

	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/clickhouse"
	"github.com/pressly/goose/v3"
)

type clickhouseMigration struct {
	dir string
}

var ClickhouseMigration *clickhouseMigration

func NewClickhouseMigration() *clickhouseMigration {
	booterConfig := config.Registry.Get("booter").(booter.Config)

	return &clickhouseMigration{
		dir: path.Join(booterConfig.RootDir, "internal/migration/clickhouse/test"),
	}
}

func (cm *clickhouseMigration) Run(callbacks ...func()) {
	clickhouseConfig := config.Registry.Get("clickhouse").(clickhouse.Config)
	conn, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s?username=%s&password=%s", clickhouseConfig.Addr[0], clickhouseConfig.Username, clickhouseConfig.Password))
	if err != nil {
		panic(err)
	}

	if _, err := conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", clickhouseConfig.Database)); err != nil {
		panic(err)
	}

	if _, err := conn.Exec(fmt.Sprintf("USE %s", clickhouseConfig.Database)); err != nil {
		panic(err)
	}

	if err := goose.SetDialect(string(goose.DialectClickHouse)); err != nil {
		panic(err)
	}

	if err := goose.Up(conn, cm.dir); err != nil {
		panic(err)
	}

	for _, callback := range callbacks {
		callback()
	}
}

func (cm *clickhouseMigration) Reset() {
	clickhouseConfig := config.Registry.Get("clickhouse").(clickhouse.Config)
	conn, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s?username=%s&password=%s", clickhouseConfig.Addr[0], clickhouseConfig.Username, clickhouseConfig.Password))
	if err != nil {
		panic(err)
	}

	if _, err := conn.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", clickhouseConfig.Database)); err != nil {
		panic(err)
	}
}
