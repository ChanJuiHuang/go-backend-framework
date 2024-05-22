package model

import "time"

type User struct {
	Id        uint `gorm:"<-:false"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Roles     []Role `gorm:"many2many:user_roles"`
}
