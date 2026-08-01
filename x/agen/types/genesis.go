package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AgenList:       []*Agen{},
		AgenCount:      0,
		ComplaintCount: 0,
		Params:         DefaultParams(),
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
	for _, agen := range gs.AgenList {
		if agen.Id == 0 {
			return ErrInvalidAgenID
		}
	}
	return nil
}

// Validate validates the Params
func (p Params) Validate() error {
	return nil
}
