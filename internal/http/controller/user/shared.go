package user

import (
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
)

type UserData struct {
	Id        uint      `json:"id" mapstructure:"id" validate:"required"`
	Name      string    `json:"name" mapstructure:"name" validate:"required"`
	Email     string    `json:"email" mapstructure:"email" validate:"required"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at" validate:"required" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at" validate:"required" format:"date-time"`
}

func (data *UserData) Fill(m *model.User) {
	data.Id = m.Id
	data.Name = m.Name
	data.Email = m.Email
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
}

type TokenData struct {
	AccessToken string `json:"access_token" mapstructure:"access_token" validate:"required"`
}

func (data *TokenData) Fill(accessToken string) {
	data.AccessToken = accessToken
}
