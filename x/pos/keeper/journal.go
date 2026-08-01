package keeper

import (
	"encoding/binary"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pos/types"
)

// GetNextJournalID returns the next journal entry ID
func (k Keeper) GetNextJournalID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.JournalCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetJournalCount sets the journal entry ID counter
func (k Keeper) SetJournalCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.JournalCountKey, bz)
}

// SetJournalEntry stores a journal entry
func (k Keeper) SetJournalEntry(ctx sdk.Context, entry types.JournalEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.JournalKey)
	bz := k.cdc.MustMarshal(&entry)
	store.Set(sdk.Uint64ToBigEndian(entry.Id), bz)

	currentCount := k.GetNextJournalID(ctx)
	if entry.Id >= currentCount {
		k.SetJournalCount(ctx, entry.Id)
	}
}

// GetJournalEntry retrieves a journal entry by ID
func (k Keeper) GetJournalEntry(ctx sdk.Context, id uint64) (entry types.JournalEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.JournalKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return entry, false
	}
	k.cdc.MustUnmarshal(bz, &entry)
	return entry, true
}

// GetAllJournalEntries returns all journal entries
func (k Keeper) GetAllJournalEntries(ctx sdk.Context) (list []*types.JournalEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.JournalKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var entry types.JournalEntry
		k.cdc.MustUnmarshal(iterator.Value(), &entry)
		list = append(list, &entry)
	}
	return
}


// FindAccountByCode returns the first account matching the given code
func (k Keeper) FindAccountByCode(ctx sdk.Context, code string) (types.Account, bool) {
	for _, account := range k.GetAllAccounts(ctx) {
		if account.Code == code {
			return *account, true
		}
	}
	return types.Account{}, false
}

// createJournalEntry creates a journal entry with the given lines and stores it.
// It validates that total debit equals total credit.
func (k Keeper) createJournalEntry(ctx sdk.Context, creator, referenceType, description string, referenceID uint64, lines []*types.JournalLine) (uint64, error) {
	if len(lines) == 0 {
		return 0, types.ErrInvalidJournalEntry
	}
	var totalDebit, totalCredit sdkmath.Int
	totalDebit = sdkmath.ZeroInt()
	totalCredit = sdkmath.ZeroInt()

	// Validate accounts exist
	for _, line := range lines {
		if _, found := k.GetAccount(ctx, line.AccountId); !found {
			return 0, types.ErrAccountNotFound
		}
		amount, ok := sdkmath.NewIntFromString(line.Debit)
		if ok {
			totalDebit = totalDebit.Add(amount)
		}
		amount, ok = sdkmath.NewIntFromString(line.Credit)
		if ok {
			totalCredit = totalCredit.Add(amount)
		}
	}

	if !totalDebit.Equal(totalCredit) {
		return 0, types.ErrInvalidJournalEntry
	}

	entryID := k.GetNextJournalID(ctx)
	entry := types.JournalEntry{
		Id:            entryID,
		ReferenceType: referenceType,
		ReferenceId:   referenceID,
		Description:   description,
		Lines:         lines,
		CreatedAt:     ctx.BlockTime().UTC().Format("2006-01-02T15:04:05Z"),
		Creator:       creator,
	}
	k.SetJournalEntry(ctx, entry)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateJournal,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyJournalID, fmt.Sprintf("%d", entryID)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator),
		),
	)

	return entryID, nil
}

