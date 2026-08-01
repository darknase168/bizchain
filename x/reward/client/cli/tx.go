package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/reward/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "reward module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdAwardReward(),
		GetCmdRedeemReward(),
	)

	return txCmd
}

// GetCmdAwardReward awards loyalty tokens to a jamaah
func GetCmdAwardReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "award-reward [jamaah] [points]",
		Short: "Award loyalty tokens to a jamaah",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rewardType, _ := cmd.Flags().GetString("reward-type")
			reason, _ := cmd.Flags().GetString("reason")

			msg := &types.MsgAwardReward{
				Creator:    clientCtx.GetFromAddress().String(),
				Jamaah:     args[0],
				Points:     args[1],
				RewardType: rewardType,
				Reason:     reason,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("reward-type", "cashback", "Reward type (cashback, discount, referral_bonus)")
	cmd.Flags().String("reason", "", "Reward reason")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRedeemReward redeems loyalty tokens
func GetCmdRedeemReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem-reward [jamaah] [points]",
		Short: "Redeem loyalty tokens",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			reason, _ := cmd.Flags().GetString("reason")

			msg := &types.MsgRedeemReward{
				Creator: clientCtx.GetFromAddress().String(),
				Jamaah:  args[0],
				Points:  args[1],
				Reason:  reason,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("reason", "", "Redemption reason")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
