package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	userModule "github.com/ChanJuiHuang/go-backend-framework/app/module/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/test/internal/test"
	"github.com/ChanJuiHuang/go-backend-framework/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type refreshTokenTestSuite struct {
	test.TestSuite
}

func (suite *refreshTokenTestSuite) TestRefreshToken() {
	requestBody, err := provider.App.Json.Marshal(user.TokenRefreshRequest{
		AccessToken:  suite.AccessToken,
		RefreshToken: suite.RefreshToken,
		Device:       "web",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("PUT", "/api/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusOK, res.Code, res)
	assert.Contains(suite.T(), resBody, "access_token")
	assert.Contains(suite.T(), resBody, "refresh_token")
}

func (suite *refreshTokenTestSuite) TestRefreshTokenError1() {
	requestBody, err := provider.App.Json.Marshal(user.TokenRefreshRequest{
		AccessToken:  suite.AccessToken,
		RefreshToken: suite.AccessToken,
		Device:       "web",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("PUT", "/api/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrJwtAuthenticationFailed, resBody["message"].(string), resBody)
	assert.EqualError(suite.T(), userModule.ErrJwtAudClaimIsNotRefreshString, resBody["previous_message"].(string), resBody)
}

func (suite *refreshTokenTestSuite) TestRefreshTokenError2() {
	util.GetRefreshTokenLimit = func() uint {
		return 0
	}
	requestBody, err := provider.App.Json.Marshal(user.TokenRefreshRequest{
		AccessToken:  suite.AccessToken,
		RefreshToken: suite.RefreshToken,
		Device:       "web",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("PUT", "/api/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrJwtAuthenticationFailed, resBody["message"].(string), resBody)
	assert.EqualError(suite.T(), userModule.ErrOverRefreshTokenLimit, resBody["previous_message"].(string), resBody)
	util.GetRefreshTokenLimit = func() uint {
		return 60
	}
}

func (suite *refreshTokenTestSuite) TestRefreshTokenCsrfError() {
	req := httptest.NewRequest("PUT", "/api/token", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
	assert.EqualError(suite.T(), response.ErrCsrfTokenMismatch, resBody["message"].(string), res.Body)
}

func (suite *refreshTokenTestSuite) TestRefreshTokenValidationError() {
	req := httptest.NewRequest("PUT", "/api/token", nil)
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrRequestValidationFailed, resBody["message"].(string), res.Body)
}

func TestRefreshTokenTestSuite(t *testing.T) {
	suite.Run(t, new(refreshTokenTestSuite))
}
