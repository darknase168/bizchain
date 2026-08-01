package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/referral/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the referral module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListReferral(),
		GetCmdShowReferral(),
		GetCmdReferralByAgent(),
	)

	return queryCmd
}

// GetCmdListReferral queries all referrals
func GetCmdListReferral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-referral",
		Short: "List all referrals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.ReferralAll(cmd.Context(), &types.QueryAllReferralRequest{})
			if err != nil {
				return err
			}

			for _, r := range resp.Referral {
				fmt.Printf("ID: %d | Agent: %s | Jamaah: %s | Paket: %d | Rate: %s | Commission: %s | Status: %s\n",
					r.Id, r.Agent, r.ReferredJamaah, r.PaketId, r.CommissionRate, r.Commission, r.Status)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowReferral queries a specific referral
func GetCmdShowReferral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-referral [id]",
		Short: "Show a referral by ID",
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
			resp, err := queryClient.Referral(cmd.Context(), &types.QueryGetReferralRequest{Id: id})
			if err != nil {
				return err
			}

			r := resp.Referral
			fmt.Printf("ID: %d\n", r.Id)
			fmt.Printf("Agent: %s\n", r.Agent)
			fmt.Printf("Referred Jamaah: %s\n", r.ReferredJamaah)
			fmt.Printf("Paket ID: %d\n", r.PaketId)
			fmt.Printf("Commission Rate: %s\n", r.CommissionRate)
			fmt.Printf("Commission: %s\n", r.Commission)
			fmt.Printf("Status: %s\n", r.Status)
			fmt.Printf("Paid At: %s\n", r.PaidAt)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdReferralByAgent queries referrals of an agent
func GetCmdReferralByAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "referral-by-agent [agent]",
		Short: "List referrals of an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.ReferralByAgent(cmd.Context(), &types.QueryByAgentRequest{Agent: args[0]})
			if err != nil {
				return err
			}

			for _, r := range resp.Referral {
				fmt.Printf("ID: %d | Jamaah: %s | Paket: %d | Commission: %s | Status: %s\n",
					r.Id, r.ReferredJamaah, r.PaketId, r.Commission, r.Status)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
