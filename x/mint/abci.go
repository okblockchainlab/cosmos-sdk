package mint

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
	"strconv"
)

func disableMining(minter *types.Minter) {
	minter.Inflation = sdk.ZeroDec()
}

var setInflationHandler func(minter *types.Minter)

// BeginBlocker mints new tokens for the previous block.
func beginBlocker(ctx sdk.Context, k keeper.Keeper) {
	logger := ctx.Logger().With("module", "mint")
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.MetricKeyBeginBlocker)

	// fetch stored minter & params
	params := k.GetParams(ctx)
	minter := k.GetMinterCustom(ctx)
	if ctx.BlockHeight() == 0 || uint64(ctx.BlockHeight()) > minter.NextBlockToUpdate {
		k.UpdateMinterCustom(ctx, &minter, params)
	}

	logger.Debug(fmt.Sprintf(
		"total supply <%v>, "+
			"annual provisions <%v>, "+
			"params <%v>, "+
			"minted this block <%v>, "+
			"next block to update minted per block <%v>, ",
		sdk.NewDecCoinFromDec(params.MintDenom, k.StakingTokenSupply(ctx)),
		sdk.NewDecCoinFromDec(params.MintDenom, minter.AnnualProvisions),
		params,
		minter.MintedPerBlock,
		minter.NextBlockToUpdate))

	err := k.MintCoins(ctx, minter.MintedPerBlock)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, minter.MintedPerBlock)
	if err != nil {
		panic(err)
	}

	if mintedCnt, err := strconv.ParseFloat(minter.MintedPerBlock.AmountOf(sdk.DefaultBondDenom).String(), 32); err != nil {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCnt), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyInflation, params.InflationRate.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, minter.MintedPerBlock.String()),
		),
	)
}

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	setInflationHandler = disableMining
	beginBlocker(ctx, k)
}
