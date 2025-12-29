package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	graphql "github.com/hasura/go-graphql-client"

	"github.com/ryanchristo/agent0-go/sdk/subgraph/model"
	"github.com/ryanchristo/agent0-go/sdk/types"
)

// SubgraphQueryOptions are the options for a subgraph query.
type SubgraphQueryOptions struct {
	Where                   map[string]any
	First                   int64
	Skip                    int64
	OrderBy                 string
	OrderDirection          OrderDirection
	IncludeRegistrationFile *bool
}

// QueryAgent is the query agent response from a subgraph query.
type QueryAgent struct {
	ID               string                      `json:"id"`
	ChainID          types.ChainID               `json:"chainId"`
	AgentID          types.AgentID               `json:"agentId"`
	Owner            types.Address               `json:"owner"`
	Operators        []types.Address             `json:"operators"`
	AgentURI         types.URI                   `json:"agentUri"`
	CreatedAt        int64                       `json:"createdAt"`
	UpdatedAt        int64                       `json:"updatedAt"`
	RegistrationFile model.AgentRegistrationFile `json:"registrationFile"`
}

// SubgraphClient is a client for querying the subgraph.
type SubgraphClient struct {
	client *graphql.Client
}

// NewSubgraphClient creates a new subgraph client.
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
	return c.client.Query(context.Background(), query, variables)
}

