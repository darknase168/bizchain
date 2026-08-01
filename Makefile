#!/usr/bin/make

# BizChain Makefile
# For business blockchain with zero fees

BINARY := bizchaind
PROJECT_NAME := BizChain

# Go version
GO_VERSION := $(shell go version 2>/dev/null | sed 's/.*go\([0-9]\.[0-9]*\).*/\1/')

# Build flags
BUILD_DIR := ./build
BUILD_FLAGS := -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=bizchain -X github.com/cosmos/cosmos-sdk/version.AppName=bizchaind -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)"

# Default target
.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(PROJECT_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/bizchaind
	@echo "Binary built at $(BUILD_DIR)/$(BINARY)"

.PHONY: install
install:
	@echo "Installing $(BINARY)..."
	go install $(BUILD_FLAGS) ./cmd/bizchaind
	@echo "$(BINARY) installed to $$GOPATH/bin"

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Done!"

.PHONY: init
init: build
	@echo "Initializing BizChain..."
	$(BUILD_DIR)/$(BINARY) init $(MONIKER) --chain-id bizchain-1
	@echo "BizChain initialized!"

.PHONY: init-zero-fees
init-zero-fees: build
	@echo "Initializing BizChain with zero fees..."
	$(BUILD_DIR)/$(BINARY) init $(MONIKER) --chain-id bizchain-1
	# Set minimum gas prices to zero for zero-fee transactions
	sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0upoint"/g' ~/.bizchain/config/app.toml
	@echo "Zero-fee BizChain initialized!"

.PHONY: start
start: build
	@echo "Starting BizChain node..."
	$(BUILD_DIR)/$(BINARY) start

.PHONY: reset
reset:
	@echo "Resetting BizChain..."
	$(BUILD_DIR)/$(BINARY) unsafe-reset-all
	@echo "Done!"

.PHONY: add-account
add-account:
	@echo "Adding genesis account..."
	$(BUILD_DIR)/$(BINARY) add-genesis-account $(ADDRESS) $(AMOUNT)upoint
	@echo "Account added!"

.PHONY: gentx
gentx:
	@echo "Generating genesis transaction..."
	$(BUILD_DIR)/$(BINARY) gentx $(NAME) $(AMOUNT)upoint --chain-id bizchain-1
	@echo "Gentx generated!"

.PHONY: collect-gentx
collect-gentx:
	$(BUILD_DIR)/$(BINARY) collect-gentxs
	@echo "Gentxs collected!"

# POS module commands
.PHONY: create-product
create-product:
	$(BUILD_DIR)/$(BINARY) tx pos create-product $(NAME) $(PRICE) $(SKU) $(CATEGORY) --from $(FROM) --chain-id bizchain-1 --fees 0upoint --yes

.PHONY: create-transaction
create-transaction:
	$(BUILD_DIR)/$(BINARY) tx pos create-transaction $(CUSTOMER) $(PRODUCT_IDS) $(QUANTITIES) $(PRICES) --from $(FROM) --chain-id bizchain-1 --fees 0upoint --yes

.PHONY: add-stock
add-stock:
	$(BUILD_DIR)/$(BINARY) tx pos add-stock $(PRODUCT_ID) $(QUANTITY) --from $(FROM) --chain-id bizchain-1 --fees 0upoint --yes

.PHONY: list-products
list-products:
	$(BUILD_DIR)/$(BINARY) query pos list-product

.PHONY: show-product
show-product:
	$(BUILD_DIR)/$(BINARY) query pos show-product $(ID)

# Wallet commands
.PHONY: create-wallet
create-wallet:
	$(BUILD_DIR)/$(BINARY) keys add $(NAME)

.PHONY: list-wallets
list-wallets:
	$(BUILD_DIR)/$(BINARY) keys list

.PHONY: show-balance
show-balance:
	$(BUILD_DIR)/$(BINARY) query bank balances $(ADDRESS)

# Testing
.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -v
	@echo "Tests completed!"

.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...
	@echo "Linting completed!"

# Help
.PHONY: help
help:
	@echo "$(PROJECT_NAME) - Feeless Retail/Education Blockchain"
	@echo ""
	@echo "Usage:"
	@echo "  make build              Build the binary"
	@echo "  make install            Install the binary"
	@echo "  make init               Initialize the chain"
	@echo "  make init-zero-fees     Initialize with zero fees"
	@echo "  make start              Start the node"
	@echo "  make reset              Reset chain data"
	@echo ""
	@echo "POS Commands:"
	@echo "  make create-product      Create a product (NAME, PRICE, SKU, CATEGORY, FROM)"
	@echo "  make create-transaction  Create POS sale (CUSTOMER, PRODUCT_IDS, QUANTITIES, PRICES, FROM)"
	@echo "  make add-stock           Add stock (PRODUCT_ID, QUANTITY, FROM)"
	@echo "  make list-products       List all products"
	@echo "  make show-product        Show product (ID)"
	@echo ""
	@echo "Wallet Commands:"
	@echo "  make create-wallet       Create a wallet (NAME)"
	@echo "  make list-wallets        List wallets"
	@echo "  make show-balance        Show balance (ADDRESS)"
	@echo ""
	@echo "Environment Variables:"
	@echo "  MONIKER=node1            Node moniker"
	@echo "  ADDRESS=point1...        Account address"
	@echo "  AMOUNT=1000000000        Token amount"
	@echo "  FROM=mykey               Key name for transactions"
	@echo "  NAME=product1            Product name"
	@echo "  PRICE=1000               Product price"
	@echo "  SKU=SKU001               Product SKU"
