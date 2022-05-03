package main

import (
	"github.com/dkirste/arbbot/arbbot"
)

func main() {
	grpcNodes := []string{"46.38.245.118:9090", "osmosis.strange.love:9090", "37.120.165.37:9090"}
	//grpcNodes = []string{"osmosis.strange.love:9090"}
	rpcNodes := []string{"tcp://46.38.245.118:26657"}
	infoMachineBaseUrl := "https://api-osmosis.imperator.co"

	address := "osmo12xjjkmefue4l655h32nxf7ducnvn4tndwqzf85"
	privateKeyArmor := "-----BEGIN TENDERMINT PRIVATE KEY-----\nkdf: bcrypt\nsalt: C0F7FAA463DFEC95EC78990480CD23CC\ntype: secp256k1\n\nMDngiOXNmCdHIrFtEdfMkvKM9mlmjhLDIeOFhZPltTw9ogjTnrQin8bcS8BJ48G2\nRLeU9FnEnzSijxk8Zg8JWvFA0TLSsqjR7Mhio7s=\n=FKgG\n-----END TENDERMINT PRIVATE KEY-----"
	privateKeyPassphrase := "aoisdhgoiuahjgklasdjkfahsdfaosdfp"

	var liquidityThreshold float64 = 1000

	bot := arbbot.ArbBot{}
	bot.Setup(grpcNodes, rpcNodes, infoMachineBaseUrl, address, privateKeyArmor, privateKeyPassphrase, liquidityThreshold)
	bot.Run()

	return
}
