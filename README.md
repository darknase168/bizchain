# BizChain 🚀

**Feeless Blockchain for Retail, Minimarket & Education Businesses**

BizChain is a custom Cosmos SDK blockchain built with CometBFT, designed specifically for retail/minimarket and education businesses. Key features:

- ✅ **Zero Transaction Fees** - Send transactions without any gas costs
- ✅ **Custom POS Module** - Point of Sales on-chain (products, transactions, inventory)
- ✅ **Custom Wallet** - Full-featured web wallet with CosmJS
- ✅ **Fast Finality** - Powered by CometBFT consensus
- ✅ **Bech32 Address** - `point1...` prefix addresses

## Architecture

```
bizchaind                    # Main blockchain daemon
├── app/                     # Cosmos SDK application setup
│   ├── app.go              # Main app initialization
│   ├── config.go           # App configuration & constants
│   └── encoding.go         # Encoding configuration
├── cmd/bizchaind/          # CLI entry point
│   ├── main.go             # Main function
│   └── cmd/root.go         # Root CLI command
├── x/pos/                  # Custom POS Module
│   ├── keeper/             # State management
│   ├── types/              # Types, messages, errors
│   ├── client/cli/         # CLI commands
│   ├── module.go           # Module registration
│   └── genesis.go          # Genesis handling
├── proto/                  # Protocol Buffer definitions
├── wallet/                 # Custom Wallet Frontend
│   ├── src/                # React + CosmJS
│   └── package.json        # Dependencies
└── Makefile                # Build & manage commands
```

## Prerequisites

- [Go](https://golang.org/dl/) 1.22+ (recommended: 1.26+)
- [Node.js](https://nodejs.org/) 18+ (for wallet)
- [npm](https://www.npmjs.com/) or [yarn](https://yarnpkg.com/)

## Quick Start - Running the Blockchain

### 1. Build the binary

```bash
# Build from source
make build

# Or install to GOPATH/bin
make install
```

### 2. Initialize the chain

```bash
# Initialize with default settings
bizchaind init my-node --chain-id bizchain-1

# Or use the Makefile
make init MONIKER=my-node
```

### 3. Configure Zero Fees (Already set by default)

The chain is configured with `minimum-gas-prices = "0upoint"` for zero-fee transactions.
This is set in `app.toml` automatically.

### 4. Create a wallet

```bash
# Create a new wallet (address: point1...)
bizchaind keys add my-wallet

# List wallets
bizchaind keys list

# Show wallet balance
bizchaind query bank balances point1...
```

### 5. Add genesis account & start

```bash
# Add genesis account with tokens
bizchaind add-genesis-account point1... 1000000000upoint

# Generate genesis transaction
bizchaind gentx my-wallet 500000000upoint --chain-id bizchain-1

# Collect genesis transactions
bizchaind collect-gentxs

# Start the chain
bizchaind start
```

## POS Module Commands

### Products

```bash
# Create a product (0 fee!)
bizchaind tx pos create-product "Indomie Goreng" 3500000 IDM-001 Makanan \
  --from my-wallet --chain-id bizchain-1 --fees 0upoint --yes

# List all products
bizchaind query pos list-product

# Show product details
bizchaind query pos show-product 1

# Add stock
bizchaind tx pos add-stock 1 100 \
  --from my-wallet --chain-id bizchain-1 --fees 0upoint --yes
```

### Transactions (Sales)

```bash
# Create a POS sale (sell 2x product 1 at 3500 POINT each)
bizchaind tx pos create-transaction "customer-address" "1" "2" "3500000" \
  --from my-wallet --chain-id bizchain-1 --fees 0upoint --yes

# List all transactions
bizchaind query pos list-transaction
```

## Wallet Frontend

### Development

```bash
cd wallet
npm install
npm run dev
```

The wallet will be available at `http://localhost:5173`

### Features

- **Create Wallet** - Generate new 24-word mnemonic wallet
- **Import Wallet** - Import from existing mnemonic
- **Check Balance** - View POINT token balance
- **Send Tokens** - Send POINT tokens to other addresses
- **POS Dashboard** - Full Point of Sales interface
- **Product Management** - Create and manage products
- **Transaction History** - View all transactions
- **Zero Fee** - All transactions are free

## Tokenomics

| Token | Denom | Decimals | Address Prefix |
|-------|-------|----------|----------------|
| POINT | upoint | 6 | point1... |

- **Bech32 Prefix:** `point`
- **Coin Type:** 118 (Cosmos standard)
- **BIP44 Purpose:** 44'

## Makefile Commands

```bash
make build              # Build the blockchain binary
make install            # Install binary to GOPATH
make init               # Initialize chain
make start              # Start the node
make create-wallet      # Create wallet (NAME=my-wallet)
make create-product     # Create product
make create-transaction # Create POS sale
make list-products      # List all products
make add-stock          # Add stock to product
```

## Project Structure Details

### Custom POS Module (`x/pos/`)

The POS module provides on-chain functionality for:

- **Product Management:** Create, update, and manage product inventory
- **POS Transactions:** Record sales with multiple items
- **Stock Management:** Track and update product stock levels
- **Zero Fee Configuration:** All transactions are feeless

### Messages

| Message | Description |
|---------|-------------|
| `MsgCreateProduct` | Create a new product listing |
| `MsgUpdateProduct` | Update product details |
| `MsgDeleteProduct` | Deactivate a product |
| `MsgCreateTransaction` | Record a POS sale |
| `MsgAddStock` | Add inventory to a product |

### Queries

| Query | Description |
|-------|-------------|
| `Product` | Get product by ID |
| `ProductAll` | List all products |
| `Transaction` | Get transaction by ID |
| `TransactionAll` | List all transactions |

## Security

- **Cosmos SDK**: Industry-standard blockchain framework
- **CometBFT**: Byzantine Fault Tolerant consensus
- **Bech32 Addresses**: Human-readable, checksummed addresses
- **BIP-39 Mnemonics**: 24-word recovery phrases

## License

MIT

## Support

Built for the retail & education sector. Zero fees, maximum efficiency. 🎯
