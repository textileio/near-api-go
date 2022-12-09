package account

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	logging "github.com/textileio/go-log/v2"
	itypes "github.com/textileio/near-api-go/internal/types"
	"github.com/textileio/near-api-go/keys"
	"github.com/textileio/near-api-go/transaction"
	"github.com/textileio/near-api-go/types"
	"github.com/textileio/near-api-go/util"
)

const (
	nonceRetryCount   = 12
	nonceRetryWait    = time.Millisecond * 500
	nonceRetryBackoff = 1.5
)

var (
	log = logging.Logger("nearclient/account")
)

// Account provides functions for a single account.
type Account struct {
	config    *types.Config
	accountID string
}

// NewAccount creates a new account.
func NewAccount(config *types.Config, accountID string) *Account {
	return &Account{
		config:    config,
		accountID: accountID,
	}
}

// ViewState queries the contract state.
func (a *Account) ViewState(ctx context.Context, opts ...ViewStateOption) (*AccountStateView, error) {
	req := &itypes.QueryRequest{
		RequestType:  "view_state",
		AccountID:    a.accountID,
		PrefixBase64: "",
	}
	for _, opt := range opts {
		opt(req)
	}
	if req.BlockID == nil && req.Finality == "" {
		return nil, fmt.Errorf("you must provide ViewStateWithBlockHeight, ViewStateWithBlockHash or ViewStateWithFinality")
	}
	if req.BlockID != nil && req.Finality != "" {
		return nil, fmt.Errorf(
			"you must provide one of ViewStateWithBlockHeight, ViewStateWithBlockHash or ViewStateWithFinality",
		)
	}
	var res AccountStateView
	if err := a.config.RPCClient.CallContext(ctx, &res, "query", rpc.NewNamedParams(req)); err != nil {
		return nil, fmt.Errorf("calling rpc: %v", util.MapRPCError(err))
	}
	return &res, nil
}

