package arbbot

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpcMachine "github.com/dkirste/arbbot/grpcmachine"
	info "github.com/dkirste/arbbot/infomachine"
	"github.com/dkirste/arbbot/poolstorage"
	"github.com/dkirste/arbbot/swaproutes"
	"github.com/dkirste/arbbot/txmachine"
	appparams "github.com/osmosis-labs/osmosis/v7/app/params"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"sync"
)

type ArbBot struct {
	EncodingConfig                appparams.EncodingConfig
	InterfaceRegistry             codectypes.InterfaceRegistry
	ProtoCodec                    *codec.ProtoCodec
	ArbAddress                    string
	im                            info.InfoMachine
	txm                           txmachine.TxMachine
	clientCtxs                    []client.Context
	grpcms                        []grpcMachine.GrpcMachine
	mempoolRPCs                   []*rpchttp.HTTP
	psMutex                       sync.Mutex
	ps                            poolstorage.PoolStorage
	maxReserve                    sdk.Coin
	reserveThreshold              sdk.Int
	whitelist                     []uint64
	executedProfRoutes            swaproutes.ProfitableArbitrages
	executedProfRoutesMutex       sync.Mutex
	sequenceNumber                uint64
	sequenceNumberMutex           sync.Mutex
	currentHeight                 uint64
	currentHeightMutex            sync.Mutex
	executedOptimisticRoutes      swaproutes.ProfitableArbitrages
	executedOptimisticRoutesMutex sync.Mutex
	analysedTxs                   []string
	analysedTxsMutex              sync.Mutex
	priceOracle                   map[string]info.TokenPriceResponse
}
