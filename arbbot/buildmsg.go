package arbbot

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dkirste/arbbot/swaproutes"
	pooltypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

func (ab *ArbBot) BuildSwapExactAmountInMsg(clientCtx client.Context, arbitrageRoutes swaproutes.SwapAmountInRoutesId, tokenIn sdk.Coin, tokenOutMinAmount sdk.Int) sdk.Msg {
	addressKey, err := clientCtx.Keyring.Key(clientCtx.From)
	if err != nil {
		fmt.Println(err)
		panic("Could not derive public key in bech32")
	}
	addressBytes, err := addressKey.GetAddress()
	if err != nil {
		fmt.Println(err)
		panic("Could not derive public key in bech32")
	}
	addressBech32, err := sdk.Bech32ifyAddressBytes("osmo", addressBytes)
	if err != nil {
		fmt.Println(err)
		panic("Could not derive public key in bech32")
	}

	swapExactAmountInMsg := pooltypes.MsgSwapExactAmountIn{
		Sender:            addressBech32,
		Routes:            arbitrageRoutes.ConvertToPooltype(),
		TokenIn:           tokenIn,
		TokenOutMinAmount: tokenOutMinAmount,
	} /*
		res, err := poolMsgClient.SwapExactAmountIn(ctx, &swapExactAmountInMsg)
		if err != nil {
			fmt.Println("Error while executing SwapExactAmountIn transaction.")
		}
		_ = res*/
	return &swapExactAmountInMsg
}
