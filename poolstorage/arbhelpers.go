package poolstorage

import (
	balancer "github.com/osmosis-labs/osmosis/v13/x/gamm/pool-models/balancer"
)

func (ps *PoolStorage) GetPoolIdInPoolsByAssets(assetId int, poolId uint64) int {
	// Check if there is no pool registered for the asset
	if len(ps.PoolsByAsset[assetId]) == 0 {
		ps.PoolsByAsset[assetId] = make([]balancer.Pool, 1)
		return 0
	}

	for i, pool := range ps.PoolsByAsset[assetId] {
		if pool.GetId() == poolId {
			return i
		}
	}

	// Pool not found in PoolsByAsset -> append new pool
	ps.PoolsByAsset[assetId] = append(ps.PoolsByAsset[assetId], balancer.Pool{})
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
