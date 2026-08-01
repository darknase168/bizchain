package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pos/types"
)

// GetNextUnitID returns the next unit ID
func (k Keeper) GetNextUnitID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UnitCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetUnitCount sets the unit ID counter
func (k Keeper) SetUnitCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.UnitCountKey, bz)
}

// SetUnit stores a unit of measure
func (k Keeper) SetUnit(ctx sdk.Context, unit types.Unit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnitKey)
	bz := k.cdc.MustMarshal(&unit)
	store.Set(sdk.Uint64ToBigEndian(unit.Id), bz)

	currentCount := k.GetNextUnitID(ctx)
	if unit.Id >= currentCount {
		k.SetUnitCount(ctx, unit.Id)
	}
}

// GetUnit retrieves a unit by ID
func (k Keeper) GetUnit(ctx sdk.Context, id uint64) (unit types.Unit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnitKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return unit, false
	}
	k.cdc.MustUnmarshal(bz, &unit)
	return unit, true
}

// GetAllUnits returns all units
func (k Keeper) GetAllUnits(ctx sdk.Context) (list []*types.Unit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnitKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var unit types.Unit
		k.cdc.MustUnmarshal(iterator.Value(), &unit)
		list = append(list, &unit)
	}
	return
}

// RemoveUnit removes a unit from the store
func (k Keeper) RemoveUnit(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnitKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}

// GetBaseUnits returns all units marked as base units
func (k Keeper) GetBaseUnits(ctx sdk.Context) (list []*types.Unit) {
	for _, unit := range k.GetAllUnits(ctx) {
		if unit.IsBase {
			list = append(list, unit)
		}
	}
	return
}
