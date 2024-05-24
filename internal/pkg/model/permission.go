package model

import "time"

type Permission struct {
	Id        uint `gorm:"<-:false"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
