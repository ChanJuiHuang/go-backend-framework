package model

import (
	"time"
)

type UserRole struct {
	Id        uint `gorm:"<-:false"`
	UserId    uint
	RoleId    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
