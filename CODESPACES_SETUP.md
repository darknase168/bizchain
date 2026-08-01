# 🚀 Setup BizChain di GitHub Codespaces

Panduan lengkap untuk menjalankan BizChain blockchain node di GitHub Codespaces.

## ✅ Langkah 1: Buat Codespace

1. Buka repository: https://github.com/darknase168/bizchain
2. Klik tombol **Code** > **Codespaces**
3. Klik **Create codespace on main**
4. Tunggu Codespaces selesai dibuat (sekitar 2-3 menit)

## 🔧 Langkah 2: Install Dependencies

```bash
# Update system
sudo apt-get update -y

# Install Go 1.21 (jika belum ada)
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
```

## 🏗️ Langkah 3: Build Blockchain

```bash
cd /workspaces/bizchain  # Path codespace Anda
go build -o build/bizchaind ./cmd/bizchaind
chmod +x build/bizchaind
```

## ⚙️ Langkah 4: Inisialisasi Chain

```bash
# Initialize node
./build/bizchaind init mynode --chain-id bizchain-1 --home ~/.bizchain

# Configure zero fees
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0upoint"/' ~/.bizchain/config/app.toml

# Enable API dan CORS untuk Codespaces
sed -i 's/enable = false/enable = true/' ~/.bizchain/config/app.toml
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.bizchain/config/app.toml
sed -i 's/swagger = false/swagger = true/' ~/.bizchain/config/app.toml

# Set API address untuk bind ke semua interface
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' ~/.bizchain/config/app.toml

# Set gRPC address
sed -i 's/address = "localhost:9090"/address = "0.0.0.0:9090"/' ~/.bizchain/config/app.toml

# Set RPC address di config.toml
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' ~/.bizchain/config/config.toml

# Set CORS untuk RPC
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' ~/.bizchain/config/config.toml
```

## 👤 Langkah 5: Buat Genesis Account

```bash
# Create alice wallet
./build/bizchaind keys add alice --keyring-backend test

# PENTING: Simpan mnemonic phrase yang muncul!

# Add alice to genesis
./build/bizchaind add-genesis-account alice 1000000000000upoint --keyring-backend test

# Update bond_denom dari stake ke upoint (PENTING!)
sed -i 's/"bond_denom": "stake"/"bond_denom": "upoint"/' ~/.bizchain/config/genesis.json

# Generate genesis transaction
./build/bizchaind gentx alice 1000000upoint --chain-id bizchain-1 --keyring-backend test

# Collect gentxs
./build/bizchaind collect-gentxs
```

## 🚀 Langkah 6: Start Node

```bash
# Start node di background dengan nohup
nohup ./build/bizchaind start > bizchain.log 2>&1 &

# Simpan PID untuk stop nanti
echo $! > bizchain.pid

# View logs real-time
tail -f bizchain.log

# Tekan Ctrl+C untuk keluar dari tail (node tetap jalan)
```

## 🌐 Langkah 7: Forward Ports

1. Buka tab **PORTS** di VS Code (biasanya di bawah terminal)
2. Codespaces akan otomatis detect ports, atau klik **Forward a Port**
3. Tambahkan ports ini:
   - **26657** - Tendermint RPC
   - **26656** - P2P
   - **1317** - REST API (Set visibility ke **Public**)
   - **9090** - gRPC
   - **5173** - Wallet (nanti)

4. Untuk port 1317, klik kanan > **Port Visibility** > **Public**

5. Salin URL port 1317, contoh:
   ```
   https://scaling-dollop-q7964qgrq65jf9764.github.dev/
   ```

## ✅ Langkah 8: Test Node

```bash
# Test dengan curl (ganti dengan URL Anda)
curl https://scaling-dollop-q7964qgrq65jf9764.github.dev/node_info

# Atau test localhost
curl http://localhost:1317/node_info

# Check balance
curl http://localhost:1317/cosmos/bank/v1beta1/balances/$(./build/bizchaind keys show alice -a --keyring-backend test)
```

## 💻 Langkah 9: Setup Web Wallet

```bash
cd /workspaces/bizchain/wallet

# Install dependencies
npm install

# Update CHAIN_CONFIG di src/types.ts
# Ganti restUrl dengan URL Codespaces port 1317 Anda
# Sudah diupdate ke: https://scaling-dollop-q7964qgrq65jf9764.github.dev/

# Start development server
npm run dev
```

Codespaces akan otomatis forward port 5173. Klik notifikasi "Open in Browser" atau buka tab PORTS dan klik URL untuk port 5173.

