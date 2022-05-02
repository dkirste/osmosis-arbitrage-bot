package poolstorage

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

func (ps *PoolStorage) swapExactAmountIn(poolId uint64, tokenIn sdk.Coin, tokenOutDenom string, tokenOutMinAmount sdk.Int) (tokenOutAmount sdk.Int, err error) {
	pool := ps.PoolsById[poolId]
	inPoolAsset, _ := pool.GetPoolAsset(tokenIn.Denom)
	outPoolAsset, _ := pool.GetPoolAsset(tokenOutDenom)

	tokenOutAmount = calcOutGivenIn(
		inPoolAsset.Token.Amount.ToDec(),
		inPoolAsset.Weight.ToDec(),
		outPoolAsset.Token.Amount.ToDec(),
		outPoolAsset.Weight.ToDec(),
		tokenIn.Amount.ToDec(),
		pool.GetPoolSwapFee(),
	).TruncateInt()
	if tokenOutAmount.LTE(sdk.ZeroInt()) {
		return sdk.Int{}, sdkerrors.Wrapf(gammtypes.ErrInvalidMathApprox, "token amount is zero or negative")
	}

	if tokenOutAmount.LT(tokenOutMinAmount) {
		return sdk.Int{}, sdkerrors.Wrapf(gammtypes.ErrLimitMinAmount, "%s token is lesser than min amount", outPoolAsset.Token.Denom)
	}

	inPoolAsset.Token.Amount = inPoolAsset.Token.Amount.Add(tokenIn.Amount)
	outPoolAsset.Token.Amount = outPoolAsset.Token.Amount.Sub(tokenOutAmount)

	//fmt.Printf("Before: %v\n", ps.PoolsById[poolId])
	// Execute update on acutal pool
	err = ps.PoolsById[poolId].UpdatePoolAssetBalances(sdk.NewCoins(
		inPoolAsset.Token,
		outPoolAsset.Token,
	))
	if err != nil {
		fmt.Println("Could not update pool balances")
		return sdk.Int{}, err
	}
	//fmt.Printf("After: %v\n", ps.PoolsById[poolId])

	return tokenOutAmount, nil
}
