package keys

// KeyStore stores KeyPairs in association with network and account IDs.
type KeyStore interface {
	SetKey(networkID, accountID string, keyPair KeyPair)
}
