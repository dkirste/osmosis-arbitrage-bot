package poolstorage

import gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"

type SwapAmountInRoutes []SwapAmountInRoute

type SwapAmountInRoute struct {
	Pool          gammtypes.PoolI
	TokenOutDenom string
}

type SwapAmountInRoutesId []SwapAmountInRouteId
type SwapAmountInRouteId struct {
	PoolId        uint64
	TokenOutDenom string
}

func (swapRoutes SwapAmountInRoutes) ConvertToPooltype() (gammtypesSwapRoutes gammtypes.SwapAmountInRoutes) {
	var gammtypesRoute gammtypes.SwapAmountInRoute
	for _, route := range swapRoutes {
		gammtypesRoute = gammtypes.SwapAmountInRoute{
			PoolId:        route.Pool.GetId(),
			TokenOutDenom: route.TokenOutDenom,
		}
		gammtypesSwapRoutes = append(gammtypesSwapRoutes, gammtypesRoute)
	}
	return
}

func (swapRoutes SwapAmountInRoutes) CheckIfPoolsAreUnique() bool {
	var poolIds = make([]uint64, 0)
	var routePoolId uint64

	for _, route := range swapRoutes {
		routePoolId = route.Pool.GetId()
		for _, poolId := range poolIds {
			if routePoolId == poolId {
				return false
			}
		}
		poolIds = append(poolIds, routePoolId)
	}
	return true
}

func (swapRoutes SwapAmountInRoutesId) ConvertToPooltype() (gammtypesSwapRoutes gammtypes.SwapAmountInRoutes) {
	var gammtypesRoute gammtypes.SwapAmountInRoute
	for _, route := range swapRoutes {
		gammtypesRoute = gammtypes.SwapAmountInRoute{
			PoolId:        route.PoolId,
			TokenOutDenom: route.TokenOutDenom,
		}
		gammtypesSwapRoutes = append(gammtypesSwapRoutes, gammtypesRoute)
	}
	return
}

func (swapRoutes SwapAmountInRoutesId) CheckIfPoolsAreUnique() bool {
	var poolIds = make([]uint64, 0)

	for _, route := range swapRoutes {
		for _, poolId := range poolIds {
			if route.PoolId == poolId {
				return false
			}
		}
		poolIds = append(poolIds, route.PoolId)
	}
	return true
}
