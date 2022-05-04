package arbbot

import (
	"fmt"
	grpcMachine "github.com/dkirste/arbbot/grpcmachine"
	"github.com/dkirste/arbbot/swaproutes"
	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

func (ab *ArbBot) GenerateAndSendToAllRPCEndpoints(profRoute swaproutes.ProfitableArbitrage) {
	arbMsg := ab.BuildSwapExactAmountInMsg(ab.clientCtxs[0], profRoute.Route, profRoute.OptimumInToken, profRoute.OptimumInToken.Amount)

	for _, clientCtx := range ab.clientCtxs {
		clientCtxPerLoop := clientCtx
		// Safe seq otherwise it gets increased beforehand
		seq := ab.sequenceNumber
		go func() {
			err := ab.txm.GenerateBroadcastTx(clientCtxPerLoop, ab.currentHeight, seq, arbMsg)

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

func (ab *ArbBot) PoolUpdateLoop(grpcm grpcMachine.GrpcMachine, heightCh chan<- uint64) (crashed bool) {
	var latestPools []gammtypes.PoolI
	var latestHeight uint64
	for {
		defer func() {
			if err := recover(); err != nil {
				crashed = true
			}
		}()
		latestPools, latestHeight = grpcm.QueryAllPools()
		if latestHeight > ab.currentHeight {
			ab.currentHeightMutex.Lock()
			ab.currentHeight = latestHeight
			ab.currentHeightMutex.Unlock()

			ab.psMutex.Lock()
			ab.ps.UpdatePools(latestPools)
			ab.psMutex.Unlock()

			ab.sequenceNumber = grpcm.QueryAccountSequence(ab.ArbAddress)
			fmt.Printf("\n%v: ", grpcm.Conn.Target())
			heightCh <- latestHeight

			ab.maxReserve = grpcm.QueryAccountBalance(ab.ArbAddress, "uosmo")
		}
	}
}
