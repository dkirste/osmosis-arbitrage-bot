package main

import (
	"github.com/dkirste/arbbot/arbbot"
	"github.com/dkirste/arbbot/mywallet"
)

func main() {
	grpcNodes := []string{"46.38.245.118:9090",
		"osmosis.strange.love:9090",
		"37.120.165.37:9090",
		"34.68.29.156:9090",
		"135.148.55.229:9090",
		"107.178.208.25:9090",
		"37.120.165.37:9090",
		//"65.108.127.166:9090", sometimes down
		//"65.108.101.108:9090", sometimes down
		//"65.21.198.130:9090", maybe atom
		//"45.34.1.114:9090", atom node
		"35.222.218.55:9090",
	}
	//grpcNodes = []string{"37.120.165.37:9090"}
	rpcNodes := []string{"tcp://46.38.245.118:26657",
		"tcp://37.120.165.37:26657",
		"https://rpc-osmosis.ecostake.com:443",
		//"https://osmosis.validator.network:443", // banning mempool scan
		"http://34.68.29.156:26657",
		"http://135.148.55.229:26657",
		"http://107.178.208.25:26657",
		"http://37.120.165.37:26657",
		"http://65.108.127.166:26657",
		"http://65.108.101.108:26657",
		"http://65.21.198.130:36657",
		//"http://45.34.1.114:26657",
		"http://35.222.218.55:26657",
	}
	//rpcNodes = []string{"tcp://46.38.245.118:26657"}
	infoMachineBaseUrl := "https://api-osmosis.imperator.co"

	address, privateKeyArmor, privateKeyPassphrase := mywallet.GetPrivateKey("prod")

	var liquidityThreshold float64 = 10

	bot := arbbot.ArbBot{}
	bot.Setup(grpcNodes, rpcNodes, infoMachineBaseUrl, address, privateKeyArmor, privateKeyPassphrase, liquidityThreshold)
	//bot.RunBlockArb(6)
	bot.RunMempoolArb(4, 1000)

	return
}
