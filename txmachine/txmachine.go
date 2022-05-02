package txmachine

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/pflag"
	"sync"
)

type TxMachine struct {
	Address             sdk.AccAddress
	Factory             tx.Factory
	SequenceNumber      uint64
	SequenceNumberMutex sync.Mutex
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
