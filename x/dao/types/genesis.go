package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ProposalList:  []*DaoProposal{},
		ProposalCount: 0,
		Params:        DefaultParams(),
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
	for _, proposal := range gs.ProposalList {
		if proposal.Id == 0 {
			return ErrInvalidProposalID
		}
	}
	return nil
}

// Validate validates the Params
func (p Params) Validate() error {
	return nil
}
