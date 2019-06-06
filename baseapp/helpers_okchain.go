package baseapp

import (
	"bytes"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getFeeFromTags(res sdk.Result) (i int, fee sdk.Coins) {
	for i, tag := range res.Tags.ToKVPairs() {
		if bytes.EqualFold(tag.Key, []byte(sdk.Fee_TagName)) {
			//fmt.Printf("%s: %s\n", string(tag.Key), string(tag.Value))
			//res.Tags = append(res.Tags[0:i], res.Tags[i+1:]...)
			return i, strToCoins(string(tag.Value))
		}
	}
	return i, sdk.Coins{}
}

func strToCoins(amount string) sdk.Coins {
	var res sdk.Coins
	coinStrs := strings.Split(amount, ",")
	for _, coinStr := range coinStrs {
		coin := strings.Split(coinStr, ":")
		if len(coin) == 2 {
			var c sdk.Coin
			c.Denom = coin[1]
			coinDec := sdk.MustNewDecFromStr(coin[0])
			c.Amount = sdk.NewIntFromBigInt(coinDec.Int)
			res = append(res, c)
		}
	}
	return res
}

func coins2str(coins sdk.Coins)string{
	if len(coins) == 0 {
		return ""
	}

	out := ""
	for _, coin := range coins {
		out += fmt.Sprintf("%v,", coin2str(coin))
	}
	return out[:len(out)-1]
}

// String provides a human-readable representation of a coin
func coin2str(coin sdk.Coin) string {
	dec := sdk.NewDecFromIntWithPrec(coin.Amount, sdk.Precision)
	return fmt.Sprintf("%s %v", dec, coin.Denom)
}
