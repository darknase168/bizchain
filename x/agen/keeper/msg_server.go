package keeper

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/agen/types"
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

// validLevel checks an agent level
func validLevel(level string) bool {
	return level == "pusat" || level == "cabang" || level == "subagen"
}

// CreateAgen registers a new agent in the network (pusat/cabang/subagen)
func (k msgServer) CreateAgen(goCtx context.Context, msg *types.MsgCreateAgen) (*types.MsgCreateAgenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	level := msg.Level
	if level == "" {
		level = "subagen"
	}
	if !validLevel(level) {
		return nil, types.ErrInvalidLevel
	}

	parentID := uint64(0)
	if msg.ParentId != "" && msg.ParentId != "0" {
		pid, err := ParseUint64(msg.ParentId)
		if err != nil {
			return nil, types.ErrInvalidParent
		}
		if _, found := k.GetAgen(ctx, pid); !found {
			return nil, types.ErrInvalidParent
		}
		parentID = pid
	}

	agenID := k.GetNextAgenID(ctx)
	agen := types.Agen{
		Id:             agenID,
		Address:        msg.Address,
		Name:           msg.Name,
		ParentId:       strconv.FormatUint(parentID, 10),
		Level:          level,
		Status:         "active",
		CommissionRate: msg.CommissionRate,
		Score:          "0",
		RatingAvg:      "0",
		TotalDownline:  0,
		TotalSales:     0,
		TotalVolume:    "0",
		Complaints:     []*types.Complaint{},
		Performance:    []*types.AgentPerformance{},
		Creator:        msg.Creator,
		CreatedAt:      nowStr(ctx),
		UpdatedAt:      nowStr(ctx),
	}
	k.SetAgen(ctx, agen)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateAgen,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAgenID, fmt.Sprintf("%d", agenID)),
			sdk.NewAttribute(types.AttributeKeyAgenName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyAgenLevel, level),
			sdk.NewAttribute(types.AttributeKeyParentID, strconv.FormatUint(parentID, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateAgenResponse{Id: agenID}, nil
}

// UpdateAgen updates agent data
func (k msgServer) UpdateAgen(goCtx context.Context, msg *types.MsgUpdateAgen) (*types.MsgUpdateAgenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	agen, found := k.GetAgen(ctx, msg.Id)
	if !found {
		return nil, types.ErrAgenNotFound
	}
	if agen.Creator != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if msg.Name != "" {
		agen.Name = msg.Name
	}
	if msg.ParentId != "" && msg.ParentId != "0" {
		pid, err := ParseUint64(msg.ParentId)
		if err != nil {
			return nil, types.ErrInvalidParent
		}
		if _, f := k.GetAgen(ctx, pid); !f {
			return nil, types.ErrInvalidParent
		}
		agen.ParentId = msg.ParentId
	}
	if msg.Level != "" {
		if !validLevel(msg.Level) {
			return nil, types.ErrInvalidLevel
		}
		agen.Level = msg.Level
	}
	if msg.CommissionRate != "" {
		agen.CommissionRate = msg.CommissionRate
	}
	if msg.Status != "" {
		agen.Status = msg.Status
	}
	agen.UpdatedAt = nowStr(ctx)
	k.SetAgen(ctx, agen)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateAgen,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAgenID, fmt.Sprintf("%d", msg.Id)),
		),
	)

	return &types.MsgUpdateAgenResponse{}, nil
}

// AddComplaint files a complaint against an agent (rekam jejak)
func (k msgServer) AddComplaint(goCtx context.Context, msg *types.MsgAddComplaint) (*types.MsgAddComplaintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	agen, found := k.GetAgen(ctx, msg.AgenId)
	if !found {
		return nil, types.ErrAgenNotFound
	}

	complaintID := k.GetNextComplaintID(ctx)
	complaint := &types.Complaint{
		Id:         complaintID,
		Reporter:   msg.Creator,
		Reason:     msg.Reason,
		Status:     "open",
		Resolution: "",
		CreatedAt:  nowStr(ctx),
		ResolvedAt: "",
	}
	agen.Complaints = append(agen.Complaints, complaint)
	agen.UpdatedAt = nowStr(ctx)
	k.SetAgen(ctx, agen)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddComplaint,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAgenID, fmt.Sprintf("%d", msg.AgenId)),
			sdk.NewAttribute(types.AttributeKeyComplaintID, fmt.Sprintf("%d", complaintID)),
		),
	)

	return &types.MsgAddComplaintResponse{ComplaintId: complaintID}, nil
}

