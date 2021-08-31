package paymentstorage

import (
	"card-warhouse/common"
	paymentmodel "card-warhouse/modules/payment/model"
	"context"
	"gorm.io/gorm"
)

type MysqlTransactionStore interface {
	Create(ctx context.Context, data *paymentmodel.TransactionCreate) error
	Update(ctx context.Context, conditions map[string]interface{}, data *paymentmodel.TransactionUpdate) error
	FindWhere(ctx context.Context, conditions map[string]interface{}) ([]paymentmodel.Transaction, error)
	FindWhereFirst(ctx context.Context, conditions map[string]interface{}) (*paymentmodel.Transaction, error)
}

type mysqlTransactionStore struct {
	db *gorm.DB
}

func NewMysqlTransactionStore(db *gorm.DB) *mysqlTransactionStore {
	return &mysqlTransactionStore{db: db}
}

func (store *mysqlTransactionStore) Create(ctx context.Context, data *paymentmodel.TransactionCreate) error {
	if err := store.db.Create(data).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *mysqlTransactionStore) Update(ctx context.Context, conditions map[string]interface{}, data *paymentmodel.TransactionUpdate) error {
	if err := store.db.Where(conditions).Updates(data).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *mysqlTransactionStore) FindWhere(ctx context.Context, conditions map[string]interface{}) ([]paymentmodel.Transaction, error) {
	var data []paymentmodel.Transaction

	if err := store.db.Where(conditions).Find(&data).Error; nil != err {
		return nil, common.NewFailResponse(err)
	}

	return data, nil
}

func (store *mysqlTransactionStore) FindWhereFirst(ctx context.Context, conditions map[string]interface{}) (*paymentmodel.Transaction, error) {
	var data paymentmodel.Transaction

	if err := store.db.Where(conditions).First(&data).Error; nil != err {
		return nil, common.NewFailResponse(err)
	}

	return &data, nil
}
