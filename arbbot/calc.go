package arbbot

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dkirste/arbbot/swaproutes"
	gammtypes "github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

func calculateSwapExactAmountIn(tokenIn sdk.Coin, tokenOutDenom string, inPoolAsset gammtypes.PoolAsset, outPoolAsset gammtypes.PoolAsset, swapFee sdk.Dec) (tokenOutAmount sdk.Int) {

	// TODO: Understand if we are handling swap fee consistently,
	// with the global swap fee and the pool swap fee

	tokenOutAmount = calcOutGivenIn(
		inPoolAsset.Token.Amount.ToDec(),
		inPoolAsset.Weight.ToDec(),
		outPoolAsset.Token.Amount.ToDec(),
		outPoolAsset.Weight.ToDec(),
		tokenIn.Amount.ToDec(),
		swapFee,
	).TruncateInt()

	return tokenOutAmount
}

func (ab *ArbBot) CalculateMultihopSwapExactAmountIn(routes swaproutes.SwapAmountInRoutesId, tokenIn sdk.Coin) (tokenOutAmount sdk.Int) {
	var inPoolAsset gammtypes.PoolAsset
	var outPoolAsset gammtypes.PoolAsset
	var swapFee sdk.Dec
	for _, route := range routes {
		inPoolAsset, _ = ab.ps.PoolsById[route.PoolId].GetPoolAsset(tokenIn.Denom)
		outPoolAsset, _ = ab.ps.PoolsById[route.PoolId].GetPoolAsset(route.TokenOutDenom)
		swapFee = ab.ps.PoolsById[route.PoolId].GetPoolSwapFee()
		tokenOutAmount = calculateSwapExactAmountIn(
			tokenIn,             // tokenIn
			route.TokenOutDenom, // tokenOutDenom
			inPoolAsset,         // inPoolAsset
			outPoolAsset,        // outPoolAsset
			swapFee,             // swapFee
		)
		tokenIn = sdk.NewCoin(route.TokenOutDenom, tokenOutAmount)
	}
	return
}

func (ab *ArbBot) FindOptimumNew(arbitrageRoutes []swaproutes.SwapAmountInRouteId) (tokenIn sdk.Coin, tokenOutAmount sdk.Int) {

	return tokenIn, tokenOutAmount
}

func (ab *ArbBot) FindOptimumFullScan(arbitrageRoutes []swaproutes.SwapAmountInRouteId) (tokenIn sdk.Coin, tokenOutAmount sdk.Int) {
	var profitsArray []sdk.Int = make([]sdk.Int, 0)
	var adjustment = sdk.NewInt(1000000)
	var tmpTokenIn sdk.Coin
	var tmpTokenOutAmount sdk.Int
	var tmpProfit sdk.Int

	tmpTokenIn = sdk.Coin{
		Denom:  "uosmo",
		Amount: adjustment,
	}
	for tmpTokenIn.IsLTE(ab.maxReserve) {
		tmpTokenOutAmount = ab.CalculateMultihopSwapExactAmountIn(arbitrageRoutes, tmpTokenIn)
		tmpProfit = tmpTokenOutAmount.Sub(tmpTokenIn.Amount)
		profitsArray = append(profitsArray, tmpProfit)

		// Increase for next iteration
		tmpTokenIn = tmpTokenIn.AddAmount(adjustment)
	}

	i := GetIndexAtMax(profitsArray)
	tokenIn = sdk.Coin{
		Denom:  "uosmo",
		Amount: adjustment.Mul(sdk.NewInt(i)).Add(adjustment),
	}
	tokenOutAmount = tokenIn.Amount.Add(profitsArray[i])

	return tokenIn, tokenOutAmount
}

func GetIndexAtMax(array []sdk.Int) int64 {
	maxValue := array[0]
	maxIndex := 0
	for i, elem := range array {
		if elem.GT(maxValue) {
			maxValue = elem
			maxIndex = i
		}
	}
	return int64(maxIndex)
}

func (ab *ArbBot) swapExactAmountIn(poolId uint64, tokenIn sdk.Coin, tokenOutDenom string, tokenOutMinAmount sdk.Int) (tokenOutAmount sdk.Int, err error) {
	pool := ab.ps.PoolsById[poolId]
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
	err = ab.ps.PoolsById[poolId].UpdatePoolAssetBalances(sdk.NewCoins(
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
