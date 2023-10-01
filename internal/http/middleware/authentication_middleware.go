package middleware

import (
	"strconv"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

func Authenticate() gin.HandlerFunc {
	logger := provider.Registry.Logger()
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer") {
			errResp := response.NewErrorResponse(response.Unauthorized, errors.New("jwt authentication failed"), nil)
			logger.Warn(response.Unauthorized, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		accessTokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
		userIdString, err := verifyAccessToken(c, accessTokenString)
		if err != nil {
			errResp := response.NewErrorResponse(response.Unauthorized, err, nil)
			logger.Warn(response.Unauthorized, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		userId, err := strconv.ParseUint(userIdString, 10, 64)
		if err != nil {
			errResp := response.NewErrorResponse(response.Unauthorized, err, nil)
			logger.Warn(response.Unauthorized, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}
		c.Set("user_id", uint(userId))

		c.Next()
	}
}

func verifyAccessToken(c *gin.Context, accessTokenString string) (string, error) {
	accessToken, err := provider.Registry.Authenticator().VerifyJwt(accessTokenString)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if !accessToken.Valid {
		return "", errors.New("access token is invalid")
	}

	claims := accessToken.Claims.(jwt.MapClaims)
	if claims["aud"].([]any)[0].(string) != "access" {
		return "", errors.New("jwt aud claim is not [access] string")
	}

	return claims["sub"].(string), nil
}
