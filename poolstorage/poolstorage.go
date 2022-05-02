package poolstorage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

type PoolStorage struct {
	PoolsById               []gammtypes.PoolI
	PoolsByAsset            [][]gammtypes.PoolI
	ThreeCurrencyRoutes     []SwapAmountInRoutesId
	FourCurrencyRoutes      []SwapAmountInRoutesId
	FiveCurrencyRoutes      []SwapAmountInRoutesId
	ThreeCurrencyRoutesById [][]SwapAmountInRoutesId
	AssetDict               AssetDict
	MaxReserve              sdk.Coin
	ReserveThreshold        sdk.Int
	Whitelist               []uint64
}

func (ps *PoolStorage) Setup(numberOfPools int) {
	ps.PoolsById = make([]gammtypes.PoolI, numberOfPools)
	ps.PoolsByAsset = make([][]gammtypes.PoolI, numberOfPools)
	ps.ThreeCurrencyRoutesById = make([][]SwapAmountInRoutesId, numberOfPools)
}
