package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/dao/types"
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

// CreateProposal creates a new DAO proposal
func (k msgServer) CreateProposal(goCtx context.Context, msg *types.MsgCreateProposal) (*types.MsgCreateProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	deadline := msg.Deadline
	if deadline == "" {
		deadline = ctx.BlockTime().Add(7 * 24 * time.Hour).UTC().Format("2006-01-02T15:04:05Z")
	}

	proposalID := k.GetNextProposalID(ctx)
	proposal := types.DaoProposal{
		Id:            proposalID,
		Title:         msg.Title,
		Description:   msg.Description,
		Options:       msg.Options,
		Votes:         []*types.Vote{},
		Deadline:      deadline,
		Status:        "active",
		ResultOption:  "",
		Creator:       msg.Creator,
		CreatedAt:     nowStr(ctx),
		ClosedAt:      "",
	}
	k.SetProposal(ctx, proposal)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateProposal,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyTitle, msg.Title),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateProposalResponse{Id: proposalID}, nil
}

// CastVote casts a vote on a proposal
func (k msgServer) CastVote(goCtx context.Context, msg *types.MsgCastVote) (*types.MsgCastVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposal, found := k.GetProposal(ctx, msg.ProposalId)
	if !found {
		return nil, types.ErrProposalNotFound
	}
	if proposal.Status != "active" {
		return nil, types.ErrProposalClosed
	}
	if proposal.Deadline != "" {
		deadline, err := time.Parse("2006-01-02T15:04:05Z", proposal.Deadline)
		if err == nil && ctx.BlockTime().After(deadline) {
			return nil, types.ErrDeadlinePassed
		}
	}

	// validate option
	validOption := false
	for _, o := range proposal.Options {
		if o == msg.Option {
			validOption = true
			break
		}
	}
	if !validOption {
		return nil, types.ErrInvalidOption
	}

	if k.HasVoted(ctx, proposal.Id, msg.Creator) {
		return nil, types.ErrAlreadyVoted
	}

	proposal.Votes = append(proposal.Votes, &types.Vote{
		Voter:    msg.Creator,
		Option:   msg.Option,
		Weight:   "1",
		VotedAt:  nowStr(ctx),
	})
	k.SetProposal(ctx, proposal)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCastVote,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyOption, msg.Option),
		),
	)

	return &types.MsgCastVoteResponse{}, nil
}

// CloseProposal closes a proposal and tallies the result
func (k msgServer) CloseProposal(goCtx context.Context, msg *types.MsgCloseProposal) (*types.MsgCloseProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposal, found := k.GetProposal(ctx, msg.ProposalId)
	if !found {
		return nil, types.ErrProposalNotFound
	}
	if proposal.Status != "active" {
		return nil, types.ErrProposalClosed
	}

	// tally votes
	counts := make(map[string]int64)
	for _, v := range proposal.Votes {
		counts[v.Option]++
	}

	winner := ""
	maxCount := int64(0)
	for option, count := range counts {
		if count > maxCount {
			maxCount = count
			winner = option
		}
	}

	if winner == "" {
		return nil, types.ErrProposalActive
	}

	proposal.Status = "passed"
	proposal.ResultOption = winner
	proposal.ClosedAt = nowStr(ctx)
	k.SetProposal(ctx, proposal)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCloseProposal,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyResult, winner),
		),
	)

	return &types.MsgCloseProposalResponse{
		ResultOption: winner,
		TotalVotes:   fmt.Sprintf("%d", maxCount),
	}, nil
}
