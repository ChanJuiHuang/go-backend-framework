package test

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
)

type adminService struct {
	userService
}

var AdminService *adminService

func NewAdminService() *adminService {
	return &adminService{
		userService: userService{
			User: &model.User{
				Name:  "admin",
				Email: "admin@test.com",
			},
			UserPassword: "abcABC123",
		},
	}
}
