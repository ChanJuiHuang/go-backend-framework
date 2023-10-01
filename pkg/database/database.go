package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySqlDatabase(config Config, gConfig *gorm.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), gConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db
}

func NewPgSqlDatabase(config Config, gConfig *gorm.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	db, err := gorm.Open(postgres.Open(dsn), gConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db
}

func NewSqliteDatabase(config Config, gConfig *gorm.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.Database), gConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db
}

func New(config Config) *gorm.DB {
	gConfig := &gorm.Config{
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
		PrepareStmt:              true,
		Logger:                   logger.Default.LogMode(GetGormLogLevel(config.LogLevel)),
	}

	var db *gorm.DB
	switch GetDriver(config.Driver) {
	case MySql:
		db = NewMySqlDatabase(config, gConfig)
	case PgSql:
		db = NewPgSqlDatabase(config, gConfig)
	case Sqlite:
		db = NewSqliteDatabase(config, gConfig)
	}

	return db
}
