package paymentmodel

import "card-warhouse/common"

var (
	TypePurchase = 0
	TypeRefund   = 1
	TypeRevert   = 2
)

var (
	TranStatusCreated    = 0
	TranStatusProcessing = 1
	TranStatusCompleted  = 3
	TranStatusFail       = 4
	TranStatusCancel     = 5
	TranStatusRefund     = 6
	TranStatusReview     = 7
)

type Transaction struct {
	common.BaseModel
	PaymentNo   string `json:"payment_no" gorm:"column:payment_no"`
	Type        int    `json:"type" gorm:"column:type;default:0"`
	Acquirer    string `json:"acquirer" gorm:"column:acquirer"`
	AcquirerTxn string `json:"acquirer_txn" gorm:"column:acquirer_txn"`
	AcquirerUrl string `json:"acquirer_url" gorm:"column:acquirer_url"`
	CardPin     string `json:"card_pin" gorm:"column:card_pin"`
	CardSerial  string `json:"card_serial" gorm:"column:card_serial"`
	Status      int    `json:"status" gorm:"status"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionCreate struct {
	common.BaseModel
	PaymentNo   string `json:"payment_no" gorm:"column:payment_no"`
	Type        int    `json:"type" gorm:"column:type;default:0"`
	Acquirer    string `json:"acquirer" gorm:"column:acquirer"`
	AcquirerTxn string `json:"acquirer_txn" gorm:"column:acquirer_txn"`
	AcquirerUrl string `json:"acquirer_url" gorm:"column:acquirer_url"`
	CardPin     string `json:"card_pin" gorm:"column:card_pin"`
	CardSerial  string `json:"card_serial" gorm:"column:card_serial"`
	Status      int    `json:"status" gorm:"status"`
}

func (TransactionCreate) TableName() string {
	return Transaction{}.TableName()
}

type TransactionUpdate struct {
	AcquirerTxn *string `json:"acquirer_txn" gorm:"column:acquirer_txn"`
	AcquirerUrl *string `json:"acquirer_url" gorm:"column:acquirer_url"`
	Status      *int    `json:"status" gorm:"status"`
}

func (TransactionUpdate) TableName() string {
	return Transaction{}.TableName()
}
