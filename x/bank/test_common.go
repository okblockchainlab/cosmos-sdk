// nolint
package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	autypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	patypes "github.com/cosmos/cosmos-sdk/x/params/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type testInput struct {
	cdc      *codec.Codec
	appCodec codec.Marshaler
	ctx      sdk.Context
	ak       authkeeper.AccountKeeper
	bk       keeper.BaseKeeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	interfaceRegistry := types.NewInterfaceRegistry()
	appCodec := codec.NewHybridCodec(cdc, interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	bankCapKey := sdk.NewKVStoreKey("bankCapKey")
	keyParams := sdk.NewKVStoreKey("subspace")
	tkeyParams := sdk.NewTransientStoreKey("transient_subspace")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(bankCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	ps := patypes.NewSubspace(appCodec, keyParams, tkeyParams, autypes.ModuleName)
	bps := patypes.NewSubspace(appCodec, keyParams, tkeyParams, banktypes.ModuleName)
	ak := authkeeper.NewAccountKeeper(appCodec, authCapKey, ps, autypes.ProtoBaseAccount, nil)
	bk := keeper.NewBaseKeeper(appCodec, bankCapKey, ak, bps, nil)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	bk.SetSendEnabled(ctx, true)
	supply := banktypes.NewSupply(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))))
	bk.SetSupply(ctx, supply)

	return testInput{appCodec: appCodec, ctx: ctx, ak: ak, bk: bk}
}
