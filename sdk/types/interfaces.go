package types

// Endpoint is an agent endpoint.
type Endpoint struct {
	// Type is the type of the endpoint.
	Type EndpointType `json:"type"`

	// Value is the value of the endpoint.
	Value string `json:"value"`

	// Meta is the metadata of the endpoint.
	Meta map[string]any `json:"meta,omitempty"`
}

// RegistrationFile is the registration file for an agent.
type RegistrationFile struct {
	// AgentID is the ID of the agent.
	AgentID AgentID `json:"agentId,omitempty"`

	// AgentURI is the URI of the published agent file.
	AgentURI URI `json:"agentUri,omitempty"`

	// Name is the name of the agent.
	Name string `json:"name"`

	// Description is the description of the agent.
	Description string `json:"description"`

	// Image is the image of the agent.
	Image URI `json:"image,omitempty"`

	// WalletAddress is the wallet address of the agent.
	WalletAddress Address `json:"walletAddress,omitempty"`

	// WalletChainID is the wallet chain ID of the agent.
	WalletChainID ChainID `json:"walletChainId,omitempty"`

	// Endpoints is the endpoints of the agent.
	Endpoints []Endpoint `json:"endpoints"`

	// TrustModels is the trust models of the agent.
	TrustModels []TrustModel `json:"trustModels"`

	// Owners is the owners of the agent.
	Owners []Address `json:"owners"`

	// Operators is the operators of the agent.
	Operators []Address `json:"operators"`

	// Active is the active status of the agent.
	Active bool `json:"active"`

	// X402Support is the X402 support status of the agent.
	X402Support bool `json:"x402Support"`

	// Metadata is the metadata of the agent.
	Metadata map[string]any `json:"metadata"`

	// UpdatedAt is the timestamp of the last update.
	UpdatedAt Timestamp `json:"updatedAt"`
}

// AgentSummary is the summary information of an agent.
type AgentSummary struct {
	// ChainID is the chain ID of the agent.
	ChainID ChainID `json:"chainId"`

	// AgentID is the ID of the agent.
	AgentID AgentID `json:"agentId"`

	// Name is the name of the agent.
	Name string `json:"name"`

	// Description is the description of the agent.
	Description string `json:"description"`

	// Image is the image of the agent.
	Image URI `json:"image,omitempty"`

	// Owners is the owners of the agent.
	Owners []Address `json:"owners"`

	// Operators is the operators of the agent.
	Operators []Address `json:"operators"`

	// MCP is the MCP support status of the agent.
	MCP bool `json:"mcp"`

	// A2A is the A2A support status of the agent.
	A2A bool `json:"a2a"`

	// ENS is the ENS name of the agent.
	ENS string `json:"ens,omitempty"`

	// DID is the DID of the agent.
	DID string `json:"did,omitempty"`

	// WalletAddress is the wallet address of the agent.
	WalletAddress Address `json:"walletAddress,omitempty"`

	// SupportedTrusts is the supported trust models of the agent.
	SupportedTrusts []string `json:"supportedTrusts"`

	// A2ASkills is the A2A skills of the agent.
	A2ASkills []string `json:"a2aSkills"`

	// MCPTools is the MCP tools of the agent.
	MCPTools []string `json:"mcpTools"`

	// MCPPrompts is the MCP prompts of the agent.
	MCPPrompts []string `json:"mcpPrompts"`

	// MCPResources is the MCP resources of the agent.
	MCPResources []string `json:"mcpResources"`

	// Active is the active status of the agent.
	Active bool `json:"active"`

	// X402Support is the X402 support status of the agent.
	X402Support bool `json:"x402Support"`

	// Extras is the extras of the agent.
	Extras map[string]any `json:"extras"`
}

// Feedback is the feedback associated with an agent.
type Feedback struct {
	// ID is the ID of the feedback.
	ID FeedbackIDTuple `json:"id"`

	// AgentID is the ID of the agent.
	AgentID AgentID `json:"agentId"`

	// Reviewer is the address of the reviewer.
	Reviewer Address `json:"reviewer"`

	// Score is the score of the feedback.
	Score int64 `json:"score,omitempty"`

	// Tags is the tags of the feedback.
	Tags []string `json:"tags"`

	// Text is the text of the feedback.
	Text string `json:"text,omitempty"`

	// Context is the context of the feedback.
	Context map[string]any `json:"context,omitempty"`

	// ProofOfPayment is the proof of payment of the feedback.
	ProofOfPayment map[string]any `json:"proofOfPayment,omitempty"`

	// FileURI is the URI of the file of the feedback.
	FileURI URI `json:"fileUri,omitempty"`

	// CreatedAt is the timestamp of the creation of the feedback.
	CreatedAt Timestamp `json:"createdAt"`

	// Answers is the answers of the feedback.
	Answers []map[string]any `json:"answers"`

	// IsRevoked is the revoked status of the feedback.
	IsRevoked bool `json:"isRevoked"`

	// Off-chain only fields (not stored on chain)

	// Capability is the MCP capability associated with the feedback.
	Capability string `json:"capability,omitempty"`

	// Name is the MCP tool/resource name associated with the feedback.
	Name string `json:"name,omitempty"`

	// Skill is the A2A skill associated with the feedback.
	Skill string `json:"skill,omitempty"`

	// Task is the A2A task associated with the feedback.
	Task string `json:"task,omitempty"`
}

