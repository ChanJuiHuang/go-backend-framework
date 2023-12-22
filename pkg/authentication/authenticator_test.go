package authentication

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticatorTestSuite struct {
	suite.Suite
	Authenticator
}

func (suite *AuthenticatorTestSuite) SetupSuite() {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	assert.ErrorIs(suite.T(), err, nil, err)
	suite.publicKey = publicKey
	suite.privateKey = privateKey
	suite.accessTokenLifeTime = time.Hour
	suite.refreshTokenLifeTime = time.Hour
}

func (suite *AuthenticatorTestSuite) TestVerifyAccessToken() {
	subject := "john"
	accessToken, err := suite.IssueAccessToken(subject)
	assert.ErrorIs(suite.T(), err, nil, err)

	token, err := suite.VerifyJwt(accessToken)
	assert.ErrorIs(suite.T(), err, nil, err)
	assert.True(suite.T(), token.Valid)

	claims := token.Claims.(jwt.MapClaims)
	assert.Equal(suite.T(), subject, claims["sub"].(string))
	assert.Equal(suite.T(), "access", claims["aud"].([]any)[0].(string))
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
	assert.ErrorIs(suite.T(), err, nil, err)
	time.Sleep(time.Second)

	token, err := suite.VerifyJwt(invalidJwt)
	assert.NotEmpty(suite.T(), err)
	assert.Empty(suite.T(), token)
}

func TestAuthenticatorTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticatorTestSuite))
}
