package arbbot

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpcMachine "github.com/dkirste/arbbot/grpcmachine"
	info "github.com/dkirste/arbbot/infomachine"
	"github.com/dkirste/arbbot/poolstorage"
	"github.com/dkirste/arbbot/txmachine"
	appparams "github.com/osmosis-labs/osmosis/v7/app/params"
	"sync"
)

func (ab *ArbBot) Setup(grpcNodes []string, rpcNodes []string, infoMachineBaseUrl string, address string, privateKeyArmor string, privateKeyPassphrase string, liquidityThreshold float64) ArbBot {
	ab.ArbAddress = address

	// Setup config, registry and protocodec
	ab.EncodingConfig = appparams.MakeEncodingConfig()
	ab.InterfaceRegistry = setupInterfaceRegistry(ab.EncodingConfig)
	ab.ProtoCodec = codec.NewProtoCodec(ab.InterfaceRegistry)

	ab.grpcms = make([]grpcMachine.GrpcMachine, 0)
	for _, grpcNode := range grpcNodes {
		grpcConn := openGRPCConn(grpcNode)
		newGrpcMachine := grpcMachine.GrpcMachine{
			Conn:              grpcConn,
			InterfaceRegistry: ab.InterfaceRegistry,
		}
		ab.grpcms = append(ab.grpcms, newGrpcMachine)
	}

	ab.clientCtxs = make([]client.Context, 0)
	for _, rpcNode := range rpcNodes {
		rpcConn := openRPCConn(rpcNode)
		newClientCtx := ab.setupClientContext(rpcConn, "", privateKeyArmor, privateKeyPassphrase)
		ab.clientCtxs = append(ab.clientCtxs, newClientCtx)
	}

	ab.txm = txmachine.TxMachine{
		Factory: tx.Factory{},
	}
	ab.txm.Setup(ab.clientCtxs[0])

	newInfoMachine := info.InfoMachine{BaseUrl: infoMachineBaseUrl}

	ab.whitelist = newInfoMachine.BuildWhitelist(liquidityThreshold)

	newPoolStorage := poolstorage.PoolStorage{
		PoolsById:           nil,
		PoolsByAsset:        nil,
		ThreeCurrencyRoutes: nil,
		FourCurrencyRoutes:  nil,
		FiveCurrencyRoutes:  nil,
		ArbRoutesById:       nil,
		AssetDict:           poolstorage.AssetDict{},
		Whitelist:           ab.whitelist,
	}

	ab.ps = newPoolStorage
	ab.ps.Setup(1000)
	initialPools, currentHeight := ab.grpcms[0].QueryAllPools()
	ab.ps.AddPools(initialPools)
	ab.ps.GenerateThreeCurrencyRoutes()
	ab.ps.GenerateFourCurrencyRoutes()
	ab.ps.GenerateFiveCurrencyRoutes()
	ab.ps.AddGeneratedThreeCurrencyRoutesById(ab.ps.ThreeCurrencyRoutes)
	ab.ps.AddGeneratedThreeCurrencyRoutesById(ab.ps.FourCurrencyRoutes)
	ab.currentHeight = currentHeight

	ab.maxReserve = sdk.Coin{}
	ab.reserveThreshold = sdk.NewInt(1000000)

	ab.currentHeight = 0

	ab.psMutex = sync.Mutex{}
	ab.executedProfRoutesMutex = sync.Mutex{}
	ab.sequenceNumberMutex = sync.Mutex{}
	ab.currentHeightMutex = sync.Mutex{}

	return ArbBot{}
}
