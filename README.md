# 🕌 BizChain - Blockchain untuk Sistem Agen Haji & Umroh

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Cosmos SDK](https://img.shields.io/badge/Cosmos%20SDK-v0.50-brightgreen.svg)](https://cosmos.network)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**BizChain** adalah blockchain khusus untuk sistem manajemen Agen Haji & Umroh, dibangun menggunakan Cosmos SDK. Mendukung fitur-fitur lengkap dari manajemen jamaah, paket, pembayaran escrow, asuransi, hingga sistem POS/Retail dengan **Multi-Cabang Support**.

## ✨ Fitur Utama

### 🕋 Haji & Umroh Management
- **Jamaah Management** - Pendaftaran dan verifikasi data jamaah dengan DID
- **Paket Management** - Pengelolaan paket haji/umroh dengan review system
- **Pembayaran Escrow** - Sistem pembayaran bertahap dengan smart escrow
- **Visa Processing** - Tracking status visa application
- **Hotel & Flight** - Manajemen booking hotel dan tiket pesawat NFT
- **Keberangkatan** - Tracking journey stages (manasik, departure, rituals, return)

### 💼 Ekosistem Agen
- **Multi-Level Agent System** - Pusat, Cabang, Sub-agen dengan commission tracking
- **Referral & Reward** - Sistem referral dengan loyalty points
- **DAO Governance** - Voting dan proposal untuk keputusan bersama
- **Insurance** - Asuransi perjalanan dengan blockchain-based claims
- **Oleh-Oleh Marketplace** - Pre-order marketplace untuk oleh-oleh haji/umroh

### 🏪 POS/Retail System
- **Multi-Satuan** - Support multiple units of measure (pcs, dus, karton)
- **Harga Level** - Tiered pricing (ecer, grosir, member)
- **Barang Gabungan** - Bundle/composite products
- **📍 Multi-Cabang** - Multi-branch inventory & transaction management ✨
- **Akuntansi Terintegrasi** - Chart of accounts, journal entries, financial reports

### 🔗 Blockchain Features
- **Zero Fees** - Transaksi tanpa biaya gas
- **IBC Ready** - Inter-Blockchain Communication support
- **Audit Trail** - Semua transaksi tercatat dan dapat diaudit
- **gRPC & REST API** - Modern API interface

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21+ 
- **Node.js** 18+ dan npm (untuk web wallet)
- **Git**

### Installation

```bash
# Clone repository
git clone https://github.com/darknase168/bizchain.git
cd bizchain

# Build blockchain
go build -o build/bizchaind ./cmd/bizchaind

# Initialize chain
./build/bizchaind init mynode --chain-id bizchain-1

# Configure zero fees (optional, edit ~/.bizchain/config/app.toml)
# minimum-gas-prices = "0upoint"

# Create genesis account
./build/bizchaind keys add alice
./build/bizchaind add-genesis-account alice 1000000000000upoint --keyring-backend test
./build/bizchaind gentx alice 1000000upoint --chain-id bizchain-1 --keyring-backend test
./build/bizchaind collect-gentxs

# Start node
./build/bizchaind start
```

### Web Wallet

```bash
cd wallet
npm install
npm run dev
```

Buka browser: `http://localhost:5173`

## 📚 Dokumentasi Lengkap

- **[MULTI_CABANG.md](MULTI_CABANG.md)** - Spesifikasi fitur Multi-Cabang
- **[IMPLEMENTASI_MULTI_CABANG.md](IMPLEMENTASI_MULTI_CABANG.md)** - Detail implementasi Multi-Cabang
- **[JALANKAN_NODE.md](JALANKAN_NODE.md)** - Panduan lengkap menjalankan node (Windows & GitHub Codespaces)
- **[Fitur_Sistem_Agen_Haji_Umroh_Blockchain_CosmosSDK.md](Fitur_Sistem_Agen_Haji_Umroh_Blockchain_CosmosSDK.md)** - Spesifikasi lengkap sistem

## 🏗️ Arsitektur

```
bizchain/
├── app/                   # Application setup
├── cmd/bizchaind/        # CLI binary
├── proto/                # Protobuf definitions
│   └── bizchain/
│       ├── agen/         # Agent module
│       ├── jamaah/       # Pilgrim module
│       ├── paket/        # Package module
│       ├── pembayaran/   # Payment escrow
│       ├── visa/         # Visa processing
│       ├── hotel/        # Hotel management
│       ├── ticket/       # Flight tickets (NFT)
│       ├── pos/          # POS/Retail (Multi-Cabang)
│       └── ...
├── x/                    # Custom modules
│   ├── agen/
│   ├── jamaah/
│   ├── paket/
│   ├── pembayaran/
│   ├── pos/             # POS module with Multi-Cabang
│   └── ...
└── wallet/              # React web wallet
    └── src/
        ├── HajiUmrohDashboard.tsx
        ├── RetailDashboard.tsx    # Multi-Cabang UI
        ├── AgenEkosistemDashboard.tsx
        └── ...
```

## 🔌 API Endpoints

### REST API (Port 1317)

**POS/Retail (Multi-Cabang):**
```
GET /bizchain/pos/product                    # All products
GET /bizchain/pos/product_by_branch/{id}     # Products by branch ✨
GET /bizchain/pos/transaction_by_branch/{id} # Transactions by branch ✨
GET /bizchain/pos/unit                       # Units of measure
GET /bizchain/pos/account                    # Chart of accounts
GET /bizchain/pos/trial_balance              # Trial balance
GET /bizchain/pos/income_statement           # Income statement
GET /bizchain/pos/balance_sheet              # Balance sheet
```

**Haji & Umroh:**
```
GET /bizchain/jamaah/jamaah         # All pilgrims
GET /bizchain/paket/paket           # All packages
GET /bizchain/pembayaran/pembayaran # All payments
GET /bizchain/visa/visa             # All visas
GET /bizchain/hotel/hotel           # All hotels
GET /bizchain/ticket/ticket         # All tickets
GET /bizchain/keberangkatan/keberangkatan # All journeys
```

**Agen & Ecosystem:**
```
GET /bizchain/agen/agen             # All agents
GET /bizchain/referral/referral     # All referrals
GET /bizchain/reward/reward         # All rewards
GET /bizchain/dao/proposal          # All DAO proposals
GET /bizchain/oleholeh/product      # Marketplace products
```

## 💡 Contoh Penggunaan

### Buat Produk dengan Multi-Cabang ✨

```bash
bizchaind tx pos create-product \
  "Indomie Goreng" \
  "3500" \
  "SKU-001" \
  "Makanan" \
  --description "Mie instan rasa goreng" \
  --branch-id "JKT" \
  --from alice \
  --chain-id bizchain-1 \
  --keyring-backend test \
  --yes
```

### Query Produk per Cabang ✨

```bash
curl http://localhost:1317/bizchain/pos/product_by_branch/JKT
```

### Registrasi Jamaah

```bash
bizchaind tx jamaah create-jamaah \
  "Ahmad Yani" \
  "081234567890" \
  "ahmad@email.com" \
  "Jakarta Selatan" \
  "A1234567" \
  --from alice \
  --chain-id bizchain-1 \
  --yes
```

### Buat Paket Umroh

```bash
bizchaind tx paket create-paket \
  "Paket Umroh Ekonomis" \
  "25000000" \
  "2024-03-15" \
  45 \
  --category "umroh" \
  --from alice \
  --chain-id bizchain-1 \
  --yes
```

## 🌐 Deploy ke GitHub Codespaces

1. Buka repository di GitHub
2. Klik **Code** > **Codespaces** > **Create codespace on main**
3. Ikuti panduan di [JALANKAN_NODE.md](JALANKAN_NODE.md#-menjalankan-di-github-codespaces)

## 🎯 Fitur Multi-Cabang

Sistem Multi-Cabang memungkinkan satu Group (Toko) mengelola banyak cabang dalam satu blockchain:

- Setiap cabang memiliki **inventory terpisah**
- Setiap cabang memiliki **transaksi terpisah**
- Owner dapat melihat **laporan seluruh cabang real-time**
- **Audit trail transparan** untuk semua aktivitas
- Support **transfer stok antar cabang**

### Contoh Struktur:
```
Group: TOKO ABC
├── Cabang Jakarta (JKT)
├── Cabang Bandung (BDG)
├── Cabang Semarang (SMG)
└── Cabang Surabaya (SUR)
```

Lihat dokumentasi lengkap di [MULTI_CABANG.md](MULTI_CABANG.md)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License.

## 🔗 Links

- **Repository**: https://github.com/darknase168/bizchain
- **Issues**: https://github.com/darknase168/bizchain/issues
- **Cosmos SDK**: https://docs.cosmos.network
- **Tendermint**: https://docs.tendermint.com

## 👨‍💻 Author

**darknase168**
- GitHub: [@darknase168](https://github.com/darknase168)

---

**Built with ❤️ using Cosmos SDK**
