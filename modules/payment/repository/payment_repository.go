package paymentrepository

import (
	"card-warhouse/common"
	"card-warhouse/components/acquirer"
	paymentmodel "card-warhouse/modules/payment/model"
	paymentstorage "card-warhouse/modules/payment/storage"
	"context"
)

type PaymentRepo interface {
	Create(ctx context.Context, data *paymentmodel.PaymentCreate) (*paymentmodel.PaymentResult, error)
	Update(ctx context.Context, conditions map[string]interface{}, data *paymentmodel.PaymentUpdate) error
	Inquire(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*paymentmodel.Payment, error)
	IsUniqueInvoiceNo(ctx context.Context, conditions map[string]interface{}) bool
}

type paymentRepository struct {
	paymentStore     paymentstorage.MysqlStoreInterface
	transactionStore paymentstorage.MysqlTransactionStore
}

func NewPaymentRepository(paymentStore paymentstorage.MysqlStoreInterface, transactionStore paymentstorage.MysqlTransactionStore) *paymentRepository {
	return &paymentRepository{paymentStore: paymentStore, transactionStore: transactionStore}
}

func (repository *paymentRepository) Create(ctx context.Context, data *paymentmodel.PaymentCreate) (*paymentmodel.PaymentResult, error) {
	acquirerFactory := acquirer.NewAcquirerFactory()

	acquirerClass := acquirerFactory.GetAcquirerByMethod(data.Method)

	if nil == acquirerClass {
		return nil, common.NewBadRequestResponse(nil, common.CodeMethodNotSupport, common.GetMessageFromCode(common.CodeMethodNotSupport))
	}

	transactionDB := repository.paymentStore.StartTransaction()

	if err := repository.paymentStore.Create(ctx, data); nil != err {
		transactionDB.Rollback()
		return nil, err
	}

	var transaction = paymentmodel.TransactionCreate{
		PaymentNo: data.PaymentNo,
		Acquirer:  acquirerClass.Name(),
		Status:    paymentmodel.TranStatusCreated,
		Type:      paymentmodel.TypePurchase,
	}

	if err := repository.transactionStore.Create(ctx, &transaction); nil != err {
		transactionDB.Rollback()
		return nil, err
	}

	if err := transactionDB.Commit().Error; nil != err {
		transactionDB.Rollback()
		return nil, common.NewFailResponse(err)
	}

	result, errAcquirer := acquirerClass.Purchase(map[string]interface{}{
		"payment_no":  data.PaymentNo,
		"description": data.Description,
		"amount":      data.Amount,
		"card_pin":    transaction.CardPin,
		"card_serial": transaction.CardSerial,
	})

	if nil != errAcquirer {
		return nil, errAcquirer
	}

	status := transFormTransactionStatus(*result)

	if (paymentmodel.TranStatusProcessing == status || paymentmodel.TranStatusCompleted == status) && result.Amount != data.Amount {
		status = paymentmodel.TranStatusReview
	}

	transactionUpdate := paymentmodel.TransactionUpdate{
		Status:      &status,
		AcquirerUrl: &result.AcquirerUrl,
		AcquirerTxn: &result.AcquirerTxnId,
	}

	transactionDB2 := repository.paymentStore.StartTransaction()

	if err := repository.transactionStore.Update(ctx, map[string]interface{}{"id": transaction.Id}, &transactionUpdate); nil != err {
		transactionDB2.Rollback()
		return nil, err
	}

	paymentStatus := transFormPaymentStatus(transactionUpdate)

	paymentUpdate := paymentmodel.PaymentUpdate{
		Status:        &paymentStatus,
		FailureReason: &result.FailReason,
	}

	if err := repository.paymentStore.Update(ctx, map[string]interface{}{"id": data.Id}, &paymentUpdate); nil != err {
		transactionDB2.Rollback()
		return nil, err
	}

	if err := transactionDB2.Commit().Error; nil != err {
		transactionDB2.Rollback()
		return nil, common.NewFailResponse(err)
	}

	return &paymentmodel.PaymentResult{
		PaymentNo:  data.PaymentNo,
		Status:     *paymentUpdate.Status,
		PaymentUrl: *transactionUpdate.AcquirerUrl,
		InvoiceNo:  data.InvoiceNo,
		CreatedAt:  data.CreatedAt,
		Amount:     data.Amount,
	}, nil
}

func (repository *paymentRepository) Inquire(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*paymentmodel.Payment, error) {
	data, err := repository.paymentStore.FindWhereFirst(ctx, conditions, moreKeys...)

	if nil != err {
		if err == common.ErrRecordNotFound {
			return nil, common.NewNotFoundResponse()
		}

		return nil, err
	}

	return data, err
}

func (repository *paymentRepository) IsUniqueInvoiceNo(ctx context.Context, conditions map[string]interface{}) bool {
	_, err := repository.paymentStore.FindWhereFirst(ctx, conditions)

	if nil != err {
		if err == common.ErrRecordNotFound {
			return true
		}

		return false
	}

	return false
}

func (repository *paymentRepository) Update(ctx context.Context, conditions map[string]interface{}, data *paymentmodel.PaymentUpdate) error {
	if err := repository.paymentStore.Update(ctx, conditions, data); nil != err {
		return err
	}

	return nil
}

func transFormTransactionStatus(result acquirer.Result) int {
	switch result.Status {
	case acquirer.StatusProcessing:
		return paymentmodel.TranStatusProcessing
	case acquirer.StatusCompleted:
		return paymentmodel.TranStatusCompleted
	case acquirer.StatusFail:
		return paymentmodel.TranStatusFail
	case acquirer.StatusCancel:
		return paymentmodel.TranStatusCancel
	case acquirer.StatusRefund:
		return paymentmodel.TranStatusRefund
	case acquirer.StatusReview:
		return paymentmodel.TranStatusReview
	}

	return paymentmodel.TranStatusFail
}

func transFormPaymentStatus(transactionUpdate paymentmodel.TransactionUpdate) int {
	switch *transactionUpdate.Status {
	case paymentmodel.TranStatusProcessing:
		return paymentmodel.PayStatusProcessing
	case paymentmodel.TranStatusCompleted:
		return paymentmodel.PayStatusNeedCompensation
	case paymentmodel.TranStatusFail:
		return paymentmodel.PayStatusFail
	case paymentmodel.TranStatusCancel:
		return paymentmodel.PayStatusCancel
	case paymentmodel.TranStatusRefund:
		return paymentmodel.PayStatusRefund
	case paymentmodel.TranStatusReview:
		return paymentmodel.PayStatusReview
	}

	return paymentmodel.PayStatusFail
}