// GetAgents queries the subgraph for agents with the given options.
func (c *SubgraphClient) GetAgents(options SubgraphQueryOptions) []types.AgentSummary {
	if options.Where == nil {
		options.Where = make(map[string]any)
	}
	if options.First == 0 {
		options.First = 100
	}
	if options.OrderBy == "" {
		options.OrderBy = "createdAt"
	}
	if options.OrderDirection == "" {
		options.OrderDirection = ORDER_DIRECTION_DESC
	}
	if options.IncludeRegistrationFile == nil {
		*options.IncludeRegistrationFile = true
	}

	// Support Agent-level filters and nested registration file filters
	supportedWhere := map[string]any{}
	if v, ok := options.Where["agentId"]; ok {
		supportedWhere["agentId"] = v
	}
	if v, ok := options.Where["owner"]; ok {
		supportedWhere["owner"] = v
	}
	if v, ok := options.Where["owner_in"]; ok {
		supportedWhere["owner_in"] = v
	}
	if v, ok := options.Where["operators_contains"]; ok {
		supportedWhere["operators_contains"] = v
	}
	if v, ok := options.Where["agentURI"]; ok {
		supportedWhere["agentURI"] = v
	}
	if v, ok := options.Where["registrationFile_not"]; ok {
		supportedWhere["registrationFile_not"] = v
	}

	// Support nested registration file filters (pushed to subgraph level)
	// Note: Python SDK uses "registrationFile_" (with underscore) for nested filters
	if v, ok := options.Where["registrationFile"]; ok {
		supportedWhere["registrationFile_"] = v
	}
	if v, ok := options.Where["registrationFile_"]; ok {
		supportedWhere["registrationFile_"] = v
	}

	// Build where clause with support for nested filters
	whereClause := ""
	if len(supportedWhere) > 0 {
		conditions := []string{}
		for k, v := range supportedWhere {
			isRegistrationFile := (k == "registrationFile" || k == "registrationFile_")
			if vMap, ok := v.(map[string]any); ok && isRegistrationFile {
				nestedConditions := []string{}
				for nk, nv := range vMap {
					if b, ok := nv.(bool); ok {
						nestedConditions = append(nestedConditions, fmt.Sprintf("%s: %v", nk, b))
					} else if s, ok := nv.(string); ok {
						nestedConditions = append(nestedConditions, fmt.Sprintf("%s: %s", nk, s))
					} else if nv == nil {
						if strings.HasSuffix(nk, "_not") {
							nestedConditions = append(nestedConditions, fmt.Sprintf("%s: null", nk))
						} else {
							nestedConditions = append(nestedConditions, fmt.Sprintf("%s_not: null", nk))
						}
					}
				}
				if len(nestedConditions) > 0 {
					conditions = append(
						conditions,
						fmt.Sprintf("registrationFile_: { %s }", strings.Join(nestedConditions, ", ")),
					)
				}
			} else if b, ok := v.(bool); ok {
				conditions = append(conditions, fmt.Sprintf("%s: %v", k, b))
			} else if s, ok := v.(string); ok {
				conditions = append(conditions, fmt.Sprintf("%s: %s", k, s))
			} else if i, ok := v.(int); ok {
				conditions = append(conditions, fmt.Sprintf("%s: %d", k, i))
			} else if a, ok := v.([]any); ok {
				conditions = append(conditions, fmt.Sprintf("%s: %v", k, a))
			} else if v == nil {
				if strings.HasSuffix(k, "_not") {
					conditions = append(conditions, fmt.Sprintf("%s: null", k))
				} else {
					conditions = append(conditions, fmt.Sprintf("%s_not: null", k))
				}
			}
		}
		if len(conditions) > 0 {
			whereClause = fmt.Sprintf("where: { \"%s\" }", strings.Join(conditions, ", "))
		}
	}

	// Build registration file fragment
	regFileFragment := ""
	if *options.IncludeRegistrationFile {
		regFileFragment = `
			registrationFile {
				id
				agentId
				name
				description
				image
				active
				x402support
				supportedTrusts
				mcpEndpoint
				mcpVersion
				a2aEndpoint
				a2aVersion
				ens
				did
				agentWallet
				agentWalletChainId
				mcpTools
				mcpPrompts
				mcpResources
				a2aSkills
			}
		`
	}

	query := fmt.Sprintf(`
		query GetAgents($first: Int!, $skip: Int!, $orderBy: Agent_orderBy!, $orderDirection: OrderDirection!) {
			agents(
				%s
				first: $first
				skip: $skip
				orderBy: $orderBy
				orderDirection: $orderDirection
			) {
				id
				chainId
				agentId
				owner
				operators
				agentURI
				createdAt
				updatedAt
				%s
			}
		}
	`, whereClause, regFileFragment)

	variables := map[string]any{
		"first":          options.First,
		"skip":           options.Skip,
		"orderBy":        options.OrderBy,
		"orderDirection": options.OrderDirection,
	}

	data := c.Query(query, variables)

	if data, ok := data.(map[string]any); ok {
		if agents, ok := data["agents"].([]any); ok {
			agentSummaries := make([]types.AgentSummary, 0, len(agents))
			for _, agent := range agents {
				agentSummary := c.transformAgent(agent.(QueryAgent))
				agentSummaries = append(agentSummaries, agentSummary)
			}
			return agentSummaries
		}
	}

	log.Fatal("Failed to get agents from subgraph")

	return []types.AgentSummary{}
}

// GetAgentByID queries the subgraph for a single agent by ID.
func (c *SubgraphClient) GetAgentByID(agentID types.AgentID) types.AgentSummary {
	query := `
		query GetAgent($agentId: String!) {
			agent(id: $agentId) {
				id
				chainId
				agentId
				owner
				operators
				agentURI
				createdAt
				updatedAt
				registrationFile {
					id
					agentId
					name
					description
					image
					active
					x402support
					supportedTrusts
					mcpEndpoint
					mcpVersion
					a2aEndpoint
					a2aVersion
					ens
					did
					agentWallet
					agentWalletChainId
					mcpTools
					mcpPrompts
					mcpResources
					a2aSkills
				}
			}
		}
	`

	data := c.Query(query, map[string]any{"agentId": agentID})
	if d, ok := data.(map[string]any); ok {
		if a, ok := d["agent"].(QueryAgent); ok {
			return c.transformAgent(a)
		}
	}

	log.Fatal("Failed to get agent from subgraph")

	return types.AgentSummary{}
}

