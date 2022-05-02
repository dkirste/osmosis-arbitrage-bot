package arbbot

import (
	"fmt"
	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

func (ab *ArbBot) Run() {
	for _, grpcm := range ab.grpcms {
		grpcmPerLoop := grpcm
		go func() {
			var latestPools []gammtypes.PoolI
			var latestHeight uint64
			for {
				latestPools, latestHeight = grpcmPerLoop.QueryAllPools()
				if latestHeight > ab.currentHeight {
					fmt.Printf("%v,", latestHeight)
					ab.currentHeightMutex.Lock()
					ab.currentHeight = latestHeight
					ab.currentHeightMutex.Unlock()

					ab.psMutex.Lock()
					ab.ps.UpdatePools(latestPools)
					ab.psMutex.Unlock()

					ab.sequenceNumber = grpcmPerLoop.QueryAccountSequence(ab.ArbAddress)

				}
			}
		}()
	}
	return
}
