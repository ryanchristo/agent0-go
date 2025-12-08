package core

import (
	"encoding/json"
	"io"
	"log"
	"maps"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ryanchristo/agent0-go/sdk/types"
	"github.com/ryanchristo/agent0-go/sdk/utils"
)

// SDKConfig is the configuration for the SDK.
type SDKConfig struct {
	ChainID           types.ChainID
	RPCURL            types.URI
	Signer            any // string, wallet, or signer
	RegistryOverrides RegistryOverrides

	// IPFS configuration

	IPFS               IPFSProvider
	IPFSNodeURL        string
	FilecoinPrivateKey string
	PinataJWT          string

	// Subgraph configuration

	SubgraphURL       string
	SubgraphOverrides SubgraphOverrides
}

// SDK is the main SDK instance.
type SDK struct {
	web3Client         *Web3Client
	ipfsClient         *IPFSClient
	subgraphClient     *SubgraphClient
	feedbackManager    *FeedbackManager
	indexer            *AgentIndexer
	identityRegistry   *Contract
	reputationRegistry *Contract
	validationRegistry *Contract
	registries         map[string]types.Address
	chainID            types.ChainID
	subgraphURLs       map[types.ChainID]string
}

// NewSDK creates a new SDK instance.
func NewSDK(cfg SDKConfig) *SDK {
	sdk := &SDK{}

	sdk.chainID = cfg.ChainID

	// Initialize web3 client
	sdk.web3Client = NewWeb3Client(cfg.RPCURL, cfg.Signer)

	// Resolve registry addresses
	mergedRegistries := make(map[string]types.Address)
	maps.Copy(mergedRegistries, DEFAULT_REGISTRIES[cfg.ChainID])
	maps.Copy(mergedRegistries, cfg.RegistryOverrides[cfg.ChainID])
	sdk.registries = mergedRegistries

	// Resolve subgraph URL
	if cfg.SubgraphOverrides != nil {
		sdk.subgraphURLs = cfg.SubgraphOverrides
	}

	resolvedSubgraphURL := ""
	if url, ok := sdk.subgraphURLs[cfg.ChainID]; ok {
		resolvedSubgraphURL = url
	} else if url, ok := DEFAULT_SUBGRAPH_URLS[cfg.ChainID]; ok {
		resolvedSubgraphURL = url
	} else {
		resolvedSubgraphURL = cfg.SubgraphURL
	}

	// Initialize subgraph client
	if resolvedSubgraphURL != "" {
		sdk.subgraphClient = NewSubgraphClient(resolvedSubgraphURL)
	}

	// Initialize indexer
	sdk.indexer = NewAgentIndexer(sdk.web3Client, sdk.subgraphClient, sdk.subgraphURLs)

	// Initialize IPFS client
	if cfg.IPFS != "" {
		sdk.ipfsClient = initalizeIPFSClient(cfg)
	}

	// Initialize feedback manager
	sdk.feedbackManager = NewFeedbackManager(
		sdk.web3Client,
		sdk.ipfsClient,
		nil, // reputationRegistry (lazy initialization)
		nil, // identityRegistry (lazy initialization)
		sdk.subgraphClient,
	)

	// Set subgraph client getter for multi-chain support
	sdk.feedbackManager.SetSubgraphClientGetter(
		func(chainID types.ChainID) *SubgraphClient {
			return sdk.GetSubgraphClient(chainID)
		},
		sdk.chainID,
	)

	return sdk
}

// initalizeIPFSClient initializes the IPFS client.
func initalizeIPFSClient(cfg SDKConfig) *IPFSClient {
	ipfsCfg := IPFSClientConfig{}

	switch cfg.IPFS {
	case IPFSProviderNode:
		if cfg.IPFSNodeURL == "" {
			log.Fatal("IPFSNodeURL is required when IPFS=\"node\"")
		}
		ipfsCfg.URL = cfg.IPFSNodeURL
	case IPFSProviderFilecoinPin:
		if cfg.FilecoinPrivateKey == "" {
			log.Fatal("FilecoinPrivateKey is required when IPFS=\"filecoinPin\"")
		}
		ipfsCfg.FilecoinPinEnabled = true
		ipfsCfg.FilecoinPrivateKey = cfg.FilecoinPrivateKey
	case IPFSProviderPinata:
		if cfg.PinataJWT == "" {
			log.Fatal("PinataJWT is required when IPFS=\"pinata\"")
		}
		ipfsCfg.PinataEnabled = true
		ipfsCfg.PinataJWT = cfg.PinataJWT
	case "":
		log.Fatal("IPFS provider not specified")
	default:
		log.Fatalf("Invalid IPFS value: %s. Must be \"node\", \"filecoinPin\", or \"pinata\"", cfg.IPFS)
	}

	return NewIPFSClient(ipfsCfg)
}

// ChainID returns the current chain ID.
func (s *SDK) ChainID() types.ChainID {
	if s.web3Client.ChainID == 0 {
		s.web3Client.Initialize()
	}
	return s.web3Client.ChainID
}

// Registries returns the resolved registry addresses for the current chain.
func (s *SDK) Registries() map[string]types.Address {
	return s.registries
}

// GetSubgraphClient returns the subgraph client for the given chain ID.
func (s *SDK) GetSubgraphClient(chainID types.ChainID) *SubgraphClient {
	targetChain := chainID
	if targetChain == 0 {
		targetChain = s.chainID
	}

	// Check if we already have a client for this chain
	if targetChain == s.chainID && s.subgraphClient != nil {
		return s.subgraphClient
	}

	// Resolve URL for target chain
	resolvedURL := ""
	if s.subgraphURLs != nil {
		if url, ok := s.subgraphURLs[targetChain]; ok {
			resolvedURL = url
		}
	}
	if resolvedURL == "" {
		if url, ok := DEFAULT_SUBGRAPH_URLS[targetChain]; ok {
			resolvedURL = url
		}
	}

	if resolvedURL != "" {
		return NewSubgraphClient(resolvedURL)
	}

	return nil
}

// GetIdentityRegistry returns the identity registry contract.
func (s *SDK) GetIdentityRegistry() *Contract {
	if s.identityRegistry == nil {
		address := s.registries["IDENTITY"]
		if address == "" {
			log.Fatalf("No identity registry address for chain %d", s.chainID)
		}
		s.identityRegistry = s.web3Client.GetContract(address, IDENTITY_REGISTRY_ABI)
	}
	return s.identityRegistry
}

// GetReputationRegistry returns the reputation registry contract.
func (s *SDK) GetReputationRegistry() *Contract {
	if s.reputationRegistry == nil {
		address := s.registries["REPUTATION"]
		if address == "" {
			log.Fatalf("No reputation registry address for chain %d", s.chainID)
		}
		s.reputationRegistry = s.web3Client.GetContract(address, REPUTATION_REGISTRY_ABI)

		// Update feedback manager
		s.feedbackManager.SetReputationRegistry(s.reputationRegistry)
	}
	return s.reputationRegistry
}

// GetValidationRegistry returns the validation registry contract.
func (s *SDK) GetValidationRegistry() *Contract {
	if s.validationRegistry == nil {
		address := s.registries["VALIDATION"]
		if address == "" {
			log.Fatalf("No validation registry address for chain %d", s.chainID)
		}
		s.validationRegistry = s.web3Client.GetContract(address, VALIDATION_REGISTRY_ABI)
	}
	return s.validationRegistry
}

// IsReadOnly checks if SDK is in read only mode (no signer).
func (s *SDK) IsReadOnly() bool {
	return s.web3Client.Address() == ""
}

// Agent lifecycle methods

// CreateAgent creates a new agent (off-chain instance in memory).
func (s *SDK) CreateAgent(name, description, imageURL string) *Agent {
	registrationFile := types.RegistrationFile{
		Name:        name,
		Description: description,
		Image:       imageURL,
		Endpoints:   []types.Endpoint{},
		TrustModels: []types.TrustModel{},
		Owners:      []string{},
		Operators:   []string{},
		Active:      false,
		X402Support: false,
		Metadata:    map[string]any{},
		UpdatedAt:   time.Now().Unix(),
	}
	return newAgent(registrationFile)
}