// transformAgent transforms the raw subgraph agent into an agent summary.
func (c *SubgraphClient) transformAgent(agent QueryAgent) types.AgentSummary {
	return types.AgentSummary{
		ChainID:         agent.ChainID,
		AgentID:         agent.AgentID,
		Name:            *agent.RegistrationFile.Name,
		Description:     *agent.RegistrationFile.Description,
		Image:           *agent.RegistrationFile.Image,
		Owners:          []types.Address{agent.Owner},
		Operators:       agent.Operators,
		MCP:             *agent.RegistrationFile.McpEndpoint != "",
		A2A:             *agent.RegistrationFile.A2aEndpoint != "",
		ENS:             *agent.RegistrationFile.Ens,
		DID:             *agent.RegistrationFile.Did,
		WalletAddress:   *agent.RegistrationFile.AgentWallet,
		SupportedTrusts: agent.RegistrationFile.SupportedTrusts,
		A2ASkills:       agent.RegistrationFile.A2aSkills,
		MCPTools:        agent.RegistrationFile.McpTools,
		MCPPrompts:      agent.RegistrationFile.McpPrompts,
		MCPResources:    agent.RegistrationFile.McpResources,
		Active:          *agent.RegistrationFile.Active,
		X402Support:     *agent.RegistrationFile.X402support,
		Extras:          map[string]any{},
	}
}

