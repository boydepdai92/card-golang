package paymentstorage

import (
	"card-warhouse/common"
	paymentmodel "card-warhouse/modules/payment/model"
	"context"
	"gorm.io/gorm"
)

type MysqlStoreInterface interface {
	Create(ctx context.Context, data *paymentmodel.PaymentCreate) error
	Update(ctx context.Context, conditions map[string]interface{}, data *paymentmodel.PaymentUpdate) error
	FindWhere(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) ([]paymentmodel.Payment, error)
	FindWhereFirst(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*paymentmodel.Payment, error)
	StartTransaction() *gorm.DB
}

type mysqlPaymentStore struct {
	db *gorm.DB
}

func NewMysqlPaymentStore(db *gorm.DB) *mysqlPaymentStore {
	return &mysqlPaymentStore{db: db}
}

func (store *mysqlPaymentStore) Create(ctx context.Context, data *paymentmodel.PaymentCreate) error {
	if err := store.db.Create(data).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *mysqlPaymentStore) Update(ctx context.Context, conditions map[string]interface{}, data *paymentmodel.PaymentUpdate) error {
	if err := store.db.Where(conditions).Updates(data).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *mysqlPaymentStore) FindWhere(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) ([]paymentmodel.Payment, error) {
	var data []paymentmodel.Payment

	if err := store.db.Where(conditions).Find(&data).Error; nil != err {
		return nil, common.NewFailResponse(err)
	}

	return data, nil
}

func (store *mysqlPaymentStore) FindWhereFirst(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*paymentmodel.Payment, error) {
	var data paymentmodel.Payment

	for i := range moreKeys {
		store.db = store.db.Preload(moreKeys[i])
	}

	if err := store.db.Where(conditions).First(&data).Error; nil != err {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.NewFailResponse(err)
	}

	return &data, nil
}

func (store *mysqlPaymentStore) StartTransaction() *gorm.DB {
	return store.db.Begin()
}
