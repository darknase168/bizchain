package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateProduct{}
	_ sdk.Msg = &MsgUpdateProduct{}
	_ sdk.Msg = &MsgDeleteProduct{}
	_ sdk.Msg = &MsgCreateTransaction{}
	_ sdk.Msg = &MsgCancelTransaction{}
	_ sdk.Msg = &MsgAddStock{}
	_ sdk.Msg = &MsgAdjustStock{}
	_ sdk.Msg = &MsgCreateUnit{}
	_ sdk.Msg = &MsgUpdateUnit{}
	_ sdk.Msg = &MsgDeleteUnit{}
	_ sdk.Msg = &MsgCreateAccount{}
	_ sdk.Msg = &MsgCreateJournalEntry{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func validateNonNegativeAmount(name, value string) error {
	if value == "" {
		return nil
	}
	amount, ok := sdkmath.NewIntFromString(value)
	if !ok || amount.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidPrice, "%s must be a valid non-negative integer", name)
	}
	return nil
}

// ---------- MsgCreateProduct ----------

func NewMsgCreateProduct(
	creator, name, description, price, costPrice, sku, category, imageURL, branchID string,
	baseUnitID, initialStock, minStock uint64, barcode string,
	priceLevels []*PriceLevel, isBundle bool, components []*BundleComponent,
) *MsgCreateProduct {
	return &MsgCreateProduct{
		Creator:      creator,
		Name:         name,
		Description:  description,
		Price:        price,
		CostPrice:    costPrice,
		Sku:          sku,
		Category:     category,
		ImageUrl:     imageURL,
		BaseUnitId:   baseUnitID,
		InitialStock: initialStock,
		MinStock:     minStock,
		Barcode:      barcode,
		PriceLevels:  priceLevels,
		IsBundle:     isBundle,
		Components:   components,
		BranchId:     branchID,
	}
}

func (msg *MsgCreateProduct) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "name cannot be empty")
	}
	if len(msg.Price) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "price cannot be empty")
	}
	price, ok := sdkmath.NewIntFromString(msg.Price)
	if !ok || !price.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidPrice, "price must be positive")
	}
	if err := validateNonNegativeAmount("cost_price", msg.CostPrice); err != nil {
		return err
	}
	for _, pl := range msg.PriceLevels {
		if pl.Level == "" || len(pl.Price) == 0 {
			return sdkerrors.Wrap(ErrInvalidPriceLevel, "price level requires a level name and price")
		}
		p, ok := sdkmath.NewIntFromString(pl.Price)
		if !ok || !p.IsPositive() {
			return sdkerrors.Wrap(ErrInvalidPriceLevel, "price level must be positive")
		}
	}
	if msg.IsBundle && len(msg.Components) == 0 {
		return sdkerrors.Wrap(ErrInvalidBundle, "bundle must have at least one component")
	}
	for _, c := range msg.Components {
		if c.ProductId == 0 || c.Quantity == 0 {
			return sdkerrors.Wrap(ErrInvalidBundle, "component requires product id and quantity")
		}
	}
	if len(msg.BranchId) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "branch_id is required")
	}
	return nil
}

func (msg *MsgCreateProduct) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateProduct ----------

func NewMsgUpdateProduct(
	creator string, id uint64, name, description, price, costPrice, sku, category, imageURL string,
	baseUnitID, minStock uint64, barcode string,
	priceLevels []*PriceLevel, isBundle bool, components []*BundleComponent,
) *MsgUpdateProduct {
	return &MsgUpdateProduct{
		Creator:     creator,
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
		CostPrice:   costPrice,
		Sku:         sku,
		Category:    category,
		ImageUrl:    imageURL,
		BaseUnitId:  baseUnitID,
		MinStock:    minStock,
		Barcode:     barcode,
		PriceLevels: priceLevels,
		IsBundle:    isBundle,
		Components:  components,
	}
}

func (msg *MsgUpdateProduct) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidProductID, "product ID cannot be zero")
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "name cannot be empty")
	}
	return nil
}

func (msg *MsgUpdateProduct) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgDeleteProduct ----------

func NewMsgDeleteProduct(creator string, id uint64) *MsgDeleteProduct {
	return &MsgDeleteProduct{Creator: creator, Id: id}
}

func (msg *MsgDeleteProduct) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidProductID, "product ID cannot be zero")
	}
	return nil
}

func (msg *MsgDeleteProduct) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgCreateTransaction ----------

func NewMsgCreateTransaction(
	creator, customerAddress, note, paymentMethod, discount, tax, branchID string, items []*Item,
) *MsgCreateTransaction {
	return &MsgCreateTransaction{
		Creator:         creator,
		CustomerAddress: customerAddress,
		Items:           items,
		Discount:        discount,
		Tax:             tax,
		PaymentMethod:   paymentMethod,
		Note:            note,
		BranchId:        branchID,
	}
}

func (msg *MsgCreateTransaction) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Items) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "items cannot be empty")
	}
	for _, item := range msg.Items {
		if item.ProductId == 0 {
			return sdkerrors.Wrap(ErrInvalidProductID, "product ID cannot be zero")
		}
		if item.Quantity == 0 {
			return sdkerrors.Wrap(ErrInvalidQuantity, "quantity cannot be zero")
		}
		if err := validateNonNegativeAmount("price", item.Price); err != nil {
			return err
		}
	}
	if err := validateNonNegativeAmount("discount", msg.Discount); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("tax", msg.Tax); err != nil {
		return err
	}
	if len(msg.BranchId) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "branch_id is required")
	}
	return nil
}

