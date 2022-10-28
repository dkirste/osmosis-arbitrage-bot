package poolstorage

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gammtypes "github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

func (ps *PoolStorage) UpdatePoolOptimistically(swapMsg gammtypes.MsgSwapExactAmountIn) (involvedPools []uint64, err error) {
	involvedPools = make([]uint64, 0)
	routes := swapMsg.Routes
	tokenOutMinAmount := swapMsg.TokenOutMinAmount
	tokenIn := swapMsg.TokenIn
	var tokenOutAmount sdk.Int
	for i, route := range routes {
		involvedPools = append(involvedPools, route.PoolId)
		_outMinAmount := sdk.NewInt(1)
		if len(routes)-1 == i {
			_outMinAmount = tokenOutMinAmount
		}

		tokenOutAmount, err = ps.swapExactAmountIn(route.PoolId, tokenIn, route.TokenOutDenom, _outMinAmount)
		if err != nil {
			fmt.Printf("!")
			return nil, err
		}
		tokenIn = sdk.NewCoin(route.TokenOutDenom, tokenOutAmount)
	}
	return involvedPools, nil

}
