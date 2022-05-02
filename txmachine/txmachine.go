package txmachine

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/pflag"
)

type TxMachine struct {
	Factory tx.Factory
}

func (txm *TxMachine) Setup(clientCtx client.Context) {
	flagSet := pflag.FlagSet{
		Usage:                nil,
		SortFlags:            false,
		ParseErrorsWhitelist: pflag.ParseErrorsWhitelist{},
	}

	txf := tx.NewFactoryCLI(clientCtx, &flagSet)
	txm.Factory = txf

	return
}