// LoadAgent loads an existing agent (hydrates from registration file if registered).
func (s *SDK) LoadAgent(agentID types.AgentID) *Agent {
	// Parse agent ID
	parsedAgentID := utils.ParseAgentID(agentID)
	chainID := parsedAgentID.ChainID
	tokenID := parsedAgentID.TokenID

	currentChainID := s.ChainID()
	if chainID != currentChainID {
		log.Fatalf("agent %s is not on current chain %d", agentID, currentChainID)
	}

	// Get token URI from contract
	tokenURI := ""
	identityRegistry := s.GetIdentityRegistry()
	if identityRegistry != nil {
		tokenURI = s.web3Client.CallContract(identityRegistry, "tokenURI", tokenID).(string)
	} else {
		log.Fatalf("identity registry not found for chain %d", currentChainID)
	}

	// Load registration file - handle empty URI (agent registered without URI)
	registrationFile := types.RegistrationFile{}
	if tokenURI == "" {
		registrationFile = s.createEmptyRegistrationFile()
	} else {
		registrationFile = s.loadRegistrationFile(tokenURI)
	}

	registrationFile.AgentID = agentID
	registrationFile.AgentURI = tokenURI

	return newAgent(registrationFile)
}

// GetAgent gets an agent summary from the subgraph (read-only).
// Supports both default chain and explicit chain specification.
func (s *SDK) GetAgent(agentID types.AgentID) types.AgentSummary {
	// Parse agentID to extract chainID if present
	// If no colon, assume it's just tokenID on default chain
	parsedChainID := types.ChainID(0)
	formattedAgentID := ""

	if strings.Contains(agentID, ":") {
		parsed := utils.ParseAgentID(agentID)
		parsedChainID = parsed.ChainID
		formattedAgentID = agentID
	} else {
		parsedChainID = s.chainID
		formattedAgentID = utils.FormattedAgentID(s.chainID, agentID)
	}

	// Determine which chain to query
	targetChainID := parsedChainID
	if targetChainID == 0 {
		targetChainID = s.chainID
	}

	// Get subgraph client for target chain (or use default)
	subgraphClient := s.subgraphClient
	if targetChainID != 0 {
		subgraphClient = s.GetSubgraphClient(targetChainID)
	}

	if subgraphClient == nil {
		log.Fatalf("Subgraph client required for getAgent on chain %d", targetChainID)
	}

	return subgraphClient.GetAgentByID(formattedAgentID)
}

// SearchAgents searches for agents matching the given query criteria.
// Supports multi-chain search when chains parameter is provided.
func (s *SDK) SearchAgents(params types.SearchParams, sort []string, pageSize int64, cursor string) AgentSearchResult {
	if pageSize == 0 {
		pageSize = utils.DEFAULTS["SEARCH_PAGE_SIZE"]
	}
	return s.indexer.SearchAgents(params, pageSize, cursor, sort)
}

// SearchAgentsByReputation searches for agents by reputation criteria.
func (s *SDK) SearchAgentsByReputation(
	agents []types.AgentID,
	tags []string,
	reviewers []types.Address,
	capabilities []string,
	skills []string,
	tasks []string,
	names []string,
	minAverageScore int64,
	includeRevoked bool,
	pageSize int64,
	cursor string,
	sort []string,
	chains []types.ChainID,
) AgentSearchResult {
	// Parse cursor to skip valie
	skip := int64(0)
	if cursor != "" {
		parsedCursor, err := strconv.ParseInt(cursor, 10, 64)
		if err == nil {
			skip = parsedCursor
		}
	}

	// Default sort
	if len(sort) == 0 {
		sort = []string{"createdAt:desc"}
	}

	return s.indexer.SearchAgentsByReputation(
		agents,
		tags,
		reviewers,
		capabilities,
		skills,
		tasks,
		names,
		minAverageScore,
		includeRevoked,
		pageSize,
		skip,
		sort,
		chains,
	)
}

// TransferAgent transfers agent ownership.
func (s *SDK) TransferAgent(agentID types.AgentID, newOwner types.Address) TransferResult {
	agent := s.LoadAgent(agentID)
	return agent.Transfer(newOwner)
}

// IsAgentOwner checks if the given address is the owner of the agent.
func (s *SDK) IsAgentOwner(agentID types.AgentID, address types.Address) bool {
	tokenID := utils.ParseAgentID(agentID).TokenID
	identityRegistry := s.GetIdentityRegistry()
	owner := s.web3Client.CallContract(identityRegistry, "ownerOf", big.NewInt(tokenID)).(string)
	return strings.EqualFold(owner, address)
}

