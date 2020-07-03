package mint

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

func TestGenesis(t *testing.T) {
	setup := newTestInput(t)
	genesisState := types.NewGenesisState(
		types.DefaultInitialMinterCustom(), types.DefaultParams())
	defaultGenesisState := types.DefaultGenesisState()
	require.Equal(t, genesisState, defaultGenesisState)
	InitGenesis(setup.ctx, setup.mintKeeper, setup.ak, defaultGenesisState)
	require.NoError(t, types.ValidateGenesis(defaultGenesisState))
	exportedState := ExportGenesis(setup.ctx, setup.mintKeeper)
	require.NotNil(t, exportedState)
	require.Equal(t, defaultGenesisState, exportedState)
}
