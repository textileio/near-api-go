package key

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mr-tron/base58"
)

type Base58PublicKey struct {
	Type  PublicKeyType
	Value string

	pk PublicKey
}

func NewBase58PublicKey(raw string) (pk Base58PublicKey, err error) {
	split := strings.SplitN(raw, ":", 2)
	if len(split) != 2 {
		return pk, ErrInvalidPublicKey
	}

	keyTypeRaw := split[0]
	encodedKey := split[1]

	keyType, ok := reverseKeyTypeMapping[keyTypeRaw]
	if !ok {
		return pk, ErrInvalidKeyType
	}

	decoded, err := base58.Decode(encodedKey)
	if err != nil {
		return pk, fmt.Errorf("failed to decode public key: %w", err)
	}

	pk.Type = keyTypes[keyType]
	pk.Value = encodedKey

	pk.pk, err = WrapRawKey(pk.Type, decoded)

	return
}

func (pk Base58PublicKey) String() string {
	return fmt.Sprintf("%s:%s", pk.Type, pk.Value)
}

func (pk Base58PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pk.String())
}

func (pk *Base58PublicKey) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}

	*pk, err = NewBase58PublicKey(s)
	return
}

// Copies Base58PublicKey to PublicKey
func (pk *Base58PublicKey) ToPublicKey() PublicKey {
	var buf PublicKey
	copy(buf[:], pk.pk[:])
	return buf
}
