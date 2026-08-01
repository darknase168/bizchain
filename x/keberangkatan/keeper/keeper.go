package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/keberangkatan/types"
)

// Keeper of the keberangkatan module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new keberangkatan Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	authority string,
) *Keeper {
	return &Keeper{
		storeKey:  storeKey,
		memKey:    memKey,
		cdc:       cdc,
		authority: authority,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) cosmossdklog.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetAuthority returns the module's authority
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetNextKeberangkatanID returns the next keberangkatan ID
func (k Keeper) GetNextKeberangkatanID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeberangkatanCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetKeberangkatanCount sets the keberangkatan ID counter
func (k Keeper) SetKeberangkatanCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.KeberangkatanCountKey, bz)
}

// GetNextBaggageID returns the next baggage ID
func (k Keeper) GetNextBaggageID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BaggageCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetBaggageCount sets the baggage ID counter
func (k Keeper) SetBaggageCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.BaggageCountKey, bz)
}

// SetKeberangkatan stores a keberangkatan
func (k Keeper) SetKeberangkatan(ctx sdk.Context, keberangkatan types.Keberangkatan) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeberangkatanKey)
	bz := k.cdc.MustMarshal(&keberangkatan)
	store.Set(sdk.Uint64ToBigEndian(keberangkatan.Id), bz)

	currentCount := k.GetNextKeberangkatanID(ctx)
	if keberangkatan.Id >= currentCount {
		k.SetKeberangkatanCount(ctx, keberangkatan.Id)
	}
}

// GetKeberangkatan retrieves a keberangkatan by ID
func (k Keeper) GetKeberangkatan(ctx sdk.Context, id uint64) (keberangkatan types.Keberangkatan, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeberangkatanKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return keberangkatan, false
	}
	k.cdc.MustUnmarshal(bz, &keberangkatan)
	return keberangkatan, true
}

// GetAllKeberangkatan returns all keberangkatan
func (k Keeper) GetAllKeberangkatan(ctx sdk.Context) (list []*types.Keberangkatan) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeberangkatanKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var keberangkatan types.Keberangkatan
		k.cdc.MustUnmarshal(iterator.Value(), &keberangkatan)
		list = append(list, &keberangkatan)
	}

	return
}

// GetKeberangkatanByJamaah returns journeys of a jamaah
func (k Keeper) GetKeberangkatanByJamaah(ctx sdk.Context, jamaah string) (list []*types.Keberangkatan) {
	for _, keberangkatan := range k.GetAllKeberangkatan(ctx) {
		if keberangkatan.Jamaah == jamaah {
			list = append(list, keberangkatan)
		}
	}
	return
}

// RemoveKeberangkatan removes a keberangkatan from the store
func (k Keeper) RemoveKeberangkatan(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeberangkatanKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
