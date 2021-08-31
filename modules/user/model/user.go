package usermodel

import (
	"card-warhouse/common"
	"errors"
	"strings"
)

var (
	ErrUsernameRequired        = errors.New("username không thể để trống")
	ErrPasswordRequired        = errors.New("password không thể để trống")
	ErrUsernameExisted         = errors.New("username đã tồn tại")
	ErrUsernameOrPasswordWrong = errors.New("username hoặc mật khẩu không chính xác")
	ErrUserDeactive            = errors.New("user đã bị khoá")
)

var (
	StatusDeactive = 0
)

type User struct {
	common.BaseModel
	Username  string `json:"username" gorm:"column:username"`
	Password  string `json:"-" gorm:"column:password"`
	Salt      string `json:"salt" gorm:"column:salt"`
	Balance   int    `json:"balance" gorm:"column:balance; default:0"`
	SecretKey string `json:"secret_key" gorm:"column:secret_key"`
	Status    int    `json:"status" gorm:"column:status; default:0"`
}

func (u *User) GetId() int {
	return u.Id
}

func (u *User) GetUsername() string {
	return u.Username
}

func (User) TableName() string {
	return "users"
}

type UserCreate struct {
	common.BaseModel
	Username  string `json:"username" gorm:"column:username"`
	Password  string `json:"password" gorm:"column:password"`
	Salt      string `json:"salt" gorm:"column:salt"`
	Balance   int    `json:"balance" gorm:"column:balance; default:0"`
	SecretKey string `json:"secret_key" gorm:"column:secret_key"`
	Status    int    `json:"status" gorm:"column:status; default:0"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

type UserUpdate struct {
	Balance *int `json:"balance" gorm:"column:balance"`
	Status  *int `json:"status" gorm:"column:status"`
}

func (UserUpdate) TableName() string {
	return User{}.TableName()
}

type UserLogin struct {
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

func (userCreate *UserCreate) Validate() error {
	userCreate.Username = strings.TrimSpace(userCreate.Username)

	if "" == userCreate.Username {
		return ErrUsernameRequired
	}

	userCreate.Password = strings.TrimSpace(userCreate.Password)

	if "" == userCreate.Password {
		return ErrPasswordRequired
	}

	return nil
}

func (userLogin *UserLogin) Validate() error {
	userLogin.Username = strings.TrimSpace(userLogin.Username)

	if "" == userLogin.Username {
		return ErrUsernameRequired
	}

	userLogin.Password = strings.TrimSpace(userLogin.Password)

	if "" == userLogin.Password {
		return ErrPasswordRequired
	}

	return nil
}

func (userCreate *UserCreate) GenerateRawPassword(pepper string) string {
	return userCreate.Password + userCreate.Salt + pepper
}

func (userLogin *UserLogin) GenerateRawPassword(salt string, pepper string) string {
	return userLogin.Password + salt + pepper
}
