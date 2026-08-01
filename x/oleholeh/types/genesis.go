package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ProductList:  []*OlehOlehProduct{},
		OrderList:    []*OlehOlehOrder{},
		ProductCount: 0,
		OrderCount:   0,
		Params:       DefaultParams(),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() *Params {
	return &Params{
		ZeroFeeEnabled: true,
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	for _, product := range gs.ProductList {
		if product.Id == 0 {
			return ErrInvalidProductID
		}
	}
	for _, order := range gs.OrderList {
		if order.Id == 0 {
			return ErrInvalidOrderID
		}
	}
	return nil
}

// Validate validates the Params
func (p Params) Validate() error {
	return nil
}
