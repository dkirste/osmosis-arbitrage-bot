package arbbot

import (
	"fmt"
	"github.com/dkirste/arbbot/swaproutes"
)

func (ab *ArbBot) RunBlockArb(numArbWorkers int) {
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
		ab.executedProfRoutes = make(swaproutes.ProfitableArbitrages, 0)
		for workerId := 0; workerId < numArbWorkers; workerId++ {
			go ab.EvaluateArbitrage(workerId, numArbWorkers)
		}
	}
	return
}
