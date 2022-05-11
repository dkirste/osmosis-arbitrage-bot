package main

import (
	"github.com/dkirste/arbbot/arbbot"
	"github.com/dkirste/arbbot/mywallet"
)

func main() {
	grpcNodes := []string{"46.38.245.118:9090",
		"osmosis.strange.love:9090",
		//"37.120.165.37:9090",
		"34.68.29.156:9090",
		"135.148.55.229:9090", //slow connection
		"107.178.208.25:9090", //slow connection
		"35.222.218.55:9090",
		"65.108.101.53:9090",
		"65.21.34.249:9090",
		//"161.35.88.220:9090",
		//"141.95.86.126:9090",
		//"34.220.67.127:9090",
		//"65.108.127.166:9090", ///sometimes down
		//"65.108.101.108:9090", //sometimes down
		//"65.21.198.130:9090", maybe atom
		//"45.34.1.114:9090", atom node
		//"54.203.14.73:9090",
		//"54.218.32.201:9090",
		//"54.185.37.230:9090",
		//"52.11.18.67:9090",
		//"54.70.120.178:9090",
		//"34.213.138.61:9090",
		//"54.186.186.120:9090",
	}
	//grpcNodes = []string{"37.120.165.37:9090"}
	rpcNodes := []string{"tcp://46.38.245.118:26657",
		//"tcp://37.120.165.37:26657",
		"https://rpc-osmosis.ecostake.com:443",
		"https://osmosis.validator.network:443", // banning mempool scan
		"http://34.68.29.156:26657",
		"http://135.148.55.229:26657",
		"http://107.178.208.25:26657",
		"http://37.120.165.37:26657",
		"http://65.108.127.166:26657", //sometimes down
		"http://65.108.101.108:26657", //sometimes down
		"http://65.21.198.130:36657",
		//"http://45.34.1.114:26657", // cosmoshub
		"http://35.222.218.55:26657", //down 05.05.22
		"http://65.108.101.53:26657",
		"http://65.21.34.249:26657",
		"http://161.35.88.220:26657",
		"http://141.95.86.126:26657",
		//"http://54.203.14.73:26657", // down
		"http://34.220.67.127:26657",
		//"http://54.218.32.201:26657",
		//"http://54.185.37.230:26657",
		//"http://52.11.18.67:26657",
		//"http://54.70.120.178:26657",
		//"http://34.213.138.61:26657",
		//"http://54.186.186.120:26657",
	}
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
