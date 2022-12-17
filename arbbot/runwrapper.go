package arbbot

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/server/rosetta"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpcMachine "github.com/dkirste/arbbot/grpcmachine"
	"github.com/dkirste/arbbot/swaproutes"
	balancer "github.com/osmosis-labs/osmosis/v13/x/gamm/pool-models/balancer"
	gammtypes "github.com/osmosis-labs/osmosis/v13/x/gamm/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"math"
)

func (ab *ArbBot) GenerateAndSendToAllRPCEndpoints(profRoute swaproutes.ProfitableArbitrage) {
	numPools := len(profRoute.Route)
	arbMsg := ab.BuildSwapExactAmountInMsg(ab.clientCtxs[0], profRoute.Route, profRoute.OptimumInToken, profRoute.OptimumInToken.Amount)

	for _, clientCtx := range ab.clientCtxs {
		clientCtxPerLoop := clientCtx
		// Safe seq otherwise it gets increased beforehand
		seq := ab.sequenceNumber
		go func() {
			err := ab.txm.GenerateBroadcastTx(clientCtxPerLoop, ab.currentHeight, seq, arbMsg, numPools)

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
	var latestPools []balancer.Pool
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

func (ab *ArbBot) ScanMempoolLoop(rpcConnMem *rpchttp.HTTP, scanThreshold int64, numArbWorkers int) (crashed bool) {
	defer func() {
		if err := recover(); err != nil {
			crashed = true
		}
	}()

	var limit = 1000
	var tokenIn sdk.Coin
	var tokenInAmountUSD sdk.Int

	var c = rosetta.NewConverter(ab.ProtoCodec, ab.InterfaceRegistry, ab.EncodingConfig.TxConfig)
	var ctx = context.Background()

	// Run forever
	for {
		txs, err := rpcConnMem.UnconfirmedTxs(ctx, &limit)
		if err != nil {
			fmt.Printf("API banned by %v\n", rpcConnMem)
			return true
		}

	unconfirmedTxLoop:
		for _, unconfirmedTx := range txs.Txs {
			for _, analysedTx := range ab.analysedTxs {
				if analysedTx == string(unconfirmedTx.Hash()) {
					continue unconfirmedTxLoop
				}
			}
			// Add tx to analysed txs
			ab.analysedTxsMutex.Lock()
			ab.analysedTxs = append(ab.analysedTxs, string(unconfirmedTx.Hash()))
			ab.analysedTxsMutex.Unlock()

			fmt.Printf(".")
			transaction, err := c.ToRosetta().Tx(unconfirmedTx, nil)
			if err != nil {
				continue
			}

			unsignedTx, _ := c.ToSDK().UnsignedTx(transaction.Operations)
			for _, msg := range unsignedTx.GetMsgs() {
				switch swapExactAmountInMsg := msg.(type) {
				case *gammtypes.MsgSwapExactAmountIn:
					tokenIn = swapExactAmountInMsg.TokenIn
					if _, ok := ab.priceOracle[tokenIn.Denom]; !ok {
						// Denom was not included in price oracle.
						fmt.Printf("\nDenom of token not found in priceOracle (%v)\n", tokenIn.Denom)
						fmt.Printf("#")
						involvedPools, failed := ab.ps.UpdatePoolOptimistically(*swapExactAmountInMsg)
						if failed == nil {
							// Scan for arb and execute!
							for workerId := 0; workerId < numArbWorkers; workerId++ {
								go ab.EvaluateOptimistic(workerId, numArbWorkers, involvedPools)
							}
						}
					}

					// Check for arbitrage txs
					if swapExactAmountInMsg.TokenIn.Denom == swapExactAmountInMsg.TokenOutDenom() {
						//fmt.Println("Arbitrage transaction in mempool spotted!!")
						continue unconfirmedTxLoop
					}

					tokenInAmountUSD = tokenIn.Amount.Mul(sdk.NewInt(ab.priceOracle[tokenIn.Denom].PriceE6)).Quo(sdk.NewInt(int64(math.Pow10(ab.priceOracle[tokenIn.Denom].Exponent + 6))))

					if tokenInAmountUSD.GT(sdk.NewInt(scanThreshold)) {
						fmt.Printf("#")
						//fmt.Printf("\n%v\n", swapExactAmountInMsg)
						involvedPools, failed := ab.ps.UpdatePoolOptimistically(*swapExactAmountInMsg)
						if failed == nil {
							// Scan for arb and execute!
							for workerId := 0; workerId < numArbWorkers; workerId++ {
								go ab.EvaluateOptimistic(workerId, numArbWorkers, involvedPools)
							}
						}
					}
				}
			}

		}
	}
}
