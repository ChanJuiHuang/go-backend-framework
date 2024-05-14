package admin_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"
)

type AdminGetPolicySubjectUserTestSuite struct {
	suite.Suite
}

func (suite *AdminGetPolicySubjectUserTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminRegister()
}

func (suite *AdminGetPolicySubjectUserTestSuite) TestGetPolicySubjectUser() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken := test.AdminLogin()
	subject := "admin"

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/admin/policy/subject/%s/user", subject), nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &admin.AdminGetPolicySubjectUserData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.UserIds)
}

func (suite *AdminGetPolicySubjectUserTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	subject := "admin"

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/admin/policy/subject/%s/user", subject), nil)
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

func (suite *AdminGetPolicySubjectUserTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminLogin()
	subject := "admin"

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/admin/policy/subject/%s/user", subject), nil)
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

func (suite *AdminGetPolicySubjectUserTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestAdminGetPolicySubjectUserTestSuite(t *testing.T) {
	suite.Run(t, new(AdminGetPolicySubjectUserTestSuite))
}
