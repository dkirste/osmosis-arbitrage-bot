package grpcMachine

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v13/x/gamm/pool-models/balancer"
)

// PoolI defines an interface for pools that hold tokens.
/*
type BPoolI interface {
	GetAddress() sdk.AccAddress
	GetId() uint64
	GetSwapFee(_ sdk.Context) sdk.Dec
	GetTotalPoolLiquidity(_ sdk.Context) sdk.Coins
	GetExitFee(_ sdk.Context) sdk.Dec
	GetPoolParams() balancer.PoolParams
	GetTotalWeight() sdk.Int
	GetTotalShares() sdk.Int
	AddTotalShares(amt sdk.Int)
	SubTotalShares(amt sdk.Int)
	SetInitialPoolAssets(PoolAssets []balancer.PoolAsset)
	//setInitialPoolParams(params PoolParams, sortedAssets []PoolAsset, curBlockTime time.Time)
	GetPoolAsset(denom string) (balancer.PoolAsset, error)
	//getPoolAssetAndIndex(denom string) (int, balancer.PoolAsset, error)
	//parsePoolAssetsByDenoms(tokenADenom, tokenBDenom string)
	//parsePoolAssets(tokensA sdk.Coins, tokenBDenom string)
	//parsePoolAssetsCoins(tokensA sdk.Coins, tokensB sdk.Coins)
	IncreaseLiquidity(sharesOut sdk.Int, coinsIn sdk.Coins)
	UpdatePoolAssetBalance(coin sdk.Coin) error
	UpdatePoolAssetBalances(coins sdk.Coins) error
	//addToPoolAssetBalances(coins sdk.Coins) error
	GetPoolAssets(denoms ...string) ([]balancer.PoolAsset, error)
	GetAllPoolAssets() []balancer.PoolAsset
	//updateAllWeights(newWeights []balancer.PoolAsset)
	PokePool(blockTime time.Time)
	GetTokenWeight(denom string) (sdk.Int, error)
	GetTokenBalance(denom string) (sdk.Int, error)
	NumAssets() int
	IsActive(ctx sdk.Context) bool
	CalcOutAmtGivenIn(
		ctx sdk.Context,
		tokensIn sdk.Coins,
		tokenOutDenom string,
		swapFee sdk.Dec,
	) (sdk.Coin, error)
	SwapOutAmtGivenIn(
		ctx sdk.Context,
		tokensIn sdk.Coins,
		tokenOutDenom string,
		swapFee sdk.Dec,
	) (
		tokenOut sdk.Coin, err error,
	)
	CalcInAmtGivenOut(
		ctx sdk.Context, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (
		tokenIn sdk.Coin, err error,
	)
	SwapInAmtGivenOut(
		ctx sdk.Context, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (
		tokenIn sdk.Coin, err error,
	)
	//applySwap(ctx sdk.Context, tokensIn sdk.Coins, tokensOut sdk.Coins) error
	SpotPrice(ctx sdk.Context, baseAsset, quoteAsset string) (spotPrice sdk.Dec, err error)
	//calcSingleAssetJoin(tokenIn sdk.Coin, swapFee sdk.Dec, tokenInPoolAsset balancer.PoolAsset, totalShares sdk.Int) (numShares sdk.Int, err error)
	JoinPool(ctx sdk.Context, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares sdk.Int, err error)
	JoinPoolNoSwap(ctx sdk.Context, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares sdk.Int, err error)
	CalcJoinPoolShares(ctx sdk.Context, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares sdk.Int, tokensJoined sdk.Coins, err error)
	CalcJoinPoolNoSwapShares(ctx sdk.Context, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares sdk.Int, tokensJoined sdk.Coins, err error)
	//calcJoinSingleAssetTokensIn(tokensIn sdk.Coins, totalShares sdk.Int, poolAssetsByDenom map[string]balancer.PoolAsset, swapFee sdk.Dec) (sdk.Int, sdk.Coins, error)
	ExitPool(ctx sdk.Context, exitingShares sdk.Int, exitFee sdk.Dec) (exitingCoins sdk.Coins, err error)
	//exitPool(ctx sdk.Context, exitingCoins sdk.Coins, exitingShares sdk.Int) error
	CalcExitPoolCoinsFromShares(ctx sdk.Context, exitingShares sdk.Int, exitFee sdk.Dec) (exitedCoins sdk.Coins, err error)
	CalcTokenInShareAmountOut(
		ctx sdk.Context,
		tokenInDenom string,
		shareOutAmount sdk.Int,
		swapFee sdk.Dec,
	) (tokenInAmount sdk.Int, err error)
	JoinPoolTokenInMaxShareAmountOut(
		ctx sdk.Context,
		tokenInDenom string,
		shareOutAmount sdk.Int,
	) (tokenInAmount sdk.Int, err error)
	ExitSwapExactAmountOut(
		ctx sdk.Context,
		tokenOut sdk.Coin,
		shareInMaxAmount sdk.Int,
	) (shareInAmount sdk.Int, err error)
}
*/

