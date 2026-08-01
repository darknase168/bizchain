package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/reward/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the reward module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdRewardBalance(),
		GetCmdListReward(),
		GetCmdShowReward(),
	)

	return queryCmd
}

// GetCmdRewardBalance queries a jamaah's loyalty balance
func GetCmdRewardBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-balance [jamaah]",
		Short: "Show a jamaah's loyalty balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Balance(cmd.Context(), &types.QueryBalanceRequest{Jamaah: args[0]})
			if err != nil {
				return err
			}

			b := resp.Balance
			fmt.Printf("Jamaah: %s\n", b.Jamaah)
			fmt.Printf("Balance: %s\n", b.Balance)
			fmt.Printf("Earned: %s\n", b.Earned)
			fmt.Printf("Redeemed: %s\n", b.Redeemed)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListReward queries all rewards
func GetCmdListReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-reward",
		Short: "List all rewards",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.RewardAll(cmd.Context(), &types.QueryAllRewardRequest{})
			if err != nil {
				return err
			}

			for _, r := range resp.Reward {
				fmt.Printf("ID: %d | Jamaah: %s | Points: %s | Type: %s | Status: %s\n",
					r.Id, r.Jamaah, r.Points, r.RewardType, r.Status)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowReward queries a specific reward
func GetCmdShowReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-reward [id]",
		Short: "Show a reward by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Reward(cmd.Context(), &types.QueryGetRewardRequest{Id: id})
			if err != nil {
				return err
			}

			r := resp.Reward
			fmt.Printf("ID: %d\n", r.Id)
			fmt.Printf("Jamaah: %s\n", r.Jamaah)
			fmt.Printf("Points: %s\n", r.Points)
			fmt.Printf("Type: %s\n", r.RewardType)
			fmt.Printf("Reason: %s\n", r.Reason)
			fmt.Printf("Status: %s\n", r.Status)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
