package account

import (
	"encoding/json"

	"github.com/textileio/near-api-go/internal/types"
)

// Value models a state key-value pair.
type Value struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AccountStateView holds information about contract state.
type AccountStateView struct {
	types.QueryResponse
	Values []Value `json:"values"`
}

// AccountView holds information about an account.
type AccountView struct {
	types.QueryResponse
	Amount        string `json:"amount"`
	Locked        string `json:"locked"`
	CodeHash      string `json:"code_hash"`
	StorageUsage  int    `json:"storage_usage"`
	StoragePaidAt int    `json:"storage_paid_at"`
}

// PermissionType specifies the type of permission.
type PermissionType int

const (
	// FullAccessPermissionType means the account has full access.
	FullAccessPermissionType PermissionType = iota
	// FunctionCallPermissionType means the account has permission to call some functions.
	FunctionCallPermissionType
)

// AccessKeyView contains information about an access key.
type AccessKeyView struct {
	types.QueryResponse
	Nonce                      uint64
	PermissionType             PermissionType
	FunctionCallPermissionView *FunctionCallPermissionView
}

// FunctionCall provides information about the allowed function call.
type FunctionCall struct {
	Allowance   string   `json:"allowance"`
	ReceiverID  string   `json:"receiver_id"`
	MethodNames []string `json:"method_names"`
}

// FunctionCallPermissionView contains a FunctionCall.
type FunctionCallPermissionView struct {
	FunctionCall FunctionCall `json:"FunctionCall"`
}

// FinalExecutionStatus is the final status of a transaction.
type FinalExecutionStatus struct {
	SuccessValue string                 `json:"SuccessValue"`
	Failure      map[string]interface{} `json:"Failure,omitempty"`
}

// FinalExecutionStatusBasic is the final status of a transaction.
type FinalExecutionStatusBasic int

const (
	// FinalExecutionStatusBasicNotStarted means the transaction hasn't started.
	FinalExecutionStatusBasicNotStarted FinalExecutionStatusBasic = iota
	// FinalExecutionStatusBasicStarted means the transaction has started.
	FinalExecutionStatusBasicStarted
	// FinalExecutionStatusBasicFailure means the transaction has failed.
	FinalExecutionStatusBasicFailure
)

// ExecutionStatus is the status of a transaction.
type ExecutionStatus struct {
	SuccessValue     string                 `json:"SuccessValue"`
	SuccessReceiptID string                 `json:"SuccessReceiptId"`
	Failure          map[string]interface{} `json:"Failure,omitempty"`
}

// ExecutionStatusBasic is the status of a transaction.
type ExecutionStatusBasic int

const (
	// ExecutionStatusBasicUnknown means the transaction status is unknown.
	ExecutionStatusBasicUnknown ExecutionStatusBasic = iota
	// ExecutionStatusBasicPending means the transaction is pending.
	ExecutionStatusBasicPending
	// ExecutionStatusBasicFailure means the transaction has failed.
	ExecutionStatusBasicFailure
)

// ExecutionOutcome is the outcome of a transaction.
type ExecutionOutcome struct {
	Logs       []string        `json:"logs"`
	ReceiptIDs []string        `json:"receipt_ids"`
	GasBurnt   int64           `json:"gas_burnt"`
	RawStatus  json.RawMessage `json:"status"`
}

// GetStatus returns a bool indicating if the status is an ExecutionStatus, and if so, the ExecutionStatus.
func (eo *ExecutionOutcome) GetStatus() (ExecutionStatus, bool) {
	var v interface{}
	if err := json.Unmarshal(eo.RawStatus, &v); err != nil {
		return ExecutionStatus{}, false
	}
	switch v.(type) {
	case map[string]interface{}:
		var s ExecutionStatus
		if err := json.Unmarshal(eo.RawStatus, &s); err != nil {
			return ExecutionStatus{}, false
		}
		return s, true
	default:
		return ExecutionStatus{}, false
	}
}

// GetStatusBasic returns a bool indicating if the status is an
// ExecutionStatusBasic, and if so, the ExecutionStatusBasic.
func (eo *ExecutionOutcome) GetStatusBasic() (ExecutionStatusBasic, bool) {
	var v interface{}
	if err := json.Unmarshal(eo.RawStatus, &v); err != nil {
		return ExecutionStatusBasicUnknown, false
	}
	switch v.(type) {
	case string:
		switch v {
		case "Unknown":
			return ExecutionStatusBasicUnknown, true
		case "Pending":
			return ExecutionStatusBasicPending, true
		case "Failure":
			return ExecutionStatusBasicFailure, true
		default:
			return ExecutionStatusBasicUnknown, false
		}
	default:
		return ExecutionStatusBasicUnknown, false
	}
}

// ExecutionOutcomeWithID provides the transaction or receipt outcome with and id.
type ExecutionOutcomeWithID struct {
	ID      string           `json:"id"`
	Outcome ExecutionOutcome `json:"outcome"`
}

// FinalExecutionOutcome is the final outcome of a transaction.
type FinalExecutionOutcome struct {
	RawStatus          json.RawMessage          `json:"status"`
	Transaction        json.RawMessage          `json:"transaction"`
	TransactionOutcome ExecutionOutcomeWithID   `json:"transaction_outcome"`
	ReceiptsOutcome    []ExecutionOutcomeWithID `json:"receipts_outcome"`
}

// GetStatus returns a bool indicating if the status is an FinalExecutionStatus, and if so, the FinalExecutionStatus.
func (feo *FinalExecutionOutcome) GetStatus() (FinalExecutionStatus, bool) {
	var v interface{}
	if err := json.Unmarshal(feo.RawStatus, &v); err != nil {
		return FinalExecutionStatus{}, false
	}
	switch v.(type) {
	case map[string]interface{}:
		var s FinalExecutionStatus
		if err := json.Unmarshal(feo.RawStatus, &s); err != nil {
			return FinalExecutionStatus{}, false
		}
		return s, true
	default:
		return FinalExecutionStatus{}, false
	}
}

// GetStatusBasic returns a bool indicating if the status is an
// FinalExecutionStatusBasic, and if so, the FinalExecutionStatusBasic.
func (feo *FinalExecutionOutcome) GetStatusBasic() (FinalExecutionStatusBasic, bool) {
	var v interface{}
	if err := json.Unmarshal(feo.RawStatus, &v); err != nil {
		return FinalExecutionStatusBasicNotStarted, false
	}
	switch v.(type) {
	case string:
		switch v {
		case "NotStarted":
			return FinalExecutionStatusBasicNotStarted, true
		case "Started":
			return FinalExecutionStatusBasicStarted, true
		case "Failure":
			return FinalExecutionStatusBasicFailure, true
		default:
			return FinalExecutionStatusBasicNotStarted, false
		}
	default:
		return FinalExecutionStatusBasicNotStarted, false
	}
}
