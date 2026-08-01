package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pembayaran/types"
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

// CreatePembayaran creates a new escrow payment with down payment and installments
func (k msgServer) CreatePembayaran(goCtx context.Context, msg *types.MsgCreatePembayaran) (*types.MsgCreatePembayaranResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	total, ok := sdkmath.NewIntFromString(msg.Total)
	if !ok || !total.IsPositive() {
		return nil, types.ErrInvalidAmount
	}
	dp, ok := sdkmath.NewIntFromString(msg.DownPayment)
	if !ok || dp.IsNegative() {
		return nil, types.ErrInvalidAmount
	}
	if dp.GT(total) {
		return nil, types.ErrInvalidAmount
	}

	remaining := total.Sub(dp)

	// Build installments
	var installments []*types.Installment
	installmentID := k.GetNextInstallmentID(ctx)
	if msg.InstallmentCount > 0 {
		dueDates := msg.InstallmentDueDates
		perInstallment := sdkmath.ZeroInt()
		if msg.InstallmentCount > 0 {
			perInstallment = remaining.QuoRaw(int64(msg.InstallmentCount))
		}
		for i := uint32(0); i < msg.InstallmentCount; i++ {
			dueDate := ""
			if int(i) < len(dueDates) {
				dueDate = dueDates[i]
			}
			installments = append(installments, &types.Installment{
				Id:      installmentID,
				Amount:  perInstallment.String(),
				DueDate: dueDate,
				Paid:    false,
				LateFee: "0",
			})
			installmentID++
		}
	}
	if uint64(installmentID-1) > k.GetNextInstallmentID(ctx)-1 {
		k.SetInstallmentCount(ctx, installmentID-1)
	}

	paymentMethod := msg.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "transfer"
	}

	pembayaranID := k.GetNextPembayaranID(ctx)
	status := "pending"
	paid := dp
	if dp.IsPositive() {
		status = "dp_paid"
	}
	pembayaran := types.Pembayaran{
		Id:            pembayaranID,
		Jamaah:        msg.Jamaah,
		PaketId:       msg.PaketId,
		Total:         total.String(),
		DownPayment:   dp.String(),
		Paid:          paid.String(),
		Remaining:     remaining.String(),
		Status:        status,
		PaymentMethod: paymentMethod,
		Installments:  installments,
		EscrowStages:  msg.EscrowStages,
		Creator:       msg.Creator,
		CreatedAt:     nowStr(ctx),
		UpdatedAt:     nowStr(ctx),
	}
	k.SetPembayaran(ctx, pembayaran)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreatePembayaran,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPembayaranID, fmt.Sprintf("%d", pembayaranID)),
			sdk.NewAttribute(types.AttributeKeyJamaah, msg.Jamaah),
			sdk.NewAttribute(types.AttributeKeyStatus, status),
		),
	)

	return &types.MsgCreatePembayaranResponse{Id: pembayaranID}, nil
}

// PayInstallment records an installment payment (cicilan)
func (k msgServer) PayInstallment(goCtx context.Context, msg *types.MsgPayInstallment) (*types.MsgPayInstallmentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pembayaran, found := k.GetPembayaran(ctx, msg.PembayaranId)
	if !found {
		return nil, types.ErrPembayaranNotFound
	}
	if pembayaran.Status == "completed" {
		return nil, types.ErrPaymentCompleted
	}
	if pembayaran.Status == "refunded" || pembayaran.Status == "cancelled" {
		return nil, types.ErrPaymentCancelled
	}

	amount, ok := sdkmath.NewIntFromString(msg.Amount)
	if !ok || !amount.IsPositive() {
		return nil, types.ErrInvalidAmount
	}

	// Find and mark the installment paid
	foundInstallment := false
	for _, inst := range pembayaran.Installments {
		if inst.Id == msg.InstallmentId {
			if inst.Paid {
				return nil, types.ErrInstallmentPaid
			}
			inst.Paid = true
			inst.PaidAt = nowStr(ctx)
			inst.Amount = msg.Amount
			foundInstallment = true
			break
		}
	}
	if !foundInstallment {
		return nil, types.ErrInstallmentNotFound
	}

	// Update paid & remaining
	paid, _ := sdkmath.NewIntFromString(pembayaran.Paid)
	remaining, _ := sdkmath.NewIntFromString(pembayaran.Remaining)
	paid = paid.Add(amount)
	remaining = remaining.Sub(amount)
	if remaining.IsNegative() {
		return nil, types.ErrInvalidAmount
	}

	pembayaran.Paid = paid.String()
	pembayaran.Remaining = remaining.String()
	pembayaran.UpdatedAt = nowStr(ctx)

	completed := false
	if remaining.IsZero() {
		pembayaran.Status = "completed"
		completed = true
	} else if pembayaran.Status == "pending" {
		pembayaran.Status = "active"
	}
	k.SetPembayaran(ctx, pembayaran)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePayInstallment,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPembayaranID, fmt.Sprintf("%d", msg.PembayaranId)),
			sdk.NewAttribute(types.AttributeKeyPaid, paid.String()),
			sdk.NewAttribute(types.AttributeKeyRemaining, remaining.String()),
		),
	)

	return &types.MsgPayInstallmentResponse{
		Paid:      paid.String(),
		Remaining: remaining.String(),
		Completed: completed,
	}, nil
}

