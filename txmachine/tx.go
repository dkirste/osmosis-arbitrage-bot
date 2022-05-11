package txmachine

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

func (txm *TxMachine) GenerateBroadcastTx(clientCtx client.Context, currentHeight uint64, seq uint64, msg sdk.Msg, numPools int) error {
	var gas uint64
	if numPools == 2 {
		gas = 300000
	} else if numPools == 3 {
		gas = 300000
	} else if numPools == 4 {
		gas = 350000
	} else if numPools == 5 {
		gas = 400000
	} else {
		gas = 0
	}
	txm.Factory = txm.Factory.
		WithGas(gas).
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
