package txmachine

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

func (txm *TxMachine) GenerateBroadcastTx(clientCtx client.Context, currentHeight uint64, seq uint64, msg sdk.Msg) error {
	txm.Factory = txm.Factory.
		WithGas(350000).
		WithGasAdjustment(1).
		WithGasPrices("0.005uosmo").
		WithTimeoutHeight(currentHeight + 2).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT).
		WithSequence(seq)
	err := tx.GenerateOrBroadcastTxWithFactory(clientCtx, txm.Factory, msg)
	if err != nil {
		fmt.Println("Error while generating or broadcasting tx.")
		fmt.Println(err)
		return err
	}
	return nil
}
