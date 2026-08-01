package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/jamaah/types"
)

// Keeper of the jamaah module
type Keeper struct {
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	cdc      codec.BinaryCodec
	authority string
}

// NewKeeper creates a new jamaah Keeper instance
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

// GetNextJamaahID returns the next jamaah ID
func (k Keeper) GetNextJamaahID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.JamaahCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetJamaahCount sets the jamaah ID counter
func (k Keeper) SetJamaahCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.JamaahCountKey, bz)
}

// SetJamaah stores a jamaah
func (k Keeper) SetJamaah(ctx sdk.Context, jamaah types.Jamaah) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.JamaahKey)
	bz := k.cdc.MustMarshal(&jamaah)
	store.Set(sdk.Uint64ToBigEndian(jamaah.Id), bz)

	currentCount := k.GetNextJamaahID(ctx)
	if jamaah.Id >= currentCount {
		k.SetJamaahCount(ctx, jamaah.Id)
	}
}

// GetJamaah retrieves a jamaah by ID
func (k Keeper) GetJamaah(ctx sdk.Context, id uint64) (jamaah types.Jamaah, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.JamaahKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return jamaah, false
	}
	k.cdc.MustUnmarshal(bz, &jamaah)
	return jamaah, true
}

// GetAllJamaah returns all jamaah
func (k Keeper) GetAllJamaah(ctx sdk.Context) (list []*types.Jamaah) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.JamaahKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var jamaah types.Jamaah
		k.cdc.MustUnmarshal(iterator.Value(), &jamaah)
		list = append(list, &jamaah)
	}

	return
}

// RemoveJamaah removes a jamaah from the store
func (k Keeper) RemoveJamaah(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.JamaahKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