// ReleaseEscrow releases escrow funds for a stage (visa, tiket, hotel)
func (k msgServer) ReleaseEscrow(goCtx context.Context, msg *types.MsgReleaseEscrow) (*types.MsgReleaseEscrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pembayaran, found := k.GetPembayaran(ctx, msg.PembayaranId)
	if !found {
		return nil, types.ErrPembayaranNotFound
	}

	for i, stage := range pembayaran.EscrowStages {
		if stage.Name == msg.StageName {
			if stage.Released {
				return nil, types.ErrEscrowReleased
			}
			// Enforce progressive release order (visa -> tiket -> hotel)
			for j := 0; j < i; j++ {
				if !pembayaran.EscrowStages[j].Released {
					return nil, types.ErrEscrowStageOrder
				}
			}
			stage.Released = true
			stage.ReleasedAt = nowStr(ctx)

			// Ensure sufficient funds
			paid, _ := sdkmath.NewIntFromString(pembayaran.Paid)
			stageAmount, _ := sdkmath.NewIntFromString(stage.Amount)
			if paid.LT(stageAmount) {
				return nil, types.ErrInvalidAmount
			}

			pembayaran.UpdatedAt = nowStr(ctx)
			k.SetPembayaran(ctx, pembayaran)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeReleaseEscrow,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
					sdk.NewAttribute(types.AttributeKeyPembayaranID, fmt.Sprintf("%d", msg.PembayaranId)),
					sdk.NewAttribute(types.AttributeKeyStage, msg.StageName),
				),
			)

			return &types.MsgReleaseEscrowResponse{}, nil
		}
	}

	return nil, types.ErrEscrowStageNotFound
}

// RefundPembayaran refunds an escrow payment
func (k msgServer) RefundPembayaran(goCtx context.Context, msg *types.MsgRefundPembayaran) (*types.MsgRefundPembayaranResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pembayaran, found := k.GetPembayaran(ctx, msg.PembayaranId)
	if !found {
		return nil, types.ErrPembayaranNotFound
	}
	if pembayaran.Status == "refunded" {
		return nil, types.ErrInvalidStatus
	}

	pembayaran.Status = "refunded"
	pembayaran.UpdatedAt = nowStr(ctx)
	k.SetPembayaran(ctx, pembayaran)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRefundPembayaran,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPembayaranID, fmt.Sprintf("%d", msg.PembayaranId)),
			sdk.NewAttribute(types.AttributeKeyStatus, "refunded"),
		),
	)

	return &types.MsgRefundPembayaranResponse{}, nil
}

// CancelPembayaran cancels a pending payment
func (k msgServer) CancelPembayaran(goCtx context.Context, msg *types.MsgCancelPembayaran) (*types.MsgCancelPembayaranResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pembayaran, found := k.GetPembayaran(ctx, msg.PembayaranId)
	if !found {
		return nil, types.ErrPembayaranNotFound
	}
	if pembayaran.Status == "completed" {
		return nil, types.ErrPaymentCompleted
	}
	if pembayaran.Status == "refunded" || pembayaran.Status == "cancelled" {
		return nil, types.ErrPaymentCancelled
	}

	pembayaran.Status = "cancelled"
	pembayaran.UpdatedAt = nowStr(ctx)
	k.SetPembayaran(ctx, pembayaran)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCancelPembayaran,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPembayaranID, fmt.Sprintf("%d", msg.PembayaranId)),
			sdk.NewAttribute(types.AttributeKeyStatus, "cancelled"),
		),
	)

	return &types.MsgCancelPembayaranResponse{}, nil
}
