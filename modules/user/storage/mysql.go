package userstorage

import (
	"card-warhouse/common"
	usermodel "card-warhouse/modules/user/model"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserStore interface {
	Create(ctx context.Context, data *usermodel.UserCreate) error
	Update(ctx context.Context, conditions map[string]interface{}, data *usermodel.UserUpdate) error
	FindWhere(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) ([]usermodel.User, error)
	FindWhereFirst(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*usermodel.User, error)
	Increment(ctx context.Context, conditions map[string]interface{}, amount int) error
	Decrement(ctx context.Context, conditions map[string]interface{}, amount int) error
	StartTransaction() *gorm.DB
}

type mysqlStore struct {
	db *gorm.DB
}

func NewMysqlStore(db *gorm.DB) *mysqlStore {
	return &mysqlStore{db: db}
}

func (store *mysqlStore) Create(ctx context.Context, data *usermodel.UserCreate) error {
	if err := store.db.Create(data).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *mysqlStore) Update(ctx context.Context, conditions map[string]interface{}, data *usermodel.UserUpdate) error {
	if err := store.db.Where(conditions).Updates(data).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *mysqlStore) FindWhere(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) ([]usermodel.User, error) {
	var users []usermodel.User

	if err := store.db.Where(conditions).Find(&users).Error; nil != err {
		return nil, common.NewFailResponse(err)
	}

	return users, nil
}

func (store *mysqlStore) FindWhereFirst(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*usermodel.User, error) {
	var user usermodel.User

	if err := store.db.Where(conditions).First(&user).Error; nil != err {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.NewFailResponse(err)
	}

	return &user, nil
}

func (store *mysqlStore) Increment(ctx context.Context, conditions map[string]interface{}, amount int) error {
	if err := store.db.Table(getTableName()).Where(conditions).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *mysqlStore) Decrement(ctx context.Context, conditions map[string]interface{}, amount int) error {
	if err := store.db.Table(getTableName()).Where(conditions).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error; nil != err {
		return common.NewFailResponse(err)
	}

	return nil
}

func (store *mysqlStore) StartTransaction() *gorm.DB {
	return store.db.Begin()
}

func getTableName() string {
	return usermodel.User{}.TableName()
}
