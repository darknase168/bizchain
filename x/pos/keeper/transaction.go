package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pos/types"
)

// GetNextTransactionID returns the next transaction ID
func (k Keeper) GetNextTransactionID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TransactionCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetTransactionCount sets the transaction ID counter
func (k Keeper) SetTransactionCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.TransactionCountKey, bz)
}

// SetTransaction stores a transaction
func (k Keeper) SetTransaction(ctx sdk.Context, transaction types.Transaction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TransactionKey)
	bz := k.cdc.MustMarshal(&transaction)
	store.Set(sdk.Uint64ToBigEndian(transaction.Id), bz)

	// Update the transaction count
	currentCount := k.GetNextTransactionID(ctx)
	if transaction.Id >= currentCount {
		k.SetTransactionCount(ctx, transaction.Id)
	}
}

// GetTransaction retrieves a transaction by ID
func (k Keeper) GetTransaction(ctx sdk.Context, id uint64) (transaction types.Transaction, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TransactionKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return transaction, false
	}
	k.cdc.MustUnmarshal(bz, &transaction)
	return transaction, true
}

// GetAllTransactions returns all transactions
func (k Keeper) GetAllTransactions(ctx sdk.Context) (list []*types.Transaction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TransactionKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var transaction types.Transaction
		k.cdc.MustUnmarshal(iterator.Value(), &transaction)
		list = append(list, &transaction)
	}

	return
}

// GetTransactionsBySeller returns all transactions by a specific seller
func (k Keeper) GetTransactionsBySeller(ctx sdk.Context, seller string) (list []*types.Transaction) {
	allTxs := k.GetAllTransactions(ctx)
	for _, tx := range allTxs {
		if tx.Seller == seller {
			list = append(list, tx)
		}
	}
	return
}
