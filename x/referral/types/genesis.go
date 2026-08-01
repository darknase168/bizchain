package types

import (
	"encoding/json"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ReferralList:  []*Referral{},
		ReferralCount: 0,
		Params:        DefaultParams(),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() *Params {
	return &Params{
		ZeroFeeEnabled:        true,
		DefaultCommissionRate: "50", // 50 basis poin = 5%
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	for _, referral := range gs.ReferralList {
		if referral.Id == 0 {
			return ErrInvalidReferralID
		}
	}
	return nil
}

// Validate validates the Params
func (p Params) Validate() error {
	return nil
}

// GetGenesisStateFromAppState returns the GenesisState from the app state
func GetGenesisStateFromAppState(rawState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if rawState[ModuleName] != nil {
		ModuleCdc.MustUnmarshalJSON(rawState[ModuleName], &genesisState)
	}
	return genesisState
}