## 📝 Perintah Berguna

### Check Node Status

```bash
# Node info
./build/bizchaind status

# List keys
./build/bizchaind keys list --keyring-backend test

# Get alice address
./build/bizchaind keys show alice -a --keyring-backend test

# Query balance
./build/bizchaind query bank balances $(./build/bizchaind keys show alice -a --keyring-backend test)
```

### Create Product (POS dengan Multi-Cabang)

```bash
./build/bizchaind tx pos create-product \
  "Indomie Goreng" \
  "3500" \
  "SKU-001" \
  "Makanan" \
  --description "Mie instan rasa goreng" \
  --branch-id "JKT" \
  --from alice \
  --chain-id bizchain-1 \
  --keyring-backend test \
  --fees 0upoint \
  --yes
```

### Query Products

```bash
# All products
curl http://localhost:1317/bizchain/pos/product

# Products by branch (Multi-Cabang feature)
curl http://localhost:1317/bizchain/pos/product_by_branch/JKT

# Atau dengan URL public Codespaces
curl https://scaling-dollop-q7964qgrq65jf9764.github.dev/bizchain/pos/product
```

### Stop Node

```bash
# Stop node
kill $(cat bizchain.pid)

# Atau cari PID manual
ps aux | grep bizchaind
kill <PID>
```

### View Logs

```bash
# View logs
tail -f bizchain.log

# View last 100 lines
tail -100 bizchain.log

# Search for errors
grep -i error bizchain.log
```

## 🔧 Troubleshooting

### Node tidak start

```bash
# Check logs
tail -100 bizchain.log

# Check process
ps aux | grep bizchaind

# Reset chain data (HATI-HATI: akan hapus semua data!)
./build/bizchaind tendermint unsafe-reset-all
```

### Port sudah digunakan

```bash
# Check ports
lsof -i :1317
lsof -i :26657

# Kill process
kill -9 <PID>
```

### Cannot connect from wallet

1. Pastikan node berjalan: `curl http://localhost:1317/node_info`
2. Pastikan port 1317 di-forward dan visibility-nya **Public**
3. Pastikan `restUrl` di `wallet/src/types.ts` sesuai dengan URL Codespaces port 1317 Anda
4. Restart wallet: `cd wallet && npm run dev`

### Chain-ID mismatch error

Ini bug yang sedang diselidiki. Jika terjadi error "invalid chain-id on InitChain", lakukan:

```bash
# Stop node
kill $(cat bizchain.pid)

# Update genesis.json - pastikan bond_denom = upoint
sed -i 's/"bond_denom": "stake"/"bond_denom": "upoint"/' ~/.bizchain/config/genesis.json

# Delete data dan reinit
rm -rf ~/.bizchain/data
mkdir ~/.bizchain/data
echo '{"height": "0", "round": 0, "step": 0}' > ~/.bizchain/data/priv_validator_state.json

# Start node
nohup ./build/bizchaind start > bizchain.log 2>&1 &
echo $! > bizchain.pid
```

## 📚 Resources

- **Repository**: https://github.com/darknase168/bizchain
- **README**: [/workspaces/bizchain/README.md](../README.md)
- **Multi-Cabang Guide**: [/workspaces/bizchain/MULTI_CABANG.md](../MULTI_CABANG.md)
- **Implementation Guide**: [/workspaces/bizchain/IMPLEMENTASI_MULTI_CABANG.md](../IMPLEMENTASI_MULTI_CABANG.md)

## 🎯 Fitur Multi-Cabang

Setelah node berjalan, Anda bisa test fitur Multi-Cabang:

```bash
# Create products di berbagai cabang
./build/bizchaind tx pos create-product "Indomie Jakarta" "3500" "JKT-001" "Makanan" \
  --description "Product Cabang Jakarta" --branch-id "JKT" \
  --from alice --chain-id bizchain-1 --keyring-backend test --fees 0upoint --yes

./build/bizchaind tx pos create-product "Indomie Bandung" "3600" "BDG-001" "Makanan" \
  --description "Product Cabang Bandung" --branch-id "BDG" \
  --from alice --chain-id bizchain-1 --keyring-backend test --fees 0upoint --yes

# Query per cabang
curl http://localhost:1317/bizchain/pos/product_by_branch/JKT
curl http://localhost:1317/bizchain/pos/product_by_branch/BDG

# View di Web Wallet - pilih cabang dari dropdown
```

---

**Happy Coding di Codespaces! 🚀**