// ResolveComplaint resolves a complaint and recomputes the agent score
func (k msgServer) ResolveComplaint(goCtx context.Context, msg *types.MsgResolveComplaint) (*types.MsgResolveComplaintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	agen, found := k.GetAgen(ctx, msg.AgenId)
	if !found {
		return nil, types.ErrAgenNotFound
	}

	foundComplaint := false
	for _, c := range agen.Complaints {
		if c.Id == msg.ComplaintId {
			if c.Status == "resolved" {
				return nil, types.ErrComplaintResolved
			}
			c.Status = "resolved"
			c.Resolution = msg.Resolution
			c.ResolvedAt = nowStr(ctx)
			foundComplaint = true
			break
		}
	}
	if !foundComplaint {
		return nil, types.ErrComplaintNotFound
	}

	// recompute aggregate score: base 100 minus complaint penalty, blended with rating
	openComplaints := 0
	for _, c := range agen.Complaints {
		if c.Status == "open" {
			openComplaints++
		}
	}
	score := 100.0 - float64(openComplaints)*10.0
	if rating, err := strconv.ParseFloat(agen.RatingAvg, 64); err == nil && rating > 0 {
		score = score*0.6 + (rating/5.0*100.0)*0.4
	}
	if score < 0 {
		score = 0
	}
	agen.Score = fmt.Sprintf("%.1f", score)
	agen.UpdatedAt = nowStr(ctx)
	k.SetAgen(ctx, agen)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeResolveComplaint,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAgenID, fmt.Sprintf("%d", msg.AgenId)),
			sdk.NewAttribute(types.AttributeKeyComplaintID, fmt.Sprintf("%d", msg.ComplaintId)),
		),
	)

	return &types.MsgResolveComplaintResponse{}, nil
}

// RecordPerformance records agent performance and recomputes the score
func (k msgServer) RecordPerformance(goCtx context.Context, msg *types.MsgRecordPerformance) (*types.MsgRecordPerformanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	agen, found := k.GetAgen(ctx, msg.AgenId)
	if !found {
		return nil, types.ErrAgenNotFound
	}

	// compute per-period performance score (0-100) from sales & volume
	perfScoreVal := 50.0
	if msg.Sales > 0 {
		perfScoreVal = 50.0 + float64(msg.Sales)*2.0 // +2 per sale, capped at 100
		if perfScoreVal > 100 {
			perfScoreVal = 100
		}
	}

	// record the performance entry
	perf := &types.AgentPerformance{
		Agent:     agen.Address,
		Period:    msg.Period,
		Sales:     msg.Sales,
		Volume:    msg.Volume,
		RatingAvg: msg.RatingAvg,
		Score:     fmt.Sprintf("%.1f", perfScoreVal),
		CreatedAt: nowStr(ctx),
	}
	agen.Performance = append(agen.Performance, perf)

	// aggregate totals
	agen.TotalSales += msg.Sales
	totalVolume, _ := strconv.ParseFloat(agen.TotalVolume, 64)
	vol, _ := strconv.ParseFloat(msg.Volume, 64)
	agen.TotalVolume = fmt.Sprintf("%.0f", totalVolume+vol)

	if msg.RatingAvg != "" {
		agen.RatingAvg = msg.RatingAvg
	}

	// recompute score: blend performance (50%) + rating (40%) - complaint penalty (10%)
	perfScore := 0.0
	if len(agen.Performance) > 0 {
		for _, p := range agen.Performance {
			s, _ := strconv.ParseFloat(p.Score, 64)
			perfScore += s
		}
		perfScore /= float64(len(agen.Performance))
	}
	rating, _ := strconv.ParseFloat(agen.RatingAvg, 64)
	ratingScore := rating / 5.0 * 100.0
	openComplaints := 0
	for _, c := range agen.Complaints {
		if c.Status == "open" {
			openComplaints++
		}
	}
	score := perfScore*0.5 + ratingScore*0.4 - float64(openComplaints)*5.0
	if score < 0 {
		score = 0
	}
	agen.Score = fmt.Sprintf("%.1f", score)
	agen.UpdatedAt = nowStr(ctx)
	k.SetAgen(ctx, agen)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRecordPerformance,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAgenID, fmt.Sprintf("%d", msg.AgenId)),
		),
	)

	return &types.MsgRecordPerformanceResponse{}, nil
}
