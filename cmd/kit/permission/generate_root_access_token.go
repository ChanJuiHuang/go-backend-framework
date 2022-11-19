package permission

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateRootAccessToken() string {
	now := time.Now()
	timeString := strconv.FormatInt(24*365, 10) + "h"
	duration, err := time.ParseDuration(timeString)
	if err != nil {
		panic(err)
	}

	emailUser := &model.EmailUser{}
	db := provider.App.DB.Where("email = ?", "root@root.com").
		First(emailUser)
	if db.Error != nil {
		panic(db.Error)
	}

	claims := jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"access"},
		Subject:   strconv.FormatUint(uint64(emailUser.UserId), 10),
		ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := util.IssueJwtWithConfig(claims)
	fmt.Printf("root jwt is [%s]\n", token)

	return token
}
