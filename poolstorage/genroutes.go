package poolstorage

import "fmt"

func (ps *PoolStorage) GenerateThreeCurrencyRoutes() {
	routes := make([]SwapAmountInRoutesId, 0)

	firstAssetId := ps.AssetDict.GetId("uosmo")
	for _, firstPool := range ps.PoolsByAsset[firstAssetId] {
		for _, secondAsset := range firstPool.GetAllPoolAssets() {
			secondAssetId := ps.AssetDict.GetId(secondAsset.Token.Denom)
			// Loop protection
			if secondAssetId == firstAssetId {
				continue
			}

			for _, secondPool := range ps.PoolsByAsset[secondAssetId] {
				for _, thirdAsset := range secondPool.GetAllPoolAssets() {
					thirdAssetId := ps.AssetDict.GetId(thirdAsset.Token.Denom)
					// Loop protection
					if thirdAssetId == secondAssetId {
						continue
					}
					if thirdAssetId == firstAssetId {
						continue
					}

					for _, thirdPool := range ps.PoolsByAsset[thirdAssetId] {
						for _, fourthAsset := range thirdPool.GetAllPoolAssets() {
							fourthAssetId := ps.AssetDict.GetId(fourthAsset.Token.Denom)
							// Loop protection
							if fourthAssetId == thirdAssetId {
								continue
							}

							if fourthAssetId == firstAssetId {
								arbitrageRoutes := make(SwapAmountInRoutesId, 3)
								arbitrageRoutes[0] = SwapAmountInRouteId{
									PoolId:        firstPool.GetId(),
									TokenOutDenom: secondAsset.Token.Denom,
								}
								arbitrageRoutes[1] = SwapAmountInRouteId{
									PoolId:        secondPool.GetId(),
									TokenOutDenom: thirdAsset.Token.Denom,
								}
								arbitrageRoutes[2] = SwapAmountInRouteId{
									PoolId:        thirdPool.GetId(),
									TokenOutDenom: fourthAsset.Token.Denom, // = firstAsset
								}
								if arbitrageRoutes.CheckIfPoolsAreUnique() {
									routes = append(routes, arbitrageRoutes)
								}
							}
						}
					}
				}
			}
		}
	}
	fmt.Println("ThreeCurrency Number: ", len(routes))
	ps.ThreeCurrencyRoutes = routes
	return
}

func (ps *PoolStorage) GenerateFourCurrencyRoutes() {
	routes := make([]SwapAmountInRoutesId, 0)

	firstAssetId := ps.AssetDict.GetId("uosmo")
	for _, firstPool := range ps.PoolsByAsset[firstAssetId] {
		for _, secondAsset := range firstPool.GetAllPoolAssets() {
			secondAssetId := ps.AssetDict.GetId(secondAsset.Token.Denom)
			// Loop protection
			if secondAssetId == firstAssetId {
				continue
			}
			if secondAssetId == firstAssetId {
				continue
			}

			for _, secondPool := range ps.PoolsByAsset[secondAssetId] {
				for _, thirdAsset := range secondPool.GetAllPoolAssets() {
					thirdAssetId := ps.AssetDict.GetId(thirdAsset.Token.Denom)
					// Loop protection
					if thirdAssetId == secondAssetId {
						continue
					}
					if thirdAssetId == firstAssetId {
						continue
					}

					for _, thirdPool := range ps.PoolsByAsset[thirdAssetId] {
						for _, fourthAsset := range thirdPool.GetAllPoolAssets() {
							fourthAssetId := ps.AssetDict.GetId(fourthAsset.Token.Denom)
							// Loop protection
							if fourthAssetId == thirdAssetId {
								continue
							}

							for _, fourthPool := range ps.PoolsByAsset[fourthAssetId] {
								for _, fifthAsset := range fourthPool.GetAllPoolAssets() {
									fifthAssetId := ps.AssetDict.GetId(fifthAsset.Token.Denom)
									// Loop protection
									if fifthAssetId == fourthAssetId {
										continue
									}

									if fifthAssetId == firstAssetId {
										arbitrageRoutes := make(SwapAmountInRoutesId, 4)
										arbitrageRoutes[0] = SwapAmountInRouteId{
											PoolId:        firstPool.GetId(),
											TokenOutDenom: secondAsset.Token.Denom,
										}
										arbitrageRoutes[1] = SwapAmountInRouteId{
											PoolId:        secondPool.GetId(),
											TokenOutDenom: thirdAsset.Token.Denom,
										}
										arbitrageRoutes[2] = SwapAmountInRouteId{
											PoolId:        thirdPool.GetId(),
											TokenOutDenom: fourthAsset.Token.Denom,
										}
										arbitrageRoutes[3] = SwapAmountInRouteId{
											PoolId:        fourthPool.GetId(),
											TokenOutDenom: fifthAsset.Token.Denom, // = firstAsset
										}
										if arbitrageRoutes.CheckIfPoolsAreUnique() {
											routes = append(routes, arbitrageRoutes)
										}

									}
								}

							}
						}
					}
				}
			}
		}
	}
	fmt.Println("FourCurrency Number: ", len(routes))
	ps.FourCurrencyRoutes = routes
	return
}

