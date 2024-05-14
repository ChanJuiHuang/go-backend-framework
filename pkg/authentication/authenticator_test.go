package authentication

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/suite"
)

type AuthenticatorTestSuite struct {
	suite.Suite
	Authenticator
}

func (suite *AuthenticatorTestSuite) SetupSuite() {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	suite.ErrorIs(err, nil, err)
	suite.publicKey = publicKey
	suite.privateKey = privateKey
	suite.accessTokenLifeTime = time.Hour
	suite.refreshTokenLifeTime = time.Hour
}

func (suite *AuthenticatorTestSuite) TestVerifyAccessToken() {
	subject := "john"
	accessToken, err := suite.IssueAccessToken(subject)
	suite.ErrorIs(err, nil, err)

	token, err := suite.VerifyJwt(accessToken)
	suite.ErrorIs(err, nil, err)
	suite.True(token.Valid)

	claims := token.Claims.(jwt.MapClaims)
	suite.Equal(subject, claims["sub"].(string))
	suite.Equal("access", claims["aud"].([]any)[0].(string))
}

func (suite *AuthenticatorTestSuite) TestInvalidJwtToken() {
	subject := "john"
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"invalid"},
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	invalidJwt, err := suite.IssueJwt(claims)
	suite.ErrorIs(err, nil, err)
	time.Sleep(time.Second)

	token, err := suite.VerifyJwt(invalidJwt)
	suite.NotEmpty(err)
	suite.Empty(token)
}

func TestAuthenticatorTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticatorTestSuite))
}