// SearchAgents searches the subgraph for agents with the given parameters.
func (c *SubgraphClient) SearchAgents(params types.SearchParams, first, skip int64) []types.AgentSummary {
	if first == 0 {
		first = 100
	}

	where := map[string]any{
		"registrationFile_not": nil, // only get agents with registration files
	}

	// Note: Most search fields are in registration file, so we need to filter after fetching
	// For now, we'll do basic filtering on Agent fields and then filter on registrationFile fields
	if params.Active != nil || params.MCP != nil || params.A2A != nil || params.X402Support != nil ||
		params.ENS != "" || params.WalletAddress != "" || params.SupportedTrust != nil || params.A2ASkills != nil ||
		params.MCPTools != nil || params.Name != "" || params.Owners != nil || params.Operators != nil {
		// Push basic filters to subgraph using nested registrationFile filters
		registrationFileFilters := map[string]any{}
		if params.Active != nil {
			registrationFileFilters["active"] = *params.Active
		}
		if params.X402Support != nil {
			registrationFileFilters["x402support"] = *params.X402Support
		}
		if params.ENS != "" {
			registrationFileFilters["ens"] = strings.ToLower(params.ENS)
		}
		if params.WalletAddress != "" {
			registrationFileFilters["agentWallet"] = strings.ToLower(params.WalletAddress)
		}
		if params.MCP != nil {
			if *params.MCP {
				registrationFileFilters["mcpEndpoint_not"] = nil
			} else {
				registrationFileFilters["mcpEndpoint"] = nil
			}
		}
		if params.A2A != nil {
			if *params.A2A {
				registrationFileFilters["a2aEndpoint_not"] = nil
			} else {
				registrationFileFilters["a2aEndpoint"] = nil
			}
		}

		whereWithFilters := map[string]any{}
		if len(registrationFileFilters) > 0 {
			// Python SDK uses "registrationFile_" (with underscore) for nested filters
			whereWithFilters["registrationFile_"] = registrationFileFilters
		}

		// Owner filtering (at Agent level, not registrationFile)
		if len(params.Owners) > 0 {
			// Normalize addresses to lowercase for case-insensitive matching
			normalizedOwners := make([]string, len(params.Owners))
			for i, owner := range params.Owners {
				normalizedOwners[i] = strings.ToLower(string(owner))
			}
			if len(normalizedOwners) == 1 {
				whereWithFilters["owner"] = normalizedOwners[0]
			} else {
				whereWithFilters["owner_in"] = normalizedOwners
			}
		}

		// Operator filtering (at Agent level, not registrationFile)
		if len(params.Operators) > 0 {
			// Normalize addresses to lowercase for case-insensitive matching
			normalizedOperators := make([]string, len(params.Operators))
			for i, operator := range params.Operators {
				normalizedOperators[i] = strings.ToLower(string(operator))
			}
			// For operators (array field), use contains to check if any operator matches
			whereWithFilters["operators_contains"] = normalizedOperators
		}

		// Fetch records with filters and pagination applied at subgraph level
		allAgents := c.GetAgents(SubgraphQueryOptions{
			Where: whereWithFilters,
			First: first,
			Skip:  skip,
		})

		// Only filter client-side for fields that can't be filtered at subgraph level
		// Fields already filtered at subgraph level: active, x402support, mcp, a2a, ens, walletAddress, owners, operators
		filteredAgents := make([]types.AgentSummary, 0, len(allAgents))
		for _, agent := range allAgents {
			// Name filtering (substring search - not supported at subgraph level)
			if params.Name != "" && !strings.Contains(strings.ToLower(agent.Name), strings.ToLower(params.Name)) {
				continue
			}
			// Array contains filtering (supportedTrust, a2aSkills, mcpTools) - these require array contains logic
			if len(params.SupportedTrust) > 0 {
				hasAllTrusts := false
				for _, trust := range agent.SupportedTrusts {
					if slices.Contains(params.SupportedTrust, types.TrustModel(trust)) {
						hasAllTrusts = true
						break
					}
				}
				if !hasAllTrusts {
					continue
				}
			}
			if len(params.A2ASkills) > 0 {
				hasAllSkills := false
				for _, skill := range agent.A2ASkills {
					if slices.Contains(params.A2ASkills, skill) {
						hasAllSkills = true
						break
					}
				}
				if !hasAllSkills {
					continue
				}
			}
			if len(params.MCPTools) > 0 {
				hasAllTools := false
				for _, tool := range agent.MCPTools {
					if slices.Contains(params.MCPTools, tool) {
						hasAllTools = true
						break
					}
				}
				if !hasAllTools {
					continue
				}
			}
			filteredAgents = append(filteredAgents, agent)
		}
		return filteredAgents
	}

	return c.GetAgents(SubgraphQueryOptions{
		Where: where,
		First: first,
		Skip:  skip,
	})
}

