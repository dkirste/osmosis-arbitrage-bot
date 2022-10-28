package poolstorage

import (
	gammtypes "github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

func (ps *PoolStorage) GetPoolIdInPoolsByAssets(assetId int, poolId uint64) int {
	// Check if there is no pool registered for the asset
	if len(ps.PoolsByAsset[assetId]) == 0 {
		ps.PoolsByAsset[assetId] = make([]gammtypes.PoolI, 1)
		return 0
	}

	for i, pool := range ps.PoolsByAsset[assetId] {
		if pool.GetId() == poolId {
			return i
		}
	}

	// Pool not found in PoolsByAsset -> append new pool
	ps.PoolsByAsset[assetId] = append(ps.PoolsByAsset[assetId], nil)
	return len(ps.PoolsByAsset[assetId]) - 1
}

func (ps *PoolStorage) CheckIfWhitelisted(poolId uint64) bool {
	for _, whitelistedPoolId := range ps.Whitelist {
		if poolId == whitelistedPoolId {
			return true
		}
	}
	return false
}
