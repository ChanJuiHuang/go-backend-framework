package provider

import (
	"fmt"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func provideDatabase(gLogger *gormLogger) *gorm.DB {
	dbConfig := config.Database()
	if !dbConfig.Enabled {
		return new(gorm.DB)
	}

	db := new(gorm.DB)
	gConfig := &gorm.Config{
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
		PrepareStmt:              true,
		Logger:                   gLogger,
	}

	switch dbConfig.Driver {
	case config.MySql:
		db = createMySqlDatabase(gConfig)
	case config.PgSql:
		db = createPgSqlDatabase(gConfig)
	case config.Sqlite:
		db = createSqliteDatabase(gConfig)
	}

	return db
}

func createMySqlDatabase(gConfig *gorm.Config) *gorm.DB {
	dbConfig := config.Database()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), gConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	return db
}

func createPgSqlDatabase(gConfig *gorm.Config) *gorm.DB {
	dbConfig := config.Database()
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)
	db, err := gorm.Open(postgres.Open(dsn), gConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	return db
}

func createSqliteDatabase(gConfig *gorm.Config) *gorm.DB {
	dbConfig := config.Database()
	db, err := gorm.Open(sqlite.Open(dbConfig.Database), gConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	return db
}
