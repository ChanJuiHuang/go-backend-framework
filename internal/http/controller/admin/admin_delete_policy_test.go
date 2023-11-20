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

type AdminDeletePolicyTestSuite struct {
	suite.Suite
}

func (suite *AdminDeletePolicyTestSuite) SetupTest() {
	test.Migration.Run()
	test.AdminRegister()
}

func (suite *AdminDeletePolicyTestSuite) TestDeletePolicy() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken, _ := test.AdminLogin()

	subject := "role1"
	enforcer := provider.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	_, err := enforcer.AddPolicies([][]string{
		{subject, "/api1", "GET"},
		{subject, "/api2", "GET"},
	})
	if err != nil {
		panic(err)
	}

	adminDeletePolicyRequest := admin.AdminDeletePolicyRequest{
		Subject: subject,
		Rules: []admin.Rule{
			{
				Object: "/api1",
				Action: "GET",
			},
		},
	}
	reqBody, err := json.Marshal(adminDeletePolicyRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("DELETE", "/api/admin/policy", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &admin.AdminCreatePolicyResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.Equal(suite.T(), subject, respBody.Subject)
	assert.Equal(suite.T(), []admin.Rule{
		{
			Object: "/api2",
			Action: "GET",
		},
	}, respBody.Rules)
}

func (suite *AdminDeletePolicyTestSuite) TestRequestValidationFailed() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken, _ := test.AdminLogin()
	reqBody, err := json.Marshal(new(admin.AdminDeletePolicyRequest))
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("DELETE", "/api/admin/policy", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(suite.T(), response.RequestValidationFailed, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.RequestValidationFailed], respBody.Code)
}

func (suite *AdminDeletePolicyTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("DELETE", "/api/admin/policy", nil)
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

func (suite *AdminDeletePolicyTestSuite) TestCsrfMismatch() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("DELETE", "/api/admin/policy", nil)
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

func (suite *AdminDeletePolicyTestSuite) TestAuthorizationFailed() {
	accessToken, _ := test.AdminLogin()
	req := httptest.NewRequest("DELETE", "/api/admin/policy", nil)
	test.AddCsrfToken(req)
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

func (suite *AdminDeletePolicyTestSuite) TearDownTest() {
	test.Migration.Reset()
}

func TestAdminDeletePolicyTestSuite(t *testing.T) {
	suite.Run(t, new(AdminDeletePolicyTestSuite))
}
