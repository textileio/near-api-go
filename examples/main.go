package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gateway-fm/near-api-go/account"

	"github.com/ethereum/go-ethereum/rpc"
	api "github.com/gateway-fm/near-api-go"
	"github.com/gateway-fm/near-api-go/config"
)

func main() {
	rpc, err := rpc.DialContext(context.Background(), "https://rpc.mainnet.near.org")
	if err != nil {
		panic(err)
	}

	cfg := &config.Config{
		RPCClient: rpc,
		NetworkID: "mainnet",
	}

	apiN, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	status, err := apiN.NodeStatus(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(status.ChainID)

	netInfo, err := apiN.NetworkInfo(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(netInfo.NumActivePeers)
	fmt.Println(netInfo.ReceivedBytesPerSec)

	acc := account.NewAccount(cfg, "andskur.near")

	view, err := acc.State(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(view.Amount)

	protocolConfig, err := apiN.GetProtocolConfig(context.Background())
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(protocolConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
