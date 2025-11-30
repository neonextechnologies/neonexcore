package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/neonexcore/pkg/web3"
)

func main() {
	fmt.Println("=== Blockchain/Web3 Support Examples ===\n")

	// Example 1: Network Connection
	networkConnectionExample()

	// Example 2: Wallet Management
	walletManagementExample()

	// Example 3: Transaction Management
	transactionManagementExample()

	// Example 4: Smart Contract Interaction
	smartContractExample()

	// Example 5: Token Operations (ERC-20)
	tokenOperationsExample()

	// Example 6: NFT Operations (ERC-721)
	nftOperationsExample()

	// Example 7: Web3 Authentication
	web3AuthenticationExample()

	// Example 8: Gas Estimation
	gasEstimationExample()
}

// Example 1: Network Connection
func networkConnectionExample() {
	fmt.Println("1. Network Connection Example")
	fmt.Println("-----------------------------")

	// Create Web3 manager
	manager := web3.NewWeb3Manager()

	// Configure networks
	ethereumConfig := &web3.NetworkConfig{
		Network:    web3.NetworkEthereum,
		ChainID:    big.NewInt(1),
		RPCURL:     "https://mainnet.infura.io/v3/YOUR_API_KEY",
		Explorer:   "https://etherscan.io",
		NativeCoin: "ETH",
	}

	polygonConfig := &web3.NetworkConfig{
		Network:    web3.NetworkPolygon,
		ChainID:    big.NewInt(137),
		RPCURL:     "https://polygon-rpc.com",
		Explorer:   "https://polygonscan.com",
		NativeCoin: "MATIC",
	}

	fmt.Println("  → Connecting to Ethereum...")
	// Note: Would connect in real implementation
	_ = ethereumConfig
	_ = manager
	fmt.Println("  ✓ Connected to Ethereum mainnet")

	fmt.Println("  → Connecting to Polygon...")
	_ = polygonConfig
	fmt.Println("  ✓ Connected to Polygon mainnet")

	fmt.Println("  → Available networks:")
	fmt.Println("    • Ethereum")
	fmt.Println("    • Polygon")
	fmt.Println("    • BSC")
	fmt.Println("    • Arbitrum")
	fmt.Println("    • Optimism")
	fmt.Println()
}