// GetAgentOwner gets the current owner address of the agent.
func (s *SDK) GetAgentOwner(agentID types.AgentID) types.Address {
	tokenID := utils.ParseAgentID(agentID).TokenID
	identityRegistry := s.GetIdentityRegistry()
	owner := s.web3Client.CallContract(identityRegistry, "ownerOf", big.NewInt(tokenID)).(string)
	return types.Address(owner)
}

// Feedback methods

// SignFeedbackAuth signs feedback authorization and returns a signed authorization token.
func (s *SDK) SignFeedbackAuth(
	agentID types.AgentID,
	clientAddress types.Address,
	indexLimit int64,
	expiryHours int64,
) string {
	if expiryHours == 0 {
		expiryHours = utils.DEFAULTS["FEEDBACK_EXPIRY_HOURS"]
	}

	// Update feedback manager with registries
	s.feedbackManager.SetReputationRegistry(s.GetReputationRegistry())
	s.feedbackManager.SetIdentityRegistry(s.GetIdentityRegistry())

	return s.feedbackManager.SignFeedbackAuth(agentID, clientAddress, indexLimit, expiryHours)
}

// PrepareFeedback prepares feedback data for submission.
func (s *SDK) PrepareFeedback(
	agentID types.AgentID,
	score int64,
	tags []string,
	text string,
	capability string,
	name string,
	skill string,
	task string,
	context map[string]any,
	proofOfPayment map[string]any,
	extra map[string]any,
) FeedbackFile {
	return s.feedbackManager.PrepareFeedback(agentID, score, tags, text, capability, name, skill, task, context, proofOfPayment, extra)
}

// GiveFeedback submits feedback on-chain.
func (s *SDK) GiveFeedback(
	agentID types.AgentID,
	feedbackFile map[string]any,
	feedbackAuth string,
) Feedback {
	// Update feedback manager with registries
	s.feedbackManager.SetReputationRegistry(s.GetReputationRegistry())
	s.feedbackManager.SetIdentityRegistry(s.GetIdentityRegistry())

	return s.feedbackManager.GiveFeedback(agentID, feedbackFile, "", feedbackAuth)
}

// GetFeedback reads feedback from the contract.
func (s *SDK) GetFeedback(
	agentID types.AgentID,
	clientAddress types.Address,
	feedbackIndex int64,
) Feedback {
	return s.feedbackManager.GetFeedback(agentID, clientAddress, feedbackIndex)
}

// SearchFeedback searches for feedback entries with the given filters.
func (s *SDK) SearchFeedback(
	agentID types.AgentID,
	tags []string,
	capabilities []string,
	skills []string,
	minScore int64,
	maxScore int64,
) []types.Feedback {
	params := types.SearchFeedbackParams{
		Agents:       []types.AgentID{agentID},
		Tags:         tags,
		Capabilities: capabilities,
		Skills:       skills,
		MinScore:     minScore,
		MaxScore:     maxScore,
	}
	return s.feedbackManager.SearchFeedback(params)
}

// AppendResponse appends a response to feedback and returns the transaction hash.
func (s *SDK) AppendResponse(
	agentID types.AgentID,
	clientAddress types.Address,
	feedbackIndex int64,
	response FeedbackResponse,
) string {
	// Update feedback manager with registries
	s.feedbackManager.SetReputationRegistry(s.GetReputationRegistry())

	return s.feedbackManager.AppendResponse(agentID, clientAddress, feedbackIndex, response.URI, response.Hash)
}

// RevokeFeedback revokes feedback.
func (s *SDK) RevokeFeedback(
	agentID types.AgentID,
	feedbackIndex int64,
) string {
	// Update feedback manager with registries
	s.feedbackManager.SetReputationRegistry(s.GetReputationRegistry())

	return s.feedbackManager.RevokeFeedback(agentID, feedbackIndex)
}

// GetReputationSummary gets the reputation summary for an agent with a specific tag.
func (s *SDK) GetReputationSummary(
	agentID types.AgentID,
	tag1 string,
	tag2 string,
) ReputationSummary {
	// Update feedback manager with registries
	s.feedbackManager.SetReputationRegistry(s.GetReputationRegistry())

	return s.feedbackManager.GetReputationSummary(agentID, tag1, tag2)
}

