package acquirer

type napasAcquirer struct{}

func NewNapasAcquirer() *napasAcquirer {
	return &napasAcquirer{}
}

func (napas *napasAcquirer) Name() string {
	return NapasAcquirer
}

func (napas *napasAcquirer) Purchase(payload map[string]interface{}) (*Result, error) {
	return &Result{}, nil
}

func (napas *napasAcquirer) Inquire(reference string) (*Result, error) {
	return &Result{}, nil
}
