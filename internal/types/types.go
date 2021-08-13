package types

// QueryRequest is used for RPC query requests.
type QueryRequest struct {
	RequestType  string      `json:"request_type"`
	Finality     string      `json:"finality,omitempty"`
	BlockID      interface{} `json:"block_id,omitempty"`
	AccountID    string      `json:"account_id,omitempty"`
	PrefixBase64 string      `json:"prefix_base64"`
	MethodName   string      `json:"method_name,omitempty"`
	ArgsBase64   string      `json:"args_base64"`
	PublicKey    string      `json:"public_key,omitempty"`
}

// QueryResponse is a base type used for responses to query requests.
type QueryResponse struct {
	BlockHash   string `json:"block_hash"`
	BlockHeight int    `json:"block_height"`
	// TODO: this property is undocumented, but appears in some API responses. Is this the right place for it?
	Error string `json:"error"`
}

// ChangesRequest is used for RPC changes requests.
type ChangesRequest struct {
	ChangesType     string      `json:"changes_type"`
	AccountIDs      []string    `json:"account_ids"`
	KeyPrefixBase64 string      `json:"key_prefix_base64"`
	Finality        string      `json:"finality,omitempty"`
	BlockID         interface{} `json:"block_id,omitempty"`
}

// BlockRequest is used for RPC block requests.
type BlockRequest struct {
	BlockID  interface{} `json:"block_id,omitempty"`
	Finality string      `json:"finality,omitempty"`
}

// BlockHeader contains information about a block header.
type BlockHeader struct {
	Height                int           `json:"height"`
	EpochID               string        `json:"epoch_id"`
	NextEpochID           string        `json:"next_epoch_id"`
	Hash                  string        `json:"hash"`
	PrevHash              string        `json:"prev_hash"`
	PrevStateRoot         string        `json:"prev_state_root"`
	ChunkReceiptsRoot     string        `json:"hunk_receipts_root"`
	ChunkHeadersRoot      string        `json:"chunk_headers_root"`
	ChunkTxRoot           string        `json:"chunk_tx_root"`
	OutcomeToot           string        `json:"outcome_root"`
	ChunksIncluded        int           `json:"chunks_included"`
	ChallengesRoot        string        `json:"challenges_root"`
	Timestamp             int           `json:"timestamp"`
	TimestampNanosec      string        `json:"timestamp_nanosec"`
	RandomValue           string        `json:"random_value"`
	ValidatorProposals    []interface{} `json:"validator_proposals"` // TODO: what with this?
	ChunkMask             []bool        `json:"chunk_mask"`
	GasPrice              string        `json:"gas_price"`
	RentPaid              string        `json:"rent_paid"`
	ValidatorReward       string        `json:"validator_reward"`
	TotalSupply           string        `json:"total_supply"`
	ChallengesResult      []interface{} `json:"challenges_result"` // TODO: what with this?
	LastFinalBlock        string        `json:"last_final_block"`
	LastDsFinalBlock      string        `json:"last_ds_final_block"`
	NextBpHash            string        `json:"next_bp_hash"`
	BlockMerkleRoot       string        `json:"block_merkle_root"`
	Approvals             []string      `json:"approvals"`
	Signature             string        `json:"signature"`
	LatestProtocolVersion int           `json:"latest_protocol_version"`
}

// Chunk contains information about a chunk.
type Chunk struct {
	ChunkHash            string        `json:"chunk_hash"`
	PrevBlockHash        string        `json:"prev_block_hash"`
	OutcomeRoot          string        `json:"outcome_root"`
	PrevStateRoot        string        `json:"prev_state_root"`
	EncodedMerkleRoot    string        `json:"encoded_merkle_root"`
	EncodedLength        int           `json:"encoded_length"`
	HeightCreated        int           `json:"height_created"`
	HeightIncluded       int           `json:"height_included"`
	ShardID              int           `json:"shard_id"`
	GasUsed              int           `json:"gas_used"`
	GasLimit             int           `json:"gas_limit"`
	RentPaid             string        `json:"rent_paid"`
	ValidatorReward      string        `json:"validator_reward"`
	BalanceBurnt         string        `json:"balance_burnt"`
	OutgoingReceiptsRoot string        `json:"outgoing_receipts_root"`
	TxRoot               string        `json:"tx_root"`
	ValidatorProposals   []interface{} `json:"validator_proposals"` // TODO: what with this?
	Signature            string        `json:"signature"`
}

// BlockResult contains information about a block result.
type BlockResult struct {
	Author string      `json:"author"`
	Header BlockHeader `json:"header"`
	Chunks []Chunk     `json:"chunks"`
}
