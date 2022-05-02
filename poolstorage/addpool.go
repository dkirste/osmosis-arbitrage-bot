package poolstorage

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

func (ps *PoolStorage) addPool(poolToAdd gammtypes.PoolI) {
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

func (ps *PoolStorage) updatePool(poolToAdd gammtypes.PoolI) {
	poolId := poolToAdd.GetId()
	if int(poolId) > len(ps.PoolsById) {
		fmt.Printf("Pool with id: %v could not be added. Index out of bound (>%v).\n", poolId, len(ps.PoolsById))
		return
	} else {
		// Add pool to PoolsById
		ps.PoolsById[poolId] = poolToAdd
	}
}

func (ps *PoolStorage) AddPools(poolsToAdd []gammtypes.PoolI) {
	for _, pool := range poolsToAdd {
		if ps.CheckIfWhitelisted(pool.GetId()) {
			ps.addPool(pool)
		}
	}
}

func (ps *PoolStorage) UpdatePools(poolsToUpdate []gammtypes.PoolI) {
	for _, pool := range poolsToUpdate {
		ps.updatePool(pool)
	}
}

func (ps *PoolStorage) UpdateMaxReserve(newReserve sdk.Coin) {
	ps.MaxReserve = newReserve
	return
}
