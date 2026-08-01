package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/pos/types"
)

// msgServer implements the MsgServer interface
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func nowStr(ctx sdk.Context) string {
	return ctx.BlockTime().UTC().Format("2006-01-02T15:04:05Z")
}

// CreateProduct handles MsgCreateProduct
func (k msgServer) CreateProduct(goCtx context.Context, msg *types.MsgCreateProduct) (*types.MsgCreateProductResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidAddress, err.Error())
	}

	// Validate bundle components (if this product is a bundle)
	if msg.IsBundle {
		if err := k.ValidateBundleComponents(ctx, msg.Components); err != nil {
			return nil, err
		}
	}

	productID := k.GetNextProductID(ctx)
	product := types.Product{
		Id:          productID,
		Name:        msg.Name,
		Description: msg.Description,
		Price:       msg.Price,
		CostPrice:   msg.CostPrice,
		Sku:         msg.Sku,
		Category:    msg.Category,
		ImageUrl:    msg.ImageUrl,
		Stock:       msg.InitialStock,
		Owner:       msg.Creator,
		Active:      true,
		CreatedAt:   nowStr(ctx),
		UpdatedAt:   nowStr(ctx),
		BaseUnitId:  msg.BaseUnitId,
		MinStock:    msg.MinStock,
		Barcode:     msg.Barcode,
		PriceLevels: msg.PriceLevels,
		IsBundle:    msg.IsBundle,
		Components:  msg.Components,
		BranchId:    msg.BranchId,
	}
	k.SetProduct(ctx, product)

	// Post stock-in journal for the initial stock (if any and cost provided)
	if msg.InitialStock > 0 {
		cost, ok := sdkmath.NewIntFromString(msg.CostPrice)
		if ok && cost.IsPositive() {
			totalCost := cost.MulRaw(int64(msg.InitialStock))
			if _, err := k.PostPurchaseJournal(ctx, msg.Creator, productID, msg.InitialStock, totalCost); err != nil {
				return nil, err
			}
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateProduct,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyProductID, fmt.Sprintf("%d", productID)),
			sdk.NewAttribute(types.AttributeKeyProductName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateProductResponse{Id: productID}, nil
}

// UpdateProduct handles MsgUpdateProduct
func (k msgServer) UpdateProduct(goCtx context.Context, msg *types.MsgUpdateProduct) (*types.MsgUpdateProductResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	product, found := k.GetProduct(ctx, msg.Id)
	if !found {
		return nil, types.ErrProductNotFound
	}
	if product.Owner != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if msg.IsBundle {
		if err := k.ValidateBundleComponents(ctx, msg.Components); err != nil {
			return nil, err
		}
	}

	product.Name = msg.Name
	product.Description = msg.Description
	product.Price = msg.Price
	product.CostPrice = msg.CostPrice
	product.Sku = msg.Sku
	product.Category = msg.Category
	product.ImageUrl = msg.ImageUrl
	product.BaseUnitId = msg.BaseUnitId
	product.MinStock = msg.MinStock
	product.Barcode = msg.Barcode
	product.PriceLevels = msg.PriceLevels
	product.IsBundle = msg.IsBundle
	product.Components = msg.Components
	product.UpdatedAt = nowStr(ctx)

	k.SetProduct(ctx, product)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateProduct,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyProductID, fmt.Sprintf("%d", msg.Id)),
		),
	)

	return &types.MsgUpdateProductResponse{}, nil
}

// DeleteProduct handles MsgDeleteProduct
func (k msgServer) DeleteProduct(goCtx context.Context, msg *types.MsgDeleteProduct) (*types.MsgDeleteProductResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	product, found := k.GetProduct(ctx, msg.Id)
	if !found {
		return nil, types.ErrProductNotFound
	}
	if product.Owner != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	product.Active = false
	product.UpdatedAt = nowStr(ctx)
	k.SetProduct(ctx, product)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteProduct,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyProductID, fmt.Sprintf("%d", msg.Id)),
		),
	)

	return &types.MsgDeleteProductResponse{}, nil
}

