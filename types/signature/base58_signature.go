package signature

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mr-tron/base58"
)

type Base58Signature struct {
	Type  SignatureType
	Value string

	//sig Signature
}

func NewBase58Signature(raw string) (pk Base58Signature, err error) {
	split := strings.SplitN(raw, ":", 2)
	if len(split) != 2 {
		return pk, ErrInvalidSignature
	}

	sigTypeRaw := split[0]
	encodedSig := split[1]

	sigType, ok := reverseSignatureMapping[sigTypeRaw]
	if !ok {
		return pk, ErrInvalidSignatureType
	}

	decoded, err := base58.Decode(encodedSig)
	if err != nil {
		return pk, fmt.Errorf("failed to decode signature: %w", err)
	}

	pk.Type = signatureTypes[sigType]
	pk.Value = encodedSig

	// TODO
	_ = decoded

	return
}

func (sig Base58Signature) String() string {
	return fmt.Sprintf("%s:%s", sig.Type, sig.Value)
}

func (sig Base58Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(sig.String())
}

func (sig *Base58Signature) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}

	*sig, err = NewBase58Signature(s)
	return
}
