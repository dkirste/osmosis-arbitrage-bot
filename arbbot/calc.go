package arbbot

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dkirste/arbbot/swaproutes"
	balancer "github.com/osmosis-labs/osmosis/v13/x/gamm/pool-models/balancer"
	gammtypes "github.com/osmosis-labs/osmosis/v13/x/gamm/types"
)

func calculateSwapExactAmountIn(tokenIn sdk.Coin, tokenOutDenom string, inPoolAsset balancer.PoolAsset, outPoolAsset balancer.PoolAsset, swapFee sdk.Dec) (tokenOutAmount sdk.Int) {

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
	var inPoolAsset balancer.PoolAsset
	var outPoolAsset balancer.PoolAsset
	var swapFee sdk.Dec
	for _, route := range routes {
		inPoolAsset, _ = ab.ps.PoolsById[route.PoolId].GetPoolAsset(tokenIn.Denom)
		outPoolAsset, _ = ab.ps.PoolsById[route.PoolId].GetPoolAsset(route.TokenOutDenom)
		swapFee = ab.ps.PoolsById[route.PoolId].GetPoolParams().GetPoolSwapFee()
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

func (ab *ArbBot) FindOptimumBisection(arbitrageRoutes []swaproutes.SwapAmountInRouteId) (tokenIn sdk.Coin, tokenOutAmount sdk.Int) {
	var h = sdk.NewInt(10000) // Bisectioning 0.01 osmos left and right (h)
	var leftIn = sdk.Coin{
		Denom:  "uosmo",
		Amount: sdk.NewInt(100000), // using 0.1 osmo for minimum input.
	}
	var rightIn = ab.maxReserve
	var midIn = sdk.Coin{
		Denom:  "uosmo",
		Amount: leftIn.Amount.Add(rightIn.Amount).Quo(sdk.NewInt(2)), // (left + right) / 2
	}

	var midLeftProfit sdk.Int
	var midRightProfit sdk.Int
	var iterations = 20
	for i := 0; i < iterations; i++ {
		midIn.Amount = leftIn.Amount.Add(rightIn.Amount).Quo(sdk.NewInt(2))
		midLeftProfit = ab.CalculateMultihopSwapExactAmountIn(arbitrageRoutes, midIn.SubAmount(h)).Sub(midIn.SubAmount(h).Amount)  // mid - h  || out - in
		midRightProfit = ab.CalculateMultihopSwapExactAmountIn(arbitrageRoutes, midIn.AddAmount(h)).Sub(midIn.AddAmount(h).Amount) // mid + h  || out - in
		if midRightProfit.GT(midLeftProfit) {
			leftIn = midIn
		} else {
			rightIn = midIn
		}
	}
	tokenIn = midIn
	tokenOutAmount = ab.CalculateMultihopSwapExactAmountIn(arbitrageRoutes, midIn)

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
	for tmpTokenIn.IsLT(ab.maxReserve) {
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
		pool.GetPoolParams().GetPoolSwapFee(),
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
