package arbbot

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dkirste/arbbot/swaproutes"
)

func (ab *ArbBot) EvaluateArbitrage(ownId int, maxId int) (profitableRoutes swaproutes.ProfitableArbitrages) {
	var gasEstimate int64 = 2000
	//var allRoutes = [][]SwapAmountInRoutesId{ab.ThreeCurrencyRoutes, ab.FourCurrencyRoutes, ab.FiveCurrencyRoutes}
	var allRoutes = [][]swaproutes.SwapAmountInRoutesId{ab.ps.ThreeCurrencyRoutes, ab.ps.FourCurrencyRoutes, ab.ps.FiveCurrencyRoutes}
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
						arbMsg := ab.BuildSwapExactAmountInMsg(clientCtx, route, optimumIn, optimumIn.Amount) // Using inputCoin.Amount to min(losses)

						err := ab.GenerateBroadcastTxEstimateGas(clientCtx, currentHeight, txm.SequenceNumber, arbMsg)
						if err != nil {
							fmt.Printf("Error:  %v\n", err)
						} else {
							ab.executedProfRoutesMutex.Lock()
							ab.executedProfRoutes = append(ab.executedProfRoutes, profitableArbitrage)
							ab.executedProfRoutesMutex.Unlock()
							ab.sequenceNumberMutex.Lock()
							ab.sequenceNumber = ab.sequenceNumber + 1
							ab.sequenceNumberMutex.Unlock()
						}
					}

				}
			}
		}
	}

	return profitableRoutes
}
