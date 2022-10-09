package mockery

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type MockedOauth struct {
	mock.Mock
	Token *oauth2.Token
	Id    string
	Email string
}

var MockedOauthCode string = "mockedOauthCode"

func (mockedOauth *MockedOauth) GetAccessToken(code string) error {
	args := mockedOauth.Called(code)
	return args.Error(0)
}

func (mockedOauth *MockedOauth) SetIdAndEmail() error {
	mockedOauth.Id = "1234567890"
	mockedOauth.Email = "test@test.com"
	args := mockedOauth.Called()

	return args.Error(0)
}
