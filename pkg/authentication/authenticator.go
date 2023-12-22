package authentication

import (
	"crypto/ed25519"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Authenticator struct {
	privateKey           ed25519.PrivateKey
	publicKey            ed25519.PublicKey
	accessTokenLifeTime  time.Duration
	refreshTokenLifeTime time.Duration
}

func NewAuthenticator(config Config) (*Authenticator, error) {
	privateKey, err := base64.RawURLEncoding.DecodeString(config.PrivateKey)
	if err != nil {
		return nil, err
	}
	publicKey, err := base64.RawURLEncoding.DecodeString(config.PublicKey)
	if err != nil {
		return nil, err
	}
	authenticator := &Authenticator{
		privateKey:           privateKey,
		publicKey:            publicKey,
		accessTokenLifeTime:  config.AccessTokenLifeTime,
		refreshTokenLifeTime: config.RefreshTokenLifeTime,
	}

	return authenticator, nil
}

func (auth *Authenticator) IssueJwt(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenString, err := token.SignedString(auth.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (auth *Authenticator) VerifyJwt(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return auth.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (auth *Authenticator) IssueAccessToken(subject string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"access"},
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(now.Add(auth.accessTokenLifeTime)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	return auth.IssueJwt(claims)
}
