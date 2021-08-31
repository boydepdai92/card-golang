package acquirer

type funcardAcquirer struct{}

func NewFuncardAcquirer() *funcardAcquirer {
	return &funcardAcquirer{}
}

func (napas *funcardAcquirer) Name() string {
	return FuncardAcquirer
}

func (funcard *funcardAcquirer) Purchase(payload map[string]interface{}) (*Result, error) {
	return &Result{
		Type:          TypeInstant,
		PaymentNo:     payload["payment_no"].(string),
		Amount:        payload["amount"].(int),
		AcquirerTxnId: "123456",
		Status:        StatusCompleted,
	}, nil
}

func (funcard *funcardAcquirer) Inquire(reference string) (*Result, error) {
	return &Result{}, nil
}
