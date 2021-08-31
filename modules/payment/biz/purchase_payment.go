package paymentbiz

import (
	"card-warhouse/common"
	paymentmodel "card-warhouse/modules/payment/model"
	paymentrepository "card-warhouse/modules/payment/repository"
	userbiz "card-warhouse/modules/user/biz"
	"context"
	"time"
)

type purchasePaymentBiz struct {
	paymentRepo    paymentrepository.PaymentRepo
	plusBalanceBiz userbiz.PlusBalance
}

func NewPurchasePaymentBiz(paymentRepo paymentrepository.PaymentRepo, plusBalanceBiz userbiz.PlusBalance) *purchasePaymentBiz {
	return &purchasePaymentBiz{paymentRepo: paymentRepo, plusBalanceBiz: plusBalanceBiz}
}

func (purchasePayBiz *purchasePaymentBiz) Purchase(ctx context.Context, data *paymentmodel.PaymentCreate) (*paymentmodel.PaymentResult, error) {
	errValidation := data.Validate()

	if nil != errValidation {
		return nil, common.NewBadRequestResponse(errValidation, common.CodeFail, errValidation.Error())
	}

	isValidInvoiceNo := purchasePayBiz.paymentRepo.IsUniqueInvoiceNo(ctx, map[string]interface{}{"invoice_no": data.InvoiceNo, "user_id": data.UserId})

	if false == isValidInvoiceNo {
		return nil, common.NewBadRequestResponse(paymentmodel.ErrInvoiceNoInvalid, common.CodeFail, paymentmodel.ErrInvoiceNoInvalid.Error())
	}

	data.PaymentNo = common.GenerateRandomNumberString()
	data.Status = paymentmodel.PayStatusCreated
	data.ExpiredAt = time.Now().Local().Add(time.Hour * time.Duration(paymentmodel.ExpiredHour))

	result, err := purchasePayBiz.paymentRepo.Create(ctx, data)

	if nil != err {
		return nil, err
	}

	if paymentmodel.PayStatusNeedCompensation == result.Status {
		var status int

		if true == purchasePayBiz.callPlusBalance(ctx, data.UserId, result.Amount, result.PaymentNo) {
			status = paymentmodel.PayStatusCompleted
		} else {
			status = paymentmodel.PayStatusNotDelivery
		}

		paymentUpdateData := paymentmodel.PaymentUpdate{Status: &status}

		if err := purchasePayBiz.paymentRepo.Update(ctx, map[string]interface{}{"payment_no": result.PaymentNo}, &paymentUpdateData); nil != err {
			return nil, err
		}

		result.Status = *paymentUpdateData.Status

		return result, nil
	}

	return result, nil
}

func (purchasePayBiz *purchasePaymentBiz) callPlusBalance(ctx context.Context, userId int, amount int, reference string) bool {
	if err := purchasePayBiz.plusBalanceBiz.Plus(ctx, userId, amount, reference); nil != err {
		return false
	}

	return true
}
