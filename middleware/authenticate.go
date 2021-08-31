package middleware

import (
	"card-warhouse/common"
	"card-warhouse/components/appCtx"
	"card-warhouse/components/tokenProvider/jwt"
	usermodel "card-warhouse/modules/user/model"
	userstorage "card-warhouse/modules/user/storage"
	"github.com/gin-gonic/gin"
	"strings"
)

func getTokenFromHeaderString(authorization string) (string, error) {
	headers := strings.Split(authorization, " ")

	if headers[0] != "Bearer" || len(headers) < 2 || "" == strings.TrimSpace(headers[1]) {
		return "", common.NewUnauthorizedResponse()
	}

	return headers[1], nil
}

func Authenticate(appCtx appCtx.AppContext, authenticationStore userstorage.UserStore) gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := getTokenFromHeaderString(context.GetHeader("Authorization"))

		if nil != err {
			panic(err)
		}

		tokenProvider := jwt.NewJwtProvider(appCtx.GetSecret())

		payload, errToken := tokenProvider.Validate(token)

		if nil != errToken {
			panic(common.NewUnauthorizedResponse())
		}

		user, errDB := authenticationStore.FindWhereFirst(context.Request.Context(), map[string]interface{}{"id": payload.UserId})

		if nil != errDB {
			panic(common.NewUnauthorizedResponse())
		}

		if usermodel.StatusDeactive == user.Status {
			panic(common.NewUnauthorizedResponse())
		}

		context.Set(common.CurrentUser, user)

		context.Next()
	}
}
