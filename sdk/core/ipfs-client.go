package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	ipfscid "github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/multiformats/go-multiaddr"

	"github.com/ryanchristo/agent0-go/sdk/types"
	"github.com/ryanchristo/agent0-go/sdk/utils"
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
	client   *rpc.HttpApi
}

// NewIPFSClient creates a new IPFS client.
func NewIPFSClient(config IPFSClientConfig) *IPFSClient {
	ipfsClient := &IPFSClient{}

	ipfsClient.config = config

	// Determine provider
	if config.PinataEnabled {
		ipfsClient.provider = IPFSProviderPinata
		ipfsClient.verifyPinataJwt()
	} else if config.FilecoinPinEnabled {
		ipfsClient.provider = IPFSProviderFilecoinPin
		ipfsClient.verifyFilecoinPrivateKey()
	} else if config.URL != "" {
		ipfsClient.provider = IPFSProviderNode
	} else {
		log.Fatal("No IPFS provider configured")
	}

	return ipfsClient
}

// ensureClient ensures the IPFS client is initialized.
func (c *IPFSClient) ensureClient() {
	if c.provider == IPFSProviderNode && c.client == nil {
		address := multiaddr.StringCast(c.config.URL)
		client, err := rpc.NewApiWithClient(address, http.DefaultClient)
		if err != nil {
			log.Fatalf("Failed to create IPFS client: %v", err)
		}
		c.client = client
	}
}

// verifyPinataJwt verifies the Pinata JWT.
func (c *IPFSClient) verifyPinataJwt() {
	if c.config.PinataJWT == "" {
		log.Fatal("PinataJWT is required when PinataEnabled=true")
	}
}

// verifyFilecoinPrivateKey verifies the Filecoin private key.
func (c *IPFSClient) verifyFilecoinPrivateKey() {
	if c.config.FilecoinPrivateKey == "" {
		log.Fatal("FilecoinPrivateKey is required when FilecoinPinEnabled=true")
	}
}

// addToPinata adds (and pins) data to IPFS via Pinata v3 API.
func (c *IPFSClient) addToPinata(data string) string {

	// TODO: implementation

	return ""
}

// addToFilecoin adds (and pins) data to IPFS via Filecoin.
func (c *IPFSClient) addToFilecoin(data string) string {

	// TODO: implementation

	return ""
}

// addToLocalIPFS adds data to the local IPFS node.
func (c *IPFSClient) addToLocalIPFS(data string) string {

	// Initialize client if not already initialized
	c.ensureClient()
	if c.client == nil {
		log.Fatal("No IPFS client available")
	}

	// Create empty context
	ctx := context.Background()

	// Create file from data
	f := files.NewBytesFile([]byte(data))

	// Add file to IPFS node
	p, err := c.client.Unixfs().Add(ctx, f)
	if err != nil {
		log.Fatalf("Failed to add the data: %v", err)
		return ""
	}

	// Return root CID
	return p.RootCid().String()
}

// Add adds data to IPFS and returns the CID.
func (c *IPFSClient) Add(data string) string {
	switch c.provider {
	case IPFSProviderPinata:
		return c.addToPinata(data)
	case IPFSProviderFilecoinPin:
		return c.addToFilecoin(data)
	case IPFSProviderNode:
		return c.addToLocalIPFS(data)
	default:
		log.Fatal("No IPFS provider configured")
		return ""
	}
}

// AddFile adds a file to IPFS and returns the CID.
func (c *IPFSClient) AddFile(filepath string) string {

	// Read file from disk
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return ""
	}

	switch c.provider {
	case IPFSProviderPinata:
		return c.addToPinata(string(data))
	case IPFSProviderFilecoinPin:
		return c.addToFilecoin(string(data))
	case IPFSProviderNode:
		return c.addToLocalIPFS(string(data))
	default:
		log.Fatal("No IPFS provider configured")
		return ""
	}
}

