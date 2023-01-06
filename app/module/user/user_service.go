package user

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
)

var ErrRefreshTokenIsInvalid = errors.New("refresh token is invalid")
var ErrJwtAudClaimIsNotRefreshString = errors.New("jwt aud claim is not [refresh] string")
var ErrOverRefreshTokenLimit = errors.New("over refresh token limit")

func VerifyRefreshToken(refreshTokenString string) (uint, int64, error) {
	refreshToken, err := util.VerifyJwt(refreshTokenString)
	if err != nil {
		return 0, 0, util.WrapError(err)
	}
	if !refreshToken.Valid {
		return 0, 0, util.WrapError(ErrRefreshTokenIsInvalid)
	}

	claims := refreshToken.Claims.(jwt.MapClaims)
	if claims["aud"].([]any)[0].(string) != "refresh" {
		return 0, 0, util.WrapError(ErrJwtAudClaimIsNotRefreshString)
	}

	userId, err := strconv.ParseUint(claims["sub"].(string), 10, 64)
	if err != nil {
		return 0, 0, util.WrapError(err)
	}

	if err := doesOverRefreshTokenLimit(uint(userId)); err != nil {
		return 0, 0, util.WrapError(err)
	}

	return uint(userId), int64(claims["exp"].(float64)), nil
}

func doesOverRefreshTokenLimit(userId uint) error {
	var count int64
	db := provider.App.DB.Model(&model.RefreshTokenRecord{}).
		Where("user_id = ?", userId).
		Count(&count)
	if util.IsDatabaseError(db) {
		return util.WrapError(db.Error)
	}
	if count >= int64(util.GetRefreshTokenLimit()) {
		return util.WrapError(ErrOverRefreshTokenLimit)
	}

	return nil
}

func RecordAccessToken(accessTokenString string, userId uint) error {
	accessToken, _, err := jwt.NewParser().ParseUnverified(accessTokenString, jwt.MapClaims{})
	if err != nil {
		return util.WrapError(err)
	}

	err = provider.App.Redis.SetArgs(
		context.Background(),
		accessTokenString,
		userId,
		redis.SetArgs{
			Mode:     "NX",
			ExpireAt: time.Unix(int64(accessToken.Claims.(jwt.MapClaims)["exp"].(float64)), 0),
		},
	).Err()
	if err != nil {
		provider.App.Logger.Error(err.Error())
	}

	return nil
}