// PostSaleJournal posts the accounting journal for a completed sale (auto-posting).
//   - Debit Kas (1100) / Piutang (1300)  = grand total
//   - Credit Penjualan (4100)            = total sales
//   - Debit HPP (5100)                   = cost of goods sold
//   - Credit Persediaan (1200)           = cost of goods sold
func (k Keeper) PostSaleJournal(ctx sdk.Context, creator string, transactionID uint64, total, cost sdkmath.Int) (uint64, error) {
	cashAccount, found := k.FindAccountByCode(ctx, "1100")
	if !found {
		cashAccount = types.Account{Id: 2, Code: "1100", Name: "Kas", Type: "asset"}
	}
	salesAccount, found := k.FindAccountByCode(ctx, "4100")
	if !found {
		salesAccount = types.Account{Id: 10, Code: "4100", Name: "Penjualan", Type: "revenue"}
	}
	hppAccount, found := k.FindAccountByCode(ctx, "5100")
	if !found {
		hppAccount = types.Account{Id: 13, Code: "5100", Name: "HPP", Type: "expense"}
	}
	inventoryAccount, found := k.FindAccountByCode(ctx, "1200")
	if !found {
		inventoryAccount = types.Account{Id: 3, Code: "1200", Name: "Persediaan", Type: "asset"}
	}

	lines := []*types.JournalLine{
		{AccountId: cashAccount.Id, Debit: total.String()},
		{AccountId: salesAccount.Id, Credit: total.String()},
		{AccountId: hppAccount.Id, Debit: cost.String()},
		{AccountId: inventoryAccount.Id, Credit: cost.String()},
	}

	return k.createJournalEntry(ctx, creator, "sale", fmt.Sprintf("Penjualan #%d", transactionID), transactionID, lines)
}

// PostCancellationJournal posts a reversing journal entry for a cancelled sale.
func (k Keeper) PostCancellationJournal(ctx sdk.Context, creator string, transactionID uint64, total, cost sdkmath.Int) (uint64, error) {
	cashAccount, _ := k.FindAccountByCode(ctx, "1100")
	if cashAccount.Id == 0 {
		cashAccount = types.Account{Id: 2, Code: "1100", Name: "Kas", Type: "asset"}
	}
	salesAccount, _ := k.FindAccountByCode(ctx, "4100")
	if salesAccount.Id == 0 {
		salesAccount = types.Account{Id: 10, Code: "4100", Name: "Penjualan", Type: "revenue"}
	}
	hppAccount, _ := k.FindAccountByCode(ctx, "5100")
	if hppAccount.Id == 0 {
		hppAccount = types.Account{Id: 13, Code: "5100", Name: "HPP", Type: "expense"}
	}
	inventoryAccount, _ := k.FindAccountByCode(ctx, "1200")
	if inventoryAccount.Id == 0 {
		inventoryAccount = types.Account{Id: 3, Code: "1200", Name: "Persediaan", Type: "asset"}
	}

	lines := []*types.JournalLine{
		{AccountId: cashAccount.Id, Credit: total.String()},
		{AccountId: salesAccount.Id, Debit: total.String()},
		{AccountId: hppAccount.Id, Credit: cost.String()},
		{AccountId: inventoryAccount.Id, Debit: cost.String()},
	}

	return k.createJournalEntry(ctx, creator, "refund", fmt.Sprintf("Pembatalan Penjualan #%d", transactionID), transactionID, lines)
}

// PostPurchaseJournal posts the accounting journal for stock purchase / stock-in.
//   - Debit Persediaan (1200)  = cost
//   - Credit Kas (1100)        = cost
func (k Keeper) PostPurchaseJournal(ctx sdk.Context, creator string, productID uint64, quantity uint64, cost sdkmath.Int) (uint64, error) {
	inventoryAccount, _ := k.FindAccountByCode(ctx, "1200")
	if inventoryAccount.Id == 0 {
		inventoryAccount = types.Account{Id: 3, Code: "1200", Name: "Persediaan", Type: "asset"}
	}
	cashAccount, _ := k.FindAccountByCode(ctx, "1100")
	if cashAccount.Id == 0 {
		cashAccount = types.Account{Id: 2, Code: "1100", Name: "Kas", Type: "asset"}
	}

	lines := []*types.JournalLine{
		{AccountId: inventoryAccount.Id, Debit: cost.String()},
		{AccountId: cashAccount.Id, Credit: cost.String()},
	}

	return k.createJournalEntry(ctx, creator, "stock_in", fmt.Sprintf("Pembelian stok produk #%d (%d unit)", productID, quantity), productID, lines)
}

