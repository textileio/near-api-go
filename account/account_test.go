package account

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"github.com/textileio/near-api-go/keys"
	"github.com/textileio/near-api-go/types"

	"testing"
)

var ctx = context.Background()

func TestIt(t *testing.T) {
	a, cleanup := makeAccount(t)
	defer cleanup()
	require.NotNil(t, a)
}

// func TestViewState(t *testing.T) {
// 	a, cleanup := makeAccount(t)
// 	defer cleanup()
// 	res, err := a.ViewState(ctx, ViewStateWithFinality("final"))
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// }

// func TestState(t *testing.T) {
// 	a, cleanup := makeAccount(t)
// 	defer cleanup()
// 	res, err := a.State(ctx, StateWithFinality("final"))
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// }

// func TestFindAccessKey(t *testing.T) {
// 	a, cleanup := makeAccount(t)
// 	defer cleanup()
// 	pubKey, accessKeyView, err := a.FindAccessKey(ctx, "", nil)
// 	require.NoError(t, err)
// 	require.NotNil(t, pubKey)
// 	require.NotNil(t, accessKeyView)
// }

func TestViewAccessKey(t *testing.T) {
	a, cleanup := makeAccount(t)
	defer cleanup()
	pkey, err := keys.NewPublicKeyFromString("ed25519:H9k5eiU4xXS3M4z8HzKJSLaZdqGdGwBG49o7orNC4eZW")
	require.NoError(t, err)
	v, err := a.ViewAccessKey(ctx, pkey)
	require.NoError(t, err)
	require.NotNil(t, v)

	pkey, err = keys.NewPublicKeyFromString("ed25519:H9k5eiU4xXS3M4zH8zKJSLaZdqGdGwBG49o7orNC4eZW")
	require.NoError(t, err)
	_, err = a.ViewAccessKey(ctx, pkey)
	require.Error(t, err)
}

// func TestSignTransaction(t *testing.T) {
// 	a, cleanup := makeAccount(t)
// 	defer cleanup()
// 	amt := big.NewInt(1000)
// 	sendAction := transaction.TransferAction(*amt)
// 	hash, signedTxn, err := a.SignTransaction(ctx, "carsonfarmer.testnet", sendAction)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, hash)
// 	require.NotNil(t, signedTxn)
// }

// func TestSignAndSendTransaction(t *testing.T) {
// 	a, cleanup := makeAccount(t)
// 	defer cleanup()
// 	amt, ok := (&big.Int{}).SetString("1000000000000000000000000", 10)
// 	require.True(t, ok)
// 	sendAction := transaction.TransferAction(*amt)
// 	res, err := a.SignAndSendTransaction(ctx, "carsonfarmer.testnet", sendAction)
// 	require.NoError(t, err)
// 	require.NotNil(t, res)

// 	status, ok := res.GetStatus()
// 	fmt.Println(status, ok)

// 	status2, ok := res.GetStatusBasic()
// 	fmt.Println(status2, ok)
// }

func makeAccount(t *testing.T) (*Account, func()) {
	rpcClient, err := rpc.DialContext(ctx, "https://rpc.testnet.near.org")
	require.NoError(t, err)

	// keys, err := keys.NewKeyPairFromString(
	// 	"ed25519:xxxx",
	// )
	// require.NoError(t, err)

	config := &types.Config{
		RPCClient: rpcClient,
		NetworkID: "testnet",
		// Signer:    keys,
	}
	a := NewAccount(config, "client.chainlink.testnet")
	return a, func() {
		rpcClient.Close()
	}
}
