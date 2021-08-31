package userbiz

import (
	"card-warhouse/common"
	usermodel "card-warhouse/modules/user/model"
	userstorage "card-warhouse/modules/user/storage"
	"context"
)

type PlusBalance interface {
	Plus(ctx context.Context, userId int, amount int, reference string) error
}

type plusBalanceBiz struct {
	store            userstorage.UserStore
	transactionStore userstorage.TransactionStore
}

func NewPlusBalanceBiz(store userstorage.UserStore, transactionStore userstorage.TransactionStore) *plusBalanceBiz {
	return &plusBalanceBiz{store: store, transactionStore: transactionStore}
}

func (plusBalance *plusBalanceBiz) Plus(ctx context.Context, userId int, amount int, reference string) error {
	transactionDB := plusBalance.store.StartTransaction()

	user, err := plusBalance.store.FindWhereFirst(ctx, map[string]interface{}{"id": userId})

	if nil != err {
		transactionDB.Rollback()
		return err
	}

	if usermodel.StatusDeactive == user.Status {
		transactionDB.Rollback()
		return common.NewBadRequestResponse(usermodel.ErrUserDeactive, common.CodeFail, usermodel.ErrUserDeactive.Error())
	}

	if err := plusBalance.store.Increment(ctx, map[string]interface{}{"id": userId}, amount); nil != err {
		transactionDB.Rollback()
		return err
	}

	transaction := usermodel.TransactionCreate{
		UserId:       userId,
		Reference:    reference,
		Type:         usermodel.TypePurchase,
		Action:       usermodel.ActionPlus,
		Amount:       amount,
		BeforeAmount: user.Balance,
		AfterAmount:  user.Balance + amount,
	}

	if err := plusBalance.transactionStore.Create(ctx, &transaction); nil != err {
		transactionDB.Rollback()
		return err
	}

	if err := transactionDB.Commit().Error; nil != err {
		transactionDB.Rollback()
		return common.NewFailResponse(err)
	}

	return nil
}
