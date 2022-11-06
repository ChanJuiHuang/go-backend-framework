package response

import (
	"strconv"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HttpErrorResponse struct {
	Message         string         `json:"message"`
	Code            string         `json:"code"`
	PreviousMessage string         `json:"previous_message,omitempty"`
	Stacktrace      []string       `json:"stacktrace"`
	Context         map[string]any `json:"context,omitempty"`
}

type SwaggerErrorResponse struct {
	Message         string `validate:"required"`
	Code            string `validate:"required" enums:"400-001,400-002,400-003,400-004,400-005,400-006,401-001,404-001,503-001"`
	PreviousMessage string
	Stacktrace      []string `validate:"required"`
	Context         map[string]any
}

func NewHttpErrorResponse(err error) *HttpErrorResponse {
	message := err.Error()
	code, ok := ErrorMessageToCode[message]
	if !ok {
		panic("response error message does not exist")
	}

	lastStack := util.Stack(2)
	stacktrace := make([]string, 0, 5)
	stacktrace = append(stacktrace, lastStack)
	r := &HttpErrorResponse{
		Message:    message,
		Code:       code,
		Stacktrace: stacktrace,
	}

	return r
}

func (r *HttpErrorResponse) MakePreviousMessage(previousError error) *HttpErrorResponse {
	messages := strings.Split(previousError.Error(), util.ErrorDelimiter)
	stackLength := len(messages) - 1
	r.Stacktrace = append(r.Stacktrace, messages[:stackLength]...)
	r.PreviousMessage = messages[stackLength]

	return r
}

func (r *HttpErrorResponse) MakeContext(context map[string]any) *HttpErrorResponse {
	r.Context = context

	return r
}

func (r *HttpErrorResponse) StatusCode() int {
	code, _ := strconv.ParseInt(
		strings.Split(r.Code, "-")[0],
		10,
		0,
	)
	return int(code)
}

func (r *HttpErrorResponse) ToJson(c *gin.Context) {
	r.log()
	r.transform()
	c.JSON(r.StatusCode(), r)
}

func (r *HttpErrorResponse) AbortWithJson(c *gin.Context) {
	r.log()
	r.transform()
	c.AbortWithStatusJSON(r.StatusCode(), r)
}

func (r *HttpErrorResponse) AbortWithStatus(c *gin.Context) {
	r.log()
	c.AbortWithStatus(r.StatusCode())
}

func (r *HttpErrorResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("code", r.Code)

	if r.PreviousMessage != "" {
		enc.AddReflected("previousMessage", r.PreviousMessage)
	}

	if r.Context != nil {
		enc.AddReflected("context", r.Context)
	}

	enc.AddReflected("stacktrace", r.Stacktrace)

	return nil
}

func (r *HttpErrorResponse) transform() {
	if !config.App().Debug {
		r.PreviousMessage = ""
		r.Stacktrace = nil
	}
}

func (r *HttpErrorResponse) log() {
	status := r.StatusCode()
	logger := provider.App.Logger

	switch {
	case status >= 400 && status < 500:
		logger.Warn(r.Message, zap.Object("errorResponse", r))
	case status >= 500:
		logger.Error(r.Message, zap.Object("errorResponse", r))
	}
}
