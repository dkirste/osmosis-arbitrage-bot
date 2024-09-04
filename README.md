# Osmosis Arbitrage Bot

This repository contains a Golang-based arbitrage bot designed for the Osmosis DEX. The bot continuously monitors osmosis liquidity pools, identifies arbitrage opportunities, and executes profitable trades.

## Project Structure
The project is organized into several key components, each responsible for a specific aspect of the arbitrage bot's functionality. Below is an overview of each component:

### Arbbot
The `Arbbot` contains the main logic of the arbitrage bot. It orchestrates the entire arbitrage process by evaluating all possible routes through liquidity pools, identifying profitable trades, and coordinating the execution of these trades. The `Arbbot` ensures that all operations are carried out efficiently and in a timely manner.

### PoolStorage
The `PoolStorage` is responsible for storing information about the liquidity pools on Osmosis. It maintains an up-to-date record of all liquidity pools, their tokens, and other relevant metrics (e.g., fees, pool ids).

### SwapRoutes
The `SwapRoutes` component determines the optimal swap paths for executing arbitrage trades. Given the complex nature of liquidity pools and token pairs, SwapRoutes calculates the most efficient route to maximize profit in exploiting arbitrage opportunities.

### GrpcMachine
The `GrpcMachine` handles communication between the arbitrage bot and the Osmosis blockchain using gRPC (gRPC Remote Procedure Calls). It manages the interaction with the blockchain, including getting latest liquidity pool states, sending arbitrage transactions, and querying data.

### InfoMachine
The `InfoMachine` is responsible for gathering and processing information of tokens traded on Osmosis using the public endpoint of `https://api-osmosis.imperator.co`. It continuously fetches the latest information about pools, tokens, and their respective prices. This information is used to identify liquidity pools with sufficient token reserves.

### TxMachine
The `TxMachine` handles the creation, signing, and broadcasting of transactions on the Osmosis blockchain.


## Getting Started

### Prerequisites

- Golang 1.18 or higher
- Osmosis SDK v0.45.1
- Osmosis node setup or access to an Osmosis gRPC endpoint

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/dkirste/osmosis-arbitrage-bot.git
    cd osmosis-arbitrage-bot
    ```

2. Install the necessary Go dependencies:

    ```bash
    go mod tidy
    ```
   
### Configuration
1. Create the `mywallet/wallet.go` configuration file and add your private key information in order to execute arbitrage trades.
```
package mywallet

func GetPrivateKey(stage string) (address string, privateKeyArmor string, privateKeyPassphrase string) {
	if stage == "prod" {
		address = "osmo..."
		privateKeyArmor = "-----BEGIN TENDERMINT PRIVATE KEY-----
		....
		-----END TENDERMINT PRIVATE KEY-----"
		privateKeyPassphrase = "..."
		return
	}
}
```

2. Add your grpcNodes and rpcNodes to the `cmd/main.go` that they can be used to access the Osmosis blockchain for retrieving pool data and sending transactions. You can add multiple grpcNodes and rpcNodes, they will all be used in parallel to decrease latency in retrieving new blocks and broadcasting transactions.


### Running the Bot

To run the bot, simply execute:

```bash
go run main.go
```

## License

This project is licensed under the Apache License 2.0 License. See the [LICENSE](https://www.apache.org/licenses/LICENSE-2.0) file for more details.

## Disclaimer

This bot is provided as-is, without any warranties or guarantees. Trading cryptocurrencies involves significant risk, including the possibility of loss. Users should exercise caution and ensure compliance with all relevant laws and regulations in their jurisdiction. The authors are not responsible for any financial losses or legal consequences that may arise from the use of this software.
