package model

import "time"

type GoogleUser struct {
	GoogleId  string
	UserId    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
