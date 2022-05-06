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
				fmt.Printf("GRPC banned by: %v\n", grpcmPerLoop.Conn.Target())
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

func (ab *ArbBot) RunMempoolArb(numArbWorkers int, scanThreshold int64) {
	heightCh := make(chan uint64)
	for _, grpcm := range ab.grpcms {
		grpcmPerLoop := grpcm
		go func() {
			for {
				_ = ab.PoolUpdateLoop(grpcmPerLoop, heightCh)
			}
		}()
	}

	for _, mempoolRPC := range ab.mempoolRPCs {
		mempoolRPCPerLoop := mempoolRPC
		go func() {
			for {
				_ = ab.ScanMempoolLoop(mempoolRPCPerLoop, scanThreshold, numArbWorkers)
			}
		}()
	}

	for height := range heightCh {
		fmt.Printf("%v,", height)

		// Clear optimistic routes every 100 blocks
		if height%100 == 0 {
			ab.executedOptimisticRoutes = make(swaproutes.ProfitableArbitrages, 0)
		}

		ab.analysedTxs = make([]string, 0)
	}
	return
}

func (ab *ArbBot) RunFullArb(numArbWorkers int, scanThreshold int64) {
	heightCh := make(chan uint64)
	for _, grpcm := range ab.grpcms {
		grpcmPerLoop := grpcm
		go func() {
			for {
				_ = ab.PoolUpdateLoop(grpcmPerLoop, heightCh)
			}
		}()
	}

	for _, mempoolRPC := range ab.mempoolRPCs {
		mempoolRPCPerLoop := mempoolRPC
		go func() {
			for {
				_ = ab.ScanMempoolLoop(mempoolRPCPerLoop, scanThreshold, numArbWorkers)
			}
		}()
	}

	for height := range heightCh {
		fmt.Printf("%v,", height)

		// Clear optimistic routes every 100 blocks
		if height%100 == 0 {
			ab.executedOptimisticRoutes = make(swaproutes.ProfitableArbitrages, 0)
		}

		ab.analysedTxs = make([]string, 0)
		for workerId := 0; workerId < numArbWorkers; workerId++ {
			go ab.EvaluateArbitrage(workerId, numArbWorkers)
		}
	}
	return
}
