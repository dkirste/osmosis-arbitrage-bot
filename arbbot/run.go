package arbbot

import (
	"fmt"

	grpcMachine "github.com/dkirste/arbbot/grpcmachine"
	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

func (ab *ArbBot) Run() {
	heightCh := make(chan uint64)
	for _, grpcm := range ab.grpcms {
		grpcmPerLoop := grpcm
		go func() {
			for {
				_ = ab.PoolUpdateLoop(grpcmPerLoop, heightCh)
			}
		}()
	}
	for height := range heightCh {
		fmt.Printf("%v,", height)
	}
	return
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
		}
	}
}
