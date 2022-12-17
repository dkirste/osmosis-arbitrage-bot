package poolstorage

import (
	"github.com/dkirste/arbbot/swaproutes"
	balancer "github.com/osmosis-labs/osmosis/v13/x/gamm/pool-models/balancer"
)

type PoolStorage struct {
	PoolsById           []balancer.Pool
	PoolsByAsset        [][]balancer.Pool
	TwoCurrencyRoutes   []swaproutes.SwapAmountInRoutesId
	ThreeCurrencyRoutes []swaproutes.SwapAmountInRoutesId
	FourCurrencyRoutes  []swaproutes.SwapAmountInRoutesId
	FiveCurrencyRoutes  []swaproutes.SwapAmountInRoutesId
	ArbRoutesById       [][]swaproutes.SwapAmountInRoutesId
	AssetDict           AssetDict
	Whitelist           []uint64
}

func (ps *PoolStorage) Setup(numberOfPools int) {
	ps.PoolsById = make([]balancer.Pool, numberOfPools)
	ps.PoolsByAsset = make([][]balancer.Pool, numberOfPools)
	ps.ArbRoutesById = make([][]swaproutes.SwapAmountInRoutesId, numberOfPools)
}
