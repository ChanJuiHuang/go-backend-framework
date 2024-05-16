package model

import (
	"time"
)

type RolePermission struct {
	Id           uint `gorm:"<-:false"`
	RoleId       uint
	PermissionId uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
