package user

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/test/internal/mockery"
	"github.com/ChanJuiHuang/go-backend-framework/test/internal/test"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"gorm.io/gorm"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	userModule "github.com/ChanJuiHuang/go-backend-framework/app/module/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type oauthLoginTestSuite struct {
	test.TestSuite
}

func (suite *oauthLoginTestSuite) TestGoogleOauthLogin1() {
	mockedOauth := new(mockery.MockedOauth)
	mockedOauth.On("GetAccessToken", mockery.MockedOauthCode).Return(nil)
	mockedOauth.On("SetIdAndEmail").Return(nil)
	util.NewOauthProvider = func(provider util.OauthProvider) util.OauthInterface {
		return mockedOauth
	}
	userModule.GetUserId = func(oauthProvider util.OauthInterface) (uint, error) {
		return 0, nil
	}
	userModule.CreateOauthUser = func(oauthProvider util.OauthInterface) (uint, error) {
		mockedOauthProvider := oauthProvider.(*mockery.MockedOauth)
		user := new(model.User)
		err := provider.App.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(user).Error; err != nil {
				return util.WrapError(err)
			}
			db := tx.Create(&model.EmailUser{
				Email:    mockedOauthProvider.Email,
				Password: util.RandomString(16),
				UserId:   user.Id,
			})
			if db.Error != nil {
				return util.WrapError(db.Error)
			}

			db = tx.Create(&model.GoogleUser{
				GoogleId: mockedOauthProvider.Id,
				UserId:   user.Id,
			})
			if db.Error != nil {
				return util.WrapError(db.Error)
			}

			return nil
		})

		return user.Id, err
	}

	requestBody, err := provider.App.Json.Marshal(user.OauthLoginRequest{
		Code: mockery.MockedOauthCode,
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/oauth/google/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	mockedOauth.AssertExpectations(suite.T())
	assert.Equal(suite.T(), http.StatusOK, res.Code, res)
	assert.Contains(suite.T(), resBody, "access_token", res.Code)
	assert.Contains(suite.T(), resBody, "refresh_token", res.Code)
}

func (suite *oauthLoginTestSuite) TestGoogleOauthLogin2() {
	mockedOauth := new(mockery.MockedOauth)
	mockedOauth.On("GetAccessToken", mockery.MockedOauthCode).Return(nil)
	mockedOauth.On("SetIdAndEmail").Return(nil)
	util.NewOauthProvider = func(provider util.OauthProvider) util.OauthInterface {
		return mockedOauth
	}
	userModule.GetUserId = func(oauthProvider util.OauthInterface) (uint, error) {
		mockedOauthProvider := oauthProvider.(*mockery.MockedOauth)
		googleUser := new(model.GoogleUser)
		db := provider.App.DB.Where("google_id = ?", mockedOauthProvider.Id).First(googleUser)

		if util.IsDatabaseError(db) {
			return 0, util.WrapError(db.Error)
		}

		return googleUser.UserId, nil
	}

	requestBody, err := provider.App.Json.Marshal(user.OauthLoginRequest{
		Code: mockery.MockedOauthCode,
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/oauth/google/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	mockedOauth.AssertExpectations(suite.T())
	assert.Equal(suite.T(), http.StatusOK, res.Code, res)
	assert.Contains(suite.T(), resBody, "access_token", res.Code)
	assert.Contains(suite.T(), resBody, "refresh_token", res.Code)
}

func (suite *oauthLoginTestSuite) TestGoogleOauthLoginWithFailOnGetAccessToken() {
	mockedOauth := new(mockery.MockedOauth)
	mockedOauth.On("GetAccessToken", mockery.MockedOauthCode).Return(errors.New("oauth login failed"))
	util.NewOauthProvider = func(provider util.OauthProvider) util.OauthInterface {
		return mockedOauth
	}

	requestBody, err := provider.App.Json.Marshal(user.OauthLoginRequest{
		Code: mockery.MockedOauthCode,
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/oauth/google/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	mockedOauth.AssertExpectations(suite.T())
	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
}

func (suite *oauthLoginTestSuite) TestGoogleOauthLoginFailOnSetIdAndEmail() {
	mockedOauth := new(mockery.MockedOauth)
	mockedOauth.On("GetAccessToken", mockery.MockedOauthCode).Return(nil)
	mockedOauth.On("SetIdAndEmail").Return(errors.New("oauth login failed"))
	util.NewOauthProvider = func(provider util.OauthProvider) util.OauthInterface {
		return mockedOauth
	}

	requestBody, err := provider.App.Json.Marshal(user.OauthLoginRequest{
		Code: mockery.MockedOauthCode,
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/oauth/google/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	mockedOauth.AssertExpectations(suite.T())
	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
}

func (suite *oauthLoginTestSuite) TestOauthLoginCsrfError() {
	req := httptest.NewRequest("POST", "/api/oauth/google/token", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
	assert.EqualError(suite.T(), response.ErrCsrfTokenMismatch, resBody["message"].(string), res.Body)
}

func (suite *oauthLoginTestSuite) TestOauthLoginValidationError() {
	requestBody, err := provider.App.Json.Marshal(user.OauthLoginRequest{})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/oauth/google/token", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusBadRequest, res.Code, res)
	assert.EqualError(suite.T(), response.ErrRequestValidationFailed, resBody["message"].(string), res.Body)
}

func TestOauthLoginTestSuite(t *testing.T) {
	suite.Run(t, new(oauthLoginTestSuite))
}
