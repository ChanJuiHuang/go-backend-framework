package util

import (
	"errors"
)

type OauthInterface interface {
	GetAccessToken(code string) error
	SetIdAndEmail() error
}

type OauthProvider string

const (
	Google OauthProvider = "google"
)

var errEmailDoesNotExist error = errors.New("email does not exist")
var errEmailDoesNotVerify error = errors.New("email does not verify")

var NewOauthProvider func(provider OauthProvider) OauthInterface = newOauthProvider

func newOauthProvider(provider OauthProvider) OauthInterface {
	switch provider {
	case Google:
		return new(GoogleOauth)
	default:
		panic("oauth provider does not exist")
	}
}
