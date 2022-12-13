package util

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

type googleOauthConfig struct {
	ClientId     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	RedirectUrl  string `split_words:"true"`
}

type GoogleOauth struct {
	Token *oauth2.Token
	Id    string
	Email string
}

var googleOauthCfg *googleOauthConfig

func init() {
	if googleOauthCfg == nil {
		googleOauthCfg = new(googleOauthConfig)
		err := envconfig.Process("google_oauth", googleOauthCfg)

		if err != nil {
			panic(err)
		}
	}
}

func (googleOauth *GoogleOauth) GetAccessToken(code string) error {
	oauth2Config := &oauth2.Config{
		ClientID:     googleOauthCfg.ClientId,
		ClientSecret: googleOauthCfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: googleOauthCfg.RedirectUrl,
	}
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return WrapError(err)
	}
	googleOauth.Token = token

	return nil
}

func (googleOauth *GoogleOauth) SetIdAndEmail() error {
	idToken := googleOauth.Token.Extra("id_token").(string)
	token, _, err := jwt.NewParser().ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return WrapError(err)
	}

	claims := token.Claims.(jwt.MapClaims)
	if !claims["email_verified"].(bool) {
		return WrapError(errEmailDoesNotVerify)
	}

	email := claims["email"].(string)
	if email == "" {
		return WrapError(errEmailDoesNotExist)
	}
	googleOauth.Id = claims["sub"].(string)
	googleOauth.Email = email

	return nil
}
