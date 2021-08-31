package userstorage

import (
	"card-warhouse/common"
	usermodel "card-warhouse/modules/user/model"
	"context"
	"gorm.io/gorm"
)

type TransactionStore interface {
	Create(ctx context.Context, data *usermodel.TransactionCreate) error
	FindWhere(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) ([]usermodel.Transaction, error)
	FindWhereFirst(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*usermodel.Transaction, error)
}

type transactionStore struct {
	db *gorm.DB
}

func NewTransactionStore(db *gorm.DB) *transactionStore {
	return &transactionStore{db: db}
}

func (store *transactionStore) Create(ctx context.Context, data *usermodel.TransactionCreate) error {
	if err := store.db.Create(data).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *transactionStore) FindWhere(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) ([]usermodel.Transaction, error) {
	var transactions []usermodel.Transaction

	if err := store.db.Where(conditions).Find(&transactions).Error; nil != err {
		return nil, common.NewFailResponse(err)
	}

	return transactions, nil
}

func (store *transactionStore) FindWhereFirst(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*usermodel.Transaction, error) {
	var transaction usermodel.Transaction

	if err := store.db.Where(conditions).First(&transaction).Error; nil != err {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.NewFailResponse(err)
	}

	return &transaction, nil
}
