package api

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"github.com/textileio/near-api-go/types"

	"testing"
)

var ctx = context.Background()

func TestIt(t *testing.T) {
	c, cleanup := makeClient(t)
	defer cleanup()
	require.NotNil(t, c)
}

// func TestViewCode(t *testing.T) {
// 	c, cleanup := makeClient(t)
// 	defer cleanup()
// 	res, err := c.ViewCode(ctx, "filecoin-bridge.testnet")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// }

// func TestDeployContract(t *testing.T) {
// 	c, cleanup := makeClient(t)
// 	defer cleanup()
// 	res, err := c.ViewCode(ctx, "filecoin-bridge.testnet")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)

// 	bytes, err := base64.StdEncoding.DecodeString(res.CodeBase64)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, bytes)

// 	outcome, err := c.Account("<account id>").DeployContract(ctx, bytes)
// 	require.NoError(t, err)
// 	require.NotNil(t, outcome)

// 	res2, err := c.ViewCode(ctx, "<account id>")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)

// 	require.Equal(t, res.Hash, res2.Hash)
// }

// func TestDataChanges(t *testing.T) {
// 	c, cleanup := makeClient(t)
// 	defer cleanup()
// 	res, err := c.DataChanges(ctx, []string{"filecoin-bridge.testnet"}, DataChangesWithFinality("final"))
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// }

func makeClient(t *testing.T) (*Client, func()) {
	rpcClient, err := rpc.DialContext(ctx, "https://rpc.testnet.near.org")
	require.NoError(t, err)

	// keys, err := keys.NewKeyPairFromString(
	// 	"ed25519:xxxx",
	// )
	// require.NoError(t, err)

	config := &types.Config{
		RPCClient: rpcClient,
		// Signer:    keys,
		NetworkID: "testnet",
	}
	c, err := NewClient(config)

	require.NoError(t, err)
	return c, func() {
		rpcClient.Close()
	}
}
