package model

import "time"

type EmailUser struct {
	Email     string
	Password  string
	UserId    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
