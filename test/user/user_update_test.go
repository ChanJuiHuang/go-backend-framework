package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/test/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type userUpdateTestSuite struct {
	test.TestSuite
}

func (suite *userUpdateTestSuite) TestUserUpdate1() {
	requestBody, err := provider.App.Json.Marshal(user.UserUpdateRequest{
		Name:     "john",
		Password: "qweQWE890",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("PUT", "/api/user", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	suite.AddAuthorizationHeader(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)

	u := new(model.User)
	db := provider.App.DB.
		Table("users").
		Select("users.*").
		Joins("join email_users on users.id = email_users.user_id").
		Where("email_users.email = ?", "go@gogogo.com").First(u)
	if db.Error != nil {
		panic(db.Error)
	}

	assert.Equal(suite.T(), http.StatusNoContent, res.Code, res)
	assert.Equal(suite.T(), "john", u.Name, res.Body)
}

func (suite *userUpdateTestSuite) TestUserUpdate2() {
	requestBody, err := provider.App.Json.Marshal(user.UserLoginRequest{
		Email:    "go@gogogo.com",
		Password: "qweQWE890",
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

	assert.Equal(suite.T(), http.StatusOK, res.Code, resBody)
	assert.Contains(suite.T(), resBody, "access_token", res.Code)
	assert.Contains(suite.T(), resBody, "refresh_token", res.Code)
}

func (suite *userUpdateTestSuite) TestUserUpdateCsrfError() {
	req := httptest.NewRequest("PUT", "/api/user", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
	assert.EqualError(suite.T(), response.ErrCsrfTokenMismatch, resBody["message"].(string), res.Body)
}

func (suite *userUpdateTestSuite) TestUpdateAuthenticationError() {
	requestBody, err := provider.App.Json.Marshal(user.UserUpdateRequest{
		Name:     "john",
		Password: "qweQWE890",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("PUT", "/api/user", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrJwtAuthenticationFailed, resBody["message"].(string), res.Body)
}

func (suite *userUpdateTestSuite) TestUserUpdateValidationError() {
	requestBody, err := provider.App.Json.Marshal(user.UserUpdateRequest{
		Name:     "john",
		Password: "123456",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("PUT", "/api/user", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	suite.AddAuthorizationHeader(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrRequestValidationFailed, resBody["message"].(string), res.Body)
}

func TestUserUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(userUpdateTestSuite))
}
