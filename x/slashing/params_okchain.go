package slashing

import sdk "github.com/cosmos/cosmos-sdk/types"

func (p *Params) ValidateKV(key string, value string) (interface{}, sdk.Error) {
	return nil, nil
}