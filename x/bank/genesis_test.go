package bank

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	sendEnabledExpected = false
)

var (
	defaultGenExportExpected = fmt.Sprintf(`{"send_enabled":true,"balances":[],"supply":[{"denom":"%s","amount":"100.00000000"}]}`, sdk.DefaultBondDenom)
	genExportExpected        = fmt.Sprintf(`{"send_enabled":%v,"balances":[],"supply":[{"denom":"%s","amount":"100.00000000"}]}`, sendEnabledExpected, sdk.DefaultBondDenom)
)

func TestInitGenesis(t *testing.T) {
	input := setupTestInput()
	ctx, accKeeper, bankKeeper := input.ctx, input.ak, input.bk
	appModule := NewAppModule(input.appCodec, bankKeeper, accKeeper)
	// 1.check default export
	require.Equal(t, defaultGenExportExpected, string(appModule.ExportGenesis(ctx, input.appCodec)))
	// 2.change context
	bankKeeper.SetSendEnabled(ctx, sendEnabledExpected)
	// 3.export again
	genExport := appModule.ExportGenesis(ctx, input.appCodec)
	require.Equal(t, genExportExpected, string(genExport))
	// 4.init again && check
	newInput := setupTestInput()
	newCtx, newAccKeeper, newBankKeeper := newInput.ctx, newInput.ak, newInput.bk
	newAppModule := NewAppModule(input.appCodec, newBankKeeper, newAccKeeper)
	newAppModule.InitGenesis(newCtx, input.appCodec, genExport)
	require.Equal(t, sendEnabledExpected, newBankKeeper.GetSendEnabled(newCtx))
}
