package grpcMachine

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	balancer "github.com/osmosis-labs/osmosis/v13/x/gamm/pool-models/balancer"
	gammtypes "github.com/osmosis-labs/osmosis/v13/x/gamm/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func (gm *GrpcMachine) QueryAccountBalance(address string, denom string) (balance sdk.Coin) {
	bankClient := banktypes.NewQueryClient(gm.Conn)
	ctx := context.Background()
	bankRes, _ := bankClient.Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
	})
	balance = sdk.Coin{
		Denom:  bankRes.Balance.Denom,
		Amount: bankRes.Balance.Amount,
	}
	return

}

func (gm *GrpcMachine) QueryNumberOfPools() (uint64, error) {
	poolClient := gammtypes.NewQueryClient(gm.Conn)

	poolRes, err := poolClient.NumPools(context.Background(), &gammtypes.QueryNumPoolsRequest{})
	if err != nil {
		//fmt.Println("Error while executing request (QueryNumberOfPools)")
		return 0, err
	}

	return poolRes.NumPools, nil
}

func (gm *GrpcMachine) QueryAllPools() ([]balancer.Pool, uint64) {
	poolClient := gammtypes.NewQueryClient(gm.Conn)

	numberOfPools, err := gm.QueryNumberOfPools()
	if err != nil {
		panic("Could not get Number of pools")
	}

	var header metadata.MD

	poolRes, err := poolClient.Pools(context.Background(), &gammtypes.QueryPoolsRequest{Pagination: &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      numberOfPools,
		CountTotal: false,
		Reverse:    false,
	}}, grpc.Header(&header))
	if err != nil {
		fmt.Println("Error while executing request")
		return nil, 0
	}

	var poolI balancer.Pool
	var pools []balancer.Pool

	for _, pool := range poolRes.Pools {
		err := gm.InterfaceRegistry.UnpackAny(pool, &poolI)
		if err != nil {
			fmt.Printf("ERROR WHILE UNPACKING %v\n", err)
			return nil, 0
		}
		pools = append(pools, poolI)
	}
	currentHeight, _ := strconv.ParseUint(header.Get(grpctypes.GRPCBlockHeightHeader)[0], 10, 64)

	return pools, currentHeight
}

func (gm *GrpcMachine) QueryAccountNumber(address string) (accNum uint64) {
	accountClient := authtypes.NewQueryClient(gm.Conn)
	account, err := accountClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var accountI authtypes.AccountI
	err = gm.InterfaceRegistry.UnpackAny(account.Account, &accountI)
	accNum = accountI.GetAccountNumber()
	return
}

func (gm *GrpcMachine) QueryAccountSequence(address string) (seq uint64) {
	accountClient := authtypes.NewQueryClient(gm.Conn)
	account, err := accountClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var accountI authtypes.AccountI
	err = gm.InterfaceRegistry.UnpackAny(account.Account, &accountI)
	seq = accountI.GetSequence()
	return
}

func (gm *GrpcMachine) QueryCurrentHeight() (currentHeight uint64) {
	poolClient := gammtypes.NewQueryClient(gm.Conn)

	// Get header from grpc via grpc option
	var header metadata.MD

	poolRes, err := poolClient.NumPools(context.Background(), &gammtypes.QueryNumPoolsRequest{}, grpc.Header(&header))
	_ = poolRes
	if err != nil {
		panic(err)
	}

	currentHeight, _ = strconv.ParseUint(header.Get(grpctypes.GRPCBlockHeightHeader)[0], 10, 64)
	return
}
