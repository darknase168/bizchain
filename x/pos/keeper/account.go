package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pos/types"
)

// GetNextAccountID returns the next account ID
func (k Keeper) GetNextAccountID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AccountCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetAccountCount sets the account ID counter
func (k Keeper) SetAccountCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.AccountCountKey, bz)
}

// SetAccount stores a chart of accounts entry
func (k Keeper) SetAccount(ctx sdk.Context, account types.Account) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountKey)
	bz := k.cdc.MustMarshal(&account)
	store.Set(sdk.Uint64ToBigEndian(account.Id), bz)

	currentCount := k.GetNextAccountID(ctx)
	if account.Id >= currentCount {
		k.SetAccountCount(ctx, account.Id)
	}
}

// GetAccount retrieves an account by ID
func (k Keeper) GetAccount(ctx sdk.Context, id uint64) (account types.Account, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return account, false
	}
	k.cdc.MustUnmarshal(bz, &account)
	return account, true
}

// GetAllAccounts returns all accounts
func (k Keeper) GetAllAccounts(ctx sdk.Context) (list []*types.Account) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var account types.Account
		k.cdc.MustUnmarshal(iterator.Value(), &account)
		list = append(list, &account)
	}
	return
}

// GetAccountsByType returns all accounts of a given type
func (k Keeper) GetAccountsByType(ctx sdk.Context, accType string) (list []*types.Account) {
	for _, account := range k.GetAllAccounts(ctx) {
		if account.Type == accType {
			list = append(list, account)
		}
	}
	return
}

// RemoveAccount removes an account from the store
func (k Keeper) RemoveAccount(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
