package main

import (
	"github.com/dkirste/arbbot/arbbot"
	"github.com/dkirste/arbbot/mywallet"
)

func main() {
	grpcNodes := []string{"192.168.2.201:9090"}
	//grpcNodes = []string{"37.120.165.37:9090"}
	rpcNodes := []string{"tcp://192.168.2.201:26657"}
	//rpcNodes = []string{"tcp://46.38.245.118:26657"}
	infoMachineBaseUrl := "https://api-osmosis.imperator.co"

	address, privateKeyArmor, privateKeyPassphrase := mywallet.GetPrivateKey("test")

	var liquidityThreshold float64 = 750

	bot := arbbot.ArbBot{}
	bot.Setup(grpcNodes, rpcNodes, infoMachineBaseUrl, address, privateKeyArmor, privateKeyPassphrase, liquidityThreshold)
	bot.RunBlockArb(6)
	//bot.RunMempoolArb(4, 1000)

	return
}
