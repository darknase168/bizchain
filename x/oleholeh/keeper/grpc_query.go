package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/oleholeh/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Product queries a product by ID
func (k Keeper) Product(c context.Context, req *types.QueryGetProductRequest) (*types.QueryGetProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	product, found := k.GetProduct(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &types.QueryGetProductResponse{Product: &product}, nil
}

// ProductAll queries all products
func (k Keeper) ProductAll(c context.Context, req *types.QueryAllProductRequest) (*types.QueryAllProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	list := k.GetAllProducts(ctx)

	return &types.QueryAllProductResponse{
		Product: list,
	}, nil
}

// Order queries an order by ID
func (k Keeper) Order(c context.Context, req *types.QueryGetOrderRequest) (*types.QueryGetOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	order, found := k.GetOrder(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	return &types.QueryGetOrderResponse{Order: &order}, nil
}

// OrderAll queries all orders
func (k Keeper) OrderAll(c context.Context, req *types.QueryAllOrderRequest) (*types.QueryAllOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	list := k.GetAllOrders(ctx)

	return &types.QueryAllOrderResponse{
		Order: list,
	}, nil
}
