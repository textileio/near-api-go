package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/textileio/near-api-go/account"
	itypes "github.com/textileio/near-api-go/internal/types"
	"github.com/textileio/near-api-go/types"
	"github.com/textileio/near-api-go/util"
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
}

// NodeStatusResponse holds information about node status.
type NodeStatusResponse struct {
	SyncInfo *SyncInfo `json:"sync_info"`
}

// Client communicates with the NEAR API.
type Client struct {
	config *types.Config
}

// NewClient creates a new Client.
func NewClient(config *types.Config) (*Client, error) {
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
