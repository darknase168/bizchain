package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/pos/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the POS module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListProduct(),
		GetCmdShowProduct(),
		GetCmdListTransaction(),
		GetCmdShowTransaction(),
		GetCmdListUnit(),
		GetCmdShowUnit(),
		GetCmdListAccount(),
		GetCmdShowAccount(),
		GetCmdListJournal(),
		GetCmdShowJournal(),
		GetCmdTrialBalance(),
		GetCmdIncomeStatement(),
		GetCmdBalanceSheet(),
		GetCmdLedger(),
		GetCmdPriceLevelReport(),
	)

	return queryCmd
}

// GetCmdListProduct queries all products
func GetCmdListProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-product",
		Short: "List all products",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.ProductAll(cmd.Context(), &types.QueryAllProductRequest{})
			if err != nil {
				return err
			}

			for _, product := range resp.Product {
				fmt.Printf("ID: %d | Name: %s | Price: %s | Cost: %s | SKU: %s | Stock: %d | Bundle: %t | Active: %t\n",
					product.Id, product.Name, product.Price, product.CostPrice, product.Sku, product.Stock, product.IsBundle, product.Active)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowProduct queries a specific product
func GetCmdShowProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-product [id]",
		Short: "Show a product by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid product ID: %s", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Product(cmd.Context(), &types.QueryGetProductRequest{Id: id})
			if err != nil {
				return err
			}

			p := resp.Product
			fmt.Printf("ID: %d\n", p.Id)
			fmt.Printf("Name: %s\n", p.Name)
			fmt.Printf("Description: %s\n", p.Description)
			fmt.Printf("Price: %s\n", p.Price)
			fmt.Printf("Cost Price: %s\n", p.CostPrice)
			fmt.Printf("SKU: %s\n", p.Sku)
			fmt.Printf("Category: %s\n", p.Category)
			fmt.Printf("Stock: %d\n", p.Stock)
			fmt.Printf("Base Unit ID: %d\n", p.BaseUnitId)
			fmt.Printf("Min Stock: %d\n", p.MinStock)
			fmt.Printf("Barcode: %s\n", p.Barcode)
			fmt.Printf("Is Bundle: %t\n", p.IsBundle)
			fmt.Printf("Owner: %s\n", p.Owner)
			fmt.Printf("Active: %t\n", p.Active)
			if len(p.PriceLevels) > 0 {
				fmt.Println("Price Levels:")
				for _, pl := range p.PriceLevels {
					fmt.Printf("  - %s: %s (min qty %d)\n", pl.Level, pl.Price, pl.MinQuantity)
				}
			}
			if len(p.Components) > 0 {
				fmt.Println("Bundle Components:")
				for _, c := range p.Components {
					fmt.Printf("  - Product %d x %d\n", c.ProductId, c.Quantity)
				}
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListTransaction queries all transactions
func GetCmdListTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-transaction",
		Short: "List all transactions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.TransactionAll(cmd.Context(), &types.QueryAllTransactionRequest{})
			if err != nil {
				return err
			}

			for _, tx := range resp.Transaction {
				fmt.Printf("ID: %d | Seller: %s | Total: %s | Grand Total: %s | Payment: %s | Status: %s | Items: %d\n",
					tx.Id, tx.Seller, tx.Total, tx.GrandTotal, tx.PaymentMethod, tx.Status, len(tx.Items))
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowTransaction queries a specific transaction
func GetCmdShowTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-transaction [id]",
		Short: "Show a transaction by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid transaction ID: %s", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Transaction(cmd.Context(), &types.QueryGetTransactionRequest{Id: id})
			if err != nil {
				return err
			}

			tx := resp.Transaction
			fmt.Printf("Transaction ID: %d\n", tx.Id)
			fmt.Printf("Seller: %s\n", tx.Seller)
			fmt.Printf("Customer: %s\n", tx.CustomerAddress)
			fmt.Printf("Total: %s\n", tx.Total)
			fmt.Printf("Discount: %s\n", tx.Discount)
			fmt.Printf("Tax: %s\n", tx.Tax)
			fmt.Printf("Grand Total: %s\n", tx.GrandTotal)
			fmt.Printf("Payment: %s\n", tx.PaymentMethod)
			fmt.Printf("Status: %s\n", tx.Status)
			fmt.Printf("Note: %s\n", tx.Note)
			fmt.Printf("Journal ID: %d\n", tx.JournalId)
			fmt.Printf("Created: %s\n", tx.CreatedAt)
			fmt.Println("\nItems:")
			for _, item := range tx.Items {
				fmt.Printf("  - Product ID: %d | Qty: %d | Unit: %d | Price: %s | Subtotal: %s | Cost: %s\n",
					item.ProductId, item.Quantity, item.UnitId, item.Price, item.Subtotal, item.Cost)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListUnit queries all units
func GetCmdListUnit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-unit",
		Short: "List all units of measure (multi satuan)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.UnitAll(cmd.Context(), &types.QueryAllUnitRequest{})
			if err != nil {
				return err
			}

			for _, unit := range resp.Unit {
				fmt.Printf("ID: %d | Name: %s | Symbol: %s | Factor: %d | Base: %t\n",
					unit.Id, unit.Name, unit.Symbol, unit.ConversionFactor, unit.IsBase)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowUnit queries a specific unit
func GetCmdShowUnit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-unit [id]",
		Short: "Show a unit by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid unit ID: %s", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Unit(cmd.Context(), &types.QueryGetUnitRequest{Id: id})
			if err != nil {
				return err
			}

			u := resp.Unit
			fmt.Printf("ID: %d\n", u.Id)
			fmt.Printf("Name: %s\n", u.Name)
			fmt.Printf("Symbol: %s\n", u.Symbol)
			fmt.Printf("Conversion Factor: %d\n", u.ConversionFactor)
			fmt.Printf("Is Base: %t\n", u.IsBase)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListAccount queries all accounts
func GetCmdListAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-account",
		Short: "List all chart of accounts (akuntansi)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.AccountAll(cmd.Context(), &types.QueryAllAccountRequest{})
			if err != nil {
				return err
			}

			for _, account := range resp.Account {
				fmt.Printf("ID: %d | Code: %s | Name: %s | Type: %s\n",
					account.Id, account.Code, account.Name, account.Type)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowAccount queries a specific account
func GetCmdShowAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-account [id]",
		Short: "Show an account by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Account(cmd.Context(), &types.QueryGetAccountRequest{Id: id})
			if err != nil {
				return err
			}

			a := resp.Account
			fmt.Printf("ID: %d\n", a.Id)
			fmt.Printf("Code: %s\n", a.Code)
			fmt.Printf("Name: %s\n", a.Name)
			fmt.Printf("Type: %s\n", a.Type)
			fmt.Printf("Description: %s\n", a.Description)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListJournal queries all journal entries
func GetCmdListJournal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-journal",
		Short: "List all journal entries (akuntansi)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.JournalEntryAll(cmd.Context(), &types.QueryAllJournalEntryRequest{})
			if err != nil {
				return err
			}

			for _, entry := range resp.JournalEntry {
				fmt.Printf("ID: %d | Type: %s | Ref: %d | Desc: %s | Lines: %d\n",
					entry.Id, entry.ReferenceType, entry.ReferenceId, entry.Description, len(entry.Lines))
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowJournal queries a specific journal entry
func GetCmdShowJournal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-journal [id]",
		Short: "Show a journal entry by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid journal ID: %s", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.JournalEntry(cmd.Context(), &types.QueryGetJournalEntryRequest{Id: id})
			if err != nil {
				return err
			}

			e := resp.JournalEntry
			fmt.Printf("ID: %d\n", e.Id)
			fmt.Printf("Type: %s\n", e.ReferenceType)
			fmt.Printf("Reference ID: %d\n", e.ReferenceId)
			fmt.Printf("Description: %s\n", e.Description)
			fmt.Printf("Created: %s\n", e.CreatedAt)
			fmt.Println("Lines:")
			for _, line := range e.Lines {
				fmt.Printf("  - Account %d | Debit: %s | Credit: %s\n", line.AccountId, line.Debit, line.Credit)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdTrialBalance queries the trial balance
func GetCmdTrialBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trial-balance",
		Short: "Show trial balance (akuntansi)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.TrialBalance(cmd.Context(), &types.QueryTrialBalanceRequest{})
			if err != nil {
				return err
			}

			fmt.Println("=== TRIAL BALANCE ===")
			for _, acc := range resp.Accounts {
				fmt.Printf("%s | %s | %s | Debit: %s | Credit: %s | Balance: %s\n",
					acc.Code, acc.Name, acc.Type, acc.Debit, acc.Credit, acc.Balance)
			}
			fmt.Printf("\nTotal Debit: %s | Total Credit: %s\n", resp.TotalDebit, resp.TotalCredit)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdIncomeStatement queries the income statement
func GetCmdIncomeStatement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "income-statement",
		Short: "Show income statement / laba rugi (akuntansi)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.IncomeStatement(cmd.Context(), &types.QueryIncomeStatementRequest{})
			if err != nil {
				return err
			}

			fmt.Println("=== LABA RUGI (INCOME STATEMENT) ===")
			fmt.Println("Pendapatan:")
			for _, acc := range resp.Revenues {
				fmt.Printf("  %s | %s | Credit: %s\n", acc.Code, acc.Name, acc.Credit)
			}
			fmt.Printf("  Total Pendapatan: %s\n", resp.TotalRevenue)
			fmt.Println("Beban:")
			for _, acc := range resp.Expenses {
				fmt.Printf("  %s | %s | Debit: %s\n", acc.Code, acc.Name, acc.Debit)
			}
			fmt.Printf("  Total Beban: %s\n", resp.TotalExpense)
			fmt.Printf("\nLABA BERSIH: %s\n", resp.NetIncome)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdBalanceSheet queries the balance sheet
func GetCmdBalanceSheet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance-sheet",
		Short: "Show balance sheet / neraca (akuntansi)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.BalanceSheet(cmd.Context(), &types.QueryBalanceSheetRequest{})
			if err != nil {
				return err
			}

			fmt.Println("=== NERACA (BALANCE SHEET) ===")
			fmt.Println("ASET:")
			for _, acc := range resp.Assets {
				fmt.Printf("  %s | %s | Balance: %s\n", acc.Code, acc.Name, acc.Balance)
			}
			fmt.Printf("  Total Aset: %s\n", resp.TotalAssets)
			fmt.Println("KEWAJIBAN:")
			for _, acc := range resp.Liabilities {
				fmt.Printf("  %s | %s | Balance: %s\n", acc.Code, acc.Name, acc.Balance)
			}
			fmt.Printf("  Total Kewajiban: %s\n", resp.TotalLiabilities)
			fmt.Println("EKUITAS:")
			for _, acc := range resp.Equities {
				fmt.Printf("  %s | %s | Balance: %s\n", acc.Code, acc.Name, acc.Balance)
			}
			fmt.Printf("  Total Ekuitas: %s\n", resp.TotalEquity)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdLedger queries the general ledger for an account
func GetCmdLedger() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ledger [account-id]",
		Short: "Show general ledger for an account (akuntansi)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			accountID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Ledger(cmd.Context(), &types.QueryLedgerRequest{AccountId: accountID})
			if err != nil {
				return err
			}

			fmt.Printf("=== BUKU BESAR: %s (%s) ===\n", resp.Account.Name, resp.Account.Code)
			for _, line := range resp.Lines {
				fmt.Printf("Jurnal #%d | %s | Debit: %s | Credit: %s | Saldo: %s | %s\n",
					line.JournalEntryId, line.Description, line.Debit, line.Credit, line.Balance, line.CreatedAt)
			}
			fmt.Printf("\nSaldo Akhir: %s\n", resp.EndingBalance)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdPriceLevelReport queries the price level report
func GetCmdPriceLevelReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "price-level-report",
		Short: "Show price level report per product (laporan harga level)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.PriceLevelReport(cmd.Context(), &types.QueryPriceLevelReportRequest{})
			if err != nil {
				return err
			}

			fmt.Println("=== LAPORAN HARGA LEVEL ===")
			for _, item := range resp.Items {
				fmt.Printf("%s | %s | Base: %s %s\n", item.ProductName, item.Sku, item.BasePrice, item.BaseUnit)
				for _, pl := range item.PriceLevels {
					fmt.Printf("  - %s: %s (min qty %d)\n", pl.Level, pl.Price, pl.MinQuantity)
				}
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
