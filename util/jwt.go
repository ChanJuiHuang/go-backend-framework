package util

import (
	"crypto/ed25519"
	"encoding/base64"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kelseyhightower/envconfig"
)

type jwtConfig struct {
	PrivateKey           ed25519.PrivateKey `ignored:"true"`
	PublicKey            ed25519.PublicKey  `ignored:"true"`
	AccessTokenLifeTime  time.Duration      `required:"true" split_words:"true"`
	RefreshTokenLifeTime time.Duration      `required:"true" split_words:"true"`
	RefreshTokenLimit    uint               `required:"true" split_words:"true"`
}

var jwtCfg *jwtConfig

func init() {
	if jwtCfg == nil {
		jwtCfg = new(jwtConfig)
		err := envconfig.Process("jwt", jwtCfg)

		if err != nil {
			panic(err)
		}
	}

	privateKeyString, doesExists := os.LookupEnv("JWT_PRIVATE_KEY")
	if !doesExists {
		panic("environment variable [JWT_PRIVATE_KEY] does not exist")
	}
	privateKey, err := base64.RawURLEncoding.DecodeString(privateKeyString)
	if err != nil {
		panic(err)
	}

	publicKeyString, doesExists := os.LookupEnv("JWT_PUBLIC_KEY")
	if !doesExists {
		panic("environment variable [JWT_PUBLIC_KEY] does not exist")
	}
	publicKey, err := base64.RawURLEncoding.DecodeString(publicKeyString)
	if err != nil {
		panic(err)
	}

	jwtCfg.PrivateKey = ed25519.PrivateKey(privateKey)
	jwtCfg.PublicKey = ed25519.PublicKey(publicKey)
}

func IssueJwtWithConfig(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenString, err := token.SignedString(jwtCfg.PrivateKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func formatSubJect(subject any) string {
	var sub string
	switch v := subject.(type) {
	case int:
		sub = strconv.FormatInt(int64(v), 10)
	case int64:
		sub = strconv.FormatInt(v, 10)
	case uint:
		sub = strconv.FormatUint(uint64(v), 10)
	case uint64:
		sub = strconv.FormatUint(v, 10)
	case string:
		sub = v
	default:
		panic("subject type is wrong")
	}

	return sub
}

func IssueAccessToken(subject any) string {
	sub := formatSubJect(subject)
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"access"},
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(jwtCfg.AccessTokenLifeTime)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	return IssueJwtWithConfig(claims)
}

func IssueRefreshToken(subject any) string {
	sub := formatSubJect(subject)
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"refresh"},
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(jwtCfg.RefreshTokenLifeTime)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	return IssueJwtWithConfig(claims)
}

func VerifyJwt(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtCfg.PublicKey, nil
	})
	if err != nil {
		err = WrapError(err)
	}

	return token, err
}

var GetRefreshTokenLimit = func() uint {
	return jwtCfg.RefreshTokenLimit
}
