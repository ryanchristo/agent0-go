package utils

import (
	"log"
	"strconv"
	"strings"

	"github.com/ryanchristo/agent0-go/sdk/types"
)

// ParsedAgentID contains the parsed components of an agent ID.
type ParsedAgentID struct {
	ChainID types.ChainID
	TokenID int64
}

// ParseAgentID parses an agent ID string and returns the components.
// The agent ID string must be in the format "chainID:tokenID".
func ParseAgentID(id types.AgentID) ParsedAgentID {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		log.Fatalf("Invalid agent ID format: %s. Expected format: chainID:tokenID", id)
	}

	chainID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.Fatalf("Invalid chain ID in agent ID %s: %v", id, err)
	}

	tokenID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Fatalf("Invalid token ID in agent ID %s: %v", id, err)
	}

	return ParsedAgentID{
		ChainID: chainID,
		TokenID: tokenID,
	}
}

// FormattedAgentID formats agent ID components into the format "chainID:tokenID".
func FormattedAgentID(chainID types.ChainID, tokenID string) types.AgentID {
	return types.AgentID(strconv.FormatInt(chainID, 10) + ":" + tokenID)
}

// ParsedFeedbackID contains the parsed components of a feedback ID.
type ParsedFeedbackID struct {
	AgentID       types.AgentID
	ClientAddress types.Address
	FeedbackIndex int64
}

// ParseFeedbackID parses a feedback ID string and returns the components.
// The feedback ID string must be in the format "agentID:clientAddress:feedbackIndex".
// Note: An agent ID may contain colons (e.g. "11155111:123"), so we split from the right.
func ParseFeedbackID(id string) ParsedFeedbackID {
	lastColonIndex := strings.LastIndex(id, ":")
	secondLastColonIndex := strings.LastIndex(id[:lastColonIndex], ":")

	if lastColonIndex == -1 || secondLastColonIndex == -1 {
		log.Fatalf("Invalid feedback ID format: %s", id)
	}

	agentID := types.AgentID(id[:secondLastColonIndex])
	clientAddress := types.Address(id[secondLastColonIndex+1 : lastColonIndex])
	feedbackIndex, err := strconv.ParseInt(id[lastColonIndex+1:], 10, 64)
	if err != nil {
		log.Fatalf("Invalid feedback index %d: %v", feedbackIndex, err)
	}

	// Normalize address to lowercase for consistency
	clientAddress = strings.ToLower(clientAddress)

	return ParsedFeedbackID{
		AgentID:       agentID,
		ClientAddress: clientAddress,
		FeedbackIndex: feedbackIndex,
	}
}

// FormattedFeedbackID formats feedback ID components into the format "agentID:clientAddress:feedbackIndex".
func FormattedFeedbackID(agentID types.AgentID, clientAddress types.Address, feedbackIndex int64) types.FeedbackID {
	// Normalize address to lowercase for consistency
	clientAddress = strings.ToLower(clientAddress)

	return types.FeedbackID(agentID + ":" + clientAddress + ":" + strconv.FormatInt(feedbackIndex, 10))
}
