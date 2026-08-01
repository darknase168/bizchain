package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "hotel"

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
	EventTypeCreateHotel = "create_hotel"
	EventTypeUpdateHotel = "update_hotel"

	AttributeKeyHotelID   = "hotel_id"
	AttributeKeyHotelName = "hotel_name"
	AttributeKeyStatus    = "status"
	AttributeKeyCreator   = "creator"
)

// KVStore keys
var (
	// HotelKey prefix for hotel store
	HotelKey = []byte{0x01}
	// HotelCountKey for the hotel ID counter
	HotelCountKey = []byte{0x02}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x03}
)

// GetHotelKey returns the store key for a hotel by ID
func GetHotelKey(id uint64) []byte {
	return append(HotelKey, sdk.Uint64ToBigEndian(id)...)
}
