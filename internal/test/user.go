package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
)

func UserRegister() {
	userRegisterRequest := user.UserRegisterRequest{
		Name:     "john",
		Email:    "john@test.com",
		Password: "abcABC123",
	}
	reqBody, err := json.Marshal(userRegisterRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(reqBody))
	AddCsrfToken(req)
	resp := httptest.NewRecorder()
	HttpHandler.ServeHTTP(resp, req)

	respBody := &user.UserRegisterResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}
}

func UserLogin() (string, string) {
	userLoginRequest := user.UserLoginRequest{
		Email:    "john@test.com",
		Password: "abcABC123",
	}
	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBody))
	AddCsrfToken(req)
	resp := httptest.NewRecorder()
	HttpHandler.ServeHTTP(resp, req)

	respBody := &user.UserLoginResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	return respBody.AccessToken, respBody.RefreshToken
}
