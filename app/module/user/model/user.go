package model

import "time"

type User struct {
	Id        uint `gorm:"<-:false"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	EmailUser EmailUser
}
