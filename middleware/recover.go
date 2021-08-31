package middleware

import (
	"card-warhouse/common"
	"card-warhouse/components/appCtx"
	"github.com/gin-gonic/gin"
)

func Recover(appCtx appCtx.AppContext) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); nil != err {
				context.Header("Content-type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					context.AbortWithStatusJSON(appErr.StatusCode, appErr)
					panic(err)
					return
				}

				appErr := common.NewFailResponse(err.(error))
				context.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err)
				return
			}
		}()

		context.Next()
	}
}
