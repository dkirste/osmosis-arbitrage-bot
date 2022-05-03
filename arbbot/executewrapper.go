package arbbot

import (
	"fmt"
	"github.com/dkirste/arbbot/swaproutes"
)

func (ab *ArbBot) GenerateAndSendToAllRPCEndpoints(profRoute swaproutes.ProfitableArbitrage) {
	arbMsg := ab.BuildSwapExactAmountInMsg(ab.clientCtxs[0], profRoute.Route, profRoute.OptimumInToken, profRoute.OptimumInToken.Amount)

	for _, clientCtx := range ab.clientCtxs {
		clientCtxPerLoop := clientCtx
		go func() {
			err := ab.txm.GenerateBroadcastTx(clientCtxPerLoop, ab.currentHeight, ab.sequenceNumber, arbMsg)

			if err != nil {
				fmt.Printf("\nCould not send tx to:  %v\n", clientCtxPerLoop.Client)
			}
		}()
	}
	ab.executedProfRoutesMutex.Lock()
	ab.executedProfRoutes = append(ab.executedProfRoutes, profRoute)
	ab.executedProfRoutesMutex.Unlock()
	ab.sequenceNumberMutex.Lock()
	ab.sequenceNumber = ab.sequenceNumber + 1
	ab.sequenceNumberMutex.Unlock()
}
