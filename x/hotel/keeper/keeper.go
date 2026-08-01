package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/hotel/types"
)

// Keeper of the hotel module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new hotel Keeper instance
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

// GetNextHotelID returns the next hotel ID
func (k Keeper) GetNextHotelID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HotelCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetHotelCount sets the hotel ID counter
func (k Keeper) SetHotelCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.HotelCountKey, bz)
}

// SetHotel stores a hotel
func (k Keeper) SetHotel(ctx sdk.Context, hotel types.Hotel) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HotelKey)
	bz := k.cdc.MustMarshal(&hotel)
	store.Set(sdk.Uint64ToBigEndian(hotel.Id), bz)

	currentCount := k.GetNextHotelID(ctx)
	if hotel.Id >= currentCount {
		k.SetHotelCount(ctx, hotel.Id)
	}
}

// GetHotel retrieves a hotel by ID
func (k Keeper) GetHotel(ctx sdk.Context, id uint64) (hotel types.Hotel, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HotelKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return hotel, false
	}
	k.cdc.MustUnmarshal(bz, &hotel)
	return hotel, true
}

// GetAllHotel returns all hotels
func (k Keeper) GetAllHotel(ctx sdk.Context) (list []*types.Hotel) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HotelKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var hotel types.Hotel
		k.cdc.MustUnmarshal(iterator.Value(), &hotel)
		list = append(list, &hotel)
	}

	return
}

// GetActiveHotel returns active hotels
func (k Keeper) GetActiveHotel(ctx sdk.Context) (list []*types.Hotel) {
	for _, hotel := range k.GetAllHotel(ctx) {
		if hotel.Status == "active" || hotel.Status == "" {
			list = append(list, hotel)
		}
	}
	return
}

// RemoveHotel removes a hotel from the store
func (k Keeper) RemoveHotel(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HotelKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
