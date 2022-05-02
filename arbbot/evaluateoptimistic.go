package arbbot

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dkirste/arbbot/swaproutes"
	"github.com/dkirste/arbbot/txmachine"
)

func (ab *ArbBot) EvaluateOptimisticAndExecute(poolIds []uint64, txm *txmachine.TxMachine, clientCtx client.Context, currentHeight uint64, seq uint64) {
	var gasEstimate int64 = 2000
	//var allRoutes = [][]SwapAmountInRoutesId{am.ThreeCurrencyRoutes, am.FourCurrencyRoutes, am.FiveCurrencyRoutes}
	var allRoutes = make([][]swaproutes.SwapAmountInRoutesId, 0)
	for _, id := range poolIds {
		allRoutes = append(allRoutes, ab.ps.ArbRoutesById[id])
	}

	for _, poolRoutes := range allRoutes {

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
					// EXECUTE ARBITRAGE ASAP

					arbMsg := ab.BuildSwapExactAmountInMsg(clientCtx, route, optimumIn, optimumIn.Amount)
					err := txm.GenerateBroadcastTxEstimateGas(clientCtx, currentHeight, seq, arbMsg)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					} else {
						seq = seq + 1
					}

				}
			}
		}

	}
}
