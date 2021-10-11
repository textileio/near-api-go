package models

// ProtocolConfig holds protocol config information
type ProtocolConfig struct {
	AvgHiddenValidatorSeatsPerShard []uint32       `json:"avg_hidden_validator_seats_per_shard"`
	BlockProducerKickoutThreshold   uint32         `json:"block_producer_kickout_threshold"`
	ChainId                         string         `json:"chain_id"`
	ChunkProducerKickoutThreshold   uint32         `json:"chunk_producer_kickout_threshold"`
	DynamicResharding               bool           `json:"dynamic_resharding"`
	EpochLength                     uint32         `json:"epoch_length"`
	FishermenThreshold              string         `json:"fishermen_threshold"`
	GasLimit                        uint64         `json:"gas_limit"`
	GasPriceAdjustmentRate          []uint32       `json:"gas_price_adjustment_rate"`
	GenesisHeight                   uint64         `json:"genesis_height"`
	GenesisTime                     string         `json:"genesis_time"`
	MaxGasPrice                     string         `json:"max_gas_price"`
	MaxInflationRate                []uint32       `json:"max_inflation_rate"`
	MinGasPrice                     string         `json:"min_gas_price"`
	MinimumStakeDivisor             uint32         `json:"minimum_stake_divisor"`
	NumBlockProducerSeats           uint32         `json:"num_block_producer_seats"`
	NumBlockProducerSeatsPerShard   []uint32       `json:"num_block_producer_seats_per_shard"`
	NumBlocksPerYear                uint64         `json:"num_blocks_per_year"`
	OnlineMaxThreshold              []uint32       `json:"online_max_threshold"`
	OnlineMinThreshold              []uint32       `json:"online_min_threshold"`
	ProtocolRewardRate              []uint32       `json:"protocol_reward_rate"`
	ProtocolTreasuryAccount         string         `json:"protocol_treasury_account"`
	ProtocolUpgradeStakeThreshold   []uint32       `json:"protocol_upgrade_stake_threshold"`
	ProtocolVersion                 uint32         `json:"protocol_version"`
	TransactionValidityPeriod       uint32         `json:"transaction_validity_period"`
	RuntimeConfig                   *RuntimeConfig `json:"runtime_config"`
}

type RuntimeConfig struct {
	AccountCreationConfig *AccountCreationConfig `json:"account_creation_config"`
	StorageAmountPerByte  string                 `json:"storage_amount_per_byte"`
	TransactionCosts      *TransactionCosts      `json:"transaction_costs"`
	WasmConfig            *WasmConfig            `json:"wasm_config"`
}

type AccountCreationConfig struct {
	MinAllowedTopLevelAccountLength uint32 `json:"min_allowed_top_level_account_length"`
	RegistrarAccountId              string `json:"registrar_account_id"`
}

type TransactionCosts struct {
	ActionCreationConfig              *ActionCreationConfig      `json:"action_creation_config"`
	ActionReceiptCreationConfig       *Cost                      `json:"action_receipt_creation_config"`
	BurntGasReward                    []uint32                   `json:"burnt_gas_reward"`
	DataReceiptCreationConfig         *DataReceiptCreationConfig `json:"data_receipt_creation_config"`
	PessimisticGasPriceInflationRatio []uint32                   `json:"pessimistic_gas_price_inflation_ratio"`
	StorageUsageConfig                *StorageUsageConfig        `json:"storage_usage_config"`
}

type ActionCreationConfig struct {
	AddKeyCost                *AddKeyCost `json:"add_key_cost"`
	CreateAccountCost         *Cost       `json:"create_account_cost"`
	DeleteAccountCost         *Cost       `json:"delete_account_cost"`
	DeleteKeyCost             *Cost       `json:"delete_key_cost"`
	DeployContractCost        *Cost       `json:"deploy_contract_cost"`
	DeployContractCostPerByte *Cost       `json:"deploy_contract_cost_per_byte"`
	FunctionCallCost          *Cost       `json:"function_call_cost"`
	FunctionCallCostPerByte   *Cost       `json:"function_call_cost_per_byte"`
	StakeCost                 *Cost       `json:"stake_cost"`
	TransferCost              *Cost       `json:"transfer_cost"`
}

type AddKeyCost struct {
	FullAccessCost          *Cost `json:"full_access_cost"`
	FunctionCallCost        *Cost `json:"function_call_cost"`
	FunctionCallCostPerByte *Cost `json:"function_call_cost_per_byte"`
}

type DataReceiptCreationConfig struct {
	BaseCost    *Cost `json:"base_cost"`
	CostPerByte *Cost `json:"cost_per_byte"`
}

type StorageUsageConfig struct {
	NumBytesAccount     uint32 `json:"num_bytes_account"`
	NumExtraBytesRecord uint32 `json:"num_extra_bytes_record"`
}

type Cost struct {
	Execution  uint64 `json:"execution"`
	SendNotSir uint64 `json:"send_not_sir"`
	SendSir    uint64 `json:"send_sir"`
}

type WasmConfig struct {
	ExtCosts      *ExtCosts    `json:"ext_costs"`
	GrowMemCost   uint32       `json:"grow_mem_cost"`
	LimitConfig   *LimitConfig `json:"limit_config"`
	RegularOpCost uint64       `json:"regular_op_cost"`
}

