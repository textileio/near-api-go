package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gateway-fm/near-api-go/models"
	itypes "github.com/gateway-fm/near-api-go/types"

	"github.com/gateway-fm/near-api-go/account"
	"github.com/gateway-fm/near-api-go/config"
	"github.com/gateway-fm/near-api-go/util"
)

// CallFunctionResponse holds information about the result of a function call.
type CallFunctionResponse struct {
	Result      []byte   `json:"result"`
	Logs        []string `json:"logs"`
	BlockHeight int      `json:"block_height"`
	BlockHash   string   `json:"block_hash"`
}

// ViewCodeResponse holds information about contract code.
type ViewCodeResponse struct {
	CodeBase64  string `json:"code_base64"`
	Hash        string `json:"hash"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
}

// Change holds information about a state change of a key-value pair.
type Change struct {
	AccountID   string `json:"account_id"`
	KeyBase64   string `json:"key_base64"`
	ValueBase64 string `json:"value_base64"`
}

// Cause holds information about the cause of a state change.
type Cause struct {
	Type        string `json:"type"`
	ReceiptHash string `json:"receipt_hash"`
}

// ChangeData holds information about a state change.
type ChangeData struct {
	Cause  Cause  `json:"cause"`
	Type   string `json:"type"`
	Change Change `json:"change"`
}

// DataChangesResponse holds information about all data changes in a block.
type DataChangesResponse struct {
	BlockHash string       `json:"block_hash"`
	Changes   []ChangeData `json:"changes"`
}

// SyncInfo holds information about the sync status of a node.
type SyncInfo struct {
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestBlockHeight int    `json:"latest_block_height"`
	// TODO: make this time.Time and use custom json conversion.
	LatestBlockTime string `json:"latest_block_time"`
	LatestStateRoot string `json:"latest_state_root"`
	Syncing         bool   `json:"syncing"`
}

// Version struct
type Version struct {
	Build   string `json:"build"`
	Version string `json:"version"`
}

// Validators array
type Validators []Validator

// ValidatorsResponse struct
type ValidatorsResponse struct {
	CurrentFishermans []Fisherman    `json:"current_fishermen,omitempty"`
	NextFishermans    []Fisherman    `json:"next_fishermen,omitempty"`
	CurrentValidators Validators     `json:"current_validators,omitempty"`
	NextValidators    Validators     `json:"next_validators,omitempty"`
	CurrentProposal   string         `json:"current_proposal,omitempty"`
	EpochStartHeight  int64          `json:"epoch_start_height"`
	PrevEpochKickout  []EpochKickout `json:"prev_epoch_kickout"`
}

// EpochKickout struct
type EpochKickout struct {
	AccountID string                 `json:"account_id"`
	Reason    map[string]interface{} `json:"reason"`
}

// Fisherman struct
type Fisherman struct {
	AccountID string `json:"account_id"`
	PublicKey string `json:"public_key"`
	Stake     string `json:"stake"`
}

// Validator struct
type Validator struct {
	AccountID         string  `json:"account_id"`
	IsSlashed         bool    `json:"is_slashed"`
	ExpectedBlocksNum int64   `json:"num_expected_blocks,omitempty"`
	ProducedBlocksNum int64   `json:"num_produced_blocks,omitempty"`
	PublicKey         string  `json:"public_key,omitempty"`
	Shards            []int64 `json:"shards,omitempty"`
	Stake             string  `json:"stake,omitempty"`
}

// NodeStatusResponse holds information about node status.
type NodeStatusResponse struct {
	ChainID    string     `json:"chain_id"`
	RPCAddr    string     `json:"rpc_addr"`
	SyncInfo   SyncInfo   `json:"sync_info"`
	Validators Validators `json:"validators"`
	Version    Version    `json:"version"`
}

// Client communicates with the NEAR API.
type Client struct {
	config *config.Config
}

// NewClient creates a new Client.
func NewClient(config *config.Config) (*Client, error) {
	return &Client{
		config: config,
	}, nil
}

// Account provides an API for the provided account ID.
func (c *Client) Account(accountID string) *account.Account {
	return account.NewAccount(c.config, accountID)
}

// CallFunctionOption controls the behavior when calling CallFunction.
type CallFunctionOption func(*itypes.QueryRequest) error

// CallFunctionWithFinality specifies the finality to be used when calling the function.
func CallFunctionWithFinality(finalaity string) CallFunctionOption {
	return func(qr *itypes.QueryRequest) error {
		qr.Finality = finalaity
		return nil
	}
}

// CallFunctionWithBlockHeight specifies the block height to call the function for.
func CallFunctionWithBlockHeight(blockHeight int) CallFunctionOption {
	return func(qr *itypes.QueryRequest) error {
		qr.BlockID = blockHeight
		return nil
	}
}

// CallFunctionWithBlockHash specifies the block hash to call the function for.
func CallFunctionWithBlockHash(blockHash string) CallFunctionOption {
	return func(qr *itypes.QueryRequest) error {
		qr.BlockID = blockHash
		return nil
	}
}

// CallFunctionWithArgs specified the args to call the function with.
// Should be a JSON encodable object.
func CallFunctionWithArgs(args interface{}) CallFunctionOption {
	return func(qr *itypes.QueryRequest) error {
		if args == nil {
			args = make(map[string]interface{})
		}
		bytes, err := json.Marshal(args)
		if err != nil {
			return err
		}
		qr.ArgsBase64 = base64.StdEncoding.EncodeToString(bytes)
		return nil
	}
}

// CallFunction calls a function on a contract.
func (c *Client) CallFunction(
	ctx context.Context,
	accountID string,
	methodName string,
	opts ...CallFunctionOption,
) (*CallFunctionResponse, error) {
	req := &itypes.QueryRequest{
		RequestType: "call_function",
		AccountID:   accountID,
		MethodName:  methodName,
	}
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, err
		}
	}
	if req.BlockID == nil && req.Finality == "" {
		return nil, fmt.Errorf(
			"you must provide ViewAccountWithBlockHeight, ViewAccountWithBlockHash or ViewAccountWithFinality",
		)
	}
	if req.BlockID != nil && req.Finality != "" {
		return nil, fmt.Errorf(
			"you must provide one of ViewAccountWithBlockHeight, ViewAccountWithBlockHash or ViewAccountWithFinality",
		)
	}
	var res CallFunctionResponse
	if err := c.config.RPCClient.CallContext(ctx, &res, "query", rpc.NewNamedParams(req)); err != nil {
		return nil, fmt.Errorf("calling query rpc: %v", util.MapRPCError(err))
	}
	return &res, nil
}

// DataChangesOption controls behavior when calling DataChanges.
type DataChangesOption func(*itypes.ChangesRequest)

// DataChangesWithPrefix sets the data key prefix to query for.
func DataChangesWithPrefix(prefix string) DataChangesOption {
	return func(cr *itypes.ChangesRequest) {
		cr.KeyPrefixBase64 = base64.StdEncoding.EncodeToString([]byte(prefix))
	}
}

// DataChangesWithFinality specifies the finality to be used when querying data changes.
func DataChangesWithFinality(finalaity string) DataChangesOption {
	return func(qr *itypes.ChangesRequest) {
		qr.Finality = finalaity
	}
}

// DataChangesWithBlockHeight specifies the block id to query data changes for.
func DataChangesWithBlockHeight(blockHeight int) DataChangesOption {
	return func(qr *itypes.ChangesRequest) {
		qr.BlockID = blockHeight
	}
}

// DataChangesWithBlockHash specifies the block id to query data changes for.
func DataChangesWithBlockHash(blockHash string) DataChangesOption {
	return func(qr *itypes.ChangesRequest) {
		qr.BlockID = blockHash
	}
}

// DataChanges queries changes to contract data changes.
func (c *Client) DataChanges(
	ctx context.Context,
	accountIDs []string,
	opts ...DataChangesOption,
) (*DataChangesResponse, error) {
	req := &itypes.ChangesRequest{
		ChangesType: "data_changes",
		AccountIDs:  accountIDs,
	}
	for _, opt := range opts {
		opt(req)
	}
	if req.BlockID == nil && req.Finality == "" {
		return nil, fmt.Errorf(
			"you must provide DataChangesWithBlockHeight, DataChangesWithBlockHash or DataChangesWithFinality",
		)
	}
	if req.BlockID != nil && req.Finality != "" {
		return nil, fmt.Errorf(
			"you must provide one of DataChangesWithBlockHeight, DataChangesWithBlockHash or DataChangesWithFinality",
		)
	}
	var res DataChangesResponse
	if err := c.config.RPCClient.CallContext(ctx, &res, "EXPERIMENTAL_changes", rpc.NewNamedParams(req)); err != nil {
		return nil, fmt.Errorf("calling changes rpc: %v", util.MapRPCError(err))
	}
	return &res, nil
}

// ViewCodeOption controls the behavior when calling ViewCode.
type ViewCodeOption func(*itypes.QueryRequest)

// ViewCodeWithFinality specifies the finality to be used when calling the function.
func ViewCodeWithFinality(finalaity string) ViewCodeOption {
	return func(qr *itypes.QueryRequest) {
		qr.Finality = finalaity
	}
}

// ViewCodeWithBlockHeight specifies the block height to call the function for.
func ViewCodeWithBlockHeight(blockHeight int) ViewCodeOption {
	return func(qr *itypes.QueryRequest) {
		qr.BlockID = blockHeight
	}
}

// ViewCodeWithBlockHash specifies the block hash to call the function for.
func ViewCodeWithBlockHash(blockHash string) ViewCodeOption {
	return func(qr *itypes.QueryRequest) {
		qr.BlockID = blockHash
	}
}

// ViewCode returns the smart contract code for the provided account id.
func (c *Client) ViewCode(ctx context.Context, accountID string, opts ...ViewCodeOption) (*ViewCodeResponse, error) {
	req := &itypes.QueryRequest{
		RequestType: "view_code",
		AccountID:   accountID,
		Finality:    "final",
	}
	for _, opt := range opts {
		opt(req)
	}
	if req.BlockID == nil && req.Finality == "" {
		return nil, fmt.Errorf(
			"you must provide ViewAccountWithBlockHeight, ViewAccountWithBlockHash or ViewAccountWithFinality",
		)
	}
	if req.BlockID != nil && req.Finality != "" {
		return nil, fmt.Errorf(
			"you must provide one of ViewAccountWithBlockHeight, ViewAccountWithBlockHash or ViewAccountWithFinality",
		)
	}
	var viewCodeRes ViewCodeResponse
	if err := c.config.RPCClient.CallContext(ctx, &viewCodeRes, "query", rpc.NewNamedParams(req)); err != nil {
		return nil, fmt.Errorf("calling query rpc: %v", util.MapRPCError(err))
	}
	return &viewCodeRes, nil
}

// NodeStatus returns the node status.
func (c *Client) NodeStatus(ctx context.Context) (*NodeStatusResponse, error) {
	var nodeStatusRes NodeStatusResponse
	if err := c.config.RPCClient.CallContext(ctx, &nodeStatusRes, "status"); err != nil {
		return nil, fmt.Errorf("calling status rpc: %v", util.MapRPCError(err))
	}
	return &nodeStatusRes, nil
}

// NetworkInfo returns the current state of node network
// connections (active peers, transmitted data, etc.)
func (c *Client) NetworkInfo(ctx context.Context) (*models.NetworkInfo, error) {
	var networkInfo models.NetworkInfo
	if err := c.config.RPCClient.CallContext(ctx, &networkInfo, "network_info"); err != nil {
		return nil, fmt.Errorf("calling network info rpc: %v", util.MapRPCError(err))
	}
	return &networkInfo, nil
}

func (c *Client) GetBlockResult(ctx context.Context) (*itypes.BlockResult, error) {
	var blockresultinfo itypes.BlockResult
	if err := c.config.RPCClient.CallContext(ctx, &blockresultinfo, "block", rpc.NewNamedParams(itypes.BlockRequest{Finality: "final"})); err != nil {
		return nil, fmt.Errorf("getting block returned an error: %w", err)
	}

	return &blockresultinfo, nil
}

func (c *Client) GetProtocolConfig(ctx context.Context) (*models.ProtocolConfig, error) {
	var protocolConfig models.ProtocolConfig
	if err := c.config.RPCClient.CallContext(ctx, &protocolConfig, "EXPERIMENTAL_protocol_config", rpc.NewNamedParams(itypes.BlockRequest{Finality: "final"})); err != nil {
		return nil, fmt.Errorf("calling protocol config rpc: %v", util.MapRPCError(err))
	}

	return &protocolConfig, nil
}

type GasPriceOpts struct {
	Blockheight int64
	Blockhash   string
}

func (c *Client) GetGasPrice(ctx context.Context, gp *GasPriceOpts) (*itypes.BlockHeader, error) {
	var gasprice itypes.BlockHeader
	if err := c.config.RPCClient.CallContext(ctx, &gasprice, "gas_price", gp); err != nil {
		return nil, fmt.Errorf("getting gasprice returned an error: %w", err)
	}
	return &gasprice, nil
}
