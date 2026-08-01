package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pos/types"
)

// GetNextProductID returns the current product ID counter
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
func (k Keeper) SetProduct(ctx sdk.Context, product types.Product) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProductKey)
	bz := k.cdc.MustMarshal(&product)
	store.Set(sdk.Uint64ToBigEndian(product.Id), bz)

	// Update the product count if this is a new product
	currentCount := k.GetNextProductID(ctx)
	if product.Id >= currentCount {
		k.SetProductCount(ctx, product.Id)
	}
}

// GetProduct retrieves a product by ID
func (k Keeper) GetProduct(ctx sdk.Context, id uint64) (product types.Product, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProductKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return product, false
	}
	k.cdc.MustUnmarshal(bz, &product)
	return product, true
}

// GetAllProducts returns all products
func (k Keeper) GetAllProducts(ctx sdk.Context) (list []*types.Product) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProductKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var product types.Product
		k.cdc.MustUnmarshal(iterator.Value(), &product)
		list = append(list, &product)
	}

	return
}

// GetActiveProducts returns all active products
func (k Keeper) GetActiveProducts(ctx sdk.Context) (list []*types.Product) {
	allProducts := k.GetAllProducts(ctx)
	for _, product := range allProducts {
		if product.Active {
			list = append(list, product)
		}
	}
	return
}

// RemoveProduct removes a product from the store
func (k Keeper) RemoveProduct(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProductKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}

// GetBaseUnitFactor returns the conversion factor for a unit relative to the product's base unit.
// Returns 1 if the product has no unit set or the unit is the base.
func (k Keeper) GetUnitFactor(ctx sdk.Context, product types.Product, unitID uint64) uint64 {
	if unitID == 0 || unitID == product.BaseUnitId {
		return 1
	}
	unit, found := k.GetUnit(ctx, unitID)
	if !found {
		return 1
	}
	return unit.ConversionFactor
}

// ConvertToBase converts a quantity expressed in the given unit to the product's base unit quantity.
func (k Keeper) ConvertToBase(ctx sdk.Context, product types.Product, unitID, quantity uint64) uint64 {
	factor := k.GetUnitFactor(ctx, product, unitID)
	return quantity * factor
}

// GetPriceForLevel returns the price for a given price level name.
// Falls back to the base product price if the level does not exist.
func (k Keeper) GetPriceForLevel(product types.Product, level string) string {
	if level != "" {
		for _, pl := range product.PriceLevels {
			if pl.Level == level {
				return pl.Price
			}
		}
	}
	return product.Price
}

// ValidateBundleComponents checks that all components exist, are active and are not themselves bundles.
func (k Keeper) ValidateBundleComponents(ctx sdk.Context, components []*types.BundleComponent) error {
	for _, comp := range components {
		product, found := k.GetProduct(ctx, comp.ProductId)
		if !found {
			return types.ErrProductNotFound
		}
		if !product.Active {
			return types.ErrProductInactive
		}
		if product.IsBundle {
			return types.ErrInvalidBundle
		}
	}
	return nil
}

// DeductBundleStock reduces the stock of every component of a bundle.
func (k Keeper) DeductBundleStock(ctx sdk.Context, bundle types.Product, quantity uint64) error {
	// Check stock availability for all components first
	for _, comp := range bundle.Components {
		component, found := k.GetProduct(ctx, comp.ProductId)
		if !found {
			return types.ErrProductNotFound
		}
		needed := comp.Quantity * quantity
		if component.Stock < needed {
			return types.ErrInsufficientStock
		}
	}
	// Deduct stock
	for _, comp := range bundle.Components {
		component, found := k.GetProduct(ctx, comp.ProductId)
		if !found {
			return types.ErrProductNotFound
		}
		component.Stock -= comp.Quantity * quantity
		component.UpdatedAt = ctx.BlockTime().UTC().Format("2006-01-02T15:04:05Z")
		k.SetProduct(ctx, component)
	}
	return nil
}

// RestoreBundleStock increases the stock of every component of a bundle (used on cancellation).
func (k Keeper) RestoreBundleStock(ctx sdk.Context, bundle types.Product, quantity uint64) {
	for _, comp := range bundle.Components {
		component, found := k.GetProduct(ctx, comp.ProductId)
		if !found {
			continue
		}
		component.Stock += comp.Quantity * quantity
		component.UpdatedAt = ctx.BlockTime().UTC().Format("2006-01-02T15:04:05Z")
		k.SetProduct(ctx, component)
	}
}