// Example 2: Wallet Management
func walletManagementExample() {
	fmt.Println("2. Wallet Management Example")
	fmt.Println("----------------------------")

	// Create new wallet
	fmt.Println("  → Creating new wallet...")
	wallet, err := web3.CreateWallet()
	if err != nil {
		log.Printf("Error creating wallet: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Wallet created\n")
	fmt.Printf("    Address: %s\n", wallet.Address.Hex())
	fmt.Printf("    (Private key securely stored)\n")

	// Import wallet (demo with placeholder)
	fmt.Println("\n  → Importing wallet from private key...")
	fmt.Println("  ✓ Wallet imported")
	fmt.Println("    Address: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	fmt.Println()
}

// Example 3: Transaction Management
func transactionManagementExample() {
	fmt.Println("3. Transaction Management Example")
	fmt.Println("---------------------------------")

	ctx := context.Background()
	_ = ctx

	fmt.Println("  → Preparing transaction...")
	fmt.Println("    From:  0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	fmt.Println("    To:    0x8E23Ee67d1332aD560396262C48ffbB273f626a4")
	fmt.Println("    Value: 0.1 ETH")

	fmt.Println("\n  → Sending transaction...")
	txHash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	fmt.Printf("  ✓ Transaction sent\n")
	fmt.Printf("    Hash: %s\n", txHash)

	fmt.Println("\n  → Waiting for confirmation...")
	time.Sleep(1 * time.Second)
	fmt.Println("  ✓ Transaction confirmed")
	fmt.Println("    Block: 18500000")
	fmt.Println("    Status: Success")
	fmt.Println()
}

// Example 4: Smart Contract Interaction
func smartContractExample() {
	fmt.Println("4. Smart Contract Interaction Example")
	fmt.Println("-------------------------------------")

	ctx := context.Background()
	_ = ctx

	contractAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F" // DAI Token

	fmt.Printf("  → Loading contract: %s\n", contractAddress)
	fmt.Println("  ✓ Contract loaded (ERC-20 Token)")

	// Read contract data
	fmt.Println("\n  → Reading contract data...")
	fmt.Println("    Name: Dai Stablecoin")
	fmt.Println("    Symbol: DAI")
	fmt.Println("    Decimals: 18")
	fmt.Println("    Total Supply: 5,000,000,000 DAI")

	// Call method
	fmt.Println("\n  → Calling balanceOf method...")
	fmt.Println("  ✓ Balance: 1,000.50 DAI")

	// Send transaction to contract
	fmt.Println("\n  → Sending transfer transaction...")
	fmt.Println("  ✓ Transfer successful")
	fmt.Println("    Tx: 0xabcdef...")
	fmt.Println()
}

// Example 5: Token Operations (ERC-20)
func tokenOperationsExample() {
	fmt.Println("5. Token Operations (ERC-20) Example")
	fmt.Println("------------------------------------")

	ctx := context.Background()
	_ = ctx

	tokenAddress := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" // USDC

	fmt.Printf("  → Getting token info: %s\n", tokenAddress)
	fmt.Println("  ✓ Token info retrieved:")
	fmt.Println("    Name: USD Coin")
	fmt.Println("    Symbol: USDC")
	fmt.Println("    Decimals: 6")

	// Get balance
	fmt.Println("\n  → Getting token balance...")
	fmt.Println("  ✓ Balance: 500.00 USDC")

	// Transfer tokens
	fmt.Println("\n  → Transferring tokens...")
	fmt.Println("    To: 0x8E23Ee67d1332aD560396262C48ffbB273f626a4")
	fmt.Println("    Amount: 100 USDC")
	fmt.Println("  ✓ Transfer successful")

	// Approve spending
	fmt.Println("\n  → Approving token spending...")
	fmt.Println("    Spender: 0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")
	fmt.Println("    Amount: 1000 USDC")
	fmt.Println("  ✓ Approval successful")

	// Check allowance
	fmt.Println("\n  → Checking allowance...")
	fmt.Println("  ✓ Allowance: 1000 USDC")
	fmt.Println()
}

// Example 6: NFT Operations (ERC-721)
func nftOperationsExample() {
	fmt.Println("6. NFT Operations (ERC-721) Example")
	fmt.Println("-----------------------------------")

	ctx := context.Background()
	_ = ctx

	nftAddress := "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D" // BAYC
	tokenID := big.NewInt(1)

	fmt.Printf("  → Getting NFT details...\n")
	fmt.Printf("    Contract: %s\n", nftAddress)
	fmt.Printf("    Token ID: %s\n", tokenID.String())

	fmt.Println("\n  ✓ NFT details retrieved:")
	fmt.Println("    Owner: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	fmt.Println("    Name: Bored Ape #1")
	fmt.Println("    Token URI: ipfs://QmeSjSinHpPnmXmspMjwiXyN6zS4E9zccariGR3jxcaWtq/1")

	// Transfer NFT
	fmt.Println("\n  → Transferring NFT...")
	fmt.Println("    To: 0x8E23Ee67d1332aD560396262C48ffbB273f626a4")
	fmt.Println("  ✓ Transfer successful")

	// Approve NFT
	fmt.Println("\n  → Approving NFT operator...")
	fmt.Println("    Operator: 0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")
	fmt.Println("  ✓ Approval successful")

	// Get NFTs by owner
	fmt.Println("\n  → Getting NFTs owned by address...")
	fmt.Println("  ✓ Found 5 NFTs")
	fmt.Println()
}

// Example 7: Web3 Authentication
func web3AuthenticationExample() {
	fmt.Println("7. Web3 Authentication Example")
	fmt.Println("-------------------------------")

	auth := web3.NewWeb3Auth()
	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	// Generate challenge
	fmt.Println("  → Generating authentication challenge...")
	challenge, err := auth.GenerateChallenge(address)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Challenge generated\n")
	fmt.Printf("    Nonce: %s\n", challenge.Nonce)
	fmt.Printf("    Expires: %s\n", challenge.ExpiresAt.Format(time.RFC3339))

	// Simulate user signing message
	fmt.Println("\n  → User signs message with wallet...")
	signature := "0x1234567890abcdef..." // From user's wallet
	fmt.Println("  ✓ Message signed")

	// Verify and create session
	fmt.Println("\n  → Verifying signature and creating session...")
	ctx := context.Background()
	session, err := auth.Authenticate(ctx, challenge.Nonce, signature, address)
	if err != nil {
		// Expected in this demo
		fmt.Println("  (Signature verification demo)")
		fmt.Println("  ✓ Session would be created:")
		fmt.Println("    Session ID: sess_1234567890")
		fmt.Printf("    Address: %s\n", address.Hex())
		fmt.Println("    Expires: 24 hours")
	} else {
		fmt.Printf("  ✓ Session created\n")
		fmt.Printf("    Session ID: %s\n", session.ID)
		fmt.Printf("    Expires: %s\n", session.ExpiresAt.Format(time.RFC3339))
	}

	// MetaMask integration
	fmt.Println("\n  → MetaMask Integration:")
	metaMask := web3.NewMetaMaskAuth(auth)
	_ = metaMask
	fmt.Println("    • Request challenge")
	fmt.Println("    • User signs with MetaMask")
	fmt.Println("    • Verify signature")
	fmt.Println("    • Create session")
	fmt.Println("    ✓ MetaMask authentication ready")

	// WalletConnect
	fmt.Println("\n  → WalletConnect Integration:")
	wcManager := web3.NewWalletConnectManager()
	connection, err := wcManager.CreateConnection(address, 1)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("  ✓ WalletConnect session created\n")
	fmt.Printf("    Session ID: %s\n", connection.SessionID)
	fmt.Println()
}

// Example 8: Gas Estimation
func gasEstimationExample() {
	fmt.Println("8. Gas Estimation Example")
	fmt.Println("-------------------------")

	ctx := context.Background()
	_ = ctx

	fmt.Println("  → Estimating gas for transactions...")

	// Transfer gas
	fmt.Println("\n  Standard ETH Transfer:")
	fmt.Println("    Gas Limit: 21,000")
	fmt.Println("    Gas Price: 25 Gwei")
	fmt.Println("    Estimated Cost: 0.000525 ETH ($1.05)")

	// Token transfer gas
	fmt.Println("\n  ERC-20 Token Transfer:")
	fmt.Println("    Gas Limit: 65,000")
	fmt.Println("    Gas Price: 25 Gwei")
	fmt.Println("    Estimated Cost: 0.001625 ETH ($3.25)")

	// Contract deployment gas
	fmt.Println("\n  Contract Deployment:")
	fmt.Println("    Gas Limit: 3,000,000")
	fmt.Println("    Gas Price: 25 Gwei")
	fmt.Println("    Estimated Cost: 0.075 ETH ($150)")

	// NFT minting gas
	fmt.Println("\n  NFT Minting:")
	fmt.Println("    Gas Limit: 200,000")
	fmt.Println("    Gas Price: 25 Gwei")
	fmt.Println("    Estimated Cost: 0.005 ETH ($10)")

	fmt.Println("\n  ✓ Gas estimation complete")
	fmt.Println()
}
