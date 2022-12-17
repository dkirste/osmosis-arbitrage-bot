package swaproutes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	balancer "github.com/osmosis-labs/osmosis/v13/x/gamm/pool-models/balancer"
)

type ProfitableArbitrages []ProfitableArbitrage
type ProfitableArbitrage struct {
	Route            SwapAmountInRoutesId
	OptimumInToken   sdk.Coin
	OptimumOutAmount sdk.Int
}

type SwapAmountInRoutes []SwapAmountInRoute
type SwapAmountInRoute struct {
	Pool          balancer.Pool
	TokenOutDenom string
}

type SwapAmountInRoutesId []SwapAmountInRouteId
type SwapAmountInRouteId struct {
	PoolId        uint64
	TokenOutDenom string
}
