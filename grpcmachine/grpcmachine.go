package grpcMachine

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"google.golang.org/grpc"
)

type GrpcMachine struct {
	Conn              *grpc.ClientConn
	InterfaceRegistry codectypes.InterfaceRegistry
}
