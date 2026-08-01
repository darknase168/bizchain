package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AsuransiList:   []*Asuransi{},
		ClaimList:      []*AsuransiClaim{},
		AsuransiCount:  0,
		ClaimCount:     0,
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
	for _, asuransi := range gs.AsuransiList {
		if asuransi.Id == 0 {
			return ErrInvalidAsuransiID
		}
	}
	for _, claim := range gs.ClaimList {
		if claim.Id == 0 {
			return ErrInvalidClaimID
		}
	}
	return nil
}

// Validate validates the Params
func (p Params) Validate() error {
	return nil
}
