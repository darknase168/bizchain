package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/oleholeh/types"
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

// CreateProduct registers a souvenir product in the marketplace
func (k msgServer) CreateProduct(goCtx context.Context, msg *types.MsgCreateProduct) (*types.MsgCreateProductResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	price, ok := sdkmath.NewIntFromString(msg.Price)
	if !ok || !price.IsPositive() {
		return nil, types.ErrInvalidAmount
	}

	productID := k.GetNextProductID(ctx)
	product := types.OlehOlehProduct{
		Id:          productID,
		Name:        msg.Name,
		Description: msg.Description,
		Price:       msg.Price,
		ImageUrl:    msg.ImageUrl,
		Stock:       msg.Stock,
		Seller:      msg.Creator,
		Category:    msg.Category,
		Status:      "active",
		Creator:     msg.Creator,
		CreatedAt:   nowStr(ctx),
		UpdatedAt:   nowStr(ctx),
	}
	k.SetProduct(ctx, product)

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

// UpdateProduct updates a souvenir product
func (k msgServer) UpdateProduct(goCtx context.Context, msg *types.MsgUpdateProduct) (*types.MsgUpdateProductResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	product, found := k.GetProduct(ctx, msg.Id)
	if !found {
		return nil, types.ErrProductNotFound
	}
	if product.Creator != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if msg.Name != "" {
		product.Name = msg.Name
	}
	if msg.Description != "" {
		product.Description = msg.Description
	}
	if msg.Price != "" {
		product.Price = msg.Price
	}
	if msg.ImageUrl != "" {
		product.ImageUrl = msg.ImageUrl
	}
	if msg.Stock > 0 {
		product.Stock = msg.Stock
	}
	if msg.Category != "" {
		product.Category = msg.Category
	}
	if msg.Status != "" {
		product.Status = msg.Status
	}
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

// OrderProduct places a pre-order for oleh-oleh (paid via wallet)
func (k msgServer) OrderProduct(goCtx context.Context, msg *types.MsgOrderProduct) (*types.MsgOrderProductResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	product, found := k.GetProduct(ctx, msg.ProductId)
	if !found {
		return nil, types.ErrProductNotFound
	}
	if product.Status == "inactive" {
		return nil, types.ErrInvalidStatus
	}
	if msg.Quantity > product.Stock {
		return nil, types.ErrInsufficientStock
	}

	// deduct stock
	product.Stock -= msg.Quantity
	product.UpdatedAt = nowStr(ctx)
	k.SetProduct(ctx, product)

	// compute total
	price, _ := sdkmath.NewIntFromString(product.Price)
	total := price.MulRaw(int64(msg.Quantity))

	orderID := k.GetNextOrderID(ctx)
	order := types.OlehOlehOrder{
		Id:               orderID,
		ProductId:        product.Id,
		ProductName:      product.Name,
		Jamaah:           msg.Creator,
		Quantity:         msg.Quantity,
		Total:            total.String(),
		Status:           "pending",
		ShippingAddress:  msg.ShippingAddress,
		Creator:          msg.Creator,
		CreatedAt:        nowStr(ctx),
		UpdatedAt:        nowStr(ctx),
	}
	k.SetOrder(ctx, order)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOrderProduct,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyOrderID, fmt.Sprintf("%d", orderID)),
			sdk.NewAttribute(types.AttributeKeyProductID, fmt.Sprintf("%d", product.Id)),
			sdk.NewAttribute(types.AttributeKeyTotal, total.String()),
		),
	)

	return &types.MsgOrderProductResponse{OrderId: orderID, Total: total.String()}, nil
}

// UpdateOrderStatus updates an order status (paid, shipped, delivered, cancelled)
func (k msgServer) UpdateOrderStatus(goCtx context.Context, msg *types.MsgUpdateOrderStatus) (*types.MsgUpdateOrderStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	order, found := k.GetOrder(ctx, msg.OrderId)
	if !found {
		return nil, types.ErrOrderNotFound
	}
	// the order creator (buyer) or the product seller may update the status
	product, productFound := k.GetProduct(ctx, order.ProductId)
	if order.Creator != msg.Creator && !(productFound && product.Creator == msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	if msg.Status == "cancelled" {
		// restore stock
		if product, f := k.GetProduct(ctx, order.ProductId); f {
			product.Stock += order.Quantity
			product.UpdatedAt = nowStr(ctx)
			k.SetProduct(ctx, product)
		}
	}
	order.Status = msg.Status
	order.UpdatedAt = nowStr(ctx)
	k.SetOrder(ctx, order)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateOrderStatus,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyOrderID, fmt.Sprintf("%d", msg.OrderId)),
			sdk.NewAttribute(types.AttributeKeyStatus, msg.Status),
		),
	)

	return &types.MsgUpdateOrderStatusResponse{}, nil
}
