package acquirer

var (
	NapasAcquirer   = "Napas"
	FuncardAcquirer = "Funcard"
)

var (
	StatusProcessing = 1
	StatusCompleted  = 2
	StatusFail       = 3
	StatusCancel     = 4
	StatusRefund     = 5
	StatusReview     = 6
)

var (
	TypeInstant    = 1
	TypeNotInstant = 2
)

type Acquirer interface {
	Name() string
	Purchase(payload map[string]interface{}) (*Result, error)
	Inquire(reference string) (*Result, error)
}

type Result struct {
	Status        int
	Type          int
	PaymentNo     string
	Amount        int
	AcquirerTxnId string
	AcquirerUrl   string
	ExtraData     map[string]interface{}
	FailReason    string
}
