package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (keeper Keeper) BankKeeper() types.BankKeeper {
	return keeper.bankKeeper
}

func (keeper Keeper) ParamSpace() types.ParamSubspace {
	return keeper.paramSpace
}

func (keeper Keeper) StoreKey() sdk.StoreKey {
	return keeper.storeKey
}

func (keeper Keeper) Cdc() codec.Marshaler {
	return keeper.cdc
}
