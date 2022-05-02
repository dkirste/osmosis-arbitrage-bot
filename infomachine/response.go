package info

type PoolResponse []TokenResponse

type TokenResponse struct {
	Symbol             string  `json:"symbol"`
	Amount             float64 `json:"amount"`
	Denom              string  `json:"denom"`
	CoingeckoId        string  `json:"coingecko_id"`
	Liquidity          float64 `json:"liquidity"`
	Liquidity24hChange float64 `json:"liquidity_24h_change"`
	Volume24h          float64 `json:"volume_24h"`
	//Volume24hChange float64 `json:"volume_24h_change"`  // Don't parse because of "" instead of 0
	Volume7d float64 `json:"volume_7d"`
	Price    float64 `json:"price"`
	Fees     string  `json:"fees"`
}

type TokenPriceResponse struct {
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Denom    string  `json:"denom"`
	Exponent int     `json:"exponent"`
	PriceE6  int64
}
