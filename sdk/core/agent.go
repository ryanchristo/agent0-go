package core

import (
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ryanchristo/agent0-go/sdk/types"
)

// Agent is an agent instance for managing individual agents.
type Agent struct {
	registrationFile     types.RegistrationFile
	endpointCrawler      EndpointCrawler
	dirtyMetadata        map[string]bool
	lastRegisteredWallet types.Address
	lastRegisteredENS    string
}

// newAgent creates a new agent instance.
func newAgent(registrationFile types.RegistrationFile) *Agent {
	return &Agent{
		registrationFile: registrationFile,
		endpointCrawler:  NewEndpointCrawler(5000),
	}
}

// AgentID returns the agent ID.
func (a *Agent) AgentID() types.AgentID {
	return a.registrationFile.AgentID
}

// AgentURI returns the agent URI.
func (a *Agent) AgentURI() types.URI {
	return a.registrationFile.AgentURI
}

// Name returns the agent name.
func (a *Agent) Name() string {
	return a.registrationFile.Name
}

// Description returns the agent description.
func (a *Agent) Description() string {
	return a.registrationFile.Description
}

// Image returns the agent image URL.
func (a *Agent) Image() types.URI {
	return a.registrationFile.Image
}

// MCPEndpoint returns the agent MCP endpoint.
func (a *Agent) MCPEndpoint() types.URI {
	for _, endpoint := range a.registrationFile.Endpoints {
		if endpoint.Type == types.ENDPOINT_TYPE_MCP {
			return endpoint.Value
		}
	}
	return ""
}

// A2AEndpoint returns the agent A2A endpoint.
func (a *Agent) A2AEndpoint() types.URI {
	for _, endpoint := range a.registrationFile.Endpoints {
		if endpoint.Type == types.ENDPOINT_TYPE_A2A {
			return endpoint.Value
		}
	}
	return ""
}

// ENSEndpoint returns the agent ENS endpoint.
func (a *Agent) ENSEndpoint() types.URI {
	for _, endpoint := range a.registrationFile.Endpoints {
		if endpoint.Type == types.ENDPOINT_TYPE_ENS {
			return endpoint.Value
		}
	}
	return ""
}

// WalletAddress returns the agent wallet address.
func (a *Agent) WalletAddress() types.Address {
	return a.registrationFile.WalletAddress
}

// MCPTools returns the agent MCP tools.
func (a *Agent) MCPTools() []string {
	for _, endpoint := range a.registrationFile.Endpoints {
		if endpoint.Type == types.ENDPOINT_TYPE_MCP {
			return endpoint.Meta["mcpTools"].([]string)
		}
	}
	return []string{}
}

// MCPPrompts returns the agent MCP prompts.
func (a *Agent) MCPPrompts() []string {
	for _, endpoint := range a.registrationFile.Endpoints {
		if endpoint.Type == types.ENDPOINT_TYPE_MCP {
			return endpoint.Meta["mcpPrompts"].([]string)
		}
	}
	return []string{}
}

// MCPResources returns the agent MCP resources.
func (a *Agent) MCPResources() []string {
	for _, endpoint := range a.registrationFile.Endpoints {
		if endpoint.Type == types.ENDPOINT_TYPE_MCP {
			return endpoint.Meta["mcpResources"].([]string)
		}
	}
	return []string{}
}

// A2ASkills returns the agent A2A skills.
func (a *Agent) A2ASkills() []string {
	for _, endpoint := range a.registrationFile.Endpoints {
		if endpoint.Type == types.ENDPOINT_TYPE_A2A {
			return endpoint.Meta["a2aSkills"].([]string)
		}
	}
	return []string{}
}

// Endpoint management

// SetMCP sets the agent MCP endpoint.
func (a *Agent) SetMCP(endpoint types.URI, version string, autoFetch bool) *Agent {
	// default version = "2025-06-18"
	// default autoFetch = true

	// TODO: implementation

	return a
}

// SetA2A sets the agent A2A endpoint.
func (a *Agent) SetA2A(agentcard, version string, autoFetch bool) *Agent {
	// default version = "0.30"
	// default autoFetch = true

	// TODO: implementation

	return a
}

// SetENS sets the agent ENS endpoint.
func (a *Agent) SetENS(name, version string) *Agent {
	// default version = "1.0"

	// TODO: implementation

	return a
}

