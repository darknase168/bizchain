package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/pos/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "POS module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateProduct(),
		GetCmdUpdateProduct(),
		GetCmdDeleteProduct(),
		GetCmdCreateTransaction(),
		GetCmdCancelTransaction(),
		GetCmdAddStock(),
		GetCmdAdjustStock(),
		GetCmdCreateUnit(),
		GetCmdUpdateUnit(),
		GetCmdDeleteUnit(),
		GetCmdCreateAccount(),
		GetCmdCreateJournalEntry(),
	)

	return txCmd
}

// GetCmdCreateProduct creates a new product
func GetCmdCreateProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-product [name] [price] [sku] [category]",
		Short: "Create a new product",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			description, _ := cmd.Flags().GetString("description")
			imageURL, _ := cmd.Flags().GetString("image-url")
			costPrice, _ := cmd.Flags().GetString("cost-price")
			baseUnitID, _ := cmd.Flags().GetUint64("base-unit")
			initialStock, _ := cmd.Flags().GetUint64("stock")
			minStock, _ := cmd.Flags().GetUint64("min-stock")
			barcode, _ := cmd.Flags().GetString("barcode")
			bundleFlag, _ := cmd.Flags().GetBool("bundle")
			componentsStr, _ := cmd.Flags().GetString("components")
			priceLevelsStr, _ := cmd.Flags().GetString("price-levels")
			branchID, _ := cmd.Flags().GetString("branch-id")

			components, err := parseBundleComponents(componentsStr)
			if err != nil {
				return err
			}
			priceLevels, err := parsePriceLevels(priceLevelsStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateProduct(
				clientCtx.GetFromAddress().String(),
				args[0], description, args[1], costPrice, args[2], args[3], imageURL,
				branchID, baseUnitID, initialStock, minStock, barcode,
				priceLevels, bundleFlag, components,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("description", "", "Product description")
	cmd.Flags().String("image-url", "", "Product image URL")
	cmd.Flags().String("cost-price", "", "Product cost price (HPP)")
	cmd.Flags().Uint64("base-unit", 1, "Base unit ID (default 1 = Pcs)")
	cmd.Flags().Uint64("stock", 0, "Initial stock")
	cmd.Flags().Uint64("min-stock", 0, "Minimum stock alert threshold")
	cmd.Flags().String("barcode", "", "Product barcode")
	cmd.Flags().Bool("bundle", false, "Is a composite/bundle product")
	cmd.Flags().String("components", "", "Bundle components as 'productID:qty,productID:qty'")
	cmd.Flags().String("price-levels", "", "Price levels as 'level:price:minQty,level:price:minQty'")
	cmd.Flags().String("branch-id", "", "Branch ID for multi-branch support")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateProduct updates an existing product
func GetCmdUpdateProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-product [id] [name] [price] [sku] [category]",
		Short: "Update an existing product",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid product ID: %s", args[0])
			}

			description, _ := cmd.Flags().GetString("description")
			imageURL, _ := cmd.Flags().GetString("image-url")
			costPrice, _ := cmd.Flags().GetString("cost-price")
			baseUnitID, _ := cmd.Flags().GetUint64("base-unit")
			minStock, _ := cmd.Flags().GetUint64("min-stock")
			barcode, _ := cmd.Flags().GetString("barcode")
			bundleFlag, _ := cmd.Flags().GetBool("bundle")
			componentsStr, _ := cmd.Flags().GetString("components")
			priceLevelsStr, _ := cmd.Flags().GetString("price-levels")

			components, err := parseBundleComponents(componentsStr)
			if err != nil {
				return err
			}
			priceLevels, err := parsePriceLevels(priceLevelsStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateProduct(
				clientCtx.GetFromAddress().String(), id,
				args[1], description, args[2], costPrice, args[3], args[4], imageURL,
				baseUnitID, minStock, barcode,
				priceLevels, bundleFlag, components,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("description", "", "Product description")
	cmd.Flags().String("image-url", "", "Product image URL")
	cmd.Flags().String("cost-price", "", "Product cost price (HPP)")
	cmd.Flags().Uint64("base-unit", 1, "Base unit ID (default 1 = Pcs)")
	cmd.Flags().Uint64("min-stock", 0, "Minimum stock alert threshold")
	cmd.Flags().String("barcode", "", "Product barcode")
	cmd.Flags().Bool("bundle", false, "Is a composite/bundle product")
	cmd.Flags().String("components", "", "Bundle components as 'productID:qty,productID:qty'")
	cmd.Flags().String("price-levels", "", "Price levels as 'level:price:minQty,level:price:minQty'")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDeleteProduct deactivates a product
func GetCmdDeleteProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-product [id]",
		Short: "Deactivate a product",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid product ID: %s", args[0])
			}

			msg := types.NewMsgDeleteProduct(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCreateTransaction creates a new POS transaction (sale)
func GetCmdCreateTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-transaction [customer-address] [product-ids] [quantities] [prices]",
		Short: "Create a new POS transaction",
		Long: `Create a POS transaction. 
Example: 
  bizchaind tx pos create-transaction point1abc... "1,2" "2,1" "1000,500" --from mykey --payment cash
  This sells: 2x product 1 at 1000 POINT each, and 1x product 2 at 500 POINT each`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			items, err := parseItems(args[1], args[2], args[3])
			if err != nil {
				return err
			}

			// Optional units: --units "unitID,unitID" aligned with product ids
			unitsStr, _ := cmd.Flags().GetString("units")
			unitIDs := splitComma(unitsStr)
			if len(unitIDs) > 0 {
				if len(unitIDs) != len(items) {
					return fmt.Errorf("units count must match items count")
				}
				for i, u := range unitIDs {
					uid, err := strconv.ParseUint(u, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid unit ID: %s", u)
					}
					items[i].UnitId = uid
				}
			}

			note, _ := cmd.Flags().GetString("tx-note")
			paymentMethod, _ := cmd.Flags().GetString("payment")
			discount, _ := cmd.Flags().GetString("discount")
			tax, _ := cmd.Flags().GetString("tax")
			branchID, _ := cmd.Flags().GetString("branch-id")

			msg := types.NewMsgCreateTransaction(
				clientCtx.GetFromAddress().String(),
				args[0], note, paymentMethod, discount, tax, branchID, items,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("tx-note", "", "Transaction note")
	cmd.Flags().String("payment", "cash", "Payment method (cash, qris, transfer, credit)")
	cmd.Flags().String("discount", "", "Discount amount")
	cmd.Flags().String("tax", "", "Tax amount")
	cmd.Flags().String("units", "", "Unit IDs per item, e.g. \"1,2\"")
	cmd.Flags().String("branch-id", "", "Branch ID for multi-branch support")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdCancelTransaction cancels/refunds a transaction
func GetCmdCancelTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-transaction [id]",
		Short: "Cancel/refund a POS transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid transaction ID: %s", args[0])
			}

			msg := types.NewMsgCancelTransaction(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdAddStock adds stock to a product
func GetCmdAddStock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-stock [product-id] [quantity]",
		Short: "Add stock to a product",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid product ID: %s", args[0])
			}

			quantity, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid quantity: %s", args[1])
			}

			unitID, _ := cmd.Flags().GetUint64("unit")
			costPrice, _ := cmd.Flags().GetString("cost-price")

			msg := types.NewMsgAddStock(clientCtx.GetFromAddress().String(), productID, quantity, unitID, costPrice)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64("unit", 0, "Unit ID (0 = base unit)")
	cmd.Flags().String("cost-price", "", "Purchase cost price (HPP)")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdAdjustStock adjusts stock (loss/damage/return)
func GetCmdAdjustStock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adjust-stock [product-id] [quantity]",
		Short: "Adjust stock (positive=in, negative=out)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid product ID: %s", args[0])
			}

			quantity, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid quantity: %s", args[1])
			}

			reason, _ := cmd.Flags().GetString("reason")

			msg := types.NewMsgAdjustStock(clientCtx.GetFromAddress().String(), productID, quantity, reason)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("reason", "", "Adjustment reason")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCreateUnit creates a unit of measure
