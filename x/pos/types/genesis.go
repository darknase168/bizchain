package types

import (
	"encoding/json"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	accounts := DefaultChartOfAccounts()
	units := DefaultUnits()
	return &GenesisState{
		ProductList:     []*Product{},
		TransactionList: []*Transaction{},
		UnitList:        units,
		AccountList:     accounts,
		JournalList:     []*JournalEntry{},
		ProductCount:    0,
		TransactionCount: 0,
		UnitCount:       uint64(len(units)),
		AccountCount:    uint64(len(accounts)),
		JournalCount:    0,
		Params:          DefaultParams(),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() *Params {
	return &Params{
		ZeroFeeEnabled: true,
	}
}

// DefaultUnits returns the default units of measure (multi satuan)
func DefaultUnits() []*Unit {
	return []*Unit{
		{Id: 1, Name: "Pcs", Symbol: "pcs", ConversionFactor: 1, IsBase: true},
		{Id: 2, Name: "Dus", Symbol: "dus", ConversionFactor: 12, IsBase: false},
		{Id: 3, Name: "Karton", Symbol: "ktn", ConversionFactor: 48, IsBase: false},
		{Id: 4, Name: "Box", Symbol: "box", ConversionFactor: 24, IsBase: false},
	}
}

// DefaultChartOfAccounts returns the default chart of accounts (akuntansi)
func DefaultChartOfAccounts() []*Account {
	now := "2026-01-01T00:00:00Z"
	return []*Account{
		{Id: 1, Code: "1000", Name: "Aset", Type: "asset", Description: "Kelompok aset", CreatedAt: now},
		{Id: 2, Code: "1100", Name: "Kas", Type: "asset", Description: "Kas & bank", CreatedAt: now},
		{Id: 3, Code: "1200", Name: "Persediaan Barang", Type: "asset", Description: "Nilai persediaan", CreatedAt: now},
		{Id: 4, Code: "1300", Name: "Piutang Usaha", Type: "asset", Description: "Piutang pelanggan", CreatedAt: now},
		{Id: 5, Code: "2000", Name: "Kewajiban", Type: "liability", Description: "Kelompok kewajiban", CreatedAt: now},
		{Id: 6, Code: "2100", Name: "Hutang Usaha", Type: "liability", Description: "Hutang ke supplier", CreatedAt: now},
		{Id: 7, Code: "3000", Name: "Ekuitas", Type: "equity", Description: "Kelompok ekuitas", CreatedAt: now},
		{Id: 8, Code: "3100", Name: "Modal", Type: "equity", Description: "Modal pemilik", CreatedAt: now},
		{Id: 9, Code: "4000", Name: "Pendapatan", Type: "revenue", Description: "Kelompok pendapatan", CreatedAt: now},
		{Id: 10, Code: "4100", Name: "Penjualan", Type: "revenue", Description: "Pendapatan penjualan", CreatedAt: now},
		{Id: 11, Code: "4200", Name: "Pendapatan Lain", Type: "revenue", Description: "Pendapatan lain-lain", CreatedAt: now},
		{Id: 12, Code: "5000", Name: "Beban", Type: "expense", Description: "Kelompok beban", CreatedAt: now},
		{Id: 13, Code: "5100", Name: "Harga Pokok Penjualan", Type: "expense", Description: "HPP", CreatedAt: now},
		{Id: 14, Code: "5200", Name: "Beban Operasional", Type: "expense", Description: "Beban operasional", CreatedAt: now},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure
func (gs GenesisState) Validate() error {
	for _, product := range gs.ProductList {
		if product.Id == 0 {
			return ErrInvalidProductID
		}
	}
	return nil
}

// Validate validates the Params
func (p Params) Validate() error {
	return nil
}

// NewGenesisState creates a new GenesisState instance
func NewGenesisState(
	productList []*Product, transactionList []*Transaction, unitList []*Unit,
	accountList []*Account, journalList []*JournalEntry,
	productCount, transactionCount, unitCount, accountCount, journalCount uint64, params *Params,
) *GenesisState {
	return &GenesisState{
		ProductList:      productList,
		TransactionList:  transactionList,
		UnitList:         unitList,
		AccountList:      accountList,
		JournalList:      journalList,
		ProductCount:     productCount,
		TransactionCount: transactionCount,
		UnitCount:        unitCount,
		AccountCount:     accountCount,
		JournalCount:     journalCount,
		Params:           params,
	}
}

// GetGenesisStateFromAppState returns the GenesisState from the app state
func GetGenesisStateFromAppState(rawState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if rawState[ModuleName] != nil {
		ModuleCdc.MustUnmarshalJSON(rawState[ModuleName], &genesisState)
	}
	return genesisState
}
