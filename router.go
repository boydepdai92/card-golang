package main

import (
	"card-warhouse/components/appCtx"
	"card-warhouse/middleware"
	paymentgin "card-warhouse/modules/payment/transport/gin"
	userstorage "card-warhouse/modules/user/storage"
	usergin "card-warhouse/modules/user/transport/gin"
	"github.com/gin-gonic/gin"
)

func Register(appCtx appCtx.AppContext) {
	userStore := userstorage.NewMysqlStore(appCtx.GetMainDatabaseConnection())

	router := gin.Default()
	router.Use(middleware.Recover(appCtx))

	version1 := router.Group("v1")
	{
		users := version1.Group("users")
		{
			users.POST("register", usergin.Register(appCtx))
			users.POST("login", usergin.Login(appCtx))
			users.GET("get-balance", middleware.Authenticate(appCtx, userStore), usergin.GetBalance(appCtx))
		}

		payments := version1.Group("payments", middleware.Authenticate(appCtx, userStore))
		{
			payments.POST("create", paymentgin.PurchasePayment(appCtx))
			payments.GET(":reference", paymentgin.InquirePayment(appCtx))
		}
	}

	router.Run().Error()
}