type BPoolI interface {

	// Backported functions
	GetPoolAssets(denoms ...string) ([]balancer.PoolAsset, error)
	GetAllPoolAssets() []balancer.PoolAsset
	GetTokenWeight(denom string) (sdk.Int, error)

	GetAddress() sdk.AccAddress
	String() string
	GetId() uint64
	// GetSwapFee returns the pool's swap fee, based on the current state.
	// Pools may choose to make their swap fees dependent upon state
	// (prior TWAPs, network downtime, other pool states, etc.)
	// hence Context is provided as an argument.
	GetSwapFee(ctx sdk.Context) sdk.Dec
	// GetExitFee returns the pool's exit fee, based on the current state.
	// Pools may choose to make their exit fees dependent upon state.
	GetExitFee(ctx sdk.Context) sdk.Dec
	// Returns whether the pool has swaps enabled at the moment
	IsActive(ctx sdk.Context) bool
	// GetTotalPoolLiquidity returns the coins in the pool owned by all LPs
	GetTotalPoolLiquidity(ctx sdk.Context) sdk.Coins
	// GetTotalShares returns the total number of LP shares in the pool
	GetTotalShares() sdk.Int

	// SwapOutAmtGivenIn swaps 'tokenIn' against the pool, for tokenOutDenom, with the provided swapFee charged.
	// Balance transfers are done in the keeper, but this method updates the internal pool state.
	SwapOutAmtGivenIn(ctx sdk.Context, tokenIn sdk.Coins, tokenOutDenom string, swapFee sdk.Dec) (tokenOut sdk.Coin, err error)
	// CalcOutAmtGivenIn returns how many coins SwapOutAmtGivenIn would return on these arguments.
	// This does not mutate the pool, or state.
	CalcOutAmtGivenIn(ctx sdk.Context, tokenIn sdk.Coins, tokenOutDenom string, swapFee sdk.Dec) (tokenOut sdk.Coin, err error)

	// SwapInAmtGivenOut swaps exactly enough tokensIn against the pool, to get the provided tokenOut amount out of the pool.
	// Balance transfers are done in the keeper, but this method updates the internal pool state.
	SwapInAmtGivenOut(ctx sdk.Context, tokenOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (tokenIn sdk.Coin, err error)
	// CalcInAmtGivenOut returns how many coins SwapInAmtGivenOut would return on these arguments.
	// This does not mutate the pool, or state.
	CalcInAmtGivenOut(ctx sdk.Context, tokenOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (tokenIn sdk.Coin, err error)

	// Returns the spot price of the 'base asset' in terms of the 'quote asset' in the pool,
	// errors if either baseAssetDenom, or quoteAssetDenom does not exist.
	// For example, if this was a UniV2 50-50 pool, with 2 ETH, and 8000 UST
	// pool.SpotPrice(ctx, "eth", "ust") = 4000.00
	SpotPrice(ctx sdk.Context, baseAssetDenom string, quoteAssetDenom string) (sdk.Dec, error)

	// JoinPool joins the pool using all of the tokensIn provided.
	// The AMM swaps to the correct internal ratio should be and returns the number of shares created.
	// This function is mutative and updates the pool's internal state if there is no error.
	// It is up to pool implementation if they support LP'ing at arbitrary ratios, or a subset of ratios.
	// Pools are expected to guarantee LP'ing at the exact ratio, and single sided LP'ing.
	JoinPool(ctx sdk.Context, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares sdk.Int, err error)

	// JoinPoolNoSwap joins the pool with an all-asset join using the maximum amount possible given the tokensIn provided.
	// This function is mutative and updates the pool's internal state if there is no error.
	// Pools are expected to guarantee LP'ing at the exact ratio.
	JoinPoolNoSwap(ctx sdk.Context, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares sdk.Int, err error)

	// CalcJoinPoolShares returns how many LP shares JoinPool would return on these arguments.
	// This does not mutate the pool, or state.
	CalcJoinPoolShares(ctx sdk.Context, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares sdk.Int, newLiquidity sdk.Coins, err error)

	// CalcJoinPoolNoSwapShares returns how many LP shares JoinPoolNoSwap would return on these arguments.
	// This does not mutate the pool, or state.
	CalcJoinPoolNoSwapShares(ctx sdk.Context, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares sdk.Int, newLiquidity sdk.Coins, err error)

	// ExitPool exits #numShares LP shares from the pool, decreases its internal liquidity & LP share totals,
	// and returns the number of coins that are being returned.
	// This mutates the pool and state.
	ExitPool(ctx sdk.Context, numShares sdk.Int, exitFee sdk.Dec) (exitedCoins sdk.Coins, err error)
	// CalcExitPoolCoinsFromShares returns how many coins ExitPool would return on these arguments.
	// This does not mutate the pool, or state.
	CalcExitPoolCoinsFromShares(ctx sdk.Context, numShares sdk.Int, exitFee sdk.Dec) (exitedCoins sdk.Coins, err error)
}
