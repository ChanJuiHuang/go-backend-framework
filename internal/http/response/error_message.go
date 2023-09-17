package response

const (
	// 400
	BadRequest = "Bad Request"

	// 401
	Unauthorized = "Unauthorized"

	// 403
	Forbidden = "Forbidden"

	// 503
	ServiceUnavailable = "Service Unavailable"
)

var MessageToCode = map[string]string{
	BadRequest:         "400-001",
	Unauthorized:       "401-001",
	Forbidden:          "403-001",
	ServiceUnavailable: "503-001",
}
