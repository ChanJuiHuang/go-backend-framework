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
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"
)

type AdminCreateGroupingPolicyTestSuite struct {
	suite.Suite
}

func (suite *AdminCreateGroupingPolicyTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminRegister()
}

func (suite *AdminCreateGroupingPolicyTestSuite) TestCreateGroupingPolicy() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken := test.AdminLogin()

	subject := "role1"
	policies := [][]string{
		{subject, "/api1", "GET"},
	}
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	_, err := enforcer.AddPolicies(policies)
	if err != nil {
		panic(err)
	}

	var userId uint = 999
	subjects := []string{subject}
	adminCreateGroupingPolicyRequest := admin.AdminCreateGroupingPolicyRequest{
		UserId:   userId,
		Subjects: subjects,
	}
	reqBody, err := json.Marshal(adminCreateGroupingPolicyRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/grouping-policy", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &admin.AdminCreateGroupingPolicyData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(userId, data.UserId)
	suite.Equal(subjects, data.Subjects)
}

func (suite *AdminCreateGroupingPolicyTestSuite) TestRequestValidationFailed() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken := test.AdminLogin()
	reqBody, err := json.Marshal(new(admin.AdminCreateGroupingPolicyRequest))
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/grouping-policy", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.RequestValidationFailed, respBody.Message)
	suite.Equal(response.MessageToCode[response.RequestValidationFailed], respBody.Code)
}

func (suite *AdminCreateGroupingPolicyTestSuite) TestOneOfGroupingPolicyIsRepeat() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken := test.AdminLogin()

	var userId uint = 1
	subjects := []string{"admin"}
	adminCreateGroupingPolicyRequest := admin.AdminCreateGroupingPolicyRequest{
		UserId:   userId,
		Subjects: subjects,
	}
	reqBody, err := json.Marshal(adminCreateGroupingPolicyRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/grouping-policy", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.OneOfGroupingPolicyIsRepeat, respBody.Message)
	suite.Equal(response.MessageToCode[response.OneOfGroupingPolicyIsRepeat], respBody.Code)
}

func (suite *AdminCreateGroupingPolicyTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("POST", "/api/admin/grouping-policy", nil)
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

func (suite *AdminCreateGroupingPolicyTestSuite) TestCsrfMismatch() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("POST", "/api/admin/grouping-policy", nil)
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

func (suite *AdminCreateGroupingPolicyTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminLogin()
	req := httptest.NewRequest("POST", "/api/admin/grouping-policy", nil)
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

func (suite *AdminCreateGroupingPolicyTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestAdminCreateGroupingPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(AdminCreateGroupingPolicyTestSuite))
}
