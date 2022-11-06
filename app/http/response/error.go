package response

import "errors"

// http status code 400
var ErrBadRequest error = errors.New("bad request")
var ErrRequestValidationFailed error = errors.New("request validation failed")
var ErrCsrfTokenMismatch error = errors.New("csrf token mismatch")
var ErrLoginEmailIsWrong error = errors.New("login email is wrong")
var ErrLoginPasswordIsWrong error = errors.New("login password is wrong")
var ErrOauthLoginFailed error = errors.New("oauth login failed")

// http status code 401
var ErrJwtAuthenticationFailed error = errors.New("jwt authentication failed")
var ErrAuthorizationFailed error = errors.New("authorization failed")

// http status code 404
var ErrRecordNotFound error = errors.New("record not found")

// http status code 503
var ErrServiceUnavailable error = errors.New("service unavailable")

var ErrorMessageToCode = map[string]string{
	ErrBadRequest.Error():              "400-001",
	ErrRequestValidationFailed.Error(): "400-002",
	ErrCsrfTokenMismatch.Error():       "400-003",
	ErrLoginEmailIsWrong.Error():       "400-004",
	ErrLoginPasswordIsWrong.Error():    "400-005",
	ErrOauthLoginFailed.Error():        "400-006",
	ErrJwtAuthenticationFailed.Error(): "401-001",
	ErrAuthorizationFailed.Error():     "401-002",
	ErrRecordNotFound.Error():          "404-001",
	ErrServiceUnavailable.Error():      "503-001",
}
