package core

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ryanchristo/agent0-go/sdk/types"
)

// TransactionOptions are the options for a transaction.
type TransactionOptions struct {
	GasLimit             big.Int
	GasPrice             big.Int
	MaxFeePerGas         big.Int
	MaxPriorityFeePerGas big.Int
}

// Web3Client is a client for interacting with the Ethereum blockchain.
type Web3Client struct {
	Provider *ethclient.Client
	Signer   bind.SignerFn
	ChainID  types.ChainID

	// TODO: signer limitations
	privateKey *ecdsa.PrivateKey
}

// NewWeb3Client creates a new Web3Client instance.
func NewWeb3Client(rpcURL string, signerOrKey any) *Web3Client {
	web3Client := &Web3Client{}

	// Create client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	// Set provider
	web3Client.Provider = client

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Set chain ID
	web3Client.ChainID = types.ChainID(chainID.Int64())

	if signerOrKey != nil {

		// Signer is a string (private key)
		if key, ok := signerOrKey.(string); ok {

			// Trim the private key (remove whitespace)
			trimmedKey := strings.TrimSpace(key)
			if trimmedKey == "" {
				log.Fatal("Private key cannot be empty")
			}

			// Parse the private key
			privateKey, err := crypto.HexToECDSA(trimmedKey)
			if err != nil {
				log.Fatalf("Failed to parse private key: %v", err)
			}

			// Create signer from private key and chain ID
			txOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			if err != nil {
				log.Fatalf("Failed to create signer: %v", err)
			}

			web3Client.Signer = txOpts.Signer

			// TODO: signer limitations
			web3Client.privateKey = privateKey

		} else if signer, ok := signerOrKey.(bind.SignerFn); ok {

			// TODO: signer limitations
			log.Print("warning: signer limitations")

			// Signer is already signer
			web3Client.Signer = signer

		} else {
			log.Fatal("Invalid signer or key")
		}
	}

	return web3Client
}

// Initialize initializes the Web3Client.
func (c *Web3Client) Initialize() {
	// do nothing, we already set the chain ID in the constructor
}

// GetContract gets a contract instance from the Web3Client.
func (c *Web3Client) GetContract(address types.Address, abi string) *Contract {

	// Parse contract ABI
	contractABI, err := ethabi.JSON(bytes.NewReader([]byte(abi)))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Create bound contract
	contract := bind.NewBoundContract(
		common.HexToAddress(address),
		contractABI,
		c.Provider, // caller
		c.Provider, // transactor
		c.Provider, // filterer
	)

	return contract
}

// CallContract calls a contract method (view/pure function) and returns the result.
func (c *Web3Client) CallContract(contract *Contract, methodName string, args ...any) any {
	ctx := context.Background()

	// Check if the contract is the identity registry and the provided method exists
	if contract.Address() == common.HexToAddress(DEFAULT_REGISTRIES[c.ChainID]["IDENTITY"]) {
		identityABI, err := ethabi.JSON(strings.NewReader(IDENTITY_REGISTRY_ABI))
		if err != nil {
			log.Fatalf("Failed to parse identity registry ABI: %v", err)
		}
		_, exists := identityABI.Methods[methodName]
		if !exists {
			log.Fatal("Method not found")
		}
	}

	// Check if the contract is the reputation registry and the provided method exists
	if contract.Address() == common.HexToAddress(DEFAULT_REGISTRIES[c.ChainID]["REPUTATION"]) {
		reputationABI, err := ethabi.JSON(strings.NewReader(REPUTATION_REGISTRY_ABI))
		if err != nil {
			log.Fatalf("Failed to parse reputation registry ABI: %v", err)
		}
		_, exists := reputationABI.Methods[methodName]
		if !exists {
			log.Fatal("Method not found")
		}
	}

	// Check if the contract is the validation registry and the provided method exists
	if contract.Address() == common.HexToAddress(DEFAULT_REGISTRIES[c.ChainID]["VALIDATION"]) {
		validationABI, err := ethabi.JSON(strings.NewReader(VALIDATION_REGISTRY_ABI))
		if err != nil {
			log.Fatalf("Failed to parse validation registry ABI: %v", err)
		}
		_, exists := validationABI.Methods[methodName]
		if !exists {
			log.Fatal("Method not found")
		}
	}

	var result []any

	// Create call options
	opts := &bind.CallOpts{
		Context: ctx,
	}

	return contract.Call(opts, &result, methodName, args...)
}

