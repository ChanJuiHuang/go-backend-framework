package model

import (
	"time"
)

type Role struct {
	Id          uint `gorm:"<-:false"`
	Name        string
	IsPublic    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Permissions []Permission `gorm:"many2many:role_permissions"`
	Users       []User       `gorm:"many2many:user_roles"`
}
