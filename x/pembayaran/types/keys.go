package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "pembayaran"

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
	EventTypeCreatePembayaran = "create_pembayaran"
	EventTypePayInstallment   = "pay_installment"
	EventTypeReleaseEscrow    = "release_escrow"
	EventTypeRefundPembayaran = "refund_pembayaran"
	EventTypeCancelPembayaran = "cancel_pembayaran"

	AttributeKeyPembayaranID = "pembayaran_id"
	AttributeKeyJamaah       = "jamaah"
	AttributeKeyStatus       = "status"
	AttributeKeyPaid         = "paid"
	AttributeKeyRemaining    = "remaining"
	AttributeKeyStage        = "stage"
	AttributeKeyCreator      = "creator"
)

// KVStore keys
var (
	// PembayaranKey prefix for pembayaran store
	PembayaranKey = []byte{0x01}
	// PembayaranCountKey for the pembayaran ID counter
	PembayaranCountKey = []byte{0x02}
	// InstallmentCountKey for the installment ID counter
	InstallmentCountKey = []byte{0x03}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x04}
)

// GetPembayaranKey returns the store key for a pembayaran by ID
func GetPembayaranKey(id uint64) []byte {
	return append(PembayaranKey, sdk.Uint64ToBigEndian(id)...)
}
