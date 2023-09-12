package response_test

import (
	"testing"

	_ "github.com/ChanJuiHuang/go-backend-framework/internal/test"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/stretchr/testify/assert"
)

func TestMessageToCode(t *testing.T) {
	assert.Equal(t, "400-001", response.MessageToCode[response.BadRequest])
	assert.Equal(t, "400-002", response.MessageToCode[response.RequestValidationFailed])
	assert.Equal(t, "401-001", response.MessageToCode[response.Unauthorized])
	assert.Equal(t, "403-001", response.MessageToCode[response.Forbidden])
	assert.Equal(t, "500-001", response.MessageToCode[response.InternalServerError])
	assert.Equal(t, "503-001", response.MessageToCode[response.ServiceUnavailable])
}