// SearchFeedback searches the subgraph for feedback with the given parameters.
func (c *SubgraphClient) SearchFeedback(
	params SearchFeedbackParams,
	first int64,
	skip int64,
	orderBy string,
	orderDirection OrderDirection,
) []any {
	if first == 0 {
		first = 100
	}
	if orderBy == "" {
		orderBy = "createdAt"
	}
	if orderDirection == "" {
		orderDirection = ORDER_DIRECTION_DESC
	}

	// Build where clause from params
	whereConditions := []string{}

	if len(params.Agents) > 0 {
		agentIDs := make([]string, len(params.Agents))
		for i, agent := range params.Agents {
			agentIDs[i] = string(agent)
		}
		agentIDsString := strings.Join(agentIDs, ", ")
		whereConditions = append(whereConditions, fmt.Sprintf("agent_in: [%s]", agentIDsString))
	}

	if len(params.Reviewers) > 0 {
		reviewerAddresses := make([]string, len(params.Reviewers))
		for i, reviewer := range params.Reviewers {
			reviewerAddresses[i] = string(reviewer)
		}
		reviewerAddressesString := strings.Join(reviewerAddresses, ", ")
		whereConditions = append(whereConditions, fmt.Sprintf("clientAddress_in: [%s]", reviewerAddressesString))
	}

	if !params.IncludeRevoked {
		whereConditions = append(whereConditions, "isRevoked: false")
	}

	// Build all non-tag conditions first
	nonTagConditions := []string{}
	for _, condition := range whereConditions {
		nonTagConditions = append(nonTagConditions, condition)
	}

	// Hanle tag filtering separately - it needs to be at the top level
	var tagFilterCondition string
	if len(params.Tags) > 0 {
		// Tag search: any of the tags must match in tag1 OR tag2
		// Build complete condition with all filters for each tag alternative
		tagWhereItems := []string{}
		for _, tag := range params.Tags {
			// For tag1 match
			allConditionsTag1 := append(nonTagConditions, fmt.Sprintf("tag1: \"%s\"", tag))
			tagWhereItems = append(tagWhereItems, fmt.Sprintf("{ %s }", strings.Join(allConditionsTag1, ", ")))
			// For tag2 match
			allConditionsTag2 := append(nonTagConditions, fmt.Sprintf("tag2: \"%s\"", tag))
			tagWhereItems = append(tagWhereItems, fmt.Sprintf("{ %s }", strings.Join(allConditionsTag2, ", ")))
		}
		// Join all tag alternatives
		tagFilterCondition = strings.Join(tagWhereItems, ", ")
	}

	if params.MinScore != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("score_gte: %d", *params.MinScore))
	}

	if params.MaxScore != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("score_lte: %d", *params.MaxScore))
	}

	// Feedback file filters
	feedbackFileFilters := []string{}

	if len(params.Capabilities) > 0 {
		feedbackFileFilters = append(feedbackFileFilters, fmt.Sprintf("capability_in: [%s]", strings.Join(params.Capabilities, ", ")))
	}

	if len(params.Skills) > 0 {
		feedbackFileFilters = append(feedbackFileFilters, fmt.Sprintf("skill_in: [%s]", strings.Join(params.Skills, ", ")))
	}

	if len(params.Tasks) > 0 {
		feedbackFileFilters = append(feedbackFileFilters, fmt.Sprintf("task_in: [%s]", strings.Join(params.Tasks, ", ")))
	}

	if len(params.Names) > 0 {
		feedbackFileFilters = append(feedbackFileFilters, fmt.Sprintf("name_in: [%s]", strings.Join(params.Names, ", ")))
	}

	if len(feedbackFileFilters) > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("feedbackFile_: { %s }", strings.Join(feedbackFileFilters, ", ")))
	}

	// Use tag_filter_condition if tags were provided, otherwise use standard where clause
	var whereClause string
	if tagFilterCondition != "" {
		whereClause = fmt.Sprintf("where: { or: [%s] }", tagFilterCondition)
	} else {
		whereClause = fmt.Sprintf("where: { %s }", strings.Join(whereConditions, ", "))
	}

	query := fmt.Sprintf(`
	  {
	    feedbacks(
		  %s
		  first: %d
		  skip: %d
		  orderBy: %s
		  orderDirection: %s
		) {
		  id
		  agent { id agentId chainId }
		  clientAddress
		  score
		  tag1
		  tag2
		  feedbackUri
		  feedbackURIType
		  feedbackHash
		  isRevoked
		  createdAt
		  revokedAt
		  feedbackFile {
			id
			feedbackId
			text
			capability
			name
			skill
			task
			context
			proofOfPaymentFromAddress
			proofOfPaymentToAddress
			proofOfPaymentChainId
			proofOfPaymentTxHash
			tag1
			tag2
			createdAt
		  }
		  responses {
			id
			responder
			responseUri
			responseHash
			createdAt
		  }
		}
	  }`, whereClause, first, skip, orderBy, orderDirection)

	result := c.Query(query, map[string]any{})
	return result.(map[string]any)["feedbacks"].([]any)
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
	minAverageScore *int64,
	includeRevoked bool,
	first int64,
	skip int64,
	orderBy string,
	orderDirection OrderDirection,
) []SearchAgentsByReputationResult {
	if first == 0 {
		first = 100
	}
	if orderBy == "" {
		orderBy = "createdAt"
	}
	if orderDirection == "" {
		orderDirection = ORDER_DIRECTION_DESC
	}

	// Build feedback filters
	feedbackFilters := []string{}

	if !includeRevoked {
		feedbackFilters = append(feedbackFilters, "isRevoked: false")
	}

	if len(tags) > 0 {
		for _, tag := range tags {
			feedbackFilters = append(feedbackFilters, fmt.Sprintf("{or: [{tag1: \"%s\"}, {tag2: \"%s\"}]}", tag, tag))
		}
		feedbackFilters = append(feedbackFilters, fmt.Sprintf("or: [%s]", strings.Join(tags, ", ")))
	}

	if len(reviewers) > 0 {
		reviewerAddresses := make([]string, len(reviewers))
		for i, reviewer := range reviewers {
			reviewerAddresses[i] = string(reviewer)
		}
		reviewerAddressesString := strings.Join(reviewerAddresses, ", ")
		feedbackFilters = append(feedbackFilters, fmt.Sprintf("clientAddress_in: [%s]", reviewerAddressesString))
	}

	// Feedback file filters
	feedbackFileFilters := []string{}

	if len(capabilities) > 0 {
		feedbackFileFilters = append(feedbackFileFilters, fmt.Sprintf("capability_in: [%s]", strings.Join(capabilities, ", ")))
	}

	if len(skills) > 0 {
		feedbackFileFilters = append(feedbackFileFilters, fmt.Sprintf("skill_in: [%s]", strings.Join(skills, ", ")))
	}

	if len(tasks) > 0 {
		feedbackFileFilters = append(feedbackFileFilters, fmt.Sprintf("task_in: [%s]", strings.Join(tasks, ", ")))
	}

	if len(names) > 0 {
		feedbackFileFilters = append(feedbackFileFilters, fmt.Sprintf("name_in: [%s]", strings.Join(names, ", ")))
	}

	if len(feedbackFileFilters) > 0 {
		feedbackFilters = append(feedbackFilters, fmt.Sprintf("feedbackFile_: { %s }", strings.Join(feedbackFileFilters, ", ")))
	}

	// If we have feedback filters, first query feedback to get agent IDs
	agentWhere := ""
	if len(tags) > 0 || len(capabilities) > 0 || len(skills) > 0 || len(tasks) > 0 || len(names) > 0 || len(reviewers) > 0 {
		feedbackWhere := ""
		if len(feedbackFilters) > 0 {
			feedbackWhere = fmt.Sprintf("{ %s }", strings.Join(feedbackFilters, ", "))
		} else {
			feedbackWhere = "{}"
		}

		feedbackQuery := fmt.Sprintf(`
			{
				feedbacks(
					where: %s
					first: 1000
					skip: 0
				) {
					agent { id }
				}
			}
		`, feedbackWhere)

		feedbackResult := c.Query(feedbackQuery, map[string]any{})
		feedbackData := feedbackResult.(map[string]any)["feedbacks"].([]any)

		// Extract unique agnet IDs
		agentIDsSet := make(map[string]bool)
		for _, feedback := range feedbackData {
			if feedback, ok := feedback.(map[string]any); ok {
				if agent, ok := feedback["agent"].(map[string]any); ok {
					if id, ok := agent["id"].(string); ok {
						agentIDsSet[id] = true
					}
				}
			}
		}

		if len(agentIDsSet) == 0 {
			// No agents have matching feedback
			return []SearchAgentsByReputationResult{}
		}

		// Apply agent filter if specified
		agentIDsList := make([]string, 0, len(agentIDsSet))
		if len(agents) > 0 {
			for _, agent := range agents {
				if _, ok := agentIDsSet[string(agent)]; ok {
					agentIDsList = append(agentIDsList, string(agent))
				}
			}
			if len(agentIDsList) == 0 {
				// If feedback query fails, return empty
				return []SearchAgentsByReputationResult{}
			}
		}

		agentIDsString := ""
		for i, agentID := range agentIDsList {
			agentIDsString += fmt.Sprintf("\"%s\"", agentID)
			if i < len(agentIDsList)-1 {
				agentIDsString += ", "
			}
		}
		agentWhere = fmt.Sprintf("where: { id_in: [%s] }", agentIDsString)

	} else {
		// No feedback filters = query agents directly
		agentFilters := []string{}
		if len(agents) > 0 {
			agentIDs := ""
			for i, agent := range agents {
				agentIDs += fmt.Sprintf("\"%s\"", agent)
				if i < len(agents)-1 {
					agentIDs += ", "
				}
			}
			agentFilters = append(agentFilters, fmt.Sprintf("id_in: [%s]", agentIDs))
		}

		if len(agentFilters) > 0 {
			agentWhere = fmt.Sprintf("where: { %s }", strings.Join(agentFilters, ", "))
		}
	}

	// Build feedback where for agent query (to calculate scores)
	feedbackWhereForAgents := ""
	if len(feedbackFilters) > 0 {
		feedbackWhereForAgents = fmt.Sprintf("{ %s }", strings.Join(feedbackFilters, ", "))
	} else {
		feedbackWhereForAgents = "{}"
	}

	query := fmt.Sprintf(`
		{
			agents(
				%s
				first: %d
				skip: %d
				orderBy: %s
				orderDirection: %s
			) {
				id
				chainId
				agentId
				agentURI
				agentURIType
				owner
				operators
				createdAt
				updatedAt
				totalFeedback
				lastActivity
				registrationFile {
					id
					name
					description
					image
					active
					x402support
					supportedTrusts
					mcpEndpoint
					mcpVersion
					a2aEndpoint
					a2aVersion
					ens
					did
					agentWallet
					agentWalletChainId
					mcpTools
					mcpPrompts
					a2aSkills
					createdAt
				}
				feedback(
					where: { %s }
				) {
					score
					isRevoked
					feedbackFile {
						capability
						skill
						task
						name
					}
				}
			}
		}
	`, agentWhere, first, skip, orderBy, orderDirection, feedbackWhereForAgents)

	result := c.Query(query, map[string]any{})
	agentsResult := result.(map[string]any)["agents"].([]any)

	// Calculate agerage scores
	agentsWithScores := []SearchAgentsByReputationResult{}
	for _, agent := range agentsResult {
		feedbacks := []map[string]any{}
		averageScore := new(int64)

		if len(feedbacks) > 0 {
			scores := []int64{}
			for _, feedback := range feedbacks {
				if score, ok := feedback["score"].(int64); ok && score > 0 {
					scores = append(scores, feedback["score"].(int64))
				}
			}
			if len(scores) > 0 {
				for _, score := range scores {
					*averageScore += score
				}
				*averageScore /= int64(len(scores))
			}
		}

		// Remove feedback array from result (not part of QueryAgent)
		delete(agent.(map[string]any), "feedback")
		agentsWithScores = append(agentsWithScores, SearchAgentsByReputationResult{
			QueryAgent:   agent.(QueryAgent),
			AverageScore: averageScore,
		})
	}

	// Filter by minAverageScore
	filteredAgents := agentsWithScores
	if minAverageScore != nil {
		filteredAgents = []SearchAgentsByReputationResult{}
		for _, agent := range agentsWithScores {
			if agent.AverageScore != nil && *agent.AverageScore >= *minAverageScore {
				filteredAgents = append(filteredAgents, agent)
			}
		}
	}

	return filteredAgents
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
	MinScore       *int64
	MaxScore       *int64
	IncludeRevoked bool
}

type SearchAgentsByReputationResult struct {
	QueryAgent
	AverageScore *int64
}
