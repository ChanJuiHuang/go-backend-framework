package response

const (
	// 400
	BadRequest              = "Bad Request"
	RequestValidationFailed = "request validation failed"
	EmailIsWrong            = "email is wrong"
	PasswordIsWrong         = "password is wrong"
	PermissionIsRepeat      = "permission is repeat"

	// 401
	Unauthorized = "Unauthorized"

	// 403
	Forbidden = "Forbidden"

	// 500
	InternalServerError = "Internal Server Error"

	// 503
	ServiceUnavailable = "Service Unavailable"
)

var MessageToCode = map[string]string{
	BadRequest:              "400-001",
	RequestValidationFailed: "400-002",
	EmailIsWrong:            "400-003",
	PasswordIsWrong:         "400-004",
	PermissionIsRepeat:      "400-005",
	Unauthorized:            "401-001",
	Forbidden:               "403-001",
	InternalServerError:     "500-001",
	ServiceUnavailable:      "503-001",
}
