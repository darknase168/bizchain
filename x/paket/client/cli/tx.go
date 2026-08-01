package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/paket/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "paket module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreatePaket(),
		GetCmdUpdatePaket(),
		GetCmdBookPaket(),
		GetCmdAddReview(),
	)

	return txCmd
}

// GetCmdCreatePaket creates a new package
func GetCmdCreatePaket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-paket [name] [price] [quota]",
		Short: "Create a new package (smart contract paket)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			quota, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			schedule, _ := cmd.Flags().GetString("schedule")
			hotel, _ := cmd.Flags().GetString("hotel")
			airline, _ := cmd.Flags().GetString("airline")
			muthawif, _ := cmd.Flags().GetString("muthawif")
			departure, _ := cmd.Flags().GetString("departure")
			returnDate, _ := cmd.Flags().GetString("return")
			category, _ := cmd.Flags().GetString("category")
			description, _ := cmd.Flags().GetString("description")
			imageURL, _ := cmd.Flags().GetString("image-url")

			msg := &types.MsgCreatePaket{
				Creator:       clientCtx.GetFromAddress().String(),
				Name:          args[0],
				Price:         args[1],
				Quota:         quota,
				Schedule:      schedule,
				Hotel:         hotel,
				Airline:       airline,
				Muthawif:      muthawif,
				DepartureDate: departure,
				ReturnDate:    returnDate,
				Category:      category,
				Description:   description,
				ImageUrl:      imageURL,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("schedule", "", "Schedule, e.g. 2026-01-15 s/d 2026-02-05")
	cmd.Flags().String("hotel", "", "Hotel info")
	cmd.Flags().String("airline", "", "Airline info")
	cmd.Flags().String("muthawif", "", "Muthawif / pembimbing")
	cmd.Flags().String("departure", "", "Departure date")
	cmd.Flags().String("return", "", "Return date")
	cmd.Flags().String("category", "umroh", "Category (haji, umroh, haji_plus)")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("image-url", "", "Image URL")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdatePaket updates an existing package
func GetCmdUpdatePaket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-paket [id] [name] [price] [quota]",
		Short: "Update an existing package",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			quota, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			schedule, _ := cmd.Flags().GetString("schedule")
			hotel, _ := cmd.Flags().GetString("hotel")
			airline, _ := cmd.Flags().GetString("airline")
			muthawif, _ := cmd.Flags().GetString("muthawif")
			departure, _ := cmd.Flags().GetString("departure")
			returnDate, _ := cmd.Flags().GetString("return")
			category, _ := cmd.Flags().GetString("category")
			description, _ := cmd.Flags().GetString("description")
			imageURL, _ := cmd.Flags().GetString("image-url")
			status, _ := cmd.Flags().GetString("status")

			msg := &types.MsgUpdatePaket{
				Creator:       clientCtx.GetFromAddress().String(),
				Id:            id,
				Name:          args[1],
				Price:         args[2],
				Quota:         quota,
				Schedule:      schedule,
				Hotel:         hotel,
				Airline:       airline,
				Muthawif:      muthawif,
				DepartureDate: departure,
				ReturnDate:    returnDate,
				Category:      category,
				Description:   description,
				ImageUrl:      imageURL,
				Status:        status,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("schedule", "", "Schedule")
	cmd.Flags().String("hotel", "", "Hotel info")
	cmd.Flags().String("airline", "", "Airline info")
	cmd.Flags().String("muthawif", "", "Muthawif")
	cmd.Flags().String("departure", "", "Departure date")
	cmd.Flags().String("return", "", "Return date")
	cmd.Flags().String("category", "", "Category (haji, umroh, haji_plus)")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("image-url", "", "Image URL")
	cmd.Flags().String("status", "", "Status (open, full, closed, departed, completed)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdBookPaket books a package
func GetCmdBookPaket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book-paket [paket-id] [quantity]",
		Short: "Book a package (auto-closes quota when full)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			paketID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			quantity, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgBookPaket{
				Creator:  clientCtx.GetFromAddress().String(),
				PaketId:  paketID,
				Quantity: quantity,
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

// GetCmdAddReview adds a rating/review
func GetCmdAddReview() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-review [paket-id] [rating] [comment]",
		Short: "Add a rating/review (1-5 stars)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			paketID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			rating, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			msg := &types.MsgAddReview{
				Creator: clientCtx.GetFromAddress().String(),
				PaketId: paketID,
				Rating:  uint32(rating),
				Comment: args[2],
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