// CreateTransaction handles MsgCreateTransaction (POS Sale) with multi-unit and bundle support
func (k msgServer) CreateTransaction(goCtx context.Context, msg *types.MsgCreateTransaction) (*types.MsgCreateTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var total sdkmath.Int
	var totalCost sdkmath.Int
	total = sdkmath.ZeroInt()
	totalCost = sdkmath.ZeroInt()
	var txItems []*types.Item

	for _, item := range msg.Items {
		product, found := k.GetProduct(ctx, item.ProductId)
		if !found {
			return nil, types.ErrProductNotFound
		}
		if !product.Active {
			return nil, types.ErrProductInactive
		}

		// Convert quantity to base unit
		baseQty := k.ConvertToBase(ctx, product, item.UnitId, item.Quantity)

		if product.IsBundle {
			// Bundle: validate & deduct component stock
			if err := k.DeductBundleStock(ctx, product, item.Quantity); err != nil {
				return nil, err
			}
		} else {
			if product.Stock < baseQty {
				return nil, fmt.Errorf("insufficient stock for product %s: have %d, want %d",
					product.Name, product.Stock, baseQty)
			}
			// Reduce stock
			product.Stock -= baseQty
			product.UpdatedAt = nowStr(ctx)
			k.SetProduct(ctx, product)
		}

		// Item price (price per sale unit)
		price, ok := sdkmath.NewIntFromString(item.Price)
		if !ok || price.IsNegative() {
			return nil, types.ErrInvalidPrice
		}
		subtotal := price.MulRaw(int64(item.Quantity))
		total = total.Add(subtotal)

		// Cost (HPP) in base unit
		costPerBase := sdkmath.ZeroInt()
		if pc, ok := sdkmath.NewIntFromString(product.CostPrice); ok && pc.IsPositive() {
			costPerBase = pc
		}
		itemCost := costPerBase.MulRaw(int64(baseQty))
		totalCost = totalCost.Add(itemCost)

		txItems = append(txItems, &types.Item{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			UnitId:    item.UnitId,
			Price:     item.Price,
			Subtotal:  subtotal.String(),
			Cost:      itemCost.String(),
		})
	}

	// Discount & tax
	discount := sdkmath.ZeroInt()
	if msg.Discount != "" {
		if d, ok := sdkmath.NewIntFromString(msg.Discount); ok && d.IsPositive() {
			discount = d
		}
	}
	tax := sdkmath.ZeroInt()
	if msg.Tax != "" {
		if t, ok := sdkmath.NewIntFromString(msg.Tax); ok && t.IsPositive() {
			tax = t
		}
	}
	grandTotal := total.Sub(discount).Add(tax)
	if grandTotal.IsNegative() {
		return nil, sdkerrors.Wrap(types.ErrInvalidPrice, "discount exceeds total")
	}

	paymentMethod := msg.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "cash"
	}

	txID := k.GetNextTransactionID(ctx)
	tx := types.Transaction{
		Id:              txID,
		Seller:          msg.Creator,
		CustomerAddress: msg.CustomerAddress,
		Items:           txItems,
		Total:           total.String(),
		Discount:        discount.String(),
		Tax:             tax.String(),
		GrandTotal:      grandTotal.String(),
		PaymentMethod:   paymentMethod,
		Status:          "completed",
		Note:            msg.Note,
		BranchId:        msg.BranchId,
		CreatedAt:       nowStr(ctx),
		UpdatedAt:       nowStr(ctx),
	}
	k.SetTransaction(ctx, tx)

	// Auto-post accounting journal (debit Kas, credit Penjualan; debit HPP, credit Persediaan)
	journalID, err := k.PostSaleJournal(ctx, msg.Creator, txID, grandTotal, totalCost)
	if err != nil {
		return nil, err
	}
	tx.JournalId = journalID
	k.SetTransaction(ctx, tx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTransaction,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTransactionID, fmt.Sprintf("%d", txID)),
			sdk.NewAttribute(types.AttributeKeyTotal, grandTotal.String()),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateTransactionResponse{
		Id:         txID,
		Total:      total.String(),
		GrandTotal: grandTotal.String(),
	}, nil
}

