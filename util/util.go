package util

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

// MapRPCError converts a RPC error with nested informatoin to a error with a useful and complete message.
func MapRPCError(rpcErr error) error {
	if e, ok := rpcErr.(rpc.DataError); ok {
		bytes, err := json.MarshalIndent(e.ErrorData(), "", "  ")
		if err != nil {
			return rpcErr
		}
		return fmt.Errorf("%s: %s", e.Error(), string(bytes))
	}
	if e, ok := rpcErr.(rpc.Error); ok {
		return fmt.Errorf("%s with code: %d", e.Error(), e.ErrorCode())
	}
	return rpcErr
}

// Retry will retry running the provided function.
func Retry(numRetries int, retryWait time.Duration, backoff float32, run func(done *bool) error) error {
	wait := retryWait
	done := false
	for i := 0; i < numRetries; i++ {
		err := run(&done)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
		if i > 0 {
			waitNanosF := float32(wait.Nanoseconds())
			newWaitNanosF := waitNanosF * backoff
			wait = time.Duration(int64(newWaitNanosF))
		}
		time.Sleep(wait)
	}
	return fmt.Errorf("failed to finish after %d retries", numRetries)
}
