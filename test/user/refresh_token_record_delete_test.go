package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/test/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type refreshTokenRecordDeleteTestSuite struct {
	test.TestSuite
}

func (suite *refreshTokenRecordDeleteTestSuite) TestRefreshTokenRecordDelete1() {
	req := httptest.NewRequest("DELETE", "/scheduler/refresh-token-record", nil)
	suite.AddRootAccessToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)

	assert.Equal(suite.T(), http.StatusNoContent, res.Code, res)
}

func (suite *refreshTokenRecordDeleteTestSuite) TestRefreshTokenRecordDelete2() {
	expireAt, err := time.Parse("2006-01-02", "2100-01-01")
	if err != nil {
		panic(err)
	}
	refreshTokenRecord := &model.RefreshTokenRecord{
		RefreshToken: "1234567890",
		UserId:       1,
		Device:       "web",
		ExpireAt:     expireAt,
	}
	provider.App.DB.Create(refreshTokenRecord)
	req := httptest.NewRequest("DELETE", "/scheduler/refresh-token-record", nil)
	suite.AddRootAccessToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)

	db := provider.App.DB.
		Where("refresh_token = ?", "1234567890").
		First(&model.RefreshTokenRecord{})

	assert.Equal(suite.T(), http.StatusNoContent, res.Code, res)
	assert.Equal(suite.T(), 1, int(db.RowsAffected))
}

func (suite *userLogoutTestSuite) TestRefreshTokenRecordDeleteAuthenticationError() {
	req := httptest.NewRequest("DELETE", "/scheduler/refresh-token-record", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusUnauthorized, res.Code, resBody)
	assert.EqualError(suite.T(), response.ErrJwtAuthenticationFailed, resBody["message"].(string), resBody)
}

func TestRefreshTokenRecordDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(refreshTokenRecordDeleteTestSuite))
}
