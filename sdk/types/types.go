package types

// AgentID is the ID of the agent (chainID:tokenID or tokenID).
type AgentID = string

// Chain ID is the chain ID.
type ChainID = int64

// Ethereum address is the address of the Ethereum account (0x...).
type Address = string

// URI is the URI of the resource (https://... or ipfs://...).
type URI = string

// CID is the CID of the resource (IPFS).
type CID = string

// Timestamp is the timestamp in seconds.
type Timestamp = int64

// IdemKey is the idempotency key for write operations.
type IdemKey = string