func (msg *MsgCreateTransaction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgCancelTransaction ----------

func NewMsgCancelTransaction(creator string, id uint64) *MsgCancelTransaction {
	return &MsgCancelTransaction{Creator: creator, Id: id}
}

func (msg *MsgCancelTransaction) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidProductID, "transaction ID cannot be zero")
	}
	return nil
}

func (msg *MsgCancelTransaction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgAddStock ----------

func NewMsgAddStock(creator string, productID, quantity, unitID uint64, costPrice string) *MsgAddStock {
	return &MsgAddStock{
		Creator:   creator,
		ProductId: productID,
		Quantity:  quantity,
		UnitId:    unitID,
		CostPrice: costPrice,
	}
}

func (msg *MsgAddStock) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.ProductId == 0 {
		return sdkerrors.Wrap(ErrInvalidProductID, "product ID cannot be zero")
	}
	if msg.Quantity == 0 {
		return sdkerrors.Wrap(ErrInvalidQuantity, "quantity cannot be zero")
	}
	if err := validateNonNegativeAmount("cost_price", msg.CostPrice); err != nil {
		return err
	}
	return nil
}

func (msg *MsgAddStock) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgAdjustStock ----------

func NewMsgAdjustStock(creator string, productID uint64, quantity int64, reason string) *MsgAdjustStock {
	return &MsgAdjustStock{
		Creator:   creator,
		ProductId: productID,
		Quantity:  quantity,
		Reason:    reason,
	}
}

func (msg *MsgAdjustStock) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.ProductId == 0 {
		return sdkerrors.Wrap(ErrInvalidProductID, "product ID cannot be zero")
	}
	if msg.Quantity == 0 {
		return sdkerrors.Wrap(ErrInvalidQuantity, "quantity cannot be zero")
	}
	return nil
}

func (msg *MsgAdjustStock) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgCreateUnit ----------

func NewMsgCreateUnit(creator, name, symbol string, conversionFactor uint64, isBase bool) *MsgCreateUnit {
	return &MsgCreateUnit{
		Creator:          creator,
		Name:             name,
		Symbol:           symbol,
		ConversionFactor: conversionFactor,
		IsBase:           isBase,
	}
}

func (msg *MsgCreateUnit) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "unit name cannot be empty")
	}
	if msg.ConversionFactor == 0 {
		return sdkerrors.Wrap(ErrInvalidQuantity, "conversion factor must be positive")
	}
	return nil
}

func (msg *MsgCreateUnit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateUnit ----------

func NewMsgUpdateUnit(creator string, id uint64, name, symbol string, conversionFactor uint64, isBase bool) *MsgUpdateUnit {
	return &MsgUpdateUnit{
		Creator:          creator,
		Id:               id,
		Name:             name,
		Symbol:           symbol,
		ConversionFactor: conversionFactor,
		IsBase:           isBase,
	}
}

func (msg *MsgUpdateUnit) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidUnitID, "unit ID cannot be zero")
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "unit name cannot be empty")
	}
	if msg.ConversionFactor == 0 {
		return sdkerrors.Wrap(ErrInvalidQuantity, "conversion factor must be positive")
	}
	return nil
}

func (msg *MsgUpdateUnit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgDeleteUnit ----------

func NewMsgDeleteUnit(creator string, id uint64) *MsgDeleteUnit {
	return &MsgDeleteUnit{Creator: creator, Id: id}
}

func (msg *MsgDeleteUnit) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidUnitID, "unit ID cannot be zero")
	}
	return nil
}

func (msg *MsgDeleteUnit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgCreateAccount ----------

func NewMsgCreateAccount(creator, code, name, accType, description string) *MsgCreateAccount {
	return &MsgCreateAccount{
		Creator:     creator,
		Code:        code,
		Name:        name,
		Type:        accType,
		Description: description,
	}
}

func (msg *MsgCreateAccount) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Code) == 0 || len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "account code and name cannot be empty")
	}
	switch msg.Type {
	case "asset", "liability", "equity", "revenue", "expense":
	default:
		return sdkerrors.Wrap(ErrInvalidAccountType, "type must be asset, liability, equity, revenue or expense")
	}
	return nil
}

func (msg *MsgCreateAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgCreateJournalEntry ----------

func NewMsgCreateJournalEntry(creator, referenceType, description string, referenceID uint64, lines []*JournalLine) *MsgCreateJournalEntry {
	return &MsgCreateJournalEntry{
		Creator:       creator,
		ReferenceType: referenceType,
		ReferenceId:   referenceID,
		Description:   description,
		Lines:         lines,
	}
}

func (msg *MsgCreateJournalEntry) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Lines) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "journal entry requires at least one line")
	}
	totalDebit := sdkmath.ZeroInt()
	totalCredit := sdkmath.ZeroInt()
	for _, line := range msg.Lines {
		if line.AccountId == 0 {
			return sdkerrors.Wrap(ErrInvalidAccountType, "account ID cannot be zero")
		}
		amount, ok := sdkmath.NewIntFromString(line.Debit)
		if !ok || amount.IsNegative() {
			return sdkerrors.Wrap(ErrInvalidPrice, "debit must be a valid non-negative integer")
		}
		totalDebit = totalDebit.Add(amount)
		amount, ok = sdkmath.NewIntFromString(line.Credit)
		if !ok || amount.IsNegative() {
			return sdkerrors.Wrap(ErrInvalidPrice, "credit must be a valid non-negative integer")
		}
		totalCredit = totalCredit.Add(amount)
	}
	if !totalDebit.Equal(totalCredit) {
		return sdkerrors.Wrap(ErrInvalidJournalEntry, "total debit must equal total credit")
	}
	return nil
}

func (msg *MsgCreateJournalEntry) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- Genesis Types ----------

const (
	ModuleVersion = 1
)