// CancelTransaction handles MsgCancelTransaction (refund)
func (k msgServer) CancelTransaction(goCtx context.Context, msg *types.MsgCancelTransaction) (*types.MsgCancelTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tx, found := k.GetTransaction(ctx, msg.Id)
	if !found {
		return nil, types.ErrTransactionNotFound
	}
	if tx.Seller != msg.Creator {
		return nil, types.ErrUnauthorized
	}
	if tx.Status == "cancelled" {
		return nil, types.ErrTransactionCancelled
	}

	// Restore stock for every item
	var totalCost sdkmath.Int
	totalCost = sdkmath.ZeroInt()
	for _, item := range tx.Items {
		product, found := k.GetProduct(ctx, item.ProductId)
		if !found {
			continue
		}
		if product.IsBundle {
			k.RestoreBundleStock(ctx, product, item.Quantity)
		} else {
			baseQty := k.ConvertToBase(ctx, product, item.UnitId, item.Quantity)
			product.Stock += baseQty
			product.UpdatedAt = nowStr(ctx)
			k.SetProduct(ctx, product)
		}
		if c, ok := sdkmath.NewIntFromString(item.Cost); ok {
			totalCost = totalCost.Add(c)
		}
	}

	// Post reversing journal
	grandTotal, _ := sdkmath.NewIntFromString(tx.GrandTotal)
	journalID, err := k.PostCancellationJournal(ctx, msg.Creator, tx.Id, grandTotal, totalCost)
	if err != nil {
		return nil, err
	}

	tx.Status = "cancelled"
	tx.UpdatedAt = nowStr(ctx)
	tx.JournalId = journalID
	k.SetTransaction(ctx, tx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCancelTransaction,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTransactionID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCancelTransactionResponse{}, nil
}

// AddStock handles MsgAddStock with multi-unit support and accounting
func (k msgServer) AddStock(goCtx context.Context, msg *types.MsgAddStock) (*types.MsgAddStockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	product, found := k.GetProduct(ctx, msg.ProductId)
	if !found {
		return nil, types.ErrProductNotFound
	}
	if product.Owner != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	baseQty := k.ConvertToBase(ctx, product, msg.UnitId, msg.Quantity)
	product.Stock += baseQty
	product.UpdatedAt = nowStr(ctx)

	// Update cost price if provided
	if msg.CostPrice != "" {
		product.CostPrice = msg.CostPrice
	}
	k.SetProduct(ctx, product)

	// Post purchase journal (debit Persediaan, credit Kas)
	costPerUnit, ok := sdkmath.NewIntFromString(msg.CostPrice)
	if !ok {
		costPerUnit, _ = sdkmath.NewIntFromString(product.CostPrice)
	}
	if costPerUnit.IsPositive() {
		totalCost := costPerUnit.MulRaw(int64(baseQty))
		if _, err := k.PostPurchaseJournal(ctx, msg.Creator, msg.ProductId, baseQty, totalCost); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddStock,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyProductID, fmt.Sprintf("%d", msg.ProductId)),
			sdk.NewAttribute(types.AttributeKeyQuantity, fmt.Sprintf("%d", baseQty)),
		),
	)

	return &types.MsgAddStockResponse{}, nil
}

