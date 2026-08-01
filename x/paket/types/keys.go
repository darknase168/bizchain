package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "paket"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// Event attribute keys
const (
	EventTypeCreatePaket = "create_paket"
	EventTypeUpdatePaket = "update_paket"
	EventTypeBookPaket   = "book_paket"
	EventTypeAddReview   = "add_review"

	AttributeKeyPaketID      = "paket_id"
	AttributeKeyPaketName    = "paket_name"
	AttributeKeyQuotaClosed  = "quota_closed"
	AttributeKeyReviewer     = "reviewer"
	AttributeKeyRating       = "rating"
	AttributeKeyCreator      = "creator"
)

// KVStore keys
var (
	// PaketKey prefix for paket store
	PaketKey = []byte{0x01}
	// PaketCountKey for the paket ID counter
	PaketCountKey = []byte{0x02}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x03}
)

// GetPaketKey returns the store key for a paket by ID
func GetPaketKey(id uint64) []byte {
	return append(PaketKey, sdk.Uint64ToBigEndian(id)...)
}
