package poolstorage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dkirste/arbbot/swaproutes"
	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

type PoolStorage struct {
	PoolsById           []gammtypes.PoolI
	PoolsByAsset        [][]gammtypes.PoolI
	ThreeCurrencyRoutes []swaproutes.SwapAmountInRoutesId
	FourCurrencyRoutes  []swaproutes.SwapAmountInRoutesId
	FiveCurrencyRoutes  []swaproutes.SwapAmountInRoutesId
	ArbRoutesById       [][]swaproutes.SwapAmountInRoutesId
	AssetDict           AssetDict
	MaxReserve          sdk.Coin
	ReserveThreshold    sdk.Int
	Whitelist           []uint64
}

func (ps *PoolStorage) Setup(numberOfPools int) {
	ps.PoolsById = make([]gammtypes.PoolI, numberOfPools)
	ps.PoolsByAsset = make([][]gammtypes.PoolI, numberOfPools)
	ps.ArbRoutesById = make([][]swaproutes.SwapAmountInRoutesId, numberOfPools)
}
