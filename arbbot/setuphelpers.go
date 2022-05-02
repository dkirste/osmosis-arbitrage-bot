package arbbot

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	appparams "github.com/osmosis-labs/osmosis/v7/app/params"
	bpool "github.com/osmosis-labs/osmosis/v7/x/gamm/pool-models/balancer"
	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"os"
)

func setupInterfaceRegistry(encodingConfig appparams.EncodingConfig) codectypes.InterfaceRegistry {
	interfaceRegistry := encodingConfig.InterfaceRegistry

	interfaceRegistry.RegisterInterface("/osmosis.gamm.v1beta1.PoolI", (*gammtypes.PoolI)(nil))
	interfaceRegistry.RegisterImplementations((*gammtypes.PoolI)(nil), &bpool.Pool{})

	interfaceRegistry.RegisterInterface("/cosmos.auth.v1beta1.BaseAccount", (*authtypes.AccountI)(nil))
	interfaceRegistry.RegisterImplementations((*authtypes.AccountI)(nil), &authtypes.BaseAccount{})
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{}, &ed25519.PubKey{})
	interfaceRegistry.RegisterImplementations((*cryptotypes.PrivKey)(nil), &secp256k1.PrivKey{})

	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	gammtypes.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &gammtypes.MsgSwapExactAmountIn{})

	return interfaceRegistry
}

func (ab *ArbBot) setupClientContext(rpcclient *rpchttp.HTTP, nodeIp string, armor string, passphrase string) client.Context {
	var accAddress sdk.AccAddress
	accAddress, err := sdk.AccAddressFromBech32(ab.ArbAddress)
	if err != nil {
		fmt.Println(err)
	}

	kb := keyring.NewInMemory(ab.ProtoCodec)
	err = kb.ImportPrivKey("arb", armor, passphrase)
	if err != nil {
		panic("Error while importing private key")
	}

	clientCtx := client.Context{
		FromAddress:       accAddress, // CUSTOM
		Client:            rpcclient,
		GRPCClient:        nil,
		ChainID:           "osmosis-1", //CUSTOM
		Codec:             ab.ProtoCodec,
		InterfaceRegistry: ab.InterfaceRegistry, // CUSTOM
		Input:             nil,
		Keyring:           kb, // CUSTOM
		KeyringOptions:    nil,
		Output:            os.Stdout,
		OutputFormat:      "",
		Height:            0,
		HomeDir:           "",
		KeyringDir:        "",
		From:              "arb",   // CUSTOM
		BroadcastMode:     "async", //CUSTOM
		FromName:          "arb",
		SignModeStr:       signing.SignMode_SIGN_MODE_DIRECT.String(), // CUSTOM
		UseLedger:         false,
		Simulate:          false,
		GenerateOnly:      false,
		Offline:           false,
		SkipConfirm:       true,                       // CUSTOM
		TxConfig:          ab.EncodingConfig.TxConfig, // CUSTOM
		AccountRetriever:  authtypes.AccountRetriever{},
		NodeURI:           "", // CUSTOM
		FeePayer:          nil,
		FeeGranter:        nil,
		Viper:             nil,
		IsAux:             false,
		LegacyAmino:       nil,
	}

	return clientCtx
}
