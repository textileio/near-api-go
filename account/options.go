package account

import (
	"encoding/base64"
	itypes "github.com/gateway-fm/near-api-go/types"
)

// ViewStateOption controls the behavior when calling ViewState.
type ViewStateOption func(*itypes.QueryRequest)

// ViewStateWithFinality specifies the finality to be used when querying the state.
func ViewStateWithFinality(finalaity string) ViewStateOption {
	return func(qr *itypes.QueryRequest) {
		qr.Finality = finalaity
	}
}

// ViewStateWithBlockHeight specifies the block height to query the state for.
func ViewStateWithBlockHeight(blockHeight int) ViewStateOption {
	return func(qr *itypes.QueryRequest) {
		qr.BlockID = blockHeight
	}
}

// ViewStateWithBlockHash specifies the block hash to query the state for.
func ViewStateWithBlockHash(blockHash string) ViewStateOption {
	return func(qr *itypes.QueryRequest) {
		qr.BlockID = blockHash
	}
}

// ViewStateWithPrefix specifies the state key prefix to query for.
func ViewStateWithPrefix(prefix string) ViewStateOption {
	return func(qr *itypes.QueryRequest) {
		qr.PrefixBase64 = base64.StdEncoding.EncodeToString([]byte(prefix))
	}
}

// StateOption controls the behavior when calling ViewAccount.
type StateOption func(*itypes.QueryRequest)

// StateWithFinality specifies the finality to be used when querying the account.
func StateWithFinality(finalaity string) StateOption {
	return func(qr *itypes.QueryRequest) {
		qr.Finality = finalaity
	}
}

// StateWithBlockHeight specifies the block height to query the account for.
func StateWithBlockHeight(blockHeight int) StateOption {
	return func(qr *itypes.QueryRequest) {
		qr.BlockID = blockHeight
	}
}

// StateWithBlockHash specifies the block hash to query the account for.
func StateWithBlockHash(blockHash string) StateOption {
	return func(qr *itypes.QueryRequest) {
		qr.BlockID = blockHash
	}
}
