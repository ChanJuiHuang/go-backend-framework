package user

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/test/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type userSearchTestSuite struct {
	test.TestSuite
}

func (suite *userSearchTestSuite) TestUserSearch1() {
	req := httptest.NewRequest("GET", "/api/user", nil)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	resBody := &user.UserSearchResponse{}
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusOK, res.Code, resBody)
	assert.Equal(suite.T(), 2, len(resBody.Data))
	assert.Equal(suite.T(), 1, resBody.Page)
	assert.Equal(suite.T(), 10, resBody.PerPage)
	assert.EqualValues(suite.T(), 1, resBody.LastPage)
	assert.EqualValues(suite.T(), 2, resBody.Total)
}

func (suite *userSearchTestSuite) TestUserSearch2() {
	req := httptest.NewRequest("GET", "/api/user", nil)
	queryStrings := url.Values{
		"page":         []string{"1"},
		"per_page":     []string{"1"},
		"filter.name":  []string{"root"},
		"filter.email": []string{"root@root.com"},
	}
	req.URL.RawQuery = queryStrings.Encode()
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	resBody := &user.UserSearchResponse{}
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(suite.T(), http.StatusOK, res.Code, resBody)
	assert.Equal(suite.T(), 1, len(resBody.Data))
	assert.Equal(suite.T(), 1, resBody.Page)
	assert.Equal(suite.T(), 1, resBody.PerPage)
	assert.EqualValues(suite.T(), 1, resBody.LastPage)
	assert.EqualValues(suite.T(), 1, resBody.Total)
	assert.Equal(suite.T(), "root", resBody.Data[0].Name)
	assert.Equal(suite.T(), "root@root.com", resBody.Data[0].EmailUser.Email)
}

func TestUserSearchTestSuite(t *testing.T) {
	suite.Run(t, new(userSearchTestSuite))
}
