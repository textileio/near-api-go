package key

import "errors"

type PublicKeyType string

const (
	RawKeyTypeED25519 byte = iota
	RawKeyTypeSECP256K1
)

const (
	KeyTypeED25519   PublicKeyType = "ed25519"
	KeyTypeSECP256K1 PublicKeyType = "secp256k1"
)

var (
	ErrInvalidPublicKey  = errors.New("invalid public key")
	ErrInvalidPrivateKey = errors.New("invalid private key")
	ErrInvalidKeyType    = errors.New("invalid key type")

	// nolint: deadcode,varcheck,unused
	keyTypes = map[byte]PublicKeyType{
		RawKeyTypeED25519:   KeyTypeED25519,
		RawKeyTypeSECP256K1: KeyTypeSECP256K1,
	}
	reverseKeyTypeMapping = map[string]byte{
		string(KeyTypeED25519):   RawKeyTypeED25519,
		string(KeyTypeSECP256K1): RawKeyTypeSECP256K1,
	}
)
