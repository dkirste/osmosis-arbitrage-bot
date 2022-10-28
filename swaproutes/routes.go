package swaproutes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gammtypes "github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

type ProfitableArbitrages []ProfitableArbitrage
type ProfitableArbitrage struct {
	Route            SwapAmountInRoutesId
	OptimumInToken   sdk.Coin
	OptimumOutAmount sdk.Int
}

type SwapAmountInRoutes []SwapAmountInRoute
type SwapAmountInRoute struct {
	Pool          gammtypes.PoolI
	TokenOutDenom string
}

type SwapAmountInRoutesId []SwapAmountInRouteId
type SwapAmountInRouteId struct {
	PoolId        uint64
	TokenOutDenom string
}
