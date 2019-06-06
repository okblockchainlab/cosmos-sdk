package auth

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubCollectedFees - sub fee from fee pool
func (fck FeeCollectionKeeper) SubCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	logger := ctx.Logger().With("module", "auth")
	oldCoins := fck.GetCollectedFees(ctx)
	newCoins, anyNeg := oldCoins.SafeSub(coins)
	if !anyNeg {
		fck.setCollectedFees(ctx, newCoins)
		logger.Debug(fmt.Sprintf("sub fee from pool, oldCoins: %v, subCoins: %v, newCoins: %v",
			oldCoins, coins, newCoins))
	} else {
		logger.Error(fmt.Sprintf("sub fee from pool failed, oldCoins: %v, subCoins: %v",
			oldCoins, coins))
	}

	return newCoins
}