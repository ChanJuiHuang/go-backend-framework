package response

const (
	// 400
	BadRequest              = "Bad Request"
	RequestValidationFailed = "request validation failed"

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
	Unauthorized:            "401-001",
	Forbidden:               "403-001",
	InternalServerError:     "500-001",
	ServiceUnavailable:      "503-001",
}
