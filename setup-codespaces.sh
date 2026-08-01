#!/bin/bash

# BizChain Setup Script for GitHub Codespaces
# This script automates the setup and running of BizChain node in Codespaces

set -e  # Exit on error

echo "🚀 BizChain Codespaces Setup Script"
echo "===================================="

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Check Go installation
echo -e "\n${BLUE}Step 1: Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Go not found. Installing Go 1.21...${NC}"
    wget -q https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    rm go1.21.0.linux-amd64.tar.gz
fi

go version
echo -e "${GREEN}✓ Go is ready${NC}"

# Step 2: Build blockchain
echo -e "\n${BLUE}Step 2: Building blockchain binary...${NC}"
go build -o build/bizchaind ./cmd/bizchaind
chmod +x build/bizchaind
echo -e "${GREEN}✓ Binary built successfully${NC}"

# Step 3: Initialize chain (if not already initialized)
echo -e "\n${BLUE}Step 3: Initializing blockchain...${NC}"
if [ ! -d "$HOME/.bizchain" ]; then
    ./build/bizchaind init mynode --chain-id bizchain-1 --home ~/.bizchain
    echo -e "${GREEN}✓ Chain initialized${NC}"
else
    echo -e "${YELLOW}Chain already initialized, skipping...${NC}"
fi

# Step 4: Configure for Codespaces (bind to 0.0.0.0)
echo -e "\n${BLUE}Step 4: Configuring for Codespaces...${NC}"

# Update config.toml
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' ~/.bizchain/config/config.toml
sed -i 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26656"/' ~/.bizchain/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' ~/.bizchain/config/config.toml

# Update app.toml
sed -i 's/enable = false/enable = true/' ~/.bizchain/config/app.toml
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' ~/.bizchain/config/app.toml
sed -i 's/address = "localhost:9090"/address = "0.0.0.0:9090"/' ~/.bizchain/config/app.toml
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.bizchain/config/app.toml

# Set zero fees
if ! grep -q 'minimum-gas-prices = "0upoint"' ~/.bizchain/config/app.toml; then
    sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0upoint"/' ~/.bizchain/config/app.toml
fi

echo -e "${GREEN}✓ Configuration updated for Codespaces${NC}"

# Step 5: Create genesis account (if not exists)
echo -e "\n${BLUE}Step 5: Setting up genesis account...${NC}"
if ! ./build/bizchaind keys show alice --keyring-backend test &> /dev/null; then
    echo -e "${YELLOW}Creating alice account...${NC}"
    echo "y" | ./build/bizchaind keys add alice --keyring-backend test 2>&1 | tee alice-key.txt
    
    ALICE_ADDRESS=$(./build/bizchaind keys show alice -a --keyring-backend test)
    
    ./build/bizchaind add-genesis-account alice 1000000000000upoint --keyring-backend test
    ./build/bizchaind gentx alice 1000000upoint --chain-id bizchain-1 --keyring-backend test
    ./build/bizchaind collect-gentxs
    
    echo -e "${GREEN}✓ Genesis account created${NC}"
    echo -e "${YELLOW}⚠️  Alice address: ${ALICE_ADDRESS}${NC}"
    echo -e "${YELLOW}⚠️  Mnemonic saved to alice-key.txt - KEEP IT SAFE!${NC}"
else
    echo -e "${YELLOW}Alice account already exists, skipping...${NC}"
    ALICE_ADDRESS=$(./build/bizchaind keys show alice -a --keyring-backend test)
    echo -e "${YELLOW}Alice address: ${ALICE_ADDRESS}${NC}"
fi

# Step 6: Start node
echo -e "\n${BLUE}Step 6: Starting blockchain node...${NC}"
echo -e "${YELLOW}Node will start in background. Use 'tail -f bizchain.log' to view logs.${NC}"

# Kill existing process if any
pkill -f bizchaind || true
sleep 2

# Start in background with nohup
nohup ./build/bizchaind start > bizchain.log 2>&1 &
NODE_PID=$!
echo $NODE_PID > bizchain.pid

echo -e "${GREEN}✓ Node started (PID: ${NODE_PID})${NC}"

# Wait for node to be ready
echo -e "\n${BLUE}Waiting for node to be ready...${NC}"
for i in {1..30}; do
    if curl -s http://localhost:26657/status > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Node is ready!${NC}"
        break
    fi
    echo -n "."
    sleep 1
done

# Step 7: Display info
echo -e "\n${GREEN}=====================================${NC}"
echo -e "${GREEN}🎉 BizChain is ready!${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""
echo -e "${BLUE}📡 Endpoints:${NC}"
echo "  - Tendermint RPC: http://localhost:26657"
echo "  - REST API:       http://localhost:1317"
echo "  - gRPC:           localhost:9090"
echo ""
echo -e "${BLUE}👤 Alice Account:${NC}"
echo "  - Address: ${ALICE_ADDRESS}"
echo "  - Mnemonic: Check alice-key.txt"
echo ""
echo -e "${BLUE}📝 Useful Commands:${NC}"
echo "  - View logs:      tail -f bizchain.log"
echo "  - Stop node:      kill \$(cat bizchain.pid)"
echo "  - Node status:    ./build/bizchaind status"
echo "  - Check balance:  ./build/bizchaind query bank balances ${ALICE_ADDRESS}"
echo ""
echo -e "${BLUE}🌐 Forward Ports in Codespaces:${NC}"
echo "  1. Go to PORTS tab in VS Code"
echo "  2. Forward ports: 26657, 1317, 9090"
echo "  3. Set port 1317 visibility to 'Public'"
echo ""
echo -e "${BLUE}🚀 Next Steps:${NC}"
echo "  1. Setup Web Wallet:"
echo "     cd wallet && npm install"
echo "     # Update CHAIN_CONFIG.restUrl in src/types.ts with your Codespaces URL"
echo "     npm run dev"
echo ""
echo "  2. Test API:"
echo "     curl http://localhost:1317/node_info"
echo ""
echo "  3. Create a product:"
echo "     ./build/bizchaind tx pos create-product \\"
echo "       'Indomie Goreng' '3500' 'SKU-001' 'Makanan' \\"
echo "       --description 'Mie instan' --branch-id 'JKT' \\"
echo "       --from alice --chain-id bizchain-1 \\"
echo "       --keyring-backend test --fees 0upoint --yes"
echo ""
echo -e "${YELLOW}📚 Documentation: README.md, JALANKAN_NODE.md${NC}"
echo ""
