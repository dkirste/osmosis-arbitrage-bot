package arbbot

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dkirste/arbbot/swaproutes"
)

func (ab *ArbBot) EvaluateOptimistic(ownId int, maxId int, poolIds []uint64) {
	var gasEstimate int64 = 2000
	var allRoutes = make([][]swaproutes.SwapAmountInRoutesId, 0)

	for _, id := range poolIds {
		allRoutes = append(allRoutes, ab.ps.ArbRoutesById[id])
	}

	for i, poolRoutes := range allRoutes {
		if i%maxId != ownId {
			continue
		}

		for _, route := range poolRoutes {
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
}
