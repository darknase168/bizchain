package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/referral/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "referral module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateReferral(),
		GetCmdSettleReferral(),
	)

	return txCmd
}

// GetCmdCreateReferral records a referral relationship
func GetCmdCreateReferral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-referral [agent] [referred-jamaah] [paket-id]",
		Short: "Record a referral relationship",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			paketID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			commissionRate, _ := cmd.Flags().GetString("commission-rate")

			msg := &types.MsgCreateReferral{
				Creator:        clientCtx.GetFromAddress().String(),
				Agent:          args[0],
				ReferredJamaah: args[1],
				PaketId:        paketID,
				CommissionRate: commissionRate,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("commission-rate", "", "Commission rate in basis points (50 = 5%), default from params")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSettleReferral pays out commission to the agent
func GetCmdSettleReferral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settle-referral [id]",
		Short: "Settle commission for a referral",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgSettleReferral{
				Creator: clientCtx.GetFromAddress().String(),
				Id:      id,
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
