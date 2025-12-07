package types

// EndpointType is a custom type for enumeration.
type EndpointType string

const (
	// ENDPOINT_TYPE_MCP is the MCP endpoint type.
	ENDPOINT_TYPE_MCP EndpointType = "MCP"

	// ENDPOINT_TYPE_A2A is the A2A endpoint type.
	ENDPOINT_TYPE_A2A EndpointType = "A2A"

	// ENDPOINT_TYPE_ENS is the ENS endpoint type.
	ENDPOINT_TYPE_ENS EndpointType = "ENS"

	// ENDPOINT_TYPE_DID is the DID endpoint type.
	ENDPOINT_TYPE_DID EndpointType = "DID"

	// ENDPOINT_TYPE_WALLET is the wallet endpoint type.
	ENDPOINT_TYPE_WALLET EndpointType = "wallet"

	// ENDPOINT_TYPE_OASF is the OASF endpoint type.
	ENDPOINT_TYPE_OASF EndpointType = "OASF"
)

// TrustModel is a custom type for enumeration.
type TrustModel string

const (
	// TRUST_MODEL_REPUTATION is the reputation trust model.
	TRUST_MODEL_REPUTATION TrustModel = "reputation"

	// TRUST_MODEL_CRYPTO_ECONOMICS is the crypto economics trust model.
	TRUST_MODEL_CRYPTO_ECONOMICS TrustModel = "crypto-economics"

	// TRUST_MODEL_TEE_ATTESTATION is the tee attestation trust model.
	TRUST_MODEL_TEE_ATTESTATION TrustModel = "tee-attestation"
)
