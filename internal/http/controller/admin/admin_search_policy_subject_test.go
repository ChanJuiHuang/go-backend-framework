package admin_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

type AdminSearchPolicySubjectTestSuite struct {
	suite.Suite
}

func (suite *AdminSearchPolicySubjectTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminRegister()
}

func (suite *AdminSearchPolicySubjectTestSuite) TestSearchPolicySubject() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken := test.AdminLogin()

	req := httptest.NewRequest("GET", "/api/admin/policy/subject", nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &admin.AdminSearchPolicySubjectData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	subjects := enforcer.GetAllSubjects()

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.Equal(suite.T(), len(subjects), len(data.Subjects))
}

func (suite *AdminSearchPolicySubjectTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()

	req := httptest.NewRequest("GET", "/api/admin/policy/subject", nil)
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

func (suite *AdminSearchPolicySubjectTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminLogin()
	req := httptest.NewRequest("GET", "/api/admin/policy/subject", nil)
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

func (suite *AdminSearchPolicySubjectTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestAdminSearchPolicySubjectTestSuite(t *testing.T) {
	suite.Run(t, new(AdminSearchPolicySubjectTestSuite))
}