// TransactContract executes a contract transaction and returns the transaction hash.
func (c *Web3Client) TransactContract(contract *Contract, methodName string, options TransactionOptions, args ...any) string {
	if c.Signer == nil {
		log.Fatal("Cannot execute transaction: SDK is in read-only mode.")
	}

	// Special handling for register() function with multiple overloads
	if methodName == "register" {
		return c.registerAgent(contract, options, args...)
	}

	var data []byte

	// Check if the contract is the identity registry and the provided method exists
	if contract.Address() == common.HexToAddress(DEFAULT_REGISTRIES[c.ChainID]["IDENTITY"]) {
		identityABI, err := ethabi.JSON(strings.NewReader(IDENTITY_REGISTRY_ABI))
		if err != nil {
			log.Fatalf("Failed to parse identity registry ABI: %v", err)
		}
		_, exists := identityABI.Methods[methodName]
		if !exists {
			log.Fatal("Method not found")
		}
		data, err = identityABI.Pack(methodName, args...)
		if err != nil {
			log.Fatalf("Failed to pack data: %v", err)
		}
	}

	// Check if the contract is the reputation registry and the provided method exists
	if contract.Address() == common.HexToAddress(DEFAULT_REGISTRIES[c.ChainID]["REPUTATION"]) {
		reputationABI, err := ethabi.JSON(strings.NewReader(REPUTATION_REGISTRY_ABI))
		if err != nil {
			log.Fatalf("Failed to parse reputation registry ABI: %v", err)
		}
		_, exists := reputationABI.Methods[methodName]
		if !exists {
			log.Fatal("Method not found")
		}
		data, err = reputationABI.Pack(methodName, args...)
		if err != nil {
			log.Fatalf("Failed to pack data: %v", err)
		}
	}

	// Check if the contract is the validation registry and the provided method exists
	if contract.Address() == common.HexToAddress(DEFAULT_REGISTRIES[c.ChainID]["VALIDATION"]) {
		validationABI, err := ethabi.JSON(strings.NewReader(VALIDATION_REGISTRY_ABI))
		if err != nil {
			log.Fatalf("Failed to parse validation registry ABI: %v", err)
		}
		_, exists := validationABI.Methods[methodName]
		if !exists {
			log.Fatal("Method not found")
		}
		data, err = validationABI.Pack(methodName, args...)
		if err != nil {
			log.Fatalf("Failed to pack data: %v", err)
		}
	}

	ctx := context.Background()

	nonce, err := c.Provider.PendingNonceAt(ctx, common.HexToAddress(c.GetAddress()))
	if err != nil {
		log.Fatalf("Failed to get pending nonce: %v", err)
	}

	if options.GasPrice.Sign() == 0 {
		suggestedGasPrice, err := c.Provider.SuggestGasPrice(ctx)
		if err != nil {
			log.Fatalf("Failed to suggest gas price: %v", err)
		}
		options.GasPrice = *suggestedGasPrice
	}

	if options.GasLimit.Sign() == 0 {
		contractAddr := contract.Address()
		estimatedGas, err := c.Provider.EstimateGas(ctx, ethereum.CallMsg{
			To:   &contractAddr,
			Data: data,
		})
		if err != nil {
			log.Fatalf("Failed to estimate gas needed: %v", err)
		}
		options.GasLimit = *big.NewInt(int64(estimatedGas))
	}

	// Build transaction options
	opts := &bind.TransactOpts{
		From:       common.HexToAddress(c.GetAddress()),
		Nonce:      big.NewInt(int64(nonce)),
		Signer:     c.Signer,
		Value:      nil,
		GasPrice:   &options.GasPrice,
		GasFeeCap:  &options.MaxFeePerGas,
		GasTipCap:  &options.MaxPriorityFeePerGas,
		GasLimit:   options.GasLimit.Uint64(),
		AccessList: nil,
		Context:    ctx,
		NoSend:     false,
	}

	// Send transaction
	tx, err := contract.Transact(opts, methodName, args...)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	return tx.Hash().Hex()
}

