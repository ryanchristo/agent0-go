package core

import (
	"context"
	"net/http"

	graphql "github.com/hasura/go-graphql-client"

	"github.com/ryanchristo/agent0-go/sdk/types"
)

// SubgraphQueryOptions are the options for a subgraph query.
type SubgraphQueryOptions struct {
	Where                   map[string]any
	First                   int64
	Skip                    int64
	OrderBy                 string
	OrderDirection          OrderDirection
	IncludeRegistrationFile bool
}

// QueryAgent is the query agent response from a subgraph query.
type QueryAgent struct {
	ID               string                 `json:"id"`
	ChainID          types.ChainID          `json:"chainId"`
	AgentID          types.AgentID          `json:"agentId"`
	Owner            types.Address          `json:"owner"`
	Operators        []types.Address        `json:"operators"`
	AgentURI         types.URI              `json:"agentUri"`
	CreatedAt        int64                  `json:"createdAt"`
	UpdatedAt        int64                  `json:"updatedAt"`
	RegistrationFile types.RegistrationFile `json:"registrationFile,omitempty"`
}

// SubgraphClient is a client for querying the subgraph.
type SubgraphClient struct {
	client *graphql.Client
}

func NewSubgraphClient(subgraphURL string) *SubgraphClient {
	client := graphql.NewClient(subgraphURL, nil)

	client = client.WithRequestModifier(func(r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
	})

	return &SubgraphClient{
		client: client,
	}
}

// Query queries the subgraph with a given query and variables.
func (c *SubgraphClient) Query(query string, variables map[string]any) any {
	ctx := context.Background()
	return c.client.Query(ctx, query, variables)
}

// GetAgents queries the subgraph for agents with the given options.
func (c *SubgraphClient) GetAgents(options SubgraphQueryOptions) []types.AgentSummary {

	// TODO: implementation

	return []types.AgentSummary{}
}

// GetAgentByID queries the subgraph for a single agent by ID.
func (c *SubgraphClient) GetAgentByID(agentID types.AgentID) types.AgentSummary {

	// TODO: implementation

	return types.AgentSummary{}
}

// transformAgent transforms the raw subgraph agent into an agent summary.
func (c *SubgraphClient) transformAgent(agent QueryAgent) types.AgentSummary {

	// TODO: implementation

	return types.AgentSummary{}
}

// SearchAgents searches the subgraph for agents with the given parameters.
func (c *SubgraphClient) SearchAgents(
	params types.SearchParams,
	first int64,
	skip int64,
) []types.AgentSummary {
	// default first = 100
	// default skip = 0

	// TODO: implementation

	return []types.AgentSummary{}
}

// SearchFeedback searches the subgraph for feedback with the given parameters.
func (c *SubgraphClient) SearchFeedback(
	params SearchFeedbackParams,
	first int64,
	skip int64,
	orderBy string,
	orderDirection OrderDirection,
) []any {
	// default first = 100
	// default skip = 0
	// default orderBy = "createdAt"
	// default orderDirection = ORDER_DIRECTION_DESC

	// TODO: implementation

	return []any{}
}

// SearchAgentsByReputation searches the subgraph for agents by reputation with the given parameters.
func (c *SubgraphClient) SearchAgentsByReputation(
	agents []string,
	tags []string,
	reviewers []string,
	capabilities []string,
	skills []string,
	tasks []string,
	names []string,
	minAverageScore int64,
	includeRevoked bool,
	first int64,
	skip int64,
	orderBy string,
	orderDirection string,
) []SearchAgentsByReputationResult {
	// default includeRevoked = false
	// default first = 100
	// default skip = 0
	// default orderBy = "createdAt"
	// default orderDirection = ORDER_DIRECTION_DESC

	// TODO: implementation

	return []SearchAgentsByReputationResult{}
}

// ...

type OrderDirection string

const (
	ORDER_DIRECTION_ASC  OrderDirection = "asc"
	ORDER_DIRECTION_DESC OrderDirection = "desc"
)

type SearchFeedbackParams struct {
	Agents         []string
	Reviewers      []string
	Tags           []string
	Capabilities   []string
	Skills         []string
	Tasks          []string
	Names          []string
	MinScore       int64
	MaxScore       int64
	IncludeRevoked bool
}

type SearchAgentsByReputationResult struct {
	QueryAgent
	AverageScore int64
}
