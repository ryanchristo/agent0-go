package core

import (
	"github.com/ryanchristo/agent0-go/sdk/types"
	"github.com/ryanchristo/agent0-go/sdk/utils"
)

// MCPCapabilities is the MCP capabilities of an agent.
type MCPCapabilities struct {
	MCPTools     []string
	MCPPrompts   []string
	MCPResources []string
}

// A2ACapabilities is the A2A capabilities of an agent.
type A2ACapabilities struct {
	A2ASkills []string
}

// CreateJSONRPCRequest creates a JSON-RPC request.
func CreateJSONRPCRequest(method string, params map[string]any, requestID int64) JSONRPCRequest {
	if requestID == 0 {
		requestID = 1
	}
	return JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		ID:      requestID,
		Params:  params,
	}
}

// EndpointCrawler is the endpoint crawler for an agent.
type EndpointCrawler struct {
	timeout int64 // milliseconds
}

// NewEndpointCrawler creates a new endpoint crawler.
func NewEndpointCrawler(timeout int64) EndpointCrawler {
	if timeout == 0 {
		timeout = utils.TIMEOUTS["ENDPOINT_CRAWLER_DEFAULT"]
	}
	return EndpointCrawler{
		timeout: timeout,
	}
}

// FetchMCPCapabilities fetches the MCP capabilities from an MCP ednpoint.
func (e *EndpointCrawler) FetchMCPCapabilities(endpoint types.URI) MCPCapabilities {

	// TODO: implementation

	return MCPCapabilities{}
}

// fetchViaJSONRPC fetches the MCP capabilities via JSON-RPC from an HTTP endpoint.
func (e *EndpointCrawler) fetchViaJSONRPC(httpURL types.URI) MCPCapabilities {

	// TODO: implementation

	return MCPCapabilities{}
}

// jsonRPCCall makes a JSON-RPC call to the given URL.
func (e *EndpointCrawler) jsonRPCCall(url types.URI, method string, params map[string]any) any {

	// TODO: implementation

	return nil
}

// parseSSEResponse parses the Server-Sent Events (SSE) response.
func (e *EndpointCrawler) parseSSEResponse(sseText string) any {

	// TODO: implementation

	return nil
}

// fetchA2ACapabilities fetches the A2A capabilities from an A2A endpoint.
func (e *EndpointCrawler) fetchA2ACapabilities(endpoint types.URI) A2ACapabilities {

	// TODO: implementation

	return A2ACapabilities{}
}

// extractA2ASkills extracts the A2A skills from the A2A agent card.
func (e *EndpointCrawler) extractA2ASkills(data any) []string {

	// TODO: implementation

	return []string{}
}

// extractList extracts a list of strings from the nested JSON data.
func (e *EndpointCrawler) extractList(data any, key string) []string {

	// TODO: implementation

	return []string{}
}

// ...

type JSONRPCRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	ID      int64          `json:"id"`
	Params  map[string]any `json:"params"`
}
