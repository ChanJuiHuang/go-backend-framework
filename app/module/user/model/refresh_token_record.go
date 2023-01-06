package model

import "time"

type RefreshTokenRecord struct {
	RefreshToken string
	UserId       uint
	Device       string
	ExpireAt     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
