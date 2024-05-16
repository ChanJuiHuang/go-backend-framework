package database

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"gorm.io/gorm"
)

func NewTx() *gorm.DB {
	return service.Registry.Get("database").(*gorm.DB)
}

func NewTxByTable(table string, associations ...string) *gorm.DB {
	database := service.Registry.Get("database").(*gorm.DB)
	tx := database.Table(table)

	for _, association := range associations {
		tx.Preload(association)
	}

	return tx
}
