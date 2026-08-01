package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "keberangkatan"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// StageLabels maps journey stages (1-9) to human readable labels
var StageLabels = map[uint64]string{
	1: "Daftar",
	2: "DP Dibayar",
	3: "Visa Diproses",
	4: "Visa Terbit",
	5: "Tiket Terbit",
	6: "Hotel Confirm",
	7: "Manasik",
	8: "Berangkat",
	9: "Pulang",
}

// Event attribute keys
const (
	EventTypeCreateKeberangkatan = "create_keberangkatan"
	EventTypeAdvanceStage        = "advance_stage"
	EventTypeUpdateDeparture     = "update_departure"
	EventTypeAddBaggage          = "add_baggage"
	EventTypeUpdateBaggageStatus = "update_baggage_status"

	AttributeKeyKeberangkatanID = "keberangkatan_id"
	AttributeKeyJamaah          = "jamaah"
	AttributeKeyStage           = "stage"
	AttributeKeyStatusLabel     = "status_label"
	AttributeKeyBaggageID       = "baggage_id"
	AttributeKeyCreator         = "creator"
)

// KVStore keys
var (
	// KeberangkatanKey prefix for keberangkatan store
	KeberangkatanKey = []byte{0x01}
	// KeberangkatanCountKey for the keberangkatan ID counter
	KeberangkatanCountKey = []byte{0x02}
	// BaggageCountKey for the baggage ID counter
	BaggageCountKey = []byte{0x03}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x04}
)

// GetKeberangkatanKey returns the store key for a keberangkatan by ID
func GetKeberangkatanKey(id uint64) []byte {
	return append(KeberangkatanKey, sdk.Uint64ToBigEndian(id)...)
}
