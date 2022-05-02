package info

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type InfoMachine struct {
	BaseUrl string
}

func (im *InfoMachine) GetAllPools() map[string]PoolResponse {
	resp, err := http.Get(im.BaseUrl + "/pools/v2/all?low_liquidity=true")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("Error while getting liquidity of pools!")
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	pools := make(map[string]PoolResponse)
	if err := json.Unmarshal(body, &pools); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
	}

	return pools
}

func (im *InfoMachine) BuildWhitelist(liquidityThreshold float64) (whitelist []uint64) {
	allPools := im.GetAllPools()
	for poolId, pool := range allPools {
		if pool[0].Liquidity > liquidityThreshold {
			poolIdUint64, err := strconv.ParseUint(poolId, 10, 64)
			if err != nil {
				fmt.Printf("Could not add %v to whitelist\n", pool)
				continue
			}
			whitelist = append(whitelist, poolIdUint64)
		}
	}
	return
}

func (im *InfoMachine) GetAllTokenPrices() map[string]TokenPriceResponse {
	resp, err := http.Get(im.BaseUrl + "/tokens/v2/all")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("Error while getting liquidity of pools!")
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	tokenPrices := make([]TokenPriceResponse, 0)
	if err := json.Unmarshal(body, &tokenPrices); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
	}

	// What does 1 base token costs in dollar
	var resultMap = make(map[string]TokenPriceResponse)
	for _, tokenPrice := range tokenPrices {
		tokenPrice.PriceE6 = int64(tokenPrice.Price * 1000000)
		resultMap[tokenPrice.Denom] = tokenPrice
	}

	return resultMap
}
