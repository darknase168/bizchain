package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/oleholeh/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "oleholeh module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateProduct(),
		GetCmdUpdateProduct(),
		GetCmdOrderProduct(),
		GetCmdUpdateOrderStatus(),
	)

	return txCmd
}

// GetCmdCreateProduct registers a souvenir product
func GetCmdCreateProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-product [name] [price]",
		Short: "Register a souvenir product in the marketplace",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			stock, _ := cmd.Flags().GetUint64("stock")
			description, _ := cmd.Flags().GetString("description")
			imageURL, _ := cmd.Flags().GetString("image-url")
			category, _ := cmd.Flags().GetString("category")

			msg := &types.MsgCreateProduct{
				Creator:     clientCtx.GetFromAddress().String(),
				Name:        args[0],
				Price:       args[1],
				Description: description,
				ImageUrl:    imageURL,
				Stock:       stock,
				Category:    category,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64("stock", 0, "Initial stock")
	cmd.Flags().String("description", "", "Product description")
	cmd.Flags().String("image-url", "", "Product image URL")
	cmd.Flags().String("category", "", "Product category")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateProduct updates a souvenir product
func GetCmdUpdateProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-product [id]",
		Short: "Update a souvenir product",
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

			name, _ := cmd.Flags().GetString("name")
			price, _ := cmd.Flags().GetString("price")
			stock, _ := cmd.Flags().GetUint64("stock")
			description, _ := cmd.Flags().GetString("description")
			category, _ := cmd.Flags().GetString("category")
			status, _ := cmd.Flags().GetString("status")

			msg := &types.MsgUpdateProduct{
				Creator:     clientCtx.GetFromAddress().String(),
				Id:          id,
				Name:        name,
				Price:       price,
				Stock:       stock,
				Description: description,
				Category:    category,
				Status:      status,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("name", "", "Product name")
	cmd.Flags().String("price", "", "Product price")
	cmd.Flags().Uint64("stock", 0, "Product stock")
	cmd.Flags().String("description", "", "Product description")
	cmd.Flags().String("category", "", "Product category")
	cmd.Flags().String("status", "", "Status (active, inactive)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdOrderProduct places a pre-order for oleh-oleh
func GetCmdOrderProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order-product [product-id] [quantity]",
		Short: "Place a pre-order for oleh-oleh (paid via wallet)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			quantity, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			shippingAddress, _ := cmd.Flags().GetString("shipping-address")

			msg := &types.MsgOrderProduct{
				Creator:          clientCtx.GetFromAddress().String(),
				ProductId:        productID,
				Quantity:         quantity,
				ShippingAddress:  shippingAddress,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("shipping-address", "", "Shipping address")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateOrderStatus updates an order status
func GetCmdUpdateOrderStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-order-status [order-id] [status]",
		Short: "Update order status (paid, shipped, delivered, cancelled)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			orderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateOrderStatus{
				Creator: clientCtx.GetFromAddress().String(),
				OrderId: orderID,
				Status:  args[1],
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