// FeedbackIDTuple is the tuple of the feedback ID.
type FeedbackIDTuple struct {
	// AgentID is the ID of the agent.
	AgentID AgentID `json:"agentId"`

	// ClientAddress is the address of the client.
	ClientAddress Address `json:"clientAddress"`

	// FeedbackIndex is the index of the feedback.
	FeedbackIndex int64 `json:"feedbackIndex"`
}

// FeedbackID is the ID of the feedback (agentID:clientAddress:feedbackIndex).
type FeedbackID string

// SearchParams is the search criteria for searching agents.
type SearchParams struct {
	// Chains is the chains to search (empty searches all chains).
	Chains []ChainID `json:"chains,omitempty"`

	// Name is the name search criteria (case-insensitive substring).
	Name string `json:"name,omitempty"`

	// Description is the description search criteria (semantic; vector distance > threshold).
	Description string `json:"description,omitempty"`

	// Owners is the owners addresses of the agent to search.
	Owners []Address `json:"owners,omitempty"`

	// Operators is the operators addresses of the agent to search.
	Operators []Address `json:"operators,omitempty"`

	// MCP is the MCP support status of the agent to search.
	MCP bool `json:"mcp,omitempty"`

	// A2A is the A2A support status of the agent to search.
	A2A bool `json:"a2a,omitempty"`

	// ENS is the ENS name of the agent to search (exact, case-insensitive).
	ENS string `json:"ens,omitempty"`

	// DID is the DID of the agent to search (exact).
	DID string `json:"did,omitempty"`

	// WalletAddress is the wallet address of the agent to search.
	WalletAddress Address `json:"walletAddress,omitempty"`

	// SupportedTrust is the supported trust models of the agent to search.
	SupportedTrust []TrustModel `json:"supportedTrust,omitempty"`

	// A2ASkills is the A2A skills of the agent to search.
	A2ASkills []string `json:"a2aSkills,omitempty"`

	// MCPTools is the MCP tools of the agent to search.
	MCPTools []string `json:"mcpTools,omitempty"`

	// MCPPrompts is the MCP prompts of the agent to search.
	MCPPrompts []string `json:"mcpPrompts,omitempty"`

	// MCPResources is the MCP resources of the agent to search.
	MCPResources []string `json:"mcpResources,omitempty"`

	// Active is the active status of the agent to search.
	Active bool `json:"active,omitempty"`

	// X402Support is the X402 support status of the agent to search.
	X402Support bool `json:"x402Support,omitempty"`
}

// SearchFeedbackParams is the search criteria for searching feedback.
type SearchFeedbackParams struct {
	// Agents is the agent IDs to search.
	Agents []AgentID `json:"agents,omitempty"`

	// Tags is the tags of the feedback to search.
	Tags []string `json:"tags,omitempty"`

	// Reviewers is the reviewer addresses to search.
	Reviewers []Address `json:"reviewers,omitempty"`

	// Capabilities is the MCP capabilities to search.
	Capabilities []string `json:"capabilities,omitempty"`

	// Skills is the A2A skills to search.
	Skills []string `json:"skills,omitempty"`

	// Tasks is the A2A tasks to search.
	Tasks []string `json:"tasks,omitempty"`

	// Names is the names of the MCP tools/resources/prompts to search.
	Names []string `json:"names,omitempty"`

	// MinScore is the minimum score to search (0-100).
	MinScore int64 `json:"minScore,omitempty"`

	// MaxScore is the maximum score to search (0-100).
	MaxScore int64 `json:"maxScore,omitempty"`

	// IncludeRevoked is the include revoked status to search.
	IncludeRevoked bool `json:"includeRevoked,omitempty"`
}

// SearchResultMeta is the metadata for multi-chain search results.
type SearchResultMeta struct {
	// Chains is the chains that were searched.
	Chains []ChainID `json:"chains"`

	// SuccessfulChains is the chains that were successful.
	SuccessfulChains []ChainID `json:"successfulChains"`

	// FailedChains is the chains that failed.
	FailedChains []ChainID `json:"failedChains"`

	// TotalResults is the total number of results.
	TotalResults int64 `json:"totalResults"`

	// Timing is the timing for multi-chain search results.
	Timing SearchResultMetaTiming `json:"timing"`
}

// SearchResultMetaTiming is the timing for multi-chain search results.
type SearchResultMetaTiming struct {
	// TotalMs is the total time in milliseconds.
	TotalMs int64 `json:"totalMs"`

	// AveragePerChainMs is the average time per chain in milliseconds.
	AveragePerChainMs int64 `json:"averagePerChainMs"`
}