// Private methods

// createEmptyRegistrationFile creates an empty registration file with default values.
func (s *SDK) createEmptyRegistrationFile() types.RegistrationFile {
	return types.RegistrationFile{
		Name:        "",
		Description: "",
		Endpoints:   []types.Endpoint{},
		TrustModels: []types.TrustModel{},
		Owners:      []string{},
		Operators:   []string{},
		Active:      false,
		X402Support: false,
		Metadata:    map[string]any{},
		UpdatedAt:   time.Now().Unix(),
	}
}

// loadRegistrationFile loads a registration file from a URI (IPFS or HTTP).
func (s *SDK) loadRegistrationFile(uri string) types.RegistrationFile {
	var rawData []byte

	// Handle IPFS URIs
	if strings.HasPrefix(uri, "ipfs://") {
		cid := strings.TrimPrefix(uri, "ipfs://")
		if s.ipfsClient != nil {
			// Use IPFS client if available
			rawData = s.ipfsClient.GetJSON(cid)
		} else {
			// Fallback to HTTP gateways if no IPFS client configured
			httpClient := &http.Client{
				// time.Duration receives nanoseconds, so we multiply by milliseconds
				Timeout: time.Duration(utils.TIMEOUTS["IPFS_GATEWAY"]) * time.Millisecond,
			}
			for _, gateway := range utils.IPFS_GATEWAYS {
				fetched := false
				url := gateway + cid
				resp, err := httpClient.Get(url)
				if err == nil {
					if resp.StatusCode == http.StatusOK {
						body, err := io.ReadAll(resp.Body)
						resp.Body.Close()
						if err == nil {
							rawData = body
							fetched = true
							break
						}
					} else {
						resp.Body.Close()
					}
				}

				if !fetched {
					log.Fatal("Failed to retrieve data from all IPFS gateways")
				}
			}
		}
	} else if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		resp, fetchErr := http.Get(uri)
		if fetchErr != nil {
			log.Fatalf("Failed to fetch registration file: HTTP %d: %v", resp.StatusCode, fetchErr)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			rawData = body
		}
	} else if strings.HasPrefix(uri, "data:") {
		// Data URIs are not supported
		log.Fatalf("Data URIs are not supported. Expected HTTP(S) or IPFS URI, got: %s", uri)
	} else if uri == "" || strings.TrimSpace(uri) == "" {
		// Return empty registration file (agent registered without URI)
		return s.createEmptyRegistrationFile()
	} else {
		log.Fatalf("Unsupported URI scheme: %s", uri)
	}

	// Validate rawData before transformation
	if rawData == nil {
		log.Fatalf("Invalid registration file format: expected an object")
	}
	var rawMap map[string]any
	if err := json.Unmarshal(rawData, &rawMap); err != nil {
		log.Fatalf("Failed to parse registration file JSON: %v", err)
	}

	return s.transformRegistrationFile(rawMap)
}

// transformRegistrationFile transforms raw registration data to a registration file.
func (s *SDK) transformRegistrationFile(rawData map[string]any) types.RegistrationFile {
	endpoints := s.transformEndpoints(rawData)
	extractedWallet := s.extractWalletInfo(rawData)
	walletAddress := extractedWallet.WalletAddress
	walletChainID := extractedWallet.ChainID

	// Extract trust models with proper type checking
	var trustModels []types.TrustModel
	if supportedTrust, ok := rawData["supportedTrust"].([]types.TrustModel); ok {
		trustModels = make([]types.TrustModel, 0, len(supportedTrust))
	} else {
		if trustModelsRaw, ok := rawData["trustModels"].([]types.TrustModel); ok {
			trustModels = make([]types.TrustModel, 0, len(trustModelsRaw))
		}
	}

	return types.RegistrationFile{
		Name:        rawData["name"].(string),
		Description: rawData["description"].(string),
		Image:       rawData["image"].(string),
		Endpoints:   endpoints,
		TrustModels: trustModels,
		Owners:      rawData["owners"].([]string),
		Operators:   rawData["operators"].([]string),
		Active:      rawData["active"].(bool),
		X402Support: rawData["x402Support"].(bool),
		Metadata: map[string]any{
			"version": rawData["version"].(string),
			"tags":    rawData["tags"].([]string),
		},
		UpdatedAt:     rawData["updatedAt"].(int64),
		WalletAddress: walletAddress,
		WalletChainID: walletChainID,
	}
}

