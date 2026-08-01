package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/pos/types"
)

// var _ types.QueryServer = Keeper{}

// Product handles the gRPC query for a single product
func (k Keeper) Product(goCtx context.Context, req *types.QueryGetProductRequest) (*types.QueryGetProductResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	product, found := k.GetProduct(ctx, req.Id)
	if !found {
		return nil, types.ErrProductNotFound
	}

	return &types.QueryGetProductResponse{Product: &product}, nil
}

// ProductAll handles the gRPC query for all products
func (k Keeper) ProductAll(goCtx context.Context, req *types.QueryAllProductRequest) (*types.QueryAllProductResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	products := k.GetAllProducts(ctx)

	return &types.QueryAllProductResponse{
		Product: products,
	}, nil
}

// Transaction handles the gRPC query for a single transaction
func (k Keeper) Transaction(goCtx context.Context, req *types.QueryGetTransactionRequest) (*types.QueryGetTransactionResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	transaction, found := k.GetTransaction(ctx, req.Id)
	if !found {
		return nil, types.ErrTransactionNotFound
	}

	return &types.QueryGetTransactionResponse{Transaction: &transaction}, nil
}

// TransactionAll handles the gRPC query for all transactions
func (k Keeper) TransactionAll(goCtx context.Context, req *types.QueryAllTransactionRequest) (*types.QueryAllTransactionResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	transactions := k.GetAllTransactions(ctx)

	return &types.QueryAllTransactionResponse{
		Transaction: transactions,
	}, nil
}

// Params handles the gRPC query for module parameters
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	params := types.DefaultParams()

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

// Unit handles the gRPC query for a single unit
func (k Keeper) Unit(goCtx context.Context, req *types.QueryGetUnitRequest) (*types.QueryGetUnitResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	unit, found := k.GetUnit(ctx, req.Id)
	if !found {
		return nil, types.ErrUnitNotFound
	}

	return &types.QueryGetUnitResponse{Unit: &unit}, nil
}

// UnitAll handles the gRPC query for all units
func (k Keeper) UnitAll(goCtx context.Context, req *types.QueryAllUnitRequest) (*types.QueryAllUnitResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	units := k.GetAllUnits(ctx)

	return &types.QueryAllUnitResponse{
		Unit: units,
	}, nil
}

// Account handles the gRPC query for a single account
func (k Keeper) Account(goCtx context.Context, req *types.QueryGetAccountRequest) (*types.QueryGetAccountResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	account, found := k.GetAccount(ctx, req.Id)
	if !found {
		return nil, types.ErrAccountNotFound
	}

	return &types.QueryGetAccountResponse{Account: &account}, nil
}

// AccountAll handles the gRPC query for all accounts
func (k Keeper) AccountAll(goCtx context.Context, req *types.QueryAllAccountRequest) (*types.QueryAllAccountResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	accounts := k.GetAllAccounts(ctx)

	return &types.QueryAllAccountResponse{
		Account: accounts,
	}, nil
}

// JournalEntry handles the gRPC query for a single journal entry
func (k Keeper) JournalEntry(goCtx context.Context, req *types.QueryGetJournalEntryRequest) (*types.QueryGetJournalEntryResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	entry, found := k.GetJournalEntry(ctx, req.Id)
	if !found {
		return nil, types.ErrJournalNotFound
	}

	return &types.QueryGetJournalEntryResponse{JournalEntry: &entry}, nil
}

// JournalEntryAll handles the gRPC query for all journal entries
func (k Keeper) JournalEntryAll(goCtx context.Context, req *types.QueryAllJournalEntryRequest) (*types.QueryAllJournalEntryResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	entries := k.GetAllJournalEntries(ctx)

	return &types.QueryAllJournalEntryResponse{
		JournalEntry: entries,
	}, nil
}

// TrialBalance returns the trial balance (akuntansi)
func (k Keeper) TrialBalance(goCtx context.Context, req *types.QueryTrialBalanceRequest) (*types.QueryTrialBalanceResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return k.GetTrialBalance(ctx)
}

// IncomeStatement returns the income statement / laba rugi (akuntansi)
func (k Keeper) IncomeStatement(goCtx context.Context, req *types.QueryIncomeStatementRequest) (*types.QueryIncomeStatementResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return k.GetIncomeStatement(ctx)
}

// BalanceSheet returns the balance sheet / neraca (akuntansi)
func (k Keeper) BalanceSheet(goCtx context.Context, req *types.QueryBalanceSheetRequest) (*types.QueryBalanceSheetResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return k.GetBalanceSheet(ctx)
}

// Ledger returns the general ledger for an account (akuntansi)
func (k Keeper) Ledger(goCtx context.Context, req *types.QueryLedgerRequest) (*types.QueryLedgerResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return k.GetLedger(ctx, req.AccountId)
}

// PriceLevelReport returns the price level report per product
func (k Keeper) PriceLevelReport(goCtx context.Context, req *types.QueryPriceLevelReportRequest) (*types.QueryPriceLevelReportResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return k.GetPriceLevelReport(ctx), nil
}

// ProductByBranch returns all products for a specific branch
func (k Keeper) ProductByBranch(goCtx context.Context, req *types.QueryProductByBranchRequest) (*types.QueryProductByBranchResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	products := k.GetProductsByBranch(ctx, req.BranchId)

	return &types.QueryProductByBranchResponse{
		Product: products,
	}, nil
}

// TransactionByBranch returns all transactions for a specific branch
func (k Keeper) TransactionByBranch(goCtx context.Context, req *types.QueryTransactionByBranchRequest) (*types.QueryTransactionByBranchResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	transactions := k.GetTransactionsByBranch(ctx, req.BranchId)

	return &types.QueryTransactionByBranchResponse{
		Transaction: transactions,
	}, nil
}
