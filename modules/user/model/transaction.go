package usermodel

import (
	"card-warhouse/common"
)

var (
	TypePurchase = "purchase"
)

var (
	ActionPlus  = "plus"
	ActionMinus = "minus"
)

type Transaction struct {
	common.BaseModel
	Reference    string `json:"reference" gorm:"column:reference"`
	UserId       int    `json:"user_id" gorm:"column:user_id"`
	Type         string `json:"type" gorm:"column:type"`
	Action       string `json:"action" gorm:"column:action"`
	Amount       int    `json:"amount" gorm:"column:amount"`
	BeforeAmount int    `json:"before_amount" gorm:"column:before_amount"`
	AfterAmount  int    `json:"after_amount" gorm:"column:after_amount"`
}

func (Transaction) TableName() string {
	return "balance_transactions"
}

type TransactionCreate struct {
	common.BaseModel
	Reference    string `json:"reference" gorm:"column:reference"`
	UserId       int    `json:"user_id" gorm:"column:user_id"`
	Type         string `json:"type" gorm:"column:type"`
	Action       string `json:"action" gorm:"column:action"`
	Amount       int    `json:"amount" gorm:"column:amount"`
	BeforeAmount int    `json:"before_amount" gorm:"column:before_amount"`
	AfterAmount  int    `json:"after_amount" gorm:"column:after_amount"`
}

func (TransactionCreate) TableName() string {
	return Transaction{}.TableName()
}
