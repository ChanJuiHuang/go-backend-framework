package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/test/internal/test"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type emailLoginTestSuite struct {
	test.TestSuite
}

func (suite *emailLoginTestSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()
	requestBody, err := provider.App.Json.Marshal(user.UserCreateRequest{
		Name:     "bob",
		Email:    "bob@gmail.com",
		Password: "abcABC123",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/user", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
}

func (suite *emailLoginTestSuite) TestEmailLogin() {
	requestBody, err := provider.App.Json.Marshal(user.UserCreateRequest{
		Email:    "bob@gmail.com",
		Password: "abcABC123",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusOK, res.Code, res)
	assert.Contains(suite.T(), resBody, "access_token", res.Code)
	assert.Contains(suite.T(), resBody, "refresh_token", res.Code)
}

func (suite *emailLoginTestSuite) TestEmailLoginCsrfError() {
	req := httptest.NewRequest("POST", "/api/token", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
	assert.EqualError(suite.T(), response.ErrCsrfTokenMismatch, resBody["message"].(string), res.Body)
}

func (suite *emailLoginTestSuite) TestEmailLoginValidationError() {
	requestBody, err := provider.App.Json.Marshal(user.UserCreateRequest{
		Email: "bob@gmail.com",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
	assert.EqualError(suite.T(), response.ErrRequestValidationFailed, resBody["message"].(string), res.Body)
}

func (suite *emailLoginTestSuite) TestEmailLoginPasswordError() {
	requestBody, err := provider.App.Json.Marshal(user.UserCreateRequest{
		Email:    "bob@gmail.com",
		Password: "123456",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
	assert.EqualError(suite.T(), response.ErrLoginPasswordIsWrong, resBody["message"].(string), res.Body)
}

func TestEmailLoginTestSuite(t *testing.T) {
	suite.Run(t, new(emailLoginTestSuite))
}
