package swaproutes

import gammtypes "github.com/osmosis-labs/osmosis/v13/x/gamm/types"

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
