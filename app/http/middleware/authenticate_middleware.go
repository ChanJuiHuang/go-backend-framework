package middleware

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var errAccessTokenIsInvalid = errors.New("access token is invalid")
var errJwtAudClaimIsNotAccessString = errors.New("jwt aud claim is not [access] string")

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer") {
			response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
				AbortWithJson(c)
			return
		}

		accessTokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
		userIdString, err := verifyAccessToken(c, accessTokenString)
		if err != nil {
			response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
				MakePreviousMessage(util.WrapError(err)).
				AbortWithJson(c)
			return
		}

		userId, err := strconv.ParseUint(userIdString, 10, 64)
		if err != nil {
			response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
				MakePreviousMessage(util.WrapError(err)).
				AbortWithJson(c)
			return
		}
		c.Set("user_id", uint(userId))

		c.Next()
	}
}

func verifyAccessToken(c *gin.Context, accessTokenString string) (string, error) {
	accessToken, err := util.VerifyJwt(accessTokenString)
	if err != nil {
		return "", util.WrapError(err)
	}
	if !accessToken.Valid {
		return "", util.WrapError(errAccessTokenIsInvalid)
	}
	claims := accessToken.Claims.(jwt.MapClaims)
	if claims["aud"].([]any)[0].(string) != "access" {
		return "", util.WrapError(errJwtAudClaimIsNotAccessString)
	}

	doesExist, err := provider.App.Redis.Exists(context.Background(), accessTokenString).Result()
	if err != nil {
		return "", util.WrapError(err)
	}
	if doesExist == 1 {
		return "", util.WrapError(errAccessTokenIsInvalid)
	}

	return claims["sub"].(string), nil
}
