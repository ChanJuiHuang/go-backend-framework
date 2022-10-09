package example

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/test/internal/test"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/golang-jwt/jwt/v4"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type exampleTestSuite struct {
	test.TestSuite
}

func (suite *exampleTestSuite) TestPing() {
	req := httptest.NewRequest("GET", "/api/ping", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusOK, res.Code, res)
	assert.Equal(suite.T(), map[string]any{"message": "pong"}, resBody, res)
}

func (suite *exampleTestSuite) TestWelcome() {
	req := httptest.NewRequest("GET", "/scheduler/welcome", nil)
	suite.AddRootAccessToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)

	assert.Equal(suite.T(), http.StatusOK, res.Code, res)
	assert.Equal(suite.T(), `{"message":"welcome"}`, res.Body.String(), res)
}

func (suite *exampleTestSuite) TestWelcomeWithoutRootToken() {
	req := httptest.NewRequest("GET", "/scheduler/welcome", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)

	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, res)
	assert.EqualError(suite.T(), response.ErrJwtAuthenticationFailed, resBody["message"].(string), res)
}

func (suite *exampleTestSuite) TestWelcomeWithoutPolicy() {
	req := httptest.NewRequest("GET", "/scheduler/welcome", nil)
	req.Header.Add("Authorization", "Bearer "+suite.RootAccessToken)
	provider.App.Casbin.RemoveGroupingPolicy("1", "root")
	provider.App.Casbin.LoadPolicy()
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)

	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, res)
	assert.EqualError(suite.T(), response.ErrAuthorizationFailed, resBody["message"].(string), res)
}

func (suite *exampleTestSuite) TestWelcomeWithExpirationToken() {
	req := httptest.NewRequest("GET", "/scheduler/welcome", nil)
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"access"},
		Subject:   "1",
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Nanosecond)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := util.IssueJwtWithConfig(claims)
	req.Header.Add("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)

	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, res)
	assert.EqualError(suite.T(), response.ErrJwtAuthenticationFailed, resBody["message"].(string), res)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(exampleTestSuite))
}