func GetCmdCreateUnit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-unit [name] [symbol] [conversion-factor]",
		Short: "Create a unit of measure (multi satuan)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			factor, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversion factor: %s", args[2])
			}

			isBase, _ := cmd.Flags().GetBool("base")

			msg := types.NewMsgCreateUnit(clientCtx.GetFromAddress().String(), args[0], args[1], factor, isBase)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("base", false, "Is a base unit")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdUpdateUnit updates a unit of measure
func GetCmdUpdateUnit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-unit [id] [name] [symbol] [conversion-factor]",
		Short: "Update a unit of measure",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid unit ID: %s", args[0])
			}

			factor, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversion factor: %s", args[3])
			}

			isBase, _ := cmd.Flags().GetBool("base")

			msg := types.NewMsgUpdateUnit(clientCtx.GetFromAddress().String(), id, args[1], args[2], factor, isBase)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("base", false, "Is a base unit")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDeleteUnit deletes a unit of measure
func GetCmdDeleteUnit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-unit [id]",
		Short: "Delete a unit of measure",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid unit ID: %s", args[0])
			}

			msg := types.NewMsgDeleteUnit(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCreateAccount creates a chart of accounts entry
func GetCmdCreateAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-account [code] [name] [type]",
		Short: "Create a chart of accounts entry (akuntansi)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			description, _ := cmd.Flags().GetString("description")

			msg := types.NewMsgCreateAccount(clientCtx.GetFromAddress().String(), args[0], args[1], args[2], description)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("description", "", "Account description")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCreateJournalEntry creates a manual journal entry
func GetCmdCreateJournalEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-journal-entry [description] [account-ids] [debits] [credits]",
		Short: "Create a manual journal entry (akuntansi)",
		Long: `Create a journal entry with balanced debit/credit lines.
Example:
  bizchaind tx pos create-journal-entry "Beban listrik" "13,14" "500000,0" "0,500000" --from mykey
  This debits account 13 by 500000 and credits account 14 by 500000.`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			accountIDs := splitComma(args[1])
			debits := splitComma(args[2])
			credits := splitComma(args[3])

			if len(accountIDs) != len(debits) || len(accountIDs) != len(credits) {
				return fmt.Errorf("account IDs, debits and credits must have the same length")
			}

			lines := make([]*types.JournalLine, len(accountIDs))
			for i := range accountIDs {
				accID, err := strconv.ParseUint(accountIDs[i], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid account ID: %s", accountIDs[i])
				}
				lines[i] = &types.JournalLine{
					AccountId: accID,
					Debit:     debits[i],
					Credit:    credits[i],
				}
			}

			referenceType, _ := cmd.Flags().GetString("reference-type")
			referenceID, _ := cmd.Flags().GetUint64("reference-id")

			msg := types.NewMsgCreateJournalEntry(
				clientCtx.GetFromAddress().String(),
				referenceType, args[0], referenceID, lines,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("reference-type", "manual", "Reference type (sale, stock_in, adjustment, manual)")
	cmd.Flags().Uint64("reference-id", 0, "Reference ID")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// parseItems parses comma-separated product IDs, quantities, and prices into Item slices
func parseItems(productIDs, quantities, prices string) ([]*types.Item, error) {
	pidStrs := splitComma(productIDs)
	qtyStrs := splitComma(quantities)
	prcStrs := splitComma(prices)

	if len(pidStrs) != len(qtyStrs) || len(pidStrs) != len(prcStrs) {
		return nil, fmt.Errorf("product IDs, quantities, and prices must have the same number of items")
	}

	items := make([]*types.Item, len(pidStrs))
	for i := range pidStrs {
		pid, err := strconv.ParseUint(pidStrs[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %s", pidStrs[i])
		}

		qty, err := strconv.ParseUint(qtyStrs[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid quantity: %s", qtyStrs[i])
		}

		items[i] = &types.Item{
			ProductId: pid,
			Quantity:  qty,
			Price:     prcStrs[i],
		}
	}

	return items, nil
}

// parseBundleComponents parses "productID:qty,productID:qty" into BundleComponents
func parseBundleComponents(s string) ([]*types.BundleComponent, error) {
	if s == "" {
		return nil, nil
	}
	parts := splitComma(s)
	components := make([]*types.BundleComponent, 0, len(parts))
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid component %q, expected productID:qty", part)
		}
		productID, err := strconv.ParseUint(kv[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid component product ID: %s", kv[0])
		}
		qty, err := strconv.ParseUint(kv[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid component quantity: %s", kv[1])
		}
		components = append(components, &types.BundleComponent{ProductId: productID, Quantity: qty})
	}
	return components, nil
}

// parsePriceLevels parses "level:price:minQty,level:price:minQty" into PriceLevels
func parsePriceLevels(s string) ([]*types.PriceLevel, error) {
	if s == "" {
		return nil, nil
	}
	parts := splitComma(s)
	levels := make([]*types.PriceLevel, 0, len(parts))
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 3)
		if len(kv) != 3 {
			return nil, fmt.Errorf("invalid price level %q, expected level:price:minQty", part)
		}
		minQty, err := strconv.ParseUint(kv[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min quantity: %s", kv[2])
		}
		levels = append(levels, &types.PriceLevel{Level: kv[0], Price: kv[1], MinQuantity: minQty})
	}
	return levels, nil
}

// splitComma splits a string by comma
func splitComma(s string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	current := ""
	for _, c := range s {
		if c == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
