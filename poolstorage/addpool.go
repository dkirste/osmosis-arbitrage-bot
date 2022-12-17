package poolstorage

import (
	"fmt"
	balancer "github.com/osmosis-labs/osmosis/v13/x/gamm/pool-models/balancer"
)

func (ps *PoolStorage) addPool(poolToAdd balancer.Pool) {
	poolId := poolToAdd.GetId()
	if int(poolId) > len(ps.PoolsById) {
		fmt.Printf("Pool with id: %v could not be added. Index out of bound (>%v).\n", poolId, len(ps.PoolsById))
		return
	} else {
		// Add pool to PoolsById
		ps.PoolsById[poolId] = poolToAdd

		// Add pool to PoolsByAsset
		for _, asset := range poolToAdd.GetAllPoolAssets() {
			assetId := ps.AssetDict.GetId(asset.Token.Denom)
			poolIdByAsset := ps.GetPoolIdInPoolsByAssets(assetId, poolId)
			ps.PoolsByAsset[assetId][poolIdByAsset] = poolToAdd
		}
	}
}

func (ps *PoolStorage) updatePool(poolToAdd balancer.Pool) {
	poolId := poolToAdd.GetId()
	if int(poolId) > len(ps.PoolsById) {
		fmt.Printf("Pool with id: %v could not be added. Index out of bound (>%v).\n", poolId, len(ps.PoolsById))
		return
	} else {
		// Add pool to PoolsById
		ps.PoolsById[poolId] = poolToAdd
	}
}

func (ps *PoolStorage) AddPools(poolsToAdd []balancer.Pool) {
	for _, pool := range poolsToAdd {
		if ps.CheckIfWhitelisted(pool.GetId()) {
			ps.addPool(pool)
		}
	}
}

func (ps *PoolStorage) UpdatePools(poolsToUpdate []balancer.Pool) {
	for _, pool := range poolsToUpdate {
		ps.updatePool(pool)
	}
}
