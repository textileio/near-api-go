package transaction

import (
	"encoding/base64"
	"testing"

	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"
	"github.com/textileio/near-api-go/keys"
)

func TestIt(t *testing.T) {
	trans := Transaction{}
	signer, err := keys.NewKeyPairFromRandom("ed25519")
	require.NoError(t, err)
	hash, signedT, err := SignTransaction(trans, signer, "", "")
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotNil(t, signedT)
	payload, err := borsh.Serialize(*signedT)
	require.NoError(t, err)
	s := base64.StdEncoding.EncodeToString(payload)
	require.NotEmpty(t, s)
}
