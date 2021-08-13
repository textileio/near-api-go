package keys

import (
	"crypto/ed25519"
	"testing"

	"github.com/mr-tron/base58/base58"
	"github.com/stretchr/testify/require"
)

func TestNewRandom(t *testing.T) {
	requireNewRandom(t)
}

func TestNewRandomInvalidCurve(t *testing.T) {
	_, err := NewKeyPairFromRandom("edXX")
	require.Error(t, err)
}

func TestNewFromStringNoPrefix(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	b58 := base58.Encode(priv)
	k, err := NewKeyPairFromString(b58)
	require.NoError(t, err)
	require.NotNil(t, k)
}

func TestNewFromStringPrefix(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	b58 := base58.Encode(priv)
	k, err := NewKeyPairFromString("ed25519:" + b58)
	require.NoError(t, err)
	require.NotNil(t, k)
}

func TestNewFromStringInvalidPrefix(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	b58 := base58.Encode(priv)
	_, err = NewKeyPairFromString("edXX:" + b58)
	require.Error(t, err)
}

func TestSign(t *testing.T) {
	k := requireNewRandom(t)
	requireSign(t, k, []byte{1, 2, 3, 4, 5})
}

func TestVerify(t *testing.T) {
	k := requireNewRandom(t)
	msg := []byte{1, 2, 3, 4, 5}
	s := requireSign(t, k, msg)
	v := k.Verify(msg, s)
	require.True(t, v)
}

func requireNewRandom(t *testing.T) KeyPair {
	kp, err := NewKeyPairFromRandom("ed25519")
	require.NoError(t, err)
	require.NotNil(t, kp)
	require.NoError(t, err)
	require.NotNil(t, kp)
	return kp
}

func requireSign(t *testing.T, k KeyPair, message []byte) []byte {
	require.NotEmpty(t, message)
	res, err := k.Sign(message)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	return res
}
