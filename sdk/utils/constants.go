package utils

// IPFS_GATEWAYS is a list of IPFS gateway URLs.
var IPFS_GATEWAYS = []string{
	"https://gateway.pinata.cloud/ipfs/",
	"https://ipfs.io/ipfs/",
	"https://dweb.link/ipfs/",
}

// TIMEOUTS is a map of timeout values in milliseconds.
var TIMEOUTS = map[string]int64{
	"IPFS_GATEWAY":             10000, // 10 seconds
	"PINATA_UPLOAD":            80000, // 80 seconds
	"TRANSACTION_WAIT":         30000, // 30 seconds
	"ENDPOINT_CRAWLER_DEFAULT": 5000,  // 5 seconds
}

// DEFAULTS is a map of default values.
var DEFAULTS = map[string]int64{
	"FEEDBACK_EXPIRY_HOURS": 24,
	"SEARCH_PAGE_SIZE":      50,
}
