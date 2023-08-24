package response

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/stacktrace"
	"go.uber.org/zap"
)

type Debug struct {
	Error      string `json:"error" example:"error message" validate:"required"`
	err        error
	Stacktrace []string `json:"stacktrace" validate:"required"`
}

type ErrorResponse struct {
	Message string `json:"message" validate:"required"`
	Code    string `json:"code" validate:"required"`
	debug   *Debug
	Debug   *Debug         `json:"debug,omitempty"`
	Context map[string]any `json:"context,omitempty"`
}

func NewErrorResponse(message string, err error, context map[string]any) *ErrorResponse {
	debug := &Debug{
		err:        err,
		Stacktrace: stacktrace.GetStackStrace(err),
	}
	if err != nil {
		debug.Error = err.Error()
	}

	globalConfig := config.Registry.Get("global").(global.Config)
	if globalConfig.Debug {
		return &ErrorResponse{
			Message: message,
			Code:    MessageToCode[message],
			Debug:   debug,
			Context: context,
		}
	}

	return &ErrorResponse{
		Message: message,
		Code:    MessageToCode[message],
		debug:   debug,
		Context: context,
	}
}

func (er *ErrorResponse) StatusCode() int {
	code, err := strconv.ParseInt(
		strings.Split(er.Code, "-")[0],
		10,
		0,
	)
	if err != nil {
		provider.Registry.Logger().Error(err.Error())
		code = http.StatusBadRequest
	}

	return int(code)
}

func (er *ErrorResponse) MakeLogFields(req *http.Request) []zap.Field {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		provider.Registry.Logger().Error(err.Error())
		requestBody = nil
	}
	if requestBody != nil && json.Valid(requestBody) {
		buffer := bytes.NewBuffer(make([]byte, 0, len(requestBody)))
		err = json.Compact(buffer, requestBody)
		if err != nil {
			provider.Registry.Logger().Error(err.Error())
			requestBody = nil
		} else {
			requestBody = buffer.Bytes()
		}
	}

	var debug *Debug
	globalConfig := config.Registry.Get("global").(global.Config)
	if globalConfig.Debug {
		debug = er.Debug
	} else {
		debug = er.debug
	}

	return []zap.Field{
		zap.String("code", er.Code),
		zap.String("error", debug.err.Error()),
		zap.Strings("stacktrace", debug.Stacktrace),
		zap.Int("status_code", er.StatusCode()),
		zap.String("path", req.URL.Path),
		zap.String("query_string", req.URL.Query().Encode()),
		zap.ByteString("request_body", requestBody),
	}
}
