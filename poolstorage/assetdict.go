package poolstorage

type AssetDict struct {
	Assets []AssetDictEntry
}

type AssetDictEntry struct {
	Id         int
	TokenDenom string
}

func (ad *AssetDict) GetId(tokenDenom string) int {
	for _, assetDictEntry := range ad.Assets {
		if assetDictEntry.TokenDenom == tokenDenom {
			return assetDictEntry.Id
		}
	}

	// Id not found -> Add Id
	var newId int
	var newAssetDictEntry AssetDictEntry
	if len(ad.Assets) == 0 {
		ad.Assets = make([]AssetDictEntry, 1)
		newId = 0
		newAssetDictEntry = AssetDictEntry{
			Id:         newId,
			TokenDenom: tokenDenom,
		}
		ad.Assets[0] = newAssetDictEntry
	} else {
		newId = ad.Assets[len(ad.Assets)-1].Id + 1
		newAssetDictEntry = AssetDictEntry{
			Id:         newId,
			TokenDenom: tokenDenom,
		}
		ad.Assets = append(ad.Assets, newAssetDictEntry)
	}
	return newAssetDictEntry.Id
}
