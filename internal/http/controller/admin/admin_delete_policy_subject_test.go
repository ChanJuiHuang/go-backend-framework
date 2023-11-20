package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AdminDeletePolicySubjectTestSuite struct {
	suite.Suite
}

func (suite *AdminDeletePolicySubjectTestSuite) SetupTest() {
	test.Migration.Run()
	test.AdminRegister()
}

func (suite *AdminDeletePolicySubjectTestSuite) TestDeletePolicySubject() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken, _ := test.AdminLogin()

	role1 := "role1"
	enforcer := provider.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	enforcer.AddPolicies([][]string{
		{role1, "/api1", "GET"},
	})
	adminDeletePolicySubjectRequest := admin.AdminDeletePolicySubjectRequest{
		Subjects: []string{role1},
	}
	reqBody, err := json.Marshal(adminDeletePolicySubjectRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("DELETE", "/api/admin/policy/subject", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &admin.AdminDeletePolicySubjectResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.Equal(suite.T(), []string{"admin"}, respBody.Subjects)
}

func (suite *AdminDeletePolicySubjectTestSuite) TestCsrfMismatch() {
	test.AdminAddPolicies()
	test.AdminAddRole()

	req := httptest.NewRequest("DELETE", "/api/admin/policy/subject", nil)
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

func (suite *AdminDeletePolicySubjectTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()

	req := httptest.NewRequest("DELETE", "/api/admin/policy/subject", nil)
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

func (suite *AdminDeletePolicySubjectTestSuite) TestAuthorizationFailed() {
	accessToken, _ := test.AdminLogin()

	req := httptest.NewRequest("DELETE", "/api/admin/policy/subject", nil)
	test.AddBearerToken(req, accessToken)
	test.AddCsrfToken(req)
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

func (suite *AdminDeletePolicySubjectTestSuite) TearDownTest() {
	test.Migration.Reset()
}

func TestAdminDeletePolicySubjectTestSuite(t *testing.T) {
	suite.Run(t, new(AdminDeletePolicySubjectTestSuite))
}