func (ps *PoolStorage) GenerateFiveCurrencyRoutes() {
	routes := make([]SwapAmountInRoutesId, 0)

	firstAssetId := ps.AssetDict.GetId("uosmo")
	for _, firstPool := range ps.PoolsByAsset[firstAssetId] {
		for _, secondAsset := range firstPool.GetAllPoolAssets() {
			secondAssetId := ps.AssetDict.GetId(secondAsset.Token.Denom)
			// Loop protection
			if secondAssetId == firstAssetId {
				continue
			}
			if secondAssetId == firstAssetId {
				continue
			}

			for _, secondPool := range ps.PoolsByAsset[secondAssetId] {
				for _, thirdAsset := range secondPool.GetAllPoolAssets() {
					thirdAssetId := ps.AssetDict.GetId(thirdAsset.Token.Denom)
					// Loop protection
					if thirdAssetId == secondAssetId {
						continue
					}
					if thirdAssetId == firstAssetId {
						continue
					}

					for _, thirdPool := range ps.PoolsByAsset[thirdAssetId] {
						for _, fourthAsset := range thirdPool.GetAllPoolAssets() {
							fourthAssetId := ps.AssetDict.GetId(fourthAsset.Token.Denom)
							// Loop protection
							if fourthAssetId == thirdAssetId {
								continue
							}

							for _, fourthPool := range ps.PoolsByAsset[fourthAssetId] {
								for _, fifthAsset := range fourthPool.GetAllPoolAssets() {
									fifthAssetId := ps.AssetDict.GetId(fifthAsset.Token.Denom)
									// Loop protection
									if fifthAssetId == fourthAssetId {
										continue
									}

									for _, fifthPool := range ps.PoolsByAsset[fifthAssetId] {
										for _, sixthAsset := range fifthPool.GetAllPoolAssets() {
											sixthAssetId := ps.AssetDict.GetId(sixthAsset.Token.Denom)
											// Loop protection
											if sixthAssetId == fifthAssetId {
												continue
											}

											if sixthAssetId == firstAssetId {
												arbitrageRoutes := make(SwapAmountInRoutesId, 5)
												arbitrageRoutes[0] = SwapAmountInRouteId{
													PoolId:        firstPool.GetId(),
													TokenOutDenom: secondAsset.Token.Denom,
												}
												arbitrageRoutes[1] = SwapAmountInRouteId{
													PoolId:        secondPool.GetId(),
													TokenOutDenom: thirdAsset.Token.Denom,
												}
												arbitrageRoutes[2] = SwapAmountInRouteId{
													PoolId:        thirdPool.GetId(),
													TokenOutDenom: fourthAsset.Token.Denom,
												}
												arbitrageRoutes[3] = SwapAmountInRouteId{
													PoolId:        fourthPool.GetId(),
													TokenOutDenom: fifthAsset.Token.Denom, // = firstAsset
												}
												arbitrageRoutes[4] = SwapAmountInRouteId{
													PoolId:        fifthPool.GetId(),
													TokenOutDenom: sixthAsset.Token.Denom, // = firstAsset
												}
												if arbitrageRoutes.CheckIfPoolsAreUnique() {
													routes = append(routes, arbitrageRoutes)
												}

											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	fmt.Println("FiveCurrency Number: ", len(routes))
	ps.FiveCurrencyRoutes = routes
	return
}

func (ps *PoolStorage) AddGeneratedThreeCurrencyRoutesById(routes []SwapAmountInRoutesId) {
	var poolId uint64
	for _, swapRoutes := range routes {
		for _, route := range swapRoutes {
			poolId = route.PoolId

			// Check if array is empty
			if len(ps.ThreeCurrencyRoutesById[poolId]) == 0 {
				ps.ThreeCurrencyRoutesById[poolId] = make([]SwapAmountInRoutesId, 0)
			}

			// Append route
			ps.ThreeCurrencyRoutesById[poolId] = append(ps.ThreeCurrencyRoutesById[poolId], swapRoutes)

		}
	}
	/* PRINT LOOP
	for id, swapRoutes := range ps.ThreeCurrencyRoutesById {
		fmt.Printf("%v: %v\n", id, len(swapRoutes))
	}
	*/

}
