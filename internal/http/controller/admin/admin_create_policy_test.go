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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AdminCreatePolicyTestSuite struct {
	suite.Suite
}

func (suite *AdminCreatePolicyTestSuite) SetupTest() {
	test.Migration.Run()
	test.AdminRegister()
}

func (suite *AdminCreatePolicyTestSuite) TestCreatePolicy() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken, _ := test.AdminLogin()

	subject := "role1"
	rules := []admin.Rule{
		{
			Object: "/api1",
			Action: "POST",
		},
		{
			Object: "/api2",
			Action: "GET",
		},
	}
	adminCreatePolicyRequest := admin.AdminCreatePolicyRequest{
		Subject: subject,
		Rules:   rules,
	}
	reqBody, err := json.Marshal(adminCreatePolicyRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/policy", bytes.NewReader(reqBody))
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
	assert.Equal(suite.T(), rules, respBody.Rules)
}

func (suite *AdminCreatePolicyTestSuite) TestRequestValidationFailed() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken, _ := test.AdminLogin()
	reqBody, err := json.Marshal(new(admin.AdminCreatePolicyRequest))
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/policy", bytes.NewReader(reqBody))
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

func (suite *AdminCreatePolicyTestSuite) TestOneOfPolicyIsRepeat() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken, _ := test.AdminLogin()
	adminCreatePolicyRequest := admin.AdminCreatePolicyRequest{
		Subject: "role1",
		Rules: []admin.Rule{
			{
				Object: "/api1",
				Action: "POST",
			},
		},
	}
	reqBody, err := json.Marshal(adminCreatePolicyRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/policy", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	req = httptest.NewRequest("POST", "/api/admin/policy", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp = httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(suite.T(), response.OneOfPolicyIsRepeat, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.OneOfPolicyIsRepeat], respBody.Code)
}

func (suite *AdminCreatePolicyTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("POST", "/api/admin/policy", nil)
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

func (suite *AdminCreatePolicyTestSuite) TestCsrfMismatch() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("POST", "/api/admin/policy", nil)
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

func (suite *AdminCreatePolicyTestSuite) TestAuthorizationFailed() {
	accessToken, _ := test.AdminLogin()
	req := httptest.NewRequest("POST", "/api/admin/policy", nil)
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

func (suite *AdminCreatePolicyTestSuite) TearDownTest() {
	test.Migration.Reset()
}

func TestAdminCreatePolicyTestSuite(t *testing.T) {
	suite.Run(t, new(AdminCreatePolicyTestSuite))
}