// State queries information about the account state.
func (a *Account) State(
	ctx context.Context,
	opts ...StateOption,
) (*AccountView, error) {
	req := &itypes.QueryRequest{
		RequestType: "view_account",
		AccountID:   a.accountID,
		Finality:    "optimistic",
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
	var res AccountView
	if err := a.config.RPCClient.CallContext(ctx, &res, "query", rpc.NewNamedParams(req)); err != nil {
		return nil, fmt.Errorf("calling rpc: %v", util.MapRPCError(err))
	}
	return &res, nil
}

// ViewAccessKey gets the access key view for the provided public key associated with the account.
func (a *Account) ViewAccessKey(ctx context.Context, pubKey *keys.PublicKey) (*AccessKeyView, error) {
	pubKeyStr, err := pubKey.ToString()
	if err != nil {
		return nil, fmt.Errorf("converting public key to string: %v", err)
	}

	req := &itypes.QueryRequest{
		RequestType: "view_access_key",
		AccountID:   a.accountID,
		PublicKey:   pubKeyStr,
		Finality:    "optimistic",
	}

	type viewAccessKeyResp struct {
		itypes.QueryResponse
		Nonce      uint64           `json:"nonce"`
		Permission *json.RawMessage `json:"permission"`
	}

	var raw json.RawMessage
	resp := &viewAccessKeyResp{Permission: &raw}

	// var res AccessKeyView
	if err := a.config.RPCClient.CallContext(ctx, &resp, "query", rpc.NewNamedParams(req)); err != nil {
		return nil, fmt.Errorf("calling rpc: %v", util.MapRPCError(err))
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("error returned in body: %s", resp.Error)
	}

	ret := &AccessKeyView{
		QueryResponse: itypes.QueryResponse{
			BlockHash:   resp.BlockHash,
			BlockHeight: resp.BlockHeight,
		},
		Nonce: resp.Nonce,
	}

	if string(raw) == "\"FullAccess\"" {
		ret.PermissionType = FullAccessPermissionType
	} else {
		var view FunctionCallPermissionView
		if err := json.Unmarshal(raw, &view); err != nil {
			return nil, fmt.Errorf("unmarshaling permission: %v", err)
		}
		ret.FunctionCallPermissionView = &view
		ret.PermissionType = FunctionCallPermissionType
	}

	return ret, nil
}

func (a *Account) SendTransaction(ctx context.Context, signedTransaction *transaction.SignedTransaction, txHash []byte) (*FinalExecutionOutcome, error) {
	var result *FinalExecutionOutcome
	if err := util.Retry(nonceRetryCount, nonceRetryWait, nonceRetryBackoff, func(done *bool) error {
		bytes, err := borsh.Serialize(*signedTransaction)
		if err != nil {
			return fmt.Errorf("serializing signed transaction: %v", err)
		}
		var res FinalExecutionOutcome
		if err := a.config.RPCClient.CallContext(
			ctx,
			&res,
			"broadcast_tx_commit",
			base64.StdEncoding.EncodeToString(bytes),
		); err != nil {
			mappedErr := util.MapRPCError(err)
			if strings.Contains(mappedErr.Error(), "InvalidNonce") {
				log.Warnf("Retrying transaction %s:%s with new nonce.", base58.Encode(txHash))
				return nil
			}
			return mappedErr
		}
		result = &res
		*done = true
		return nil
	}); err != nil {
		return nil, fmt.Errorf("signing and sending transaction: %v", err)
	}
	if result == nil {
		return nil, fmt.Errorf("failed to send transaction, but no error was returned")
	}

	status, ok := result.GetStatus()
	if ok && status.Failure != nil {
		errorMessage, hasErrorMessage := status.Failure["error_message"]
		errorType, hasErrorType := status.Failure["error_type"]
		if hasErrorMessage && hasErrorType {
			return nil, fmt.Errorf(
				"transaction %s failed with message < %v > and type < %v > ",
				result.TransactionOutcome.ID,
				errorMessage,
				errorType,
			)
		}

		bytes, err := json.MarshalIndent(status.Failure, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshaling failure to json: %v", err)
		}
		return nil, fmt.Errorf("status failure: %s", string(bytes))
	}

	return result, nil
}

// // SignTransaction creates and signs a transaction from the supplied actions.
// func (a *Account) SignTransaction(
// 	ctx context.Context,
// 	receiverID string,
// 	actions ...transaction.Action,
// ) ([]byte, *transaction.SignedTransaction, error) {
// 	if a.config.Signer == nil {
// 		return nil, nil, fmt.Errorf("no signer configured")
// 	}

// 	_, accessKeyView, err := a.FindAccessKey(ctx, receiverID, actions)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("finding access key: %v", err)
// 	}
// 	if accessKeyView == nil {
// 		return nil, nil, fmt.Errorf("no access key view found") // TODO: Better error message.
// 	}
// 	var res itypes.BlockResult
// 	if err := a.config.RPCClient.CallContext(
// 		ctx,
// 		&res,
// 		"block",
// 		rpc.NewNamedParams(itypes.BlockRequest{Finality: "final"}),
// 	); err != nil {
// 		return nil, nil, fmt.Errorf("calling block rpc: %v", util.MapRPCError(err))
// 	}
// 	blockHash, err := base58.Decode(res.Header.Hash)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("decoding hash: %v", err)
// 	}
// 	var blockHashArr [32]byte
// 	copy(blockHashArr[:], blockHash)
// 	nonce := accessKeyView.Nonce + 1

// 	pk := a.config.Signer.GetPublicKey()
// 	var dataArr [32]byte
// 	copy(dataArr[:], pk.Data)

// 	t := transaction.Transaction{
// 		SignerID: a.accountID,
// 		PublicKey: transaction.PublicKey{
// 			KeyType: uint8(pk.Type),
// 			Data:    dataArr,
// 		},
// 		Nonce:      nonce,
// 		ReceiverID: receiverID,
// 		BlockHash:  blockHashArr,
// 		Actions:    actions,
// 	}
// 	hash, signedTransaction, err := transaction.SignTransaction(t, a.config.Signer, a.accountID, a.config.NetworkID)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("signing transaction: %v", err)
// 	}
// 	return hash, signedTransaction, nil
// }

// // SignAndSendTransaction creates, signs and sends a tranaction for the supplied actions.
// func (a *Account) SignAndSendTransaction(
// 	ctx context.Context,
// 	receiverID string,
// 	actions ...transaction.Action,
// ) (*FinalExecutionOutcome, error) {
// 	var result *FinalExecutionOutcome
// 	if err := util.Retry(nonceRetryCount, nonceRetryWait, nonceRetryBackoff, func(done *bool) error {
// 		txHash, signedTransaction, err := a.SignTransaction(ctx, receiverID, actions...)
// 		if err != nil {
// 			return fmt.Errorf("signing transaction: %v", err)
// 		}
// 		bytes, err := borsh.Serialize(*signedTransaction)
// 		if err != nil {
// 			return fmt.Errorf("serializing signed transaction: %v", err)
// 		}
// 		var res FinalExecutionOutcome
// 		if err := a.config.RPCClient.CallContext(
// 			ctx,
// 			&res,
// 			"broadcast_tx_commit",
// 			base64.StdEncoding.EncodeToString(bytes),
// 		); err != nil {
// 			mappedErr := util.MapRPCError(err)
// 			if strings.Contains(mappedErr.Error(), "InvalidNonce") {
// 				// Swallow the error and let Retry continue.
// 				log.Warnf("Retrying transaction %s:%s with new nonce.", receiverID, base58.Encode(txHash))
// 				return nil
// 			}
// 			return mappedErr
// 		}
// 		result = &res
// 		*done = true
// 		return nil
// 	}); err != nil {
// 		return nil, fmt.Errorf("signing and sending transaction: %v", err)
// 	}
// 	if result == nil {
// 		return nil, fmt.Errorf("failed to send transaction, but no error was returned")
// 	}

// 	// TODO: Log logs and failures.

// 	// if (typeof result.status === 'object' && typeof result.status.Failure === 'object') {
// 	// 	// if error data has error_message and error_type properties, node returned an error in the old format
// 	// 	if (result.status.Failure.error_message && result.status.Failure.error_type) {
// 	// 			throw new TypedError(
// 	// 					`Transaction ${result.transaction_outcome.id} failed. ${result.status.Failure.error_message}`,
// 	// 					result.status.Failure.error_type);
// 	// 	} else {
// 	// 			throw parseResultError(result);
// 	// 	}
// 	// }
// 	status, ok := result.GetStatus()
// 	if ok && status.Failure != nil {
// 		errorMessage, hasErrorMessage := status.Failure["error_message"]
// 		errorType, hasErrorType := status.Failure["error_type"]
// 		if hasErrorMessage && hasErrorType {
// 			return nil, fmt.Errorf(
// 				"transaction %s failed with message < %v > and type < %v > ",
// 				result.TransactionOutcome.ID,
// 				errorMessage,
// 				errorType,
// 			)
// 		}
// 		// TODO: the parse result error thing
// 		// The JS client code looks wrong/conflicting because the TS types say that status.Failure is always
// 		// and object of Execution error type. But then they have code that reads like it can be some complex
// 		// schema type. Just going to marshal the Failure map for now.
// 		bytes, err := json.MarshalIndent(status.Failure, "", "  ")
// 		if err != nil {
// 			return nil, fmt.Errorf("marshaling failure to json: %v", err)
// 		}
// 		return nil, fmt.Errorf("status failure: %s", string(bytes))
// 	}

// 	return result, nil
// }

// // FunctionCall calls a smart contract function.
// func (a *Account) FunctionCall(
// 	ctx context.Context,
// 	contractID,
// 	methodName string,
// 	opts ...transaction.FunctionCallOpton,
// ) (*FinalExecutionOutcome, error) {
// 	action, err := transaction.FunctionCallAction(methodName, opts...)
// 	if err != nil {
// 		return nil, fmt.Errorf("creating function call action: %v", err)
// 	}
// 	res, err := a.SignAndSendTransaction(ctx, contractID, *action)
// 	if err != nil {
// 		return nil, fmt.Errorf("signing and sending transaction: %v", err)
// 	}
// 	return res, nil
// }

// // DeployContract deploys contract code to the account.
// func (a *Account) DeployContract(ctx context.Context, code []byte) (*FinalExecutionOutcome, error) {
// 	action := transaction.DeployContractAction(code)
// 	res, err := a.SignAndSendTransaction(ctx, a.accountID, action)
// 	if err != nil {
// 		return nil, fmt.Errorf("signing and sending transaction: %v", err)
// 	}
// 	return res, nil
// }