// AdjustStock handles MsgAdjustStock (loss/damage/return)
func (k msgServer) AdjustStock(goCtx context.Context, msg *types.MsgAdjustStock) (*types.MsgAdjustStockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	product, found := k.GetProduct(ctx, msg.ProductId)
	if !found {
		return nil, types.ErrProductNotFound
	}
	if product.Owner != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if msg.Quantity < 0 {
		outQty := uint64(-msg.Quantity)
		if product.Stock < outQty {
			return nil, types.ErrInsufficientStock
		}
		product.Stock -= outQty
	} else {
		product.Stock += uint64(msg.Quantity)
	}
	product.UpdatedAt = nowStr(ctx)
	k.SetProduct(ctx, product)

	// Post adjustment journal if cost is available
	costPerUnit, _ := sdkmath.NewIntFromString(product.CostPrice)
	if costPerUnit.IsPositive() {
		absQty := msg.Quantity
		if absQty < 0 {
			absQty = -absQty
		}
		value := costPerUnit.MulRaw(absQty)
		if _, err := k.PostAdjustmentJournal(ctx, msg.Creator, msg.ProductId, msg.Quantity, value); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAdjustStock,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyProductID, fmt.Sprintf("%d", msg.ProductId)),
			sdk.NewAttribute(types.AttributeKeyQuantity, fmt.Sprintf("%d", msg.Quantity)),
		),
	)

	return &types.MsgAdjustStockResponse{}, nil
}

// CreateUnit handles MsgCreateUnit
func (k msgServer) CreateUnit(goCtx context.Context, msg *types.MsgCreateUnit) (*types.MsgCreateUnitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	unitID := k.GetNextUnitID(ctx)
	unit := types.Unit{
		Id:               unitID,
		Name:             msg.Name,
		Symbol:           msg.Symbol,
		ConversionFactor: msg.ConversionFactor,
		IsBase:           msg.IsBase,
	}
	k.SetUnit(ctx, unit)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateUnit,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyUnitID, fmt.Sprintf("%d", unitID)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateUnitResponse{Id: unitID}, nil
}

// UpdateUnit handles MsgUpdateUnit
func (k msgServer) UpdateUnit(goCtx context.Context, msg *types.MsgUpdateUnit) (*types.MsgUpdateUnitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	unit, found := k.GetUnit(ctx, msg.Id)
	if !found {
		return nil, types.ErrUnitNotFound
	}

	unit.Name = msg.Name
	unit.Symbol = msg.Symbol
	unit.ConversionFactor = msg.ConversionFactor
	unit.IsBase = msg.IsBase
	k.SetUnit(ctx, unit)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateUnit,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyUnitID, fmt.Sprintf("%d", msg.Id)),
		),
	)

	return &types.MsgUpdateUnitResponse{}, nil
}

// DeleteUnit handles MsgDeleteUnit
func (k msgServer) DeleteUnit(goCtx context.Context, msg *types.MsgDeleteUnit) (*types.MsgDeleteUnitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := k.GetUnit(ctx, msg.Id)
	if !found {
		return nil, types.ErrUnitNotFound
	}
	k.RemoveUnit(ctx, msg.Id)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteUnit,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyUnitID, fmt.Sprintf("%d", msg.Id)),
		),
	)

	return &types.MsgDeleteUnitResponse{}, nil
}

// CreateAccount handles MsgCreateAccount (chart of accounts)
func (k msgServer) CreateAccount(goCtx context.Context, msg *types.MsgCreateAccount) (*types.MsgCreateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	accountID := k.GetNextAccountID(ctx)
	account := types.Account{
		Id:          accountID,
		Code:        msg.Code,
		Name:        msg.Name,
		Type:        msg.Type,
		Description: msg.Description,
		CreatedAt:   nowStr(ctx),
	}
	k.SetAccount(ctx, account)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateAccount,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAccountID, fmt.Sprintf("%d", accountID)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateAccountResponse{Id: accountID}, nil
}

// CreateJournalEntry handles MsgCreateJournalEntry (manual journal entry)
func (k msgServer) CreateJournalEntry(goCtx context.Context, msg *types.MsgCreateJournalEntry) (*types.MsgCreateJournalEntryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	journalID, err := k.createJournalEntry(ctx, msg.Creator, msg.ReferenceType, msg.Description, msg.ReferenceId, msg.Lines)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateJournalEntryResponse{Id: journalID}, nil
}
