package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinterCustom returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinterCustom(nextBlockToUpdate uint64, annualProvisions sdk.Dec, mintedPerBlock sdk.DecCoins) MinterCustom {
	return MinterCustom{
		NextBlockToUpdate: nextBlockToUpdate,
		AnnualProvisions:  annualProvisions,
		MintedPerBlock:    mintedPerBlock,
	}
}

// InitialMinterCustom returns an initial Minter object with a given inflation value.
func InitialMinterCustom(inflation sdk.Dec) MinterCustom {
	return NewMinterCustom(
		0,
		sdk.NewDec(0),
		sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.ZeroInt())},
	)
}

// DefaultInitialMinterCustom returns a default initial MinterCustom object for a new chain
// which uses an inflation rate of 1%.
func DefaultInitialMinterCustom() MinterCustom {
	return InitialMinterCustom(
		sdk.NewDecWithPrec(1, 2),
	)
}

// ValidateMinterCustom validate minter
func ValidateMinterCustom(minter MinterCustom) error {
	if len(minter.MintedPerBlock) != 1 || minter.MintedPerBlock[0].Denom != sdk.DefaultBondDenom {
		return fmt.Errorf(" MintedPerBlock must contain only %s", sdk.DefaultBondDenom)
	}
	return nil
}