// registerAgent is a router wrapper for register() function overloads.
// The function selected the correct overload based on the following arguments:
// - register() - no arguments
// - register(string tokenUri) - one argument
// - register(string tokenUri, tuple[] metadata) - two arguments
func (c *Web3Client) registerAgent(contract *Contract, options TransactionOptions, args ...any) string {
	if c.Signer == nil {
		log.Fatal("No signer available for transaction")
	}

	// Determine which overload to use based on arguments
	var methodName string
	var callArgs []any

	if len(args) == 0 {
		methodName = "register()"
		callArgs = []any{}
	} else if len(args) == 1 {
		methodName = "register(string)"
		callArgs = []any{args[0]}
	} else if len(args) == 2 {
		methodName = "register(string,(string,bytes)[])"
		callArgs = []any{args[0], args[1]}
	} else {
		log.Fatalf("Invalid number of arguments for register() function: %d", len(args))
	}

	// Parse the JSON into an ABI object
	parsedABI, err := ethabi.JSON(strings.NewReader(IDENTITY_REGISTRY_ABI))
	if err != nil {
		log.Fatal(err)
	}

	// Access the function fragment
	_, exists := parsedABI.Methods[methodName]
	if !exists {
		log.Fatal("Method not found")
	}

	// Encode function data to avoid ambiguity - this bypasses function resolution
	data, err := parsedABI.Pack(methodName, callArgs...)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	ctx := context.Background()

	nonce, err := c.Provider.PendingNonceAt(ctx, common.HexToAddress(c.GetAddress()))
	if err != nil {
		log.Fatalf("Failed to get pending nonce: %v", err)
	}

	if options.GasPrice.Sign() == 0 {
		suggestedGasPrice, err := c.Provider.SuggestGasPrice(ctx)
		if err != nil {
			log.Fatalf("Failed to suggest gas price: %v", err)
		}
		options.GasPrice = *suggestedGasPrice
	}

	if options.GasLimit.Sign() == 0 {
		contractAddr := contract.Address()
		estimatedGas, err := c.Provider.EstimateGas(ctx, ethereum.CallMsg{
			To:   &contractAddr,
			Data: data,
		})
		if err != nil {
			log.Fatalf("Failed to estimate gas needed: %v", err)
		}
		options.GasLimit = *big.NewInt(int64(estimatedGas))
	}

	// Build transaction options
	opts := &bind.TransactOpts{
		From:       common.HexToAddress(c.GetAddress()),
		Nonce:      big.NewInt(int64(nonce)),
		Signer:     c.Signer,
		Value:      nil,
		GasPrice:   &options.GasPrice,
		GasFeeCap:  &options.MaxFeePerGas,
		GasTipCap:  &options.MaxPriorityFeePerGas,
		GasLimit:   options.GasLimit.Uint64(),
		AccessList: nil,
		Context:    ctx,
		NoSend:     false,
	}

	// Send transaction directly with encoded data (no function call resolution needed)
	tx, err := contract.RawTransact(opts, data)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	return tx.Hash().Hex()
}

// WaitForTransaction waits for a transaction to be mined and returns the receipt.
func (c *Web3Client) WaitForTransaction(txHash string, timeout int64) *ethtypes.Receipt {
	if timeout == 0 {
		timeout = 60000 // why not utils.TIMEOUTS["TRANSACTION_WAIT"] ?
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()

	receipt, err := bind.WaitMinedHash(ctx, c.Provider, common.HexToHash(txHash))
	if err != nil {
		log.Fatalf("Failed to wait for transaction to be mined: %v", err)
	}

	return receipt
}

// GetEvents gets the events from the contract and returns a list of logs.
func (c *Web3Client) GetEvents(contract *Contract, eventName string, fromBlock, toBlock int64) []ethtypes.Log {
	ctx := context.Background()
	end := uint64(toBlock)

	// Create filter options
	opts := &bind.FilterOpts{
		Start:   uint64(fromBlock),
		Context: ctx,
	}
	if end != 0 {
		opts.End = &end
	}

	// Get contract logs
	logs, _, err := contract.FilterLogs(opts, eventName)
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}

	// Convert chan to slice
	var events []ethtypes.Log
	for l := range logs {
		events = append(events, l)
	}

	return events
}

