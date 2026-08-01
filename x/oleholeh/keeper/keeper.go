package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/oleholeh/types"
)

// Keeper of the oleholeh module
type Keeper struct {
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	cdc      codec.BinaryCodec
	authority string
}

// NewKeeper creates a new oleholeh Keeper instance
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

// GetNextProductID returns the next product ID
func (k Keeper) GetNextProductID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProductCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetProductCount sets the product ID counter
func (k Keeper) SetProductCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.ProductCountKey, bz)
}

// SetProduct stores a product
func (k Keeper) SetProduct(ctx sdk.Context, product types.OlehOlehProduct) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProductKey)
	bz := k.cdc.MustMarshal(&product)
	store.Set(sdk.Uint64ToBigEndian(product.Id), bz)

	currentCount := k.GetNextProductID(ctx)
	if product.Id >= currentCount {
		k.SetProductCount(ctx, product.Id)
	}
}

// GetProduct retrieves a product by ID
func (k Keeper) GetProduct(ctx sdk.Context, id uint64) (product types.OlehOlehProduct, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProductKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return product, false
	}
	k.cdc.MustUnmarshal(bz, &product)
	return product, true
}

// GetAllProducts returns all products
func (k Keeper) GetAllProducts(ctx sdk.Context) (list []*types.OlehOlehProduct) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProductKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var product types.OlehOlehProduct
		k.cdc.MustUnmarshal(iterator.Value(), &product)
		list = append(list, &product)
	}

	return
}

// GetNextOrderID returns the next order ID
func (k Keeper) GetNextOrderID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OrderCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetOrderCount sets the order ID counter
func (k Keeper) SetOrderCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.OrderCountKey, bz)
}

// SetOrder stores an order
func (k Keeper) SetOrder(ctx sdk.Context, order types.OlehOlehOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OrderKey)
	bz := k.cdc.MustMarshal(&order)
	store.Set(sdk.Uint64ToBigEndian(order.Id), bz)

	currentCount := k.GetNextOrderID(ctx)
	if order.Id >= currentCount {
		k.SetOrderCount(ctx, order.Id)
	}
}

// GetOrder retrieves an order by ID
func (k Keeper) GetOrder(ctx sdk.Context, id uint64) (order types.OlehOlehOrder, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OrderKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return order, false
	}
	k.cdc.MustUnmarshal(bz, &order)
	return order, true
}

// GetAllOrders returns all orders
func (k Keeper) GetAllOrders(ctx sdk.Context) (list []*types.OlehOlehOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var order types.OlehOlehOrder
		k.cdc.MustUnmarshal(iterator.Value(), &order)
		list = append(list, &order)
	}

	return
}