// Get gets data from IPFS by CID.
func (c *IPFSClient) Get(cid string) string {

	// Remove "ipfs://" prefix if present
	cleanCID, _ := strings.CutPrefix(cid, "ipfs://")

	switch c.provider {
	case IPFSProviderPinata, IPFSProviderFilecoinPin:

		// TODO: try all gateways in parallel

		httpClient := &http.Client{
			// time.Duration receives nanoseconds, so we multiply by milliseconds
			Timeout: time.Duration(utils.TIMEOUTS["IPFS_GATEWAY"]) * time.Millisecond,
		}

		for _, gateway := range utils.IPFS_GATEWAYS {
			url := gateway + cleanCID
			resp, err := httpClient.Get(url)
			if err != nil {
				log.Printf("Failed to get data from gateway: %v", err)
				continue
			}

			// Read the response
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read response: %v", err)
				return ""
			}

			// Return the response body
			return string(respBody)
		}

		// If no response was successful, return an error
		log.Fatal("Failed to get data from all IPFS gateways")
		return ""

	case IPFSProviderNode:

		// Initialize client if not already initialized
		c.ensureClient()
		if c.client == nil {
			log.Fatal("No IPFS client available")
			return ""
		}

		// Create empty context
		ctx := context.Background()

		// Create decoded CID from CID string
		decodedCID, err := ipfscid.Decode(cleanCID)
		if err != nil {
			log.Fatalf("Failed to decode CID: %v", err)
			return ""
		}

		// Create path from decoded CID
		path := path.FromCid(decodedCID)

		// Get the file node from the path
		node, err := c.client.Unixfs().Get(ctx, path)
		if err != nil {
			log.Fatalf("Failed to get node: %v", err)
			return ""
		}

		// Get the file from the node
		file := files.ToFile(node)

		// Read the file to bytes
		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Failed to read file: %v", err)
			return ""
		}

		// Return file content
		return string(bytes)

	default:
		log.Fatal("No IPFS provider configured")
		return ""
	}
}

// GetJSON gets JSON data from IPFS by CID.
func (c *IPFSClient) GetJSON(cid string) map[string]any {

	// Get data from IPFS
	data := c.Get(cid)

	// Convert data to raw map
	var rawMap map[string]any
	if err := json.Unmarshal([]byte(data), &rawMap); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
		return nil
	}

	return rawMap
}

// Pin pins data to a local IPFS node and returns the result.
func (c *IPFSClient) Pin(cid string) PinResult {
	switch c.provider {
	case IPFSProviderPinata:
		// addToPinata adds and pins data so we add again to be safe
		pinnedCID := c.addToPinata(cid)
		return PinResult{
			Pinned: []string{pinnedCID},
		}
	case IPFSProviderFilecoinPin:
		// addToFilecoin adds and pins data so we add again to be safe
		pinnedCID := c.addToFilecoin(cid)
		return PinResult{
			Pinned: []string{pinnedCID},
		}
	case IPFSProviderNode:

		// Initialize client if not already initialized
		c.ensureClient()
		if c.client == nil {
			log.Fatal("No IPFS client available")
			return PinResult{}
		}

		// Create empty context
		ctx := context.Background()

		// Create decoded CID from CID string
		dcid, err := ipfscid.Decode(cid)
		if err != nil {
			log.Fatalf("Failed to decode CID: %v", err)
			return PinResult{}
		}

		// Create path from decoded CID
		path := path.FromCid(dcid)

		// Pin the data to the local IPFS node
		err = c.client.Pin().Add(ctx, path)
		if err != nil {
			log.Fatalf("Failed to pin the data: %v", err)
			return PinResult{}
		}

		// Return the pinned CID
		return PinResult{
			Pinned: []string{cid},
		}

	default:
		log.Fatal("No IPFS provider configured")
		return PinResult{}
	}
}

// Unpin unpins data from a local IPFS node and returns the result.
func (c *IPFSClient) Unpin(cid string) UnpinResult {
	switch c.provider {
	case IPFSProviderPinata:

		// TODO: implementation

		return UnpinResult{
			Unpinned: []string{},
		}
	case IPFSProviderFilecoinPin:

		// TODO: implementation

		return UnpinResult{
			Unpinned: []string{},
		}
	case IPFSProviderNode:

		// Initialize client if not already initialized
		c.ensureClient()
		if c.client == nil {
			log.Fatal("No IPFS client available")
			return UnpinResult{}
		}

		// Create empty context
		ctx := context.Background()

		// Create decoded CID from CID string
		dcid, err := ipfscid.Decode(cid)
		if err != nil {
			log.Fatalf("Failed to decode CID: %v", err)
			return UnpinResult{}
		}

		// Create path from decoded CID
		path := path.FromCid(dcid)

		// Unpin the data from the local IPFS node
		err = c.client.Pin().Rm(ctx, path)
		if err != nil {
			log.Fatalf("Failed to unpin the data: %v", err)
			return UnpinResult{}
		}

		// Return the unpinned CID
		return UnpinResult{
			Unpinned: []string{cid},
		}

	default:
		log.Fatal("No IPFS provider configured")
		return UnpinResult{}
	}
}

