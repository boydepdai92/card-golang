package paymentgin

import (
	"card-warhouse/common"
	"card-warhouse/components/appCtx"
	paymentbiz "card-warhouse/modules/payment/biz"
	paymentrepository "card-warhouse/modules/payment/repository"
	paymentstorage "card-warhouse/modules/payment/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InquirePayment(appCtx appCtx.AppContext) gin.HandlerFunc {
	return func(context *gin.Context) {
		reference := context.Param("reference")

		currentUser := context.MustGet(common.CurrentUser).(common.Requester)

		paymentStore := paymentstorage.NewMysqlPaymentStore(appCtx.GetMainDatabaseConnection())
		transactionStore := paymentstorage.NewMysqlTransactionStore(appCtx.GetMainDatabaseConnection())

		payRepository := paymentrepository.NewPaymentRepository(paymentStore, transactionStore)

		inquireBiz := paymentbiz.NewInquirePaymentBiz(payRepository)

		data, err := inquireBiz.Inquire(context.Request.Context(), reference, currentUser.GetId())

		if nil != err {
			panic(err)
		}

		context.JSON(http.StatusOK, common.NewSuccessResponse(data))
	}
}
