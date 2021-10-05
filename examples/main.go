package main

import (
	"context"
	"fmt"

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
}
