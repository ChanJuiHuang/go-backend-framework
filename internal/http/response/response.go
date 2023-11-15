package response

type Response struct {
	Data any `json:"data" validate:"required"`
}

func NewResponse(data any) *Response {
	return &Response{
		Data: data,
	}
}
