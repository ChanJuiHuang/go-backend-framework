package permission_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/stretchr/testify/suite"
)

type PermissionReloadTestSuite struct {
	suite.Suite
}

func (suite *PermissionReloadTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminRegister()
}

func (suite *PermissionReloadTestSuite) Test() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken := test.AdminLogin()

	req := httptest.NewRequest("POST", "/api/admin/permission/reload", nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	suite.Equal(http.StatusNoContent, resp.Code)
}

func (suite *PermissionReloadTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("POST", "/api/admin/permission/reload", nil)
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *PermissionReloadTestSuite) TestCsrfMismatch() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("POST", "/api/admin/permission/reload", nil)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *PermissionReloadTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminLogin()
	req := httptest.NewRequest("POST", "/api/admin/permission/reload", nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *PermissionReloadTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestPermissionReloadTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionReloadTestSuite))
}
