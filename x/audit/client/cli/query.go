package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/audit/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the audit module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListAuditLog(),
		GetCmdShowAuditLog(),
		GetCmdAuditLogByModule(),
	)

	return queryCmd
}

// GetCmdListAuditLog queries all audit logs
func GetCmdListAuditLog() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-audit-log",
		Short: "List all audit logs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.AuditLogAll(cmd.Context(), &types.QueryAllAuditLogRequest{})
			if err != nil {
				return err
			}

			for _, l := range resp.AuditLog {
				fmt.Printf("ID: %d | Module: %s | Action: %s | Actor: %s | Target: %s | %s\n",
					l.Id, l.Module, l.Action, l.Actor, l.TargetId, l.CreatedAt)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowAuditLog queries a specific audit log
func GetCmdShowAuditLog() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-audit-log [id]",
		Short: "Show an audit log by ID",
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
			resp, err := queryClient.AuditLog(cmd.Context(), &types.QueryGetAuditLogRequest{Id: id})
			if err != nil {
				return err
			}

			l := resp.AuditLog
			fmt.Printf("ID: %d\n", l.Id)
			fmt.Printf("Module: %s\n", l.Module)
			fmt.Printf("Action: %s\n", l.Action)
			fmt.Printf("Actor: %s\n", l.Actor)
			fmt.Printf("Target ID: %s\n", l.TargetId)
			fmt.Printf("Data Hash: %s\n", l.DataHash)
			fmt.Printf("Metadata: %s\n", l.Metadata)
			fmt.Printf("Created At: %s\n", l.CreatedAt)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdAuditLogByModule queries audit logs filtered by module
func GetCmdAuditLogByModule() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit-log-by-module [module]",
		Short: "List audit logs of a module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.AuditLogByModule(cmd.Context(), &types.QueryAuditLogByModuleRequest{Module: args[0]})
			if err != nil {
				return err
			}

			for _, l := range resp.AuditLog {
				fmt.Printf("ID: %d | Action: %s | Actor: %s | Target: %s | %s\n",
					l.Id, l.Action, l.Actor, l.TargetId, l.CreatedAt)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
