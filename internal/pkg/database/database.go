package database

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"gorm.io/gorm"
)

func NewTx(associations ...string) *gorm.DB {
	tx := service.Registry.Get("database").(*gorm.DB)
	for _, association := range associations {
		tx.Preload(association)
	}

	return tx
}

func NewTxByTable(table string, associations ...string) *gorm.DB {
	database := service.Registry.Get("database").(*gorm.DB)
	tx := database.Table(table)

	for _, association := range associations {
		tx.Preload(association)
	}

	return tx
}
