package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValidateMsgHandler func(ctx sdk.Context, msgs []sdk.Msg) error
type IsSystemFreeHandler func(ctx sdk.Context, msgs []sdk.Msg) bool
