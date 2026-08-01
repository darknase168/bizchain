package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "pos"

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
	EventTypeDeleteProduct     = "delete_product"
	EventTypeCreateTransaction = "create_transaction"
	EventTypeCancelTransaction = "cancel_transaction"
	EventTypeAddStock          = "add_stock"
	EventTypeAdjustStock       = "adjust_stock"
	EventTypeCreateUnit        = "create_unit"
	EventTypeUpdateUnit        = "update_unit"
	EventTypeDeleteUnit        = "delete_unit"
	EventTypeCreateAccount     = "create_account"
	EventTypeCreateJournal     = "create_journal"

	AttributeKeyProductID      = "product_id"
	AttributeKeyProductName    = "product_name"
	AttributeKeyTransactionID  = "transaction_id"
	AttributeKeyTotal          = "total"
	AttributeKeyQuantity       = "quantity"
	AttributeKeyCreator        = "creator"
	AttributeKeyUnitID         = "unit_id"
	AttributeKeyAccountID      = "account_id"
	AttributeKeyJournalID      = "journal_id"
)

// KVStore keys
var (
	// ProductKey prefix for product store
	ProductKey = []byte{0x01}
	// TransactionKey prefix for transaction store
	TransactionKey = []byte{0x02}
	// ProductCountKey for the product ID counter
	ProductCountKey = []byte{0x03}
	// TransactionCountKey for the transaction ID counter
	TransactionCountKey = []byte{0x04}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x05}
	// UnitKey prefix for unit of measure store
	UnitKey = []byte{0x06}
	// AccountKey prefix for chart of accounts store
	AccountKey = []byte{0x07}
	// JournalKey prefix for journal entry store
	JournalKey = []byte{0x08}
	// UnitCountKey for the unit ID counter
	UnitCountKey = []byte{0x09}
	// AccountCountKey for the account ID counter
	AccountCountKey = []byte{0x0A}
	// JournalCountKey for the journal entry ID counter
	JournalCountKey = []byte{0x0B}
)

// GetProductKey returns the store key for a product by ID
func GetProductKey(id uint64) []byte {
	return append(ProductKey, sdk.Uint64ToBigEndian(id)...)
}

// GetTransactionKey returns the store key for a transaction by ID
func GetTransactionKey(id uint64) []byte {
	return append(TransactionKey, sdk.Uint64ToBigEndian(id)...)
}

// GetUnitKey returns the store key for a unit by ID
func GetUnitKey(id uint64) []byte {
	return append(UnitKey, sdk.Uint64ToBigEndian(id)...)
}

// GetAccountKey returns the store key for an account by ID
func GetAccountKey(id uint64) []byte {
	return append(AccountKey, sdk.Uint64ToBigEndian(id)...)
}

// GetJournalKey returns the store key for a journal entry by ID
func GetJournalKey(id uint64) []byte {
	return append(JournalKey, sdk.Uint64ToBigEndian(id)...)
}
