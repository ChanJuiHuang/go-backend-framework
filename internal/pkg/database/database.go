package database

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"gorm.io/gorm"
)

func NewTx(table string, preloads ...string) *gorm.DB {
	database := service.Registry.Get("database").(*gorm.DB)
	tx := database.Table(table)

	for _, preload := range preloads {
		tx.Preload(preload)
	}

	return tx
}
