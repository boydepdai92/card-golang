package acquirer

import paymentmodel "card-warhouse/modules/payment/model"

type Factory interface {
	GetAcquirerByMethod(method string) Acquirer
	GetAcquirerByName(name string) Acquirer
}

type factory struct{}

func NewAcquirerFactory() *factory {
	return &factory{}
}

func (f *factory) GetAcquirerByMethod(method string) Acquirer {
	switch method {
	case paymentmodel.AtmCard:
	case paymentmodel.CreditCard:
		return NewNapasAcquirer()
	case paymentmodel.FunCard:
		return NewFuncardAcquirer()
	}

	return nil
}

func (f *factory) GetAcquirerByName(method string) Acquirer {
	switch method {
	case NapasAcquirer:
		return NewNapasAcquirer()
	case FuncardAcquirer:
		return NewFuncardAcquirer()
	}

	return nil
}
