package usergin

import (
	"card-warhouse/common"
	"card-warhouse/components/appCtx"
	hasher2 "card-warhouse/components/hasher"
	userbiz "card-warhouse/modules/user/biz"
	usermodel "card-warhouse/modules/user/model"
	userstorage "card-warhouse/modules/user/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(appCtx appCtx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user usermodel.UserCreate

		if err := c.ShouldBind(&user); nil != err {
			panic(common.NewFailResponse(err))
		}

		userStore := userstorage.NewMysqlStore(appCtx.GetMainDatabaseConnection())
		hasher := hasher2.NewMd5Hash()

		userRegisterBiz := userbiz.NewUserRegisterBiz(userStore, hasher, appCtx.GetPepper())

		if err := userRegisterBiz.Register(c.Request.Context(), &user); nil != err {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(map[string]interface{}{"id": user.Id}))
	}
}
