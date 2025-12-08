package core

import (
	"math/big"

	"github.com/ryanchristo/agent0-go/sdk/types"
	"github.com/ryanchristo/agent0-go/sdk/utils"
)

// FeedbackAuth contains the authorization information for feedback.
type FeedbackAuth struct {
	AgentID          big.Int
	ClientAddress    types.Address
	IndexLimit       big.Int
	Expiry           big.Int
	ChainID          big.Int
	IdentityRegistry types.Address
	SignerAddress    types.Address
}

// FeedbackManager is the manager for feedback operations.
type FeedbackManager struct {
	web3Client         *Web3Client
	ipfsClient         *IPFSClient
	reputationRegistry *Contract
	identityRegistry   *Contract
	subgraphClient     *SubgraphClient

	// properties set after initialization

	getSubgraphClientForChain func(chainID types.ChainID) *SubgraphClient
	defaultChainID            types.ChainID
}

// NewFeedbackManager creates a new FeedbackManager instance.
func NewFeedbackManager(
	web3Client *Web3Client,
	ipfsClient *IPFSClient,
	reputationRegistry *Contract,
	identityRegistry *Contract,
	subgraphClient *SubgraphClient,
) *FeedbackManager {
	return &FeedbackManager{
		web3Client:         web3Client,
		ipfsClient:         ipfsClient,
		reputationRegistry: reputationRegistry,
		identityRegistry:   identityRegistry,
		subgraphClient:     subgraphClient,
	}
}

// SetSubgraphClientGetter sets the getter function for the subgraph client.
// The getter function gets the subgraph client for a specific chain.
func (f *FeedbackManager) SetSubgraphClientGetter(
	getter func(chainID types.ChainID) *SubgraphClient,
	defaultChainID types.ChainID,
) {
	f.getSubgraphClientForChain = getter
	f.defaultChainID = defaultChainID
}

// SetReputationRegistry sets the reputation registry contract (for lazy initialization).
func (f *FeedbackManager) SetReputationRegistry(registry *Contract) {
	f.reputationRegistry = registry
}

// SetIdentityRegistry sets the identity registry contract (for lazy initialization).
func (f *FeedbackManager) SetIdentityRegistry(registry *Contract) {
	f.identityRegistry = registry
}

// SignFeedbackAuth signs the feedback authorization information for a client.
func (f *FeedbackManager) SignFeedbackAuth(
	agentID types.AgentID,
	clientAddress types.Address,
	indexLimit int64,
	expiryHours int64,
) string {
	if expiryHours == 0 {
		expiryHours = utils.DEFAULTS["FEEDBACK_EXPIRY_HOURS"]
	}

	// TODO: implementation

	return "0x"
}

// PrepareFeedback prepares a feedback file for submission.
func (f *FeedbackManager) PrepareFeedback(
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
) map[string]any {

	// TODO: implementation

	return map[string]any{}
}

// GiveFeedback submits feedback (maps 8004 endpoint).
func (f *FeedbackManager) GiveFeedback(
	agentID types.AgentID,
	feedbackFile map[string]any,
	idem types.IdemKey,
	feedbackAuth string,
) Feedback {

	// TODO: implementation

	return Feedback{}
}

// GetFeedback gets a single feedback entry (currently only supports blockchain query).
func (f *FeedbackManager) GetFeedback(
	agentID types.AgentID,
	clientAddress types.Address,
	feedbackIndex int64,
) Feedback {

	// TODO: implementation

	return Feedback{}
}

// getFeedbackFromBlockchain gets a single feedback entry from the blockchain.
func (f *FeedbackManager) getFeedbackFromBlockchain(
	agentID types.AgentID,
	clientAddress types.Address,
	feedbackIndex int64,
) Feedback {

	// TODO: implementation

	return Feedback{}
}

// SearchFeedback searches feedback entries with filters (uses subgraph if available).
func (f *FeedbackManager) SearchFeedback(params types.SearchFeedbackParams) []types.Feedback {

	// TODO: implementation

	return []types.Feedback{}
}

// mapSubgraphFeedbackToModel maps the feedback data from subgraph to the feedback model.
func (f *FeedbackManager) mapSubgraphFeedbackToModel(
	feedbackData any,
	agentID types.AgentID,
	clientAddress types.Address,
	feedbackIndex int64,
) Feedback {

	// TODO: implementation

	return Feedback{}
}

// hexBytes32ToTakes converts two hex bytes32 strings to plain strings.
// The subgraph now stores tags as human-readable strings (not hex), so
// this method handles both formats for backwards compatibility.
func (f *FeedbackManager) hexBytes32ToTakes(tag1 string, tag2 string) []string {

	// TODO: implementation

	return []string{}
}

// AppendResponse appends a response to feedback.
func (f *FeedbackManager) AppendResponse(
	agentID types.AgentID,
	clientAddress types.Address,
	feedbackIndex int64,
	responseURI types.URI,
	responseHash string,
) string {

	// TODO: implementation

	return ""
}

// RevokeFeedback revokes feedback.
func (f *FeedbackManager) RevokeFeedback(agentID types.AgentID, feedbackIndex int64) string {

	// TODO: implementation

	return ""
}

// stringToBytes32 converts a string to bytes32 for blockchain storage.
func (f *FeedbackManager) stringToBytes32(text string) string {

	// TODO: implementation

	return ""
}

// bytes32ToTags converts bytes32 tags back to plain strings.
func (f *FeedbackManager) bytes32ToTags(tag1Bytes, tag2Bytes string) []string {

	// TODO: implementation

	return []string{}
}

// GetReputationSummary gets the reputation summary for an agent.
func (f *FeedbackManager) GetReputationSummary(
	agentID types.AgentID,
	tag1 string,
	tag2 string,
) ReputationSummary {

	// TODO: implementation

	return ReputationSummary{}
}

// ...

type Feedback struct {
	ID    []string
	Score int64
	Tags  []string
}

type FeedbackFile = map[string]any

type FeedbackResponse struct {
	URI  string
	Hash string
}

type ReputationSummary struct {
	AverageScore int64
	Count        int64
}
