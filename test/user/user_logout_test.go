package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/test/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type userLogoutTestSuite struct {
	test.TestSuite
}

func (suite *userLogoutTestSuite) TestUserLogout() {
	requestBody, err := provider.App.Json.Marshal(user.UserLogoutRequest{
		RefreshToken: suite.RefreshToken,
		Device:       "web",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("DELETE", "/api/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	suite.AddAuthorizationHeader(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusNoContent, res.Code, resBody)
}

func (suite *userLogoutTestSuite) TestUserLogoutRefreshTokenError() {
	requestBody, err := provider.App.Json.Marshal(user.UserLogoutRequest{
		Device: "web",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("DELETE", "/api/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	suite.AddAuthorizationHeader(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrJwtAuthenticationFailed, resBody["message"].(string), resBody)
}

func (suite *userLogoutTestSuite) TestRefreshTokenCsrfError() {
	req := httptest.NewRequest("DELETE", "/api/token", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
	assert.EqualError(suite.T(), response.ErrCsrfTokenMismatch, resBody["message"].(string), res.Body)
}

func (suite *userLogoutTestSuite) TestUserLogoutAuthenticationError() {
	req := httptest.NewRequest("DELETE", "/api/token", nil)
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrJwtAuthenticationFailed, resBody["message"].(string), resBody)
}

func (suite *userLogoutTestSuite) TestRefreshTokenValidationError() {
	req := httptest.NewRequest("DELETE", "/api/token", nil)
	suite.AddCsrfToken(req)
	suite.AddAuthorizationHeader(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrRequestValidationFailed, resBody["message"].(string), res.Body)
}

func TestUserLogoutTestSuite(t *testing.T) {
	suite.Run(t, new(userLogoutTestSuite))
}
