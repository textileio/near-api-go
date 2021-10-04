package key_test

import (
	"testing"

	"github.com/eteu-technologies/near-api-go/pkg/types/key"
)

func TestED25519Key(t *testing.T) {
	expected := `ed25519:DcA2MzgpJbrUATQLLceocVckhhAqrkingax4oJ9kZ847`

	parsed, err := key.NewBase58PublicKey(expected)
	if err != nil {
		t.Errorf("failed to parse public key: %s", err)
	}

	if s := parsed.String(); s != expected {
		t.Errorf("%s != %s", s, expected)
	}
}

func TestED25519Key_Base58_And_Back(t *testing.T) {
	expected := `ed25519:3xCFas58RKvD5UpF9GqvEb6q9rvgfbEJPhLf85zc4HpC`

	parsed, err := key.NewBase58PublicKey(expected)
	if err != nil {
		t.Errorf("failed to parse public key: %s", err)
	}

	publicKey := parsed.ToPublicKey()
	converted := publicKey.ToBase58PublicKey()

	if s := converted.String(); s != expected {
		t.Errorf("%s != %s", s, expected)
	}
}
