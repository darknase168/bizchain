package types

import (
	"time"
)

// NewProduct creates a new Product instance (generated type)
func NewProduct(
	id uint64, name, description, price, costPrice, sku, category, imageURL, owner, branchID string,
	stock, baseUnitID, minStock uint64, barcode string, priceLevels []*PriceLevel,
	isBundle bool, components []*BundleComponent,
) Product {
	now := time.Now().UTC().Format(time.RFC3339)
	return Product{
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
		CostPrice:   costPrice,
		Sku:         sku,
		Category:    category,
		ImageUrl:    imageURL,
		Stock:       stock,
		Owner:       owner,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
		BaseUnitId:  baseUnitID,
		MinStock:    minStock,
		Barcode:     barcode,
		PriceLevels: priceLevels,
		IsBundle:    isBundle,
		Components:  components,
		BranchId:    branchID,
	}
}

// NewTransaction creates a new Transaction instance (generated type)
func NewTransaction(
	id uint64, seller, customerAddress, branchID string, items []*Item, total, discount, tax, grandTotal, paymentMethod, note string,
) Transaction {
	now := time.Now().UTC().Format(time.RFC3339)
	return Transaction{
		Id:              id,
		Seller:          seller,
		CustomerAddress: customerAddress,
		Items:           items,
		Total:           total,
		Discount:        discount,
		Tax:             tax,
		GrandTotal:      grandTotal,
		PaymentMethod:   paymentMethod,
		Status:          "completed",
		Note:            note,
		BranchId:        branchID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
