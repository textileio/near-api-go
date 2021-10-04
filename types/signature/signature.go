package signature

import (
	"crypto/ed25519"
)

// TODO: SECP256K1 support
type Signature [1 + ed25519.SignatureSize]byte

func NewSignatureED25519(data []byte) Signature {
	var buf Signature
	buf[0] = RawSignatureTypeED25519
	copy(buf[1:], data[0:ed25519.SignatureSize])
	return buf
}

func (s Signature) Type() byte {
	return s[0]
}

func (s Signature) Value() []byte {
	return s[1:]
}
