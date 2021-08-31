package userbiz

import (
	"card-warhouse/common"
	"card-warhouse/components/hasher"
	usermodel "card-warhouse/modules/user/model"
	userstorage "card-warhouse/modules/user/storage"
	"context"
)

type userRegisterBiz struct {
	store  userstorage.UserStore
	hasher hasher.Hasher
	pepper string
}

func NewUserRegisterBiz(store userstorage.UserStore, hasher hasher.Hasher, pepper string) *userRegisterBiz {
	return &userRegisterBiz{store: store, hasher: hasher, pepper: pepper}
}

func (userBiz *userRegisterBiz) Register(ctx context.Context, data *usermodel.UserCreate) error {
	if err := data.Validate(); nil != err {
		return common.NewBadRequestResponse(err, common.CodeFail, err.Error())
	}

	user, _ := userBiz.store.FindWhereFirst(ctx, map[string]interface{}{"username": data.Username})

	if nil != user {
		return common.NewBadRequestResponse(usermodel.ErrUsernameExisted, common.CodeUserExisted, common.GetMessageFromCode(common.CodeUserExisted))
	}

	data.Salt = common.GenerateSalt(50)
	data.SecretKey = common.GenerateSecretKey(30)
	data.Password = userBiz.hasher.Hash(data.GenerateRawPassword(userBiz.pepper))

	if err := userBiz.store.Create(ctx, data); nil != err {
		return err
	}

	return nil
}
