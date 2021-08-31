package usergin

import (
	"card-warhouse/common"
	"card-warhouse/components/appCtx"
	userstorage "card-warhouse/modules/user/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBalance(appCtx appCtx.AppContext) gin.HandlerFunc {
	return func(context *gin.Context) {
		currentUser := context.MustGet(common.CurrentUser).(common.Requester)

		userStore := userstorage.NewMysqlStore(appCtx.GetMainDatabaseConnection())

		user, err := userStore.FindWhereFirst(context.Request.Context(), map[string]interface{}{"id": currentUser.GetId()})

		if nil != err {
			panic(err)
		}

		context.JSON(http.StatusOK, common.NewSuccessResponse(map[string]interface{}{"balance": user.Balance}))
	}
}
