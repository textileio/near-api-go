package types

import (
	"github.com/ethereum/go-ethereum/rpc"
)

// Config configures the NEAR client.
type Config struct {
	NetworkID string
	RPCClient *rpc.Client

	// /**
	//  * {@link https://github.com/near/near-contract-helper | NEAR Contract Helper} url used
	//  * to create accounts if no master account is provided
	//  * @see {@link UrlAccountCreator}
	//  */
	// helperUrl string

	// /**
	// * The balance transferred from the {@link NearConfig.masterAccount | masterAccount} to a created account
	// * @see {@link LocalAccountCreator}
	//  */
	// initialBalance string

	// /**
	// * The account to use when creating new accounts
	// * @see {@link LocalAccountCreator}
	//  */
	// masterAccount string

	// /**
	// * NEAR wallet url used to redirect users to their wallet in browser applications.
	// * @see {@link https://docs.near.org/docs/tools/near-wallet}
	//  */
	// walletUrl string
}
