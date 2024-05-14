package permission_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin/permission"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	pkgPermission "github.com/ChanJuiHuang/go-backend-framework/internal/pkg/permission"
	"github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/pagination"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PermissionSearchTestSuite struct {
	suite.Suite
}

func (suite *PermissionSearchTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.AdminRegister()
}

func (suite *PermissionSearchTestSuite) Test() {
	permissionModel := &model.Permission{Name: "permission1"}
	casbinRules := []gormadapter.CasbinRule{
		{
			Ptype: "p",
			V0:    permissionModel.Name,
			V1:    "/api/test-api-1",
			V2:    "GET",
		},
	}

	err := database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := pkgPermission.Create(tx, permissionModel); err != nil {
			return err
		}

		if err := pkgPermission.CreateCasbinRule(tx, casbinRules); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken := test.AdminLogin()

	searchRequest := permission.PermissionSearchRequest{
		Name: permissionModel.Name,
		PaginationRequest: pagination.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}
	queryString := url.Values{}
	encoder := schema.NewEncoder()
	if err := encoder.Encode(searchRequest, queryString); err != nil {
		panic(err)
	}

	req := httptest.NewRequest("GET", "/api/admin/permission?"+queryString.Encode(), nil)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	decoder := service.Registry.Get("mapstructureDecoder").(func(any, any) error)
	data := &permission.PermissionSearchData{}
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Total)
	suite.NotEmpty(data.LastPage)
	suite.NotEmpty(data.Permissions[0].Id)
	suite.NotEmpty(data.Permissions[0].Name)
	suite.NotEmpty(data.Permissions[0].CreatedAt)
	suite.NotEmpty(data.Permissions[0].UpdatedAt)
}

func (suite *PermissionSearchTestSuite) TestRequestValidationFailed() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	accessToken := test.AdminLogin()

	req := httptest.NewRequest("GET", "/api/admin/permission", nil)
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

func (suite *PermissionSearchTestSuite) TestWrongAccessToken() {
	test.AdminAddPolicies()
	test.AdminAddRole()
	req := httptest.NewRequest("GET", "/api/admin/permission", nil)
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

func (suite *PermissionSearchTestSuite) TestAuthorizationFailed() {
	accessToken := test.AdminLogin()
	req := httptest.NewRequest("GET", "/api/admin/permission", nil)
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

func (suite *PermissionSearchTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestPermissionSearchTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionSearchTestSuite))
}
