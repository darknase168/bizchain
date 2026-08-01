package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "oleholeh"

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
	EventTypeCreateProduct     = "create_product"
	EventTypeUpdateProduct     = "update_product"
	EventTypeOrderProduct      = "order_product"
	EventTypeUpdateOrderStatus = "update_order_status"

	AttributeKeyProductID   = "product_id"
	AttributeKeyProductName = "product_name"
	AttributeKeyOrderID     = "order_id"
	AttributeKeyTotal       = "total"
	AttributeKeyStatus      = "status"
	AttributeKeyCreator     = "creator"
)

// KVStore keys
var (
	// ProductKey prefix for product store
	ProductKey = []byte{0x01}
	// ProductCountKey for the product ID counter
	ProductCountKey = []byte{0x02}
	// OrderKey prefix for order store
	OrderKey = []byte{0x03}
	// OrderCountKey for the order ID counter
	OrderCountKey = []byte{0x04}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x05}
)

// GetProductKey returns the store key for a product by ID
func GetProductKey(id uint64) []byte {
	return append(ProductKey, sdk.Uint64ToBigEndian(id)...)
}

// GetOrderKey returns the store key for an order by ID
func GetOrderKey(id uint64) []byte {
	return append(OrderKey, sdk.Uint64ToBigEndian(id)...)
}
