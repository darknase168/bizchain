package cli

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/dao/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "dao module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateProposal(),
		GetCmdCastVote(),
		GetCmdCloseProposal(),
	)

	return txCmd
}

// GetCmdCreateProposal creates a new DAO proposal
func GetCmdCreateProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-proposal [title] [options]",
		Short: "Create a DAO proposal (options comma-separated, e.g. 'yes,no')",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			options := strings.Split(args[1], ",")
			description, _ := cmd.Flags().GetString("description")
			deadline, _ := cmd.Flags().GetString("deadline")

			msg := &types.MsgCreateProposal{
				Creator:     clientCtx.GetFromAddress().String(),
				Title:       args[0],
				Description: description,
				Options:     options,
				Deadline:    deadline,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("description", "", "Proposal description")
	cmd.Flags().String("deadline", "", "Voting deadline (RFC3339)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdCastVote casts a vote on a proposal
func GetCmdCastVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cast-vote [proposal-id] [option]",
		Short: "Cast a vote on a DAO proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgCastVote{
				Creator:    clientCtx.GetFromAddress().String(),
				ProposalId: proposalID,
				Option:     args[1],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdCloseProposal closes a proposal and tallies the result
func GetCmdCloseProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-proposal [proposal-id]",
		Short: "Close a DAO proposal and tally the result",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgCloseProposal{
				Creator:    clientCtx.GetFromAddress().String(),
				ProposalId: proposalID,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