// transformEndpoints transforms the endpoints from the raw data to the RegistrationFile format.
func (s *SDK) transformEndpoints(rawData map[string]any) []types.Endpoint {
	var endpoints []types.Endpoint

	if _, ok := rawData["endpoints"].([]any); !ok {
		return []types.Endpoint{}
	}

	for _, endpoint := range rawData["endpoints"].([]map[string]any) {
		// Check if endpoint is already in new format
		if endpoint["type"] != "" && endpoint["value"] != "" {
			endpoints = append(endpoints, types.Endpoint{
				Type:  endpoint["type"].(types.EndpointType),
				Value: endpoint["value"].(string),
				Meta:  endpoint["meta"].(map[string]any),
			})
		} else {
			// Transform endpoint from legacy format to new format
			transformed := s.transformEndpointLegacy(endpoint, rawData)
			if transformed != nil {
				endpoints = append(endpoints, *transformed)
			}
		}
	}

	return endpoints
}

// transformEndpointLegacy transforms an endpoint from the legacy format to the new format.
func (s *SDK) transformEndpointLegacy(endpoint map[string]any, rawData map[string]any) *types.Endpoint {
	name := endpoint["name"].(string)
	value := endpoint["value"].(string)
	version := endpoint["version"].(string)

	// Map endpoint names to types using case-insensitive lookup
	nameLower := strings.ToLower(name)
	var ENDPOINT_TYPE_MAP = map[string]types.EndpointType{
		"mcp":         types.ENDPOINT_TYPE_MCP,
		"a2a":         types.ENDPOINT_TYPE_A2A,
		"ens":         types.ENDPOINT_TYPE_ENS,
		"did":         types.ENDPOINT_TYPE_DID,
		"agentWallet": types.ENDPOINT_TYPE_WALLET,
		"wallet":      types.ENDPOINT_TYPE_WALLET,
	}

	var endpointType string
	if mappedType, ok := ENDPOINT_TYPE_MAP[nameLower]; ok {
		endpointType = string(mappedType)

		// Special handling for wallet endpoints - parse eip155 format
		if endpointType == string(types.ENDPOINT_TYPE_WALLET) {
			walletMatch := regexp.MustCompile(`eip155:(\d+):(0x[a-fA-F0-9]{40})`).FindStringSubmatch(value)
			if walletMatch != nil {
				walletChainID, err := strconv.ParseInt(walletMatch[1], 10, 64)
				if err != nil {
					log.Fatalf("Failed to parse wallet chain ID: %v", err)
				}
				rawData["walletAddress"] = walletMatch[2]
				rawData["walletChainID"] = walletChainID
			}
		}
	} else {
		endpointType = name // fallback to name as type
	}

	meta := map[string]any{}
	if version != "" {
		meta["version"] = version
	}

	return &types.Endpoint{
		Type:  types.EndpointType(endpointType),
		Value: value,
		Meta:  meta,
	}
}

// extractWalletInfo extracts the wallet address and chain ID from raw data.
func (s *SDK) extractWalletInfo(rawData map[string]any) WalletInfo {
	if walletAddress, ok := rawData["walletAddress"].(string); ok {
		if walletChainID, ok := rawData["walletChainID"].(int64); ok {
			return WalletInfo{
				WalletAddress: walletAddress,
				ChainID:       walletChainID,
			}
		}
	}
	return WalletInfo{}
}

// Expose clients for advanced usage

// Web3Client returns the web3 client.
func (s *SDK) Web3Client() *Web3Client {
	return s.web3Client
}

// IPFSClient returns the IPFS client.
func (s *SDK) IPFSClient() *IPFSClient {
	return s.ipfsClient
}

// SubgraphClient returns the subgraph client.
func (s *SDK) SubgraphClient() *SubgraphClient {
	return s.subgraphClient
}

// ...

type RegistryOverrides = map[types.ChainID]map[string]types.Address

type SubgraphOverrides = map[types.ChainID]string

type WalletInfo struct {
	WalletAddress string
	ChainID       int64
}
