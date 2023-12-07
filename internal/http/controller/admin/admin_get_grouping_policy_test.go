package admin_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AdminGetGroupingPolicyTestSuite struct {
	suite.Suite
}

func (suite *AdminGetGroupingPolicyTestSuite) SetupTest() {
	test.Migration.Run()
	test.AdminRegister()
}

func (suite *AdminGetGroupingPolicyTestSuite) TestSearchPolicySubject() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken, _ := test.AdminLogin()
	userId := "1"

	req := httptest.NewRequest("GET", "/api/admin/grouping-policy/"+userId, nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &admin.AdminGetGroupingPolicyData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	subjects := enforcer.GetFilteredGroupingPolicy(0, userId)
	id, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.Equal(suite.T(), uint(id), data.UserId)
	assert.Equal(suite.T(), len(subjects), len(data.Subjects))
}

func (suite *AdminGetGroupingPolicyTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	userId := "1"

	req := httptest.NewRequest("GET", "/api/admin/grouping-policy/"+userId, nil)
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(suite.T(), response.Unauthorized, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *AdminGetGroupingPolicyTestSuite) TestAuthorizationFailed() {
	accessToken, _ := test.AdminLogin()
	userId := "1"

	req := httptest.NewRequest("GET", "/api/admin/grouping-policy/"+userId, nil)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusForbidden, resp.Code)
	assert.Equal(suite.T(), response.Forbidden, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *AdminGetGroupingPolicyTestSuite) TearDownTest() {
	test.Migration.Reset()
}

func TestAdminGetGroupingPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(AdminGetGroupingPolicyTestSuite))
}
