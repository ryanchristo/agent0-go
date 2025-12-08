package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

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
	Provider *ethrpc.Server
	Signer   Signer // wallet or signer
	ChainID  types.ChainID
}

// NewWeb3Client creates a new Web3Client instance.
func NewWeb3Client(rpcURL string, signerOrKey any) *Web3Client {
	web3Client := &Web3Client{}

	// Create a new RPC server instance.
	web3Client.Provider = ethrpc.NewServer()

	// TODO: implementation

	return web3Client
}

// Initialize initializes the Web3Client instance.
func (c *Web3Client) Initialize() {
	// TODO: implementation
}

// GetContract gets a contract instance from the Web3Client.
func (c *Web3Client) GetContract(address types.Address, abi string) *Contract {

	// TODO: implementation

	return &Contract{}
}

// CallContract calls a contract method (view/pure function) and returns the result.
func (c *Web3Client) CallContract(contract *Contract, methodName string, args ...interface{}) any {

	// TODO: implementation

	return nil
}

// TransactContract executes a contract transaction and returns the transaction hash.
func (c *Web3Client) TransactContract(contract *Contract, methodName string, options TransactionOptions, args ...interface{}) string {

	// TODO: implementation

	return ""
}

// registerAgent is a router wrapper for register() function overloads.
// The function selected the correct overload based on the following arguments:
// - register() - no arguments
// - register(string tokenUri) - one argument
// - register(string tokenUri, tuple[] metadata) - two arguments
func (c *Web3Client) registerAgent(contract *Contract, options TransactionOptions, args ...interface{}) string {

	// TODO: implementation

	return ""
}

// WaitForTransaction waits for a transaction to be mined and returns the receipt.
func (c *Web3Client) WaitForTransaction(txHash string, timeout int64) ethtypes.Receipt {
	// default timeout = 60000

	// TODO: implementation

	return ethtypes.Receipt{}
}

// GetEvents gets the events from the contract and returns a list of logs.
func (c *Web3Client) GetEvents(contract *Contract, eventName string, fromBlock, toBlock int64) []ethtypes.Log {
	// default fromBlock = 0

	// TODO: implementation

	return []ethtypes.Log{}
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

	// TODO: implementation

	return ""
}

// SignMessage signs a message with the account's private key.
func (c *Web3Client) SignMessage(message string) string {

	// TODO: implementation

	return ""
}

// RecoverAddress recovers the address from a message and signature.
func (c *Web3Client) RecoverAddress(message, signature string) string {

	// TODO: implementation

	return ""
}

// Keccak256 computes the Keccak-256 hash of the input data.
func (c *Web3Client) Keccak256(data string) string {

	// TODO: implementation

	return ""
}

// ToChecksumAddress converts an address to checksum format.
func (c *Web3Client) ToChecksumAddress(address string) string {

	// TODO: implementation

	return ""
}

// IsAddress checks if the address is a valid Ethereum address.
func (c *Web3Client) IsAddress(address string) bool {

	// TODO: implementation

	return false
}

// GetBalance gets the ETH balance of the address.
func (c *Web3Client) GetBalance(address string) big.Int {

	// TODO: implementation

	return big.Int{}
}

// GetTransactionCount gets the transaction count of the address.
func (c *Web3Client) GetTransactionCount(address string) int64 {

	// TODO: implementation

	return 0
}

// Address gets the account address (if signer is available).
func (c *Web3Client) Address() types.Address {

	// TODO: implementation

	return ""
}

// GetAddress gets the account address (if signer is available).
func (c *Web3Client) GetAddress() types.Address {

	// TODO: implementation

	return ""
}

// ...

type Signer any // wallet or signer

type Contract bind.BoundContract
