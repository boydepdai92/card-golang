package userbiz

import (
	"card-warhouse/common"
	"card-warhouse/components/hasher"
	"card-warhouse/components/tokenProvider"
	usermodel "card-warhouse/modules/user/model"
	userstorage "card-warhouse/modules/user/storage"
	"context"
)

type userBizLogin struct {
	store         userstorage.UserStore
	hasher        hasher.Hasher
	pepper        string
	tokenProvider tokenProvider.Provider
	expiresIn     int
}

func NewUserBizLogin(store userstorage.UserStore, hasher hasher.Hasher, pepper string, tokenProvider tokenProvider.Provider, expiresIn int) *userBizLogin {
	return &userBizLogin{store: store, hasher: hasher, pepper: pepper, tokenProvider: tokenProvider, expiresIn: expiresIn}
}

func (userBiz *userBizLogin) Login(ctx context.Context, data usermodel.UserLogin) (*tokenProvider.Token, error) {
	if err := data.Validate(); nil != err {
		return nil, common.NewBadRequestResponse(err, common.CodeFail, err.Error())
	}

	user, _ := userBiz.store.FindWhereFirst(ctx, map[string]interface{}{"username": data.Username})

	if nil == user {
		return nil, common.NewBadRequestResponse(usermodel.ErrUsernameOrPasswordWrong, common.CodeFail, usermodel.ErrUsernameOrPasswordWrong.Error())
	}

	password := userBiz.hasher.Hash(data.GenerateRawPassword(user.Salt, userBiz.pepper))

	if user.Password != password {
		return nil, common.NewBadRequestResponse(usermodel.ErrUsernameOrPasswordWrong, common.CodeFail, usermodel.ErrUsernameOrPasswordWrong.Error())
	}

	accessToken, err := userBiz.tokenProvider.Generate(tokenProvider.TokenPayload{UserId: user.Id}, userBiz.expiresIn)

	if nil != err {
		return nil, common.NewFailResponse(err)
	}

	return accessToken, nil
}
