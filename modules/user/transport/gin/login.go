package usergin

import (
	"card-warhouse/common"
	"card-warhouse/components/appCtx"
	hasher2 "card-warhouse/components/hasher"
	"card-warhouse/components/tokenProvider/jwt"
	userbiz "card-warhouse/modules/user/biz"
	usermodel "card-warhouse/modules/user/model"
	userstorage "card-warhouse/modules/user/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(appCtx appCtx.AppContext) gin.HandlerFunc {
	return func(context *gin.Context) {
		var data usermodel.UserLogin

		if err := context.ShouldBind(&data); nil != err {
			panic(common.NewFailResponse(err))
		}

		userStore := userstorage.NewMysqlStore(appCtx.GetMainDatabaseConnection())
		hasher := hasher2.NewMd5Hash()
		tokenProvider := jwt.NewJwtProvider(appCtx.GetSecret())

		userBiz := userbiz.NewUserBizLogin(userStore, hasher, appCtx.GetPepper(), tokenProvider, appCtx.GetExpiresIn())

		token, err := userBiz.Login(context.Request.Context(), data)

		if nil != err {
			panic(err)
		}

		context.JSON(http.StatusOK, common.NewSuccessResponse(token))
	}
}
