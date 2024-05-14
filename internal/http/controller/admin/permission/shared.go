package permission

import (
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type PermissionData struct {
	Id        uint      `json:"id" mapstructure:"id" validate:"required"`
	Name      string    `json:"name" mapstructure:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at" validate:"required" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at" validate:"required" format:"date-time"`
}

func (data *PermissionData) Fill(m *model.Permission) {
	data.Id = m.Id
	data.Name = m.Name
	data.CreatedAt = m.CreatedAt
	data.UpdatedAt = m.UpdatedAt
}

type HttpApiData struct {
	Path   string `json:"path" mapstructure:"path" validate:"required"`
	Method string `json:"method" mapstructure:"method" validate:"required"`
}

func (data *HttpApiData) Fill(casbinRule gormadapter.CasbinRule) {
	data.Path = casbinRule.V1
	data.Method = casbinRule.V2
}
