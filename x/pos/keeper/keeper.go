package keeper

import (
	cosmossdklog "cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pos/types"
)

// Keeper of the POS module
type Keeper struct {
	storeKey    storetypes.StoreKey
	memKey      storetypes.StoreKey
	cdc         codec.BinaryCodec
	authKeeper  types.AccountKeeper
	bankKeeper  types.BankKeeper
	authority   string
}

// NewKeeper creates a new POS Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		storeKey:   storeKey,
		memKey:     memKey,
		cdc:        cdc,
		authKeeper: authKeeper,
		bankKeeper: bankKeeper,
		authority:  authority,
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

// StoreKey returns the module's store key
func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}

// Cdc returns the module's codec
func (k Keeper) Cdc() codec.BinaryCodec {
	return k.cdc
}

// GetProductsByBranch returns all products for a specific branch
func (k Keeper) GetProductsByBranch(ctx sdk.Context, branchID string) (list []*types.Product) {
	allProducts := k.GetAllProducts(ctx)
	for _, product := range allProducts {
		if product.BranchId == branchID {
			list = append(list, product)
		}
	}
	return
}

// GetTransactionsByBranch returns all transactions for a specific branch
func (k Keeper) GetTransactionsByBranch(ctx sdk.Context, branchID string) (list []*types.Transaction) {
	allTransactions := k.GetAllTransactions(ctx)
	for _, transaction := range allTransactions {
		if transaction.BranchId == branchID {
			list = append(list, transaction)
		}
	}
	return
}