// EncodeFeedbackAuth encodes the feedback authorization data for a client.
func (c *Web3Client) EncodeFeedbackAuth(
	agentID big.Int,
	clientAddress string,
	indexLimit big.Int,
	expiry big.Int,
	chainID big.Int,
	identityRegistry string,
	signerAddress string,
) string {

	// Helper function to create ABI type
	mustABIType := func(typeString string) ethabi.Type {
		t, err := ethabi.NewType(typeString, "", nil)
		if err != nil {
			log.Fatalf("Failed to create ABI type: %v", err)
		}
		return t
	}

	// Create ABI arguments
	arguments := ethabi.Arguments{
		{
			Type: mustABIType("uint256"),
			Name: "agentId",
		},
		{
			Type: mustABIType("address"),
			Name: "clientAddress",
		},
		{
			Type: mustABIType("uint256"),
			Name: "indexLimit",
		},
		{
			Type: mustABIType("uint256"),
			Name: "expiry",
		},
		{
			Type: mustABIType("uint256"),
			Name: "chainId",
		},
		{
			Type: mustABIType("address"),
			Name: "identityRegistry",
		},
		{
			Type: mustABIType("address"),
			Name: "signerAddress",
		},
	}

	// Pack arguments to get the encoded data
	encoded, err := arguments.Pack(agentID, clientAddress, indexLimit, expiry, chainID, identityRegistry, signerAddress)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	// Convert bytes to hash and return hex string
	return common.BytesToHash(encoded).Hex()
}

// SignMessage signs a message with the account's private key.
func (c *Web3Client) SignMessage(message string) string {

	// Check if private key is available
	if c.privateKey != nil {

		// Hash the message
		hash := crypto.Keccak256Hash([]byte(message))

		// Sign the message
		signature, err := crypto.Sign(hash.Bytes(), c.privateKey)
		if err != nil {
			log.Fatalf("Failed to sign message: %v", err)
		}

		// Encode the signature
		return hexutil.Encode(signature)
	}

	// TODO: signer limitations
	log.Fatal("signer limitations")

	return ""
}

// RecoverAddress recovers the address from a message and signature.
func (c *Web3Client) RecoverAddress(message, signature string) string {
	hash := crypto.Keccak256Hash([]byte(message))

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), []byte(signature))
	if err != nil {
		log.Fatalf("Failed to recover public key: %v", err)
	}

	return common.BytesToAddress(sigPublicKey).Hex()
}

// Keccak256 computes the Keccak-256 hash of the input data.
func (c *Web3Client) Keccak256(data any) string {
	if dataStr, ok := data.(string); ok {
		return crypto.Keccak256Hash([]byte(dataStr)).Hex()
	}

	if dataBytes, ok := data.([]byte); ok {
		return crypto.Keccak256Hash(dataBytes).Hex()
	}

	log.Fatalf("Invalid data type")

	return ""
}

// ToChecksumAddress converts an address to checksum format.
func (c *Web3Client) ToChecksumAddress(address string) string {
	return common.HexToAddress(address).Hex()
}

// IsAddress checks if the address is a valid Ethereum address.
func (c *Web3Client) IsAddress(address string) bool {
	return common.IsHexAddress(address)
}

// GetBalance gets the ETH balance of the address.
func (c *Web3Client) GetBalance(address string) *big.Int {
	balance, err := c.Provider.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	return balance
}

// GetTransactionCount gets the transaction count of the address.
func (c *Web3Client) GetTransactionCount(address string) int64 {
	nonce, err := c.Provider.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to get transaction count: %v", err)
	}
	return int64(nonce)
}

// Address gets the account address (if signer is available).
func (c *Web3Client) Address() types.Address {
	if c.privateKey != nil {
		publicKey := c.privateKey.Public().(*ecdsa.PublicKey)
		publicKeyBytes := crypto.FromECDSAPub(publicKey)
		return common.BytesToAddress(publicKeyBytes).Hex()
	}

	// TODO: signer limitations
	log.Fatal("signer limitations")

	return ""
}

// GetAddress gets the account address (if signer is available).
func (c *Web3Client) GetAddress() types.Address {
	if c.privateKey != nil {
		publicKey := c.privateKey.Public().(*ecdsa.PublicKey)
		publicKeyBytes := crypto.FromECDSAPub(publicKey)
		return common.BytesToAddress(publicKeyBytes).Hex()
	}

	// TODO: signer limitations
	log.Fatal("signer limitations")

	return ""
}

// ...

type Contract = bind.BoundContract
