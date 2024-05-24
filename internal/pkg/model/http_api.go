package model

import "time"

type HttpApi struct {
	Id        uint `gorm:"<-:false"`
	Method    string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
