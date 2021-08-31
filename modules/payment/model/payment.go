package paymentmodel

import (
	"card-warhouse/common"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvoiceNoCanNotBeEmpty = errors.New("invoice_no can not be empty")
	ErrAmountInvalid          = errors.New("amount is invalid")
	ErrMethodCanNotBeEmpty    = errors.New("method can not be empty")
	ErrMethodInvalid          = errors.New("method is invalid")
	ErrMaxLenDescription      = errors.New("maximum description is 255 characters")
	ErrInvoiceNoInvalid       = errors.New("invoice no is in in used")
)

var (
	PayStatusCreated          = 0
	PayStatusProcessing       = 1
	PayStatusNotDelivery      = 2
	PayStatusCompleted        = 3
	PayStatusFail             = 4
	PayStatusCancel           = 5
	PayStatusRefund           = 6
	PayStatusReview           = 7
	PayStatusNeedCompensation = 8
)

var MinAmount = 10000
var MaxAmount = 10000000

var AtmCard = "ATM_CARD"
var CreditCard = "CREDIT_CARD"
var FunCard = "FUNCARD"

var Method = []string{AtmCard, CreditCard, FunCard}
var MaxLenDescription = 255
var ExpiredHour = 1

type Payment struct {
	common.BaseModel
	PaymentNo     string       `json:"payment_no" gorm:"column:payment_no"`
	UserId        int          `json:"user_id" gorm:"column:user_id"`
	InvoiceNo     string       `json:"invoice_no" gorm:"column:invoice_no"`
	Amount        int          `json:"amount" gorm:"column:amount;default:0"`
	Description   string       `json:"description" gorm:"column:description"`
	ClientIp      string       `json:"client_ip" gorm:"column:client_ip"`
	Method        string       `json:"method" gorm:"column:method"`
	Status        int          `json:"status" gorm:"column:status"`
	ExpiredAt     time.Time    `json:"expired_at" gorm:"column:expired_at"`
	FailureReason string       `json:"failure_reason" gorm:"column:failure_reason"`
	Transaction   *Transaction `json:"transaction" gorm:"foreignKey:PaymentNo;references:PaymentNo"`
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentCreate struct {
	common.BaseModel
	PaymentNo     string    `json:"payment_no" gorm:"column:payment_no"`
	UserId        int       `json:"user_id" gorm:"column:user_id"`
	InvoiceNo     string    `json:"invoice_no" gorm:"column:invoice_no"`
	Amount        int       `json:"amount" gorm:"column:amount"`
	Description   string    `json:"description" gorm:"column:description"`
	ClientIp      string    `json:"client_ip" gorm:"column:client_ip"`
	Method        string    `json:"method" gorm:"column:method"`
	Status        int       `json:"status" gorm:"column:status"`
	ExpiredAt     time.Time `json:"expired_at" gorm:"column:expired_at"`
	FailureReason string    `json:"failure_reason" gorm:"column:failure_reason"`
}

func (PaymentCreate) TableName() string {
	return Payment{}.TableName()
}

type PaymentUpdate struct {
	common.BaseModel
	Status        *int       `json:"status" gorm:"column:status"`
	ExpiredAt     *time.Time `json:"expired_at" gorm:"column:expired_at"`
	FailureReason *string    `json:"failure_reason" gorm:"column:failure_reason"`
}

func (PaymentUpdate) TableName() string {
	return Payment{}.TableName()
}

type PaymentResult struct {
	PaymentNo  string    `json:"payment_no" gorm:"column:payment_no"`
	InvoiceNo  string    `json:"invoice_no" gorm:"column:invoice_no"`
	Amount     int       `json:"amount" gorm:"column:amount"`
	Status     int       `json:"status" gorm:"column:status"`
	PaymentUrl string    `json:"payment_url"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
}

func (payment *PaymentCreate) Validate() error {
	payment.InvoiceNo = strings.TrimSpace(payment.InvoiceNo)

	if "" == payment.InvoiceNo {
		return ErrInvoiceNoCanNotBeEmpty
	}

	if payment.Amount < MinAmount || payment.Amount >= MaxAmount {
		return ErrAmountInvalid
	}

	payment.Method = strings.TrimSpace(payment.Method)

	if "" == payment.Method {
		return ErrMethodCanNotBeEmpty
	}

	if false == common.FindElement(Method, payment.Method) {
		return ErrMethodInvalid
	}

	if len(payment.Description) > MaxLenDescription {
		return ErrMaxLenDescription
	}

	return nil
}
