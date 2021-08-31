package paymentbiz

import (
	paymentmodel "card-warhouse/modules/payment/model"
	paymentrepository "card-warhouse/modules/payment/repository"
	"context"
)

type inquirePaymentBiz struct {
	paymentRepo paymentrepository.PaymentRepo
}

func NewInquirePaymentBiz(paymentRepo paymentrepository.PaymentRepo) *inquirePaymentBiz {
	return &inquirePaymentBiz{paymentRepo: paymentRepo}
}

func (inquirePaymentBiz *inquirePaymentBiz) Inquire(ctx context.Context, reference string, userId int) (*paymentmodel.Payment, error) {
	data, err := inquirePaymentBiz.paymentRepo.Inquire(ctx, map[string]interface{}{"invoice_no": reference, "user_id": userId}, "Transaction")

	if nil != err {
		return nil, err
	}

	return data, err
}
