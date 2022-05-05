package arbbot

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dkirste/arbbot/swaproutes"
)

func (ab *ArbBot) EvaluateArbitrage(ownId int, maxId int) (profitableRoutes swaproutes.ProfitableArbitrages) {
	var gasEstimate int64 = 1750
	//var allRoutes = [][]SwapAmountInRoutesId{ab.ThreeCurrencyRoutes, ab.FourCurrencyRoutes, ab.FiveCurrencyRoutes}
	var allRoutes = [][]swaproutes.SwapAmountInRoutesId{ab.ps.TwoCurrencyRoutes, ab.ps.ThreeCurrencyRoutes, ab.ps.FourCurrencyRoutes}
	profitableRoutes = make([]swaproutes.ProfitableArbitrage, 0)

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
				Amount: ab.CalculateMultihopSwapExactAmountIn(route, inputCoin),
			}

			if outputCoin.IsGTE(inputCoin) {
				optimumIn, optimumOut := ab.FindOptimumFullScan(route)
				if optimumOut.GT(optimumIn.Amount.Add(sdk.NewInt(gasEstimate))) {
					profitableArbitrage := swaproutes.ProfitableArbitrage{
						Route:            route,
						OptimumInToken:   optimumIn,
						OptimumOutAmount: optimumOut,
					}

					// Check if there was no more profitable arb already executed.
					if !ab.executedProfRoutes.CheckIfMoreProfitableRouteWasAlreadyExecuted(profitableArbitrage) {
						// EXECUTE ARBITRAGE
						ab.GenerateAndSendToAllRPCEndpoints(profitableArbitrage)
					}

				}
			}
		}
	}

	return profitableRoutes
}
