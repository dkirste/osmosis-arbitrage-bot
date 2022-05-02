package poolstorage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProfitableArbitrages []ProfitableArbitrage

type ProfitableArbitrage struct {
	Route            SwapAmountInRoutesId
	OptimumInToken   sdk.Coin
	OptimumOutAmount sdk.Int
}

// EvaluateArbitrage
// ownId has to start at 0
// maxId needs to be the number of Ifs
func (ps *PoolStorage) EvaluateArbitrage(ownId int, maxId int) (profitableRoutes ProfitableArbitrages) {
	var gasEstimate int64 = 2000
	//var allRoutes = [][]SwapAmountInRoutesId{ps.ThreeCurrencyRoutes, ps.FourCurrencyRoutes, ps.FiveCurrencyRoutes}
	var allRoutes = [][]SwapAmountInRoutesId{ps.ThreeCurrencyRoutes, ps.FourCurrencyRoutes, ps.FiveCurrencyRoutes}
	profitableRoutes = make([]ProfitableArbitrage, 0)

	for i, xCurrencyRoutes := range allRoutes {
		if i%maxId != ownId {
			continue
		}
		for _, route := range xCurrencyRoutes {
			//if route[0].PoolId == 481 && route[1].PoolId == 482 && route[2].PoolId == 464 && route[3].PoolId == 463 {
			//	fmt.Printf("Good.")
			//}
			if route == nil {
				continue
			}
			var inputCoin = sdk.Coin{
				Denom:  "uosmo",
				Amount: sdk.NewInt(1000000),
			}
			var outputCoin = sdk.Coin{
				Denom:  "uosmo",
				Amount: ps.CalculateMultihopSwapExactAmountIn(route, inputCoin),
			}

			if outputCoin.IsGTE(inputCoin) {
				optimumIn, optimumOut := ps.FindOptimumFullScan(route)
				if optimumOut.GT(optimumIn.Amount.Add(sdk.NewInt(gasEstimate))) {
					//fmt.Printf("\n--------------------------\n")
					//fmt.Printf("OptimumIn: %v\nOptimumOut: %vuosmo\n", optimumIn, optimumOut)
					profitableArbitrage := ProfitableArbitrage{
						Route:            route,
						OptimumInToken:   optimumIn,
						OptimumOutAmount: optimumOut,
					}
					profitableRoutes = append(profitableRoutes, profitableArbitrage)
				}
			}
		}
	}

	return profitableRoutes
}

func (parbs ProfitableArbitrages) CheckForConflict() ProfitableArbitrages {
	var cleanedRoutes ProfitableArbitrages = make([]ProfitableArbitrage, 0)
	var profRouteProfit sdk.Int
	var cleanedRouteProfit sdk.Int

RouteLoop:
	for _, profRoute := range parbs {
		profRouteProfit = profRoute.OptimumOutAmount.Sub(profRoute.OptimumInToken.Amount)

		for i, cleanedRoute := range cleanedRoutes {
			if cleanedRoute.Route == nil {
				continue
			}
			cleanedRouteProfit = cleanedRoute.OptimumOutAmount.Sub(cleanedRoute.OptimumInToken.Amount)

		PoolLoop:
			for _, profRoutePool := range profRoute.Route {
				profRoutePoolId := profRoutePool.PoolId

				for _, cleanedRoutePool := range cleanedRoute.Route {
					cleanedRoutePoolId := cleanedRoutePool.PoolId

					// Check if routes have to be nilled...
					if profRoutePoolId == cleanedRoutePoolId {
						if profRouteProfit.GT(cleanedRouteProfit) {
							cleanedRoutes[i] = ProfitableArbitrage{}
							break PoolLoop
						} else {
							// Continue with next profRoute, since a more profitable route is already included
							continue RouteLoop
						}
					}
				}
			}
		}
		// All routes nilled
		// Since there was no more profitable route, include this route
		cleanedRoutes = append(cleanedRoutes, profRoute)

	}
	//if len(parbs) != 0 {
	//	fmt.Printf("\n\nProfitableRoutes: %v\n\n", parbs)
	//	fmt.Printf("\n\nCleanedRoutes: %v\n\n", cleanedRoutes)
	//}
	return cleanedRoutes
}
