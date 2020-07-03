package slashing

import (
	"encoding/hex"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// TODO remove dependencies on staking (should only refer to validator set type from sdk)

var (
	pks = []crypto.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
	}
	addrs = []sdk.ValAddress{
		sdk.ValAddress(pks[0].Address()),
		sdk.ValAddress(pks[1].Address()),
		sdk.ValAddress(pks[2].Address()),
	}
	initTokens = sdk.TokensFromConsensusPower(200)
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
)

func createTestInput(t *testing.T, defaults types.Params) (codec.Marshaler, sdk.Context, authkeeper.AccountKeeper, bankkeeper.Keeper, stakingkeeper.Keeper,
	paramstypes.Subspace, keeper.Keeper) {
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keySlashing := sdk.NewKVStoreKey(types.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)

	db := dbm.NewMemDB()

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySlashing, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewNopLogger())
	cdc := codec.New()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewHybridCodec(cdc, interfaceRegistry)
	std.RegisterCodec(cdc)
	std.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterCodec(cdc)

	feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[notBondedPool.String()] = true
	blacklistedAddrs[bondPool.String()] = true

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: []string{authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    []string{authtypes.Burner, authtypes.Staking},
	}

	paramsKeeper := paramskeeper.NewKeeper(appCodec, keyParams, tkeyParams)
	accountKeeper := authkeeper.NewAccountKeeper(appCodec, keyAcc, paramsKeeper.Subspace(authtypes.DefaultParamspace), authtypes.ProtoBaseAccount, maccPerms)

	bk := bankkeeper.NewBaseKeeper(appCodec, keyBank, accountKeeper, paramsKeeper.Subspace(banktypes.DefaultParamspace), blacklistedAddrs)

	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(int64(len(addrs)))))
	bk.SetSupply(ctx, banktypes.NewSupply(totalSupply))

	sk := stakingkeeper.NewKeeper(appCodec, keyStaking, accountKeeper, bk, paramsKeeper.Subspace(stakingtypes.DefaultParamspace))
	genesis := stakingtypes.DefaultGenesisState()

	// set module accounts
	accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	accountKeeper.SetModuleAccount(ctx, bondPool)
	accountKeeper.SetModuleAccount(ctx, notBondedPool)

	_ = staking.InitGenesis(ctx, sk, accountKeeper, bk, genesis)

	for _, addr := range addrs {
		_, err = bk.AddCoins(ctx, sdk.AccAddress(addr), initCoins)
	}
	require.Nil(t, err)
	paramstore := paramsKeeper.Subspace(types.DefaultParamspace)
	keeper := keeper.NewKeeper(appCodec, keySlashing, &sk, paramstore)
	sk.SetHooks(keeper.Hooks())

	require.NotPanics(t, func() {
		InitGenesis(ctx, keeper, sk, types.GenesisState{defaults, nil, nil})
	})

	return appCodec, ctx, accountKeeper, bk, sk, paramstore, keeper
}

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey crypto.PubKey, amt sdk.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	return stakingtypes.NewMsgCreateValidator(
		address, pubKey, sdk.NewCoin(sdk.DefaultBondDenom, amt),
		stakingtypes.Description{}, commission, sdk.OneInt(),
	)
}

func newTestMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, delAmount sdk.Int) *stakingtypes.MsgDelegate {
	amount := sdk.NewCoin(sdk.DefaultBondDenom, delAmount)
	return stakingtypes.NewMsgDelegate(delAddr, valAddr, amount)
}

func TestInitGenesis(t *testing.T) {
	cdc, ctx, ak, bk, stakingKeeper, _, slashingKeeper := createTestInput(t, types.DefaultParams())
	appModule := NewAppModule(cdc, slashingKeeper, ak, bk, stakingKeeper)
	var exportState types.GenesisState
	// 1.check default export
	cdc.MustUnmarshalJSON(appModule.ExportGenesis(ctx, cdc), &exportState)
	require.Equal(t, types.DefaultGenesisState(), exportState)
	// 2.change params and check again
	initParams := types.NewParams(1000, sdk.MustNewDecFromStr("0.05"), 600000000000, sdk.ZeroDec(), sdk.ZeroDec())
	genesisState := types.NewGenesisState(initParams, exportState.SigningInfos, exportState.MissedBlocks)
	appModule.InitGenesis(ctx, cdc, cdc.MustMarshalJSON(genesisState))
	cdc.MustUnmarshalJSON(appModule.ExportGenesis(ctx, cdc), &exportState)
	require.Equal(t, genesisState, exportState)
	// 3.change the state.SigningInfos and state.MissedBlocks info and check again
	conAddress := sdk.GetConsAddress(pks[0])
	slashingKeeper.AddPubkey(ctx, pks[0])
	sigingInfo := types.NewValidatorSigningInfo(conAddress, 10, 1, time.Now().Add(10000), true, 5, types.Destroying)
	slashingKeeper.SetValidatorSigningInfo(ctx, conAddress, sigingInfo)
	slashingKeeper.HandleValidatorSignature(ctx, pks[0].Address(), 100, false)
	cdc.MustUnmarshalJSON(appModule.ExportGenesis(ctx, cdc), &exportState)
	require.Equal(t, 1, len(exportState.SigningInfos))
	require.Equal(t, 1, len(exportState.MissedBlocks))
}