// getOrCreateOASFEndpoint gets or creates the OASF endpoint.
func (a *Agent) getOrCreateOASFEndpoint() types.Endpoint {

	// TODO: implementation

	return types.Endpoint{}
}

// AddSkills adds a skill to the OASF endpoint.
func (a *Agent) AddSkill(slug string, validateOASF bool) *Agent {
	// default validateOASF = false

	// TODO: implementation

	return a
}

// RemoveSkill removes a skill from the OASF endpoint.
func (a *Agent) RemoveSkill(slug string) *Agent {

	// TODO: implementation

	return a
}

// AddDomain adds a domain to the OASF endpoint.
func (a *Agent) AddDomain(slug string, validateOASF bool) *Agent {
	// default validateOASF = false

	// TODO: implementation

	return a
}

// RemoveDomain removes a domain from the OASF endpoint.
func (a *Agent) RemoveDomain(slug string) *Agent {

	// TODO: implementation

	return a
}

// SetAgentWallet sets the agent wallet address and the associated chain ID.
func (a *Agent) SetAgentWallet(address types.Address, chainID int64) *Agent {

	// TODO: implementation

	return a
}

// SetActive sets the active status of the agent.
func (a *Agent) SetActive(active bool) *Agent {

	// TODO: implementation

	return a
}

// SetX402Support sets the X402 support status of the agent.
func (a *Agent) SetX402Support(x402Support bool) *Agent {

	// TODO: implementation

	return a
}

// SetTrust sets the trust models of the agent.
func (a *Agent) SetTrust(reputation, cryptoEconomics, teeAttestation bool) *Agent {
	// default reputation = false
	// default cryptoEconomics = false
	// default teeAttestation = false

	// TODO: implementation

	return a
}

// SetMetadata sets the metadata of the agent.
func (a *Agent) SetMetadata(kv map[string]any) *Agent {

	// TODO: implementation

	return a
}

// GetMetadata gets the metadata of the agent.
func (a *Agent) GetMetadata(kv map[string]any) map[string]any {

	// TODO: implementation

	return nil
}

// DelMetadata deletes metadata from the agent.
func (a *Agent) DelMetadata(key string) *Agent {

	// TODO: implementation

	return a
}

// GetRegistrationFile gets the registration file of the agent.
func (a *Agent) GetRegistrationFile() types.RegistrationFile {
	return a.registrationFile
}

// UpdateInfo updates the information of the agent.
func (a *Agent) UpdateInfo(name, description string, image types.URI) *Agent {

	// TODO: implementation

	return a
}

// RegisterIPFS registers the agent on chain using the IPFS workflow.
func (a *Agent) RegisterIPFS() types.RegistrationFile {

	// TODO: implementation

	return a.registrationFile
}

// RegisterHTTP registers the agent on chain using the HTTP workflow.
func (a *Agent) RegisterHTTP(agentURI types.URI) types.RegistrationFile {

	// TODO: implementation

	return a.registrationFile
}

// SetAgentURI sets the agent URI (used for updating the agent).
func (a *Agent) SetAgentURI(agentURI types.URI) {

	// TODO: implementation

}

// Transfer transfers the agent ownership to a new owner.
func (a *Agent) Transfer(newOwner types.Address) TransferResult {

	// TODO: implementation

	return TransferResult{}
}

// Private helper methods

// registerWithoutURI registers the agent without a URI.
func (a *Agent) registerWithoutURI() {

	// TODO: implementation

}

// registerWithURI registers the agent with a URI.
func (a *Agent) registerWithURI(agentURI types.URI) types.RegistrationFile {

	// TODO: implementation

	return a.registrationFile
}

// updateMetadataOnChain updates the metadata of the agent on chain.
func (a *Agent) updateMetadataOnChain() {

	// TODO: implementation

}

// collectMetadataForRegistration collects the metadata for registration.
func (a *Agent) collectMetadataForRegistration() map[string][]byte {

	// TODO: implementation

	return nil
}

// extractAgentIDFromReceipt extracts the agent ID from the receipt.
func (a *Agent) extractAgentIDFromReceipt(receipt ethtypes.Receipt) big.Int {

	// TODO: implementation

	return big.Int{}
}

// ...

type TransferResult struct {
	TXHash  string
	From    types.Address
	To      types.Address
	AgentID types.AgentID
}