// PostAdjustmentJournal posts the accounting journal for a stock adjustment.
//   - Positive (in): Debit Persediaan (1200), Credit Pendapatan Lain (4200)
//   - Negative (out): Debit Beban Operasional (5200), Credit Persediaan (1200)
func (k Keeper) PostAdjustmentJournal(ctx sdk.Context, creator string, productID uint64, quantity int64, value sdkmath.Int) (uint64, error) {
	inventoryAccount, _ := k.FindAccountByCode(ctx, "1200")
	if inventoryAccount.Id == 0 {
		inventoryAccount = types.Account{Id: 3, Code: "1200", Name: "Persediaan", Type: "asset"}
	}

	var lines []*types.JournalLine
	if quantity >= 0 {
		otherAccount, _ := k.FindAccountByCode(ctx, "4200")
		if otherAccount.Id == 0 {
			otherAccount = types.Account{Id: 11, Code: "4200", Name: "Pendapatan Lain", Type: "revenue"}
		}
		lines = []*types.JournalLine{
			{AccountId: inventoryAccount.Id, Debit: value.String()},
			{AccountId: otherAccount.Id, Credit: value.String()},
		}
	} else {
		expenseAccount, _ := k.FindAccountByCode(ctx, "5200")
		if expenseAccount.Id == 0 {
			expenseAccount = types.Account{Id: 14, Code: "5200", Name: "Beban Operasional", Type: "expense"}
		}
		lines = []*types.JournalLine{
			{AccountId: expenseAccount.Id, Debit: value.String()},
			{AccountId: inventoryAccount.Id, Credit: value.String()},
		}
	}

	return k.createJournalEntry(ctx, creator, "adjustment", fmt.Sprintf("Penyesuaian stok produk #%d", productID), productID, lines)
}

// GetAccountBalances computes the net debit/credit/balance for every account.
func (k Keeper) GetAccountBalances(ctx sdk.Context) map[uint64]*types.AccountBalance {
	balances := make(map[uint64]*types.AccountBalance)

	// Initialize with all accounts
	for _, account := range k.GetAllAccounts(ctx) {
		balances[account.Id] = &types.AccountBalance{
			AccountId: account.Id,
			Code:      account.Code,
			Name:      account.Name,
			Type:      account.Type,
			Debit:     "0",
			Credit:    "0",
			Balance:   "0",
		}
	}

	// Sum all journal lines
	for _, entry := range k.GetAllJournalEntries(ctx) {
		for _, line := range entry.Lines {
			bal, exists := balances[line.AccountId]
			if !exists {
				bal = &types.AccountBalance{AccountId: line.AccountId, Debit: "0", Credit: "0", Balance: "0"}
				balances[line.AccountId] = bal
			}
			balDebit, _ := sdkmath.NewIntFromString(bal.Debit)
			balCredit, _ := sdkmath.NewIntFromString(bal.Credit)
			lineDebit, _ := sdkmath.NewIntFromString(line.Debit)
			lineCredit, _ := sdkmath.NewIntFromString(line.Credit)
			bal.Debit = balDebit.Add(lineDebit).String()
			bal.Credit = balCredit.Add(lineCredit).String()
			bal.Balance = balDebit.Add(lineDebit).Sub(balCredit.Add(lineCredit)).String()
		}
	}

	return balances
}

// GetTrialBalance returns the trial balance for all accounts.
func (k Keeper) GetTrialBalance(ctx sdk.Context) (*types.QueryTrialBalanceResponse, error) {
	balances := k.GetAccountBalances(ctx)
	accounts := make([]*types.AccountBalance, 0, len(balances))
	totalDebit := sdkmath.ZeroInt()
	totalCredit := sdkmath.ZeroInt()

	for _, bal := range balances {
		accounts = append(accounts, bal)
		balDebit, _ := sdkmath.NewIntFromString(bal.Debit)
		balCredit, _ := sdkmath.NewIntFromString(bal.Credit)
		totalDebit = totalDebit.Add(balDebit)
		totalCredit = totalCredit.Add(balCredit)
	}

	return &types.QueryTrialBalanceResponse{
		Accounts:    accounts,
		TotalDebit:  totalDebit.String(),
		TotalCredit: totalCredit.String(),
	}, nil
}

