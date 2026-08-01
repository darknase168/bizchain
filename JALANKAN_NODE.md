# Panduan Menjalankan BizChain Node

## 🚀 Menjalankan di Windows (Lokal)

### 1. Inisialisasi Chain (Pertama Kali)

```bash
cd d:\bizchain
.\build\bizchaind.exe init mynode --chain-id bizchain-1
```

### 2. Konfigurasi (Optional - Zero Fees)

Edit file `~/.bizchain/config/app.toml`:
```toml
minimum-gas-prices = "0upoint"
```

### 3. Buat Genesis Account

```bash
# Buat wallet
.\build\bizchaind.exe keys add alice

# Tambahkan ke genesis
.\build\bizchaind.exe add-genesis-account alice 1000000000000upoint --keyring-backend test

# Generate genesis transaction
.\build\bizchaind.exe gentx alice 1000000upoint --chain-id bizchain-1 --keyring-backend test

# Collect gentxs
.\build\bizchaind.exe collect-gentxs
```

### 4. Start Node

```bash
.\build\bizchaind.exe start
```

Node akan berjalan di:
- **gRPC**: `localhost:9090`
- **REST API**: `localhost:1317`
- **Tendermint RPC**: `localhost:26657`

---

## 🌐 Menjalankan di GitHub Codespaces

### 1. Buka Repository di Codespaces

1. Buka repository GitHub Anda
2. Klik tombol **Code** > **Codespaces** > **Create codespace on main**
3. Tunggu codespace selesai dibuat

### 2. Install Dependencies

```bash
# Update system
sudo apt-get update

# Install Go (jika belum ada)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Verify Go installation
go version
```

### 3. Build Blockchain

```bash
cd /workspaces/bizchain  # atau path workspace Anda
go build -o build/bizchaind ./cmd/bizchaind
```

### 4. Inisialisasi Chain

```bash
./build/bizchaind init mynode --chain-id bizchain-1 --home ~/.bizchain
```

### 5. Konfigurasi Ports untuk Codespaces

Edit `~/.bizchain/config/config.toml`:

```toml
[rpc]
laddr = "tcp://0.0.0.0:26657"

[p2p]
laddr = "tcp://0.0.0.0:26656"
```

Edit `~/.bizchain/config/app.toml`:

```toml
[api]
enable = true
address = "tcp://0.0.0.0:1317"

[grpc]
address = "0.0.0.0:9090"

minimum-gas-prices = "0upoint"
```

### 6. Buat Genesis Account

```bash
# Buat wallet
./build/bizchaind keys add alice --keyring-backend test

# Simpan mnemonic yang muncul!

# Tambahkan ke genesis
./build/bizchaind add-genesis-account alice 1000000000000upoint --keyring-backend test

# Generate genesis transaction
./build/bizchaind gentx alice 1000000upoint --chain-id bizchain-1 --keyring-backend test

# Collect gentxs
./build/bizchaind collect-gentxs
```

### 7. Start Node di Background

```bash
# Gunakan nohup untuk background process
nohup ./build/bizchaind start > bizchain.log 2>&1 &

# Atau gunakan screen
screen -S bizchain
./build/bizchaind start
# Tekan Ctrl+A, lalu D untuk detach

# Lihat log
tail -f bizchain.log

# Atau attach kembali ke screen
screen -r bizchain
```

### 8. Forward Ports di Codespaces

Codespaces akan otomatis mendeteksi ports yang terbuka. Atau Anda bisa manual forward:

1. Buka tab **PORTS** di VS Code
2. Klik **Forward a Port**
3. Tambahkan ports:
   - `26657` - Tendermint RPC
   - `26656` - P2P
   - `1317` - REST API
   - `9090` - gRPC

4. Set **Port Visibility** ke **Public** untuk REST API (1317)

### 9. Akses Node

Setelah ports di-forward, Anda akan mendapatkan URL public seperti:

```
https://[codespace-name]-1317.app.github.dev
```

Test dengan curl:
```bash
curl https://[codespace-name]-1317.app.github.dev/node_info
```

---

## 🖥️ Menjalankan Web Wallet

### Di Windows (Lokal)

```bash
cd d:\bizchain\wallet

# Install dependencies (pertama kali)
npm install

# Update CHAIN_CONFIG di src/types.ts
# Ganti restUrl dengan: http://localhost:1317

# Start development server
npm run dev
```

Buka browser: `http://localhost:5173`

### Di GitHub Codespaces

```bash
cd /workspaces/bizchain/wallet

# Install dependencies
npm install

# Update CHAIN_CONFIG di src/types.ts
# Ganti restUrl dengan URL codespace Anda
# Contoh: https://[codespace-name]-1317.app.github.dev

# Start development server
npm run dev
```

Codespaces akan otomatis forward port 5173. Klik notifikasi "Open in Browser" atau buka tab PORTS dan klik URL.

---

## 📝 Perintah Berguna

### Check Node Status

```bash
# Node info
./build/bizchaind status

# Query account balance
./build/bizchaind query bank balances [address]

# List keys
./build/bizchaind keys list --keyring-backend test
```

### Create Product (POS)

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

# Products by branch
curl http://localhost:1317/bizchain/pos/product_by_branch/JKT
```

### Stop Node

```bash
# Jika menggunakan nohup
ps aux | grep bizchaind
kill [PID]

# Jika menggunakan screen
screen -r bizchain
# Tekan Ctrl+C

# View logs
tail -f ~/.bizchain/bizchaind.log
# atau
tail -f bizchain.log
```

---

## 🔧 Troubleshooting

### Port Already in Use

```bash
# Cek port yang digunakan
netstat -ano | findstr :1317  # Windows
lsof -i :1317                 # Linux/Mac

# Kill process
taskkill /PID [PID] /F        # Windows
kill -9 [PID]                 # Linux/Mac
```

### Node Not Responding

```bash
# Reset chain data (HATI-HATI: akan hapus semua data!)
./build/bizchaind unsafe-reset-all

# Reinitialize
./build/bizchaind init mynode --chain-id bizchain-1
```

### Cannot Connect from Wallet

1. Pastikan node berjalan: `curl http://localhost:1317/node_info`
2. Periksa CORS di `~/.bizchain/config/app.toml`:
   ```toml
   [api]
   enable = true
   swagger = true
   enabled-unsafe-cors = true
   ```
3. Restart node setelah perubahan config

---

## 📚 Resources

- **Cosmos SDK Docs**: https://docs.cosmos.network
- **Tendermint Docs**: https://docs.tendermint.com
- **BizChain README**: `/workspaces/bizchain/README.md`
- **Multi-Cabang Guide**: `/workspaces/bizchain/MULTI_CABANG.md`

---

**Happy Coding! 🚀**
