package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/audit/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "audit module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdLogAction(),
	)

	return txCmd
}

// GetCmdLogAction records an audit log entry
func GetCmdLogAction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log-action [module] [action] [actor] [target-id]",
		Short: "Record an audit log entry (immutable trail)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			dataHash, _ := cmd.Flags().GetString("data-hash")
			metadata, _ := cmd.Flags().GetString("metadata")

			msg := &types.MsgLogAction{
				Creator:  clientCtx.GetFromAddress().String(),
				Module:   args[0],
				Action:   args[1],
				Actor:    args[2],
				TargetId: args[3],
				DataHash: dataHash,
				Metadata: metadata,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("data-hash", "", "Hash of the data for verification")
	cmd.Flags().String("metadata", "", "JSON metadata")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
