package main

import (
	"github.com/dkirste/arbbot/arbbot"
	"github.com/dkirste/arbbot/mywallet"
)

func main() {
	grpcNodes := []string{"x.x.x.x:9090"}
	rpcNodes := []string{"tcp://x.x.x.x:26657"}
	infoMachineBaseUrl := "https://api-osmosis.imperator.co"

	address, privateKeyArmor, privateKeyPassphrase := mywallet.GetPrivateKey("test")

	var liquidityThreshold float64 = 750

	bot := arbbot.ArbBot{}
	bot.Setup(grpcNodes, rpcNodes, infoMachineBaseUrl, address, privateKeyArmor, privateKeyPassphrase, liquidityThreshold)
	bot.RunBlockArb(6)
	return
}
