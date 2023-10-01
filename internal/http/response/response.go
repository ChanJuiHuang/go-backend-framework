package response

type Response struct {
	Data any `json:"data" validate:"required"`
}
