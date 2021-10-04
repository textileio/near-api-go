package signature

import (
	"errors"
)

type SignatureType string

const (
	RawSignatureTypeED25519 byte = iota
	RawSignatureTypeSECP256K1
)

const (
	SignatureTypeED25519   = SignatureType("ed25519")
	SignatureTypeSECP256K1 = SignatureType("secp256k1")
)

var (
	ErrInvalidSignature     = errors.New("invalid signature")
	ErrInvalidSignatureType = errors.New("invalid signature type")

	signatureTypes = map[byte]SignatureType{
		RawSignatureTypeED25519:   SignatureTypeED25519,
		RawSignatureTypeSECP256K1: SignatureTypeSECP256K1,
	}
	reverseSignatureMapping = map[string]byte{
		string(SignatureTypeED25519):   RawSignatureTypeED25519,
		string(SignatureTypeSECP256K1): RawSignatureTypeSECP256K1,
	}
)
