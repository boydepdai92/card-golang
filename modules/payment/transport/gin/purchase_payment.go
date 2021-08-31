package paymentgin

import (
	"card-warhouse/common"
	"card-warhouse/components/appCtx"
	paymentbiz "card-warhouse/modules/payment/biz"
	paymentmodel "card-warhouse/modules/payment/model"
	paymentrepository "card-warhouse/modules/payment/repository"
	paymentstorage "card-warhouse/modules/payment/storage"
	userbiz "card-warhouse/modules/user/biz"
	userstorage "card-warhouse/modules/user/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PurchasePayment(appCtx appCtx.AppContext) gin.HandlerFunc {
	return func(context *gin.Context) {
		var data paymentmodel.PaymentCreate

		if err := context.ShouldBind(&data); nil != err {
			panic(common.NewFailResponse(err))
		}

		currentUser := context.MustGet(common.CurrentUser).(common.Requester)

		data.UserId = currentUser.GetId()
		data.ClientIp = context.ClientIP()

		paymentStore := paymentstorage.NewMysqlPaymentStore(appCtx.GetMainDatabaseConnection())
		transactionStore := paymentstorage.NewMysqlTransactionStore(appCtx.GetMainDatabaseConnection())

		payRepository := paymentrepository.NewPaymentRepository(paymentStore, transactionStore)

		userStore := userstorage.NewMysqlStore(appCtx.GetMainDatabaseConnection())
		transactionBalanceStore := userstorage.NewTransactionStore(appCtx.GetMainDatabaseConnection())

		plusBalanceBiz := userbiz.NewPlusBalanceBiz(userStore, transactionBalanceStore)

		paymentBiz := paymentbiz.NewPurchasePaymentBiz(payRepository, plusBalanceBiz)

		result, err := paymentBiz.Purchase(context.Request.Context(), &data)

		if nil != err {
			panic(err)
		}

		context.JSON(http.StatusOK, common.NewSuccessResponse(result))
	}
}
