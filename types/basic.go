package types

// Account identifier. Provides access to user's state.
type AccountID = string

// Gas is a type for storing amounts of gas.
type Gas = uint64

// Nonce for transactions.
type Nonce = uint64

// Time nanoseconds fit into uint128. Using existing Balance type which
// implements JSON marshal/unmarshal
type TimeNanos = Balance

// BlockHeight is used for height of the block
type BlockHeight = uint64

// ShardID is used for a shard index, from 0 to NUM_SHARDS - 1.
type ShardID = uint64

// StorageUsage is used to count the amount of storage used by a contract.
type StorageUsage = uint64

// NumBlocks holds number of blocks in current group.
type NumBlocks = uint64
