# near-api-go

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

A NEAR client written in Go

The goal of this project is to provide a fully featured NEAR cleint in Go. There is support for most NEAR RPC requests, including those that use signed transactions. Of course, there is room for improvement, especially with integration testing and a fully featured `KeyStore`, so please give it a spin and feel free to open a PR to help us improve the library. 

We're currently relying on [our fork of go-ethereum's JSON RPC client](https://github.com/textileio/go-ethereum) that adds support for named RPC parameters. That work is [pending PR](https://github.com/ethereum/go-ethereum/pull/22656) merge into their master branch.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Install

```
go get github.com/textileio/near-api-go
```

## Usage

Import the required modules.

```golang
import (
  api "github.com/textileio/near-api-go"
  "github.com/textileio/near-api-go/keys"
  "github.com/textileio/near-api-go/transaction"
  "github.com/textileio/near-api-go/types"
  "github.com/ethereum/go-ethereum/rpc"
)

```

Configure and create an API client.

```golang
rpcClient, err := rpc.DialContext(ctx, "https://rpc.testnet.near.org")

keyPair, err := keys.NewKeyPairFromString(
  "ed25519:...",
)

config := &types.Config{
  RPCClient: rpcClient,
  Signer:    keyPair, // Currently we use a key pair directly as a signer.
  NetworkID: "testnet",
}

client, err := api.NewClient(config)
```

Interact with top level functions like `CallFunction`, for example. It can be used for calling non-signed "view" functions.

```golang
res, err := client.CallFunction(
  ctx,
  "<account id>",
  "myFunction",
  api.CallFunctionWithFinality("final"),
)
```

Most other functionality is provided by the `Account` sub module. For example, you can call state-modifying functions that are sent as signed transactions, and even include a deposit while you're at it.

```golang
deposit, ok := (&big.Int{}).SetString("1000000000000000000000000", 10)
res, err := client.Account("<client account id>").FunctionCall(
  ctx,
  <contract account id>,
  "myTxnFunction",
  transaction.FunctionCallWithArgs(map[string]interface{}{
    "arg1": value1, 
    "arg2": value2
  }),
  transaction.FunctionCallWithDeposit(*deposit),
)
```

Check out the [API docs](https://pkg.go.dev/github.com/textileio/near-api-go) to see all that is possible.

## API

[https://pkg.go.dev/github.com/textileio/near-api-go](https://pkg.go.dev/github.com/textileio/near-api-go)

## Maintainers

[@asutula](https://github.com/asutula)

## Contributing

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2021 Textile
