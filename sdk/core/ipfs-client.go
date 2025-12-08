package core

import (
	"log"

	"github.com/ryanchristo/agent0-go/sdk/types"
)

// IPFSClientConfig is the configuration for the IPFS client.
type IPFSClientConfig struct {
	URL                string
	FilecoinPinEnabled bool
	FilecoinPrivateKey string
	PinataEnabled      bool
	PinataJWT          string
}

// IPFSClient is IPFS operations supporting multiple providers.
type IPFSClient struct {
	provider IPFSProvider
	config   IPFSClientConfig
	client   any // TODO: Add IPFS client
}

func NewIPFSClient(config IPFSClientConfig) *IPFSClient {
	ipfsClient := &IPFSClient{}

	ipfsClient.config = config

	// Determine provider
	if config.PinataEnabled {
		ipfsClient.provider = IPFSProviderPinata
		ipfsClient.verifyPinataJwt()
	} else if config.FilecoinPinEnabled {
		ipfsClient.provider = IPFSProviderFilecoinPin
	} else if config.URL != "" {
		ipfsClient.provider = IPFSProviderNode
	} else {
		log.Fatal("No IPFS provider configured")
	}

	return ipfsClient
}

// ensureClient ensures the IPFS client is initialized.
func (c *IPFSClient) ensureClient() {
	// TODO: implementation
}

// verifyPinataJwt verifies the Pinata JWT.
func (c *IPFSClient) verifyPinataJwt() {
	// TODO: implementation
}

// pinToPinata pins data to Pinata using v3 API.
func (c *IPFSClient) pinToPinata(data string) string {

	// TODO: implementation

	return ""
}

// pinToFilecoin pins data to Filecoin.
func (c *IPFSClient) pinToFilecoin(data string) string {

	// TODO: implementation

	return ""
}

// pinToLocalIPFS pins data to the local IPFS node.
func (c *IPFSClient) pinToLocalIPFS(data string) string {

	// TODO: implementation

	return ""
}

// Add adds data to IPFS and returns the CID.
func (c *IPFSClient) Add(data string) string {

	// TODO: implementation

	return ""
}

// AddFile adds a file to IPFS and returns the CID.
func (c *IPFSClient) AddFile(filepath string) string {

	// TODO: implementation

	return ""
}

// Get gets data from IPFS by CID.
func (c *IPFSClient) Get(cid string) string {

	// TODO: implementation

	return ""
}

// GetJSON gets JSON data from IPFS by CID.
func (c *IPFSClient) GetJSON(cid string) []byte {

	// TODO: implementation

	return []byte{}
}

// Pin pins data to a local IPFS node and returns the result.
func (c *IPFSClient) Pin(cid string) PinResult {

	// TODO: implementation

	return PinResult{}
}

// Unpin unpins data from a local IPFS node and returns the result.
func (c *IPFSClient) Unpin(cid string) UnpinResult {

	// TODO: implementation

	return UnpinResult{}
}

// AddJSON adds JSON data to IPFS and returns the CID.
func (c *IPFSClient) AddJSON(data map[string]any) string {

	// TODO: implementation

	return ""
}

// AddRegistrationFile adds a registration file to IPFS and returns the CID.
func (c *IPFSClient) AddRegistrationFile(
	registrationFile types.RegistrationFile,
	chainID types.ChainID,
	identityRegistryAddress types.Address,
) string {

	// TODO: implementation

	return ""
}

// GetRegistrationFile gets a registration file from IPFS by CID.
func (c *IPFSClient) GetRegistrationFile(cid string) types.RegistrationFile {

	// TODO: implementation

	return types.RegistrationFile{}
}

// Close closes the IPFS client connection.
func (c *IPFSClient) Close() {
	// TODO: implementation
}

// ...

type IPFSProvider string

const (
	IPFSProviderNode        IPFSProvider = "node"
	IPFSProviderFilecoinPin IPFSProvider = "filecoinPin"
	IPFSProviderPinata      IPFSProvider = "pinata"
)

type PinResult struct {
	Pinned []string
}

type UnpinResult struct {
	Unpinned []string
}
