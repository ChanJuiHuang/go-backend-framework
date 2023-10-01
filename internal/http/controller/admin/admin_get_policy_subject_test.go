package admin_test

import (
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

type AdminGetPolicySubjectTestSuite struct {
	suite.Suite
}

func (suite *AdminGetPolicySubjectTestSuite) SetupTest() {
	test.Migration.Run()
	test.AdminRegister()
}

func (suite *AdminGetPolicySubjectTestSuite) TestGetPolicySubject() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken, _ := test.AdminLogin()
	subject := "admin"

	req := httptest.NewRequest("GET", "/api/admin/policy/subject/"+subject, nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &admin.AdminGetPolicySubjectResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.Equal(suite.T(), subject, respBody.Subject)
	assert.NotEmpty(suite.T(), respBody.Rules)
}

func (suite *AdminGetPolicySubjectTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	subject := "admin"

	req := httptest.NewRequest("GET", "/api/admin/policy/subject/"+subject, nil)
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

func (suite *AdminGetPolicySubjectTestSuite) TestAuthorizationFailed() {
	accessToken, _ := test.AdminLogin()
	subject := "admin"

	req := httptest.NewRequest("GET", "/api/admin/policy/subject/"+subject, nil)
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

func (suite *AdminGetPolicySubjectTestSuite) TearDownTest() {
	test.Migration.Reset()
}

func TestAdminGetPolicySubjectTestSuite(t *testing.T) {
	suite.Run(t, new(AdminGetPolicySubjectTestSuite))
}