type ExtCosts struct {
	Base                        uint64 `json:"base"`
	ContractCompileBase         uint64 `json:"contract_compile_base"`
	ContractCompileBytes        uint64 `json:"contract_compile_bytes"`
	EcrecoverBase               uint64 `json:"ecrecover_base"`
	Keccak256Base               uint64 `json:"keccak_256_base"`
	Keccak256Byte               uint64 `json:"keccak_256_byte"`
	Keccak512Base               uint64 `json:"keccak_512_base"`
	Keccak512Byte               uint64 `json:"keccak_512_byte"`
	LogBase                     uint64 `json:"log_base"`
	LogByte                     uint64 `json:"log_byte"`
	PromiseAndBase              uint64 `json:"promise_and_base"`
	PromiseAndPerPromise        uint64 `json:"promise_and_per_promise"`
	PromiseReturn               uint64 `json:"promise_return"`
	ReadMemoryBase              uint64 `json:"read_memory_base"`
	ReadMemoryByte              uint64 `json:"read_memory_byte"`
	ReadRegisterBase            uint64 `json:"read_register_base"`
	ReadRegisterByte            uint64 `json:"read_register_byte"`
	Ripemd160Base               uint64 `json:"ripemd_160_base"`
	Ripemd160Block              uint64 `json:"ripemd_160_block"`
	Sha256Base                  uint64 `json:"sha_256_base"`
	Sha256Byte                  uint64 `json:"sha_256_byte"`
	StorageHasKeyBase           uint64 `json:"storage_has_key_base"`
	StorageHasKeyByte           uint64 `json:"storage_has_key_byte"`
	StorageIterCreateFromByte   uint64 `json:"storage_iter_create_from_byte"`
	StorageIterCreatePrefixBase uint64 `json:"storage_iter_create_prefix_base"`
	StorageIterCreatePrefixByte uint64 `json:"storage_iter_create_prefix_byte"`
	StorageIterCreateRangeBase  uint64 `json:"storage_iter_create_range_base"`
	StorageIterCreateToByte     uint64 `json:"storage_iter_create_to_byte"`
	StorageIterNextBase         uint64 `json:"storage_iter_next_base"`
	StorageIterNextKeyByte      uint64 `json:"storage_iter_next_key_byte"`
	StorageIterNextValueByte    uint64 `json:"storage_iter_next_value_byte"`
	StorageReadBase             uint64 `json:"storage_read_base"`
	StorageReadKeyByte          uint64 `json:"storage_read_key_byte"`
	StorageReadValueByte        uint64 `json:"storage_read_value_byte"`
	StorageRemoveBase           uint64 `json:"storage_remove_base"`
	StorageRemoveKeyByte        uint64 `json:"storage_remove_key_byte"`
	StorageRemoveRetValueByte   uint64 `json:"storage_remove_ret_value_byte"`
	StorageWriteBase            uint64 `json:"storage_write_base"`
	StoragewriteEvictedByte     uint64 `json:"storagewrite_evicted_byte"`
	StorageWriteKeyByte         uint64 `json:"storage_write_key_byte"`
	StorageWriteValueByte       uint64 `json:"storage_write_value_byte"`
	TouchingTrieNode            uint64 `json:"touching_trie_node"`
	Utf16DecodingBase           uint64 `json:"utf_16_decoding_base"`
	Utf16DecodingByte           uint64 `json:"utf_16_decoding_byte"`
	Utf8DecodingBase            uint64 `json:"utf_8_decoding_base"`
	Utf8DecodingByte            uint64 `json:"utf_8_decoding_byte"`
	ValidatorStakeBase          uint64 `json:"validator_stake_base"`
	ValidatorTotalStakeBase     uint64 `json:"validator_total_stake_base"`
	WriteMemoryBase             uint64 `json:"write_memory_base"`
	WriteMemoryByte             uint64 `json:"write_memory_byte"`
	WriteRegisterBase           uint64 `json:"write_register_base"`
	WriteRegisterByte           uint64 `json:"write_register_byte"`
}

type LimitConfig struct {
	InitialMemoryPages               uint64 `json:"initial_memory_pages"`
	MaxActionsPerReceipt             uint64 `json:"max_actions_per_receipt"`
	MaxArgumentsLength               uint64 `json:"max_arguments_length"`
	MaxContractSize                  uint64 `json:"max_contract_size"`
	MaxGasBurnt                      uint64 `json:"max_gas_burnt"`
	MaxGasBurntView                  uint64 `json:"max_gas_burnt_view"`
	MaxLengthMethodName              uint64 `json:"max_length_method_name"`
	MaxLengthReturnedData            uint64 `json:"max_length_returned_data"`
	MaxLengthStorageKey              uint64 `json:"max_length_storage_key"`
	MaxLengthStorageValue            uint64 `json:"max_length_storage_value"`
	MaxMemoryPages                   uint64 `json:"max_memory_pages"`
	MaxNumberBytesMethodNames        uint64 `json:"max_number_bytes_method_names"`
	MaxNumberInputDataDependencies   uint64 `json:"max_number_input_data_dependencies"`
	MaxNumberLogs                    uint64 `json:"max_number_logs"`
	MaxNumberRegisters               uint64 `json:"max_number_registers"`
	MaxPromisesPerFunctionCallAction uint64 `json:"max_promises_per_function_call_action"`
	MaxRegisterSize                  uint64 `json:"max_register_size"`
	MaxStackHeight                   uint64 `json:"max_stack_height"`
	MaxTotalLogLength                uint64 `json:"max_total_log_length"`
	MaxTotalPrepaidGas               uint64 `json:"max_total_prepaid_gas"`
	MaxTransactionSize               uint64 `json:"max_transaction_size"`
	RegistersMemoryLimit             uint64 `json:"registers_memory_limit"`
}
