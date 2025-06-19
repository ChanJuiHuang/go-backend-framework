package database

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/service"
	"gorm.io/gorm"
)

func NewTx(associations ...string) *gorm.DB {
	tx := service.Registry.Get("database").(*gorm.DB)
	for _, association := range associations {
		tx = tx.Preload(association)
	}

	return tx
}

func NewTxByTable(table string, associations ...string) *gorm.DB {
	database := service.Registry.Get("database").(*gorm.DB)
	tx := database.Table(table)

	for _, association := range associations {
		tx = tx.Preload(association)
	}

	return tx
}
