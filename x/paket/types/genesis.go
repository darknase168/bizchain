package types

import (
	"encoding/json"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PaketList:  []*Paket{},
		PaketCount: 0,
		Params:     DefaultParams(),
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
	for _, paket := range gs.PaketList {
		if paket.Id == 0 {
			return ErrInvalidPaketID
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
