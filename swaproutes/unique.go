package swaproutes

import sdk "github.com/cosmos/cosmos-sdk/types"

func (parbs ProfitableArbitrages) CheckIfMoreProfitableRouteWasAlreadyExecuted(toCeckRoute ProfitableArbitrage) (alreadExecuted bool) {
	// Return false since there are no routes
	if len(parbs) == 0 {
		alreadExecuted = false
		return
	}

	var profitThreshold = sdk.NewInt(5000)
	var profRouteProfit sdk.Int
	var toCheckRouteProfit sdk.Int

	// Calculate profit
	toCheckRouteProfit = toCeckRoute.OptimumOutAmount.Sub(toCeckRoute.OptimumInToken.Amount)

RouteLoop:
	for _, profRoute := range parbs {
		profRouteProfit = profRoute.OptimumOutAmount.Sub(profRoute.OptimumInToken.Amount)

		for _, profRoutePool := range profRoute.Route {
			profRoutePoolId := profRoutePool.PoolId

			for _, cleanedRoutePool := range toCeckRoute.Route {
				cleanedRoutePoolId := cleanedRoutePool.PoolId

				if profRoutePoolId == cleanedRoutePoolId {
					if profRouteProfit.Add(profitThreshold).GT(toCheckRouteProfit) {
						alreadExecuted = true
						return
					} else {
						// Continue with next profRoute, since the toCheck one is currently better
						continue RouteLoop
					}
				}
			}
		}
	}

	alreadExecuted = false
	return
}

func (swapRoutes SwapAmountInRoutes) CheckIfPoolsAreUnique() bool {
	var poolIds = make([]uint64, 0)
	var routePoolId uint64

	for _, route := range swapRoutes {
		routePoolId = route.Pool.GetId()
		for _, poolId := range poolIds {
			if routePoolId == poolId {
				return false
			}
		}
		poolIds = append(poolIds, routePoolId)
	}
	return true
}

func (swapRoutes SwapAmountInRoutesId) CheckIfPoolsAreUnique() bool {
	var poolIds = make([]uint64, 0)

	for _, route := range swapRoutes {
		for _, poolId := range poolIds {
			if route.PoolId == poolId {
				return false
			}
		}
		poolIds = append(poolIds, route.PoolId)
	}
	return true
}
