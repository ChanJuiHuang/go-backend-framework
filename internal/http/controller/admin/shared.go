package admin

type Rule struct {
	Object string `json:"object" binding:"required,startswith=/"`
	Action string `json:"action" binding:"required,oneof=GET POST PUT PATCH DELETE"`
}