// GetIncomeStatement returns the income statement (laba rugi).
func (k Keeper) GetIncomeStatement(ctx sdk.Context) (*types.QueryIncomeStatementResponse, error) {
	balances := k.GetAccountBalances(ctx)

	revenues := make([]*types.AccountBalance, 0)
	expenses := make([]*types.AccountBalance, 0)
	totalRevenue := sdkmath.ZeroInt()
	totalExpense := sdkmath.ZeroInt()

	for _, bal := range balances {
		// Revenue accounts have credit balances; expense accounts have debit balances
		balCredit, _ := sdkmath.NewIntFromString(bal.Credit)
		balDebit, _ := sdkmath.NewIntFromString(bal.Debit)
		switch bal.Type {
		case "revenue":
			revenues = append(revenues, bal)
			totalRevenue = totalRevenue.Add(balCredit)
		case "expense":
			expenses = append(expenses, bal)
			totalExpense = totalExpense.Add(balDebit)
		}
	}

	return &types.QueryIncomeStatementResponse{
		Revenues:     revenues,
		Expenses:     expenses,
		TotalRevenue: totalRevenue.String(),
		TotalExpense: totalExpense.String(),
		NetIncome:    totalRevenue.Sub(totalExpense).String(),
	}, nil
}

// GetBalanceSheet returns the balance sheet (neraca).
func (k Keeper) GetBalanceSheet(ctx sdk.Context) (*types.QueryBalanceSheetResponse, error) {
	balances := k.GetAccountBalances(ctx)

	assets := make([]*types.AccountBalance, 0)
	liabilities := make([]*types.AccountBalance, 0)
	equities := make([]*types.AccountBalance, 0)
	totalAssets := sdkmath.ZeroInt()
	totalLiabilities := sdkmath.ZeroInt()
	totalEquity := sdkmath.ZeroInt()

	for _, bal := range balances {
		balBalance, _ := sdkmath.NewIntFromString(bal.Balance)
		switch bal.Type {
		case "asset":
			assets = append(assets, bal)
			totalAssets = totalAssets.Add(balBalance)
		case "liability":
			liabilities = append(liabilities, bal)
			totalLiabilities = totalLiabilities.Add(balBalance)
		case "equity":
			equities = append(equities, bal)
			totalEquity = totalEquity.Add(balBalance)
		}
	}

	return &types.QueryBalanceSheetResponse{
		Assets:            assets,
		Liabilities:       liabilities,
		Equities:          equities,
		TotalAssets:       totalAssets.String(),
		TotalLiabilities:  totalLiabilities.String(),
		TotalEquity:       totalEquity.String(),
	}, nil
}

// GetLedger returns the general ledger for an account.
func (k Keeper) GetLedger(ctx sdk.Context, accountID uint64) (*types.QueryLedgerResponse, error) {
	account, found := k.GetAccount(ctx, accountID)
	if !found {
		return nil, types.ErrAccountNotFound
	}

	lines := make([]*types.LedgerLine, 0)
	running := sdkmath.ZeroInt()

	for _, entry := range k.GetAllJournalEntries(ctx) {
		for _, line := range entry.Lines {
			if line.AccountId != accountID {
				continue
			}
			lineDebit, _ := sdkmath.NewIntFromString(line.Debit)
			lineCredit, _ := sdkmath.NewIntFromString(line.Credit)
			running = running.Add(lineDebit).Sub(lineCredit)

			lines = append(lines, &types.LedgerLine{
				JournalEntryId: entry.Id,
				ReferenceType:  entry.ReferenceType,
				ReferenceId:    entry.ReferenceId,
				Description:    entry.Description,
				Debit:          line.Debit,
				Credit:         line.Credit,
				Balance:        running.String(),
				CreatedAt:      entry.CreatedAt,
			})
		}
	}

	return &types.QueryLedgerResponse{
		Account:       &account,
		Lines:         lines,
		EndingBalance: running.String(),
	}, nil
}

// GetPriceLevelReport returns the price level report for all active products.
func (k Keeper) GetPriceLevelReport(ctx sdk.Context) *types.QueryPriceLevelReportResponse {
	items := make([]*types.PriceLevelReportItem, 0)
	for _, product := range k.GetActiveProducts(ctx) {
		baseUnit := ""
		if product.BaseUnitId != 0 {
			if unit, found := k.GetUnit(ctx, product.BaseUnitId); found {
				baseUnit = unit.Symbol
			}
		}
		items = append(items, &types.PriceLevelReportItem{
			ProductId:   product.Id,
			ProductName: product.Name,
			Sku:         product.Sku,
			BasePrice:   product.Price,
			BaseUnit:    baseUnit,
			PriceLevels: product.PriceLevels,
		})
	}
	return &types.QueryPriceLevelReportResponse{Items: items}
}
