package poolstorage

import (
	"github.com/dkirste/arbbot/swaproutes"
	gammtypes "github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

type PoolStorage struct {
	PoolsById           []gammtypes.PoolI
	PoolsByAsset        [][]gammtypes.PoolI
	TwoCurrencyRoutes   []swaproutes.SwapAmountInRoutesId
	ThreeCurrencyRoutes []swaproutes.SwapAmountInRoutesId
	FourCurrencyRoutes  []swaproutes.SwapAmountInRoutesId
	FiveCurrencyRoutes  []swaproutes.SwapAmountInRoutesId
	ArbRoutesById       [][]swaproutes.SwapAmountInRoutesId
	AssetDict           AssetDict
	Whitelist           []uint64
}

func (ps *PoolStorage) Setup(numberOfPools int) {
	ps.PoolsById = make([]gammtypes.PoolI, numberOfPools)
	ps.PoolsByAsset = make([][]gammtypes.PoolI, numberOfPools)
	ps.ArbRoutesById = make([][]swaproutes.SwapAmountInRoutesId, numberOfPools)
}
