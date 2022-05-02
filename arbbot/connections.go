package arbbot

import (
	"fmt"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"
)

func openGRPCConn(targetIp, targetPort string) *grpc.ClientConn {
	grpcConn, _ := grpc.Dial(
		(targetIp + ":" + targetPort), // Or your gRPC server address.
		grpc.WithInsecure(),           // The SDK doesn't support any transport security mechanism.
	)
	fmt.Println("Successfully opened grpc")
	return grpcConn
}

func closeGRPCConn(grpcConn *grpc.ClientConn) {
	err := grpcConn.Close()
	if err != nil {
		panic("Could not close grpc connection")
	}
	fmt.Println("Successfully closed grpc")
}

func openRPCConn(targetIp string, targetPort string) *rpchttp.HTTP {
	rpcConn, err := rpchttp.New("tcp://" + targetIp + ":" + targetPort)
	if err != nil {
		fmt.Println(err)
	}
	return rpcConn
}

func closeRPCConn(rpcConn *rpchttp.HTTP) {
	err := rpcConn.Stop()
	if err != nil {
		return
	}
}
