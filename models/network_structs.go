package models

import (
	"encoding/json"
	"time"

	"github.com/gateway-fm/near-api-go/types"
	"github.com/gateway-fm/near-api-go/types/hash"
	"github.com/gateway-fm/near-api-go/types/key"
	"github.com/gateway-fm/near-api-go/types/signature"
)

// NetworkInfo holds network information
type NetworkInfo struct {
	ActivePeers         []*PeerInfo     `json:"active_peers"`
	NumActivePeers      uint            `json:"num_active_peers"`
	PeerMaxCount        uint32          `json:"peer_max_count"`
	HighestHeightPeers  []*PeerInfo     `json:"highest_height_peers"`
	SentBytesPerSec     uint64          `json:"sent_bytes_per_sec"`
	ReceivedBytesPerSec uint64          `json:"received_bytes_per_sec"`
	KnownProducers      []*PeerInfo     `json:"known_producers"`
	MetricRecorder      *MetricRecorder `json:"metric_recorder"`
	PeerCounter         uint            `json:"peer_counter"`
}

// PeerInfo holds peer information
type PeerInfo struct {
	ID        string          `json:"id"`
	Addr      string          `json:"addr"`
	AccountID types.AccountID `json:"account_id"`
}

// PeerChainInfo contains peer chain information. This is derived from PeerCHainInfoV2 in nearcore
type PeerChainInfo struct {
	// Chain Id and hash of genesis block.
	GenesisID GenesisID `json:"genesis_id"`
	// Last known chain height of the peer.
	Height types.BlockHeight `json:"height"`
	// Shards that the peer is tracking.
	TrackedShards []types.ShardID `json:"tracked_shards"`
	// Denote if a node is running in archival mode or not.
	Archival bool `json:"archival"`
}

// EdgeInfo contains information that will be ultimately used to create a new edge. It contains nonce proposed for the edge with signature from peer.
type EdgeInfo struct {
	Nonce     types.Nonce         `json:"nonce"`
	Signature signature.Signature `json:"signature"`
}

// KnownProducer is basically PeerInfo, but AccountID is known
type KnownProducer struct {
	AccountID types.AccountID `json:"account_id"`
	Addr      *string         `json:"addr"`
	PeerID    key.PeerID      `json:"peer_id"`
}

// TODO: chain/network/src/recorder.rs
type MetricRecorder = json.RawMessage

type GenesisID struct {
	// Chain Id
	ChainID string `json:"chain_id"`
	// Hash of genesis block
	Hash hash.CryptoHash `json:"hash"`
}

type StatusResponse struct {
	// Binary version
	Version NodeVersion `json:"version"`
	// Unique chain id.
	ChainID string `json:"chain_id"`
	// Currently active protocol version.
	ProtocolVersion uint32 `json:"protocol_version"`
	// Latest protocol version that this client supports.
	LatestProtocolVersion uint32 `json:"latest_protocol_version"`
	// Address for RPC server.
	RPCAddr string `json:"rpc_addr"`
	// Current epoch validators.
	Validators []ValidatorInfo `json:"validators"`
	// Sync status of the node.
	SyncInfo StatusSyncInfo `json:"sync_info"`
	// Validator id of the node
	ValidatorAccountID *types.AccountID `json:"validator_account_id"`
}

type NodeVersion struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

type ValidatorInfo struct {
	AccountID types.AccountID `json:"account_id"`
	Slashed   bool            `json:"is_slashed"`
}

type StatusSyncInfo struct {
	LatestBlockHash   hash.CryptoHash   `json:"latest_block_hash"`
	LatestBlockHeight types.BlockHeight `json:"latest_block_height"`
	LatestBlockTime   time.Time         `json:"latest_block_time"`
	Syncing           bool              `json:"syncing"`
}

type ValidatorsResponse struct {
	CurrentValidators []CurrentEpochValidatorInfo `json:"current_validator"`
}

type CurrentEpochValidatorInfo struct {
	ValidatorInfo
	PublicKey         key.PublicKey   `json:"public_key"`
	Stake             types.Balance   `json:"stake"`
	Shards            []types.ShardID `json:"shards"`
	NumProducedBlocks types.NumBlocks `json:"num_produced_blocks"`
	NumExpectedBlocks types.NumBlocks `json:"num_expected_blocks"`
}