// AddJSON adds JSON data to IPFS and returns the CID.
func (c *IPFSClient) AddJSON(data map[string]any) string {

	// Convert data to JSON bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
		return ""
	}

	// Add JSON data to IPFS
	return c.Add(string(jsonData))
}

// AddRegistrationFile adds a registration file to IPFS and returns the CID.
func (c *IPFSClient) AddRegistrationFile(
	registrationFile types.RegistrationFile,
	chainID types.ChainID,
	identityRegistryAddress types.Address,
) string {

	// Convert the endpoints array data from the internal format { type, value, meta }
	// to the ERC-8004 format { name, endpoint, version }
	var endpoints []map[string]any
	for _, endpoint := range registrationFile.Endpoints {
		endpointDict := map[string]any{
			"name":     endpoint.Type,
			"endpoint": endpoint.Value,
		}

		// Spread meta fields (version, mcpTools, mcpPrompts, etc.)
		if endpoint.Meta != nil {
			maps.Copy(endpointDict, endpoint.Meta)
		}

		endpoints = append(endpoints, endpointDict)
	}

	// Add walletAddress as an endpoint if present
	if registrationFile.WalletAddress != "" {
		walletChainID := registrationFile.WalletChainID
		if walletChainID == 0 {
			walletChainID = chainID
		}
		if walletChainID == 0 {
			walletChainID = 1
		}
		endpoints = append(endpoints, map[string]any{
			"name":     "agentWallet",
			"endpoint": fmt.Sprintf("eip155:%d:%s", walletChainID, registrationFile.WalletAddress),
		})
	}

	// Build registrations array
	var registrations []map[string]any
	if registrationFile.AgentID != "" {
		parts := strings.Split(registrationFile.AgentID, ":")
		var agentRegistry string
		if chainID != 0 {
			agentRegistry = fmt.Sprintf("eip155:%d:%s", chainID, identityRegistryAddress)
		} else {
			agentRegistry = fmt.Sprintf("eip155:1:%s", identityRegistryAddress)
		}
		agentID, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			log.Fatalf("Failed to parse agent ID: %v", err)
			return ""
		}
		registrations = append(registrations, map[string]any{
			"agentId":       agentID,
			"agentRegistry": agentRegistry,
		})
	}

	// Build ERC-8004 compliant registration file
	data := map[string]any{
		"type":        "https://eips.ethereum.org/EIPS/eip-8004#registration-v1",
		"name":        registrationFile.Name,
		"description": registrationFile.Description,
		"active":      registrationFile.Active,
		"x402Support": registrationFile.X402Support,
	}

	// Conditionally add fields only if they're not empty
	if registrationFile.Image != "" {
		data["image"] = registrationFile.Image
	}
	if len(endpoints) > 0 {
		data["endpoints"] = endpoints
	}
	if len(registrations) > 0 {
		data["registrations"] = registrations
	}
	if len(registrationFile.TrustModels) > 0 {
		data["supportedTrusts"] = registrationFile.TrustModels
	}

	return c.AddJSON(data)
}

// GetRegistrationFile gets a registration file from IPFS by CID.
func (c *IPFSClient) GetRegistrationFile(cid string) types.RegistrationFile {

	// Get data from IPFS
	data := c.Get(cid)

	// Convert data to registration file
	var registrationFile types.RegistrationFile
	if err := json.Unmarshal([]byte(data), &registrationFile); err != nil {
		log.Printf("Failed to parse registration file: %v", err)
		return types.RegistrationFile{}
	}

	return registrationFile
}

// Close closes the IPFS client connection.
func (c *IPFSClient) Close() {
	if c.client != nil {
		c.client = nil
	}
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
