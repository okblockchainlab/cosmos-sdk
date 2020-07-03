// nolint
package auth

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	patypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type testInput struct {
	cdc      *codec.Codec
	appCodec codec.Marshaler
	ctx      sdk.Context
	ak       keeper.AccountKeeper
	sk       types.BankKeeper
}

// moduleAccount defines an account for modules that holds coins on a pool
type moduleAccount struct {
	*types.BaseAccount
	name        string   `json:"name" yaml:"name"`              // name of the module
	permissions []string `json:"permissions" yaml"permissions"` // permissions of module account
}

// HasPermission returns whether or not the module account has permission.
func (ma moduleAccount) HasPermission(permission string) bool {
	for _, perm := range ma.permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// GetName returns the the name of the holder's module
func (ma moduleAccount) GetName() string {
	return ma.name
}

// GetPermissions returns permissions granted to the module account
func (ma moduleAccount) GetPermissions() []string {
	return ma.permissions
}

func setupTestInput() testInput {
	maccPerms := map[string][]string{
		types.FeeCollectorName: nil,
	}
	db := dbm.NewMemDB()

	cdc := codec.New()
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	appCodec := codec.NewHybridCodec(cdc, interfaceRegistry)
	types.RegisterCodec(cdc)
	//cdc.RegisterInterface((*types.ModuleAccountI)(nil), nil)
	//cdc.RegisterConcrete(&moduleAccount{}, "cosmos-sdk/ModuleAccount", nil)
	cryptocodec.RegisterCrypto(cdc)

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

	ps := patypes.NewSubspace(appCodec, keyParams, tkeyParams, types.ModuleName)
	bps := patypes.NewSubspace(appCodec, keyParams, tkeyParams, banktypes.ModuleName)
	ak := keeper.NewAccountKeeper(appCodec, authCapKey, ps, types.ProtoBaseAccount, maccPerms)
	bk := bankkeeper.NewBaseKeeper(appCodec, bankCapKey, ak, bps, nil)
	sk := NewDummySupplyKeeper(ak, bk)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	ak.SetParams(ctx, types.DefaultParams())

	return testInput{cdc: cdc, appCodec: appCodec, ctx: ctx, ak: ak, sk: sk}
}

// DummySupplyKeeper defines a supply keeper used only for testing to avoid
// circle dependencies
type DummySupplyKeeper struct {
	ak keeper.AccountKeeper
	bk bankkeeper.BaseKeeper
}

// NewDummySupplyKeeper creates a DummySupplyKeeper instance
func NewDummySupplyKeeper(ak keeper.AccountKeeper, bk bankkeeper.BaseKeeper) DummySupplyKeeper {
	return DummySupplyKeeper{ak, bk}
}

// SendCoinsFromAccountToModule for the dummy supply keeper
func (sk DummySupplyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, fromAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {

	fromAcc := sk.ak.GetAccount(ctx, fromAddr)
	moduleAcc := sk.GetModuleAccount(ctx, recipientModule)

	newFromCoins, hasNeg := sk.bk.GetAllBalances(ctx, fromAddr).SafeSub(amt)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, sk.bk.GetAllBalances(ctx, fromAddr).String())
	}

	newToCoins := sk.bk.GetAllBalances(ctx, moduleAcc.GetAddress()).Add(amt...)

	if err := sk.bk.SetBalances(ctx, fromAddr, newFromCoins); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInternal, err.Error())
	}

	if err := sk.bk.SetBalances(ctx, moduleAcc.GetAddress(), newToCoins); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInternal, err.Error())
	}

	sk.ak.SetAccount(ctx, fromAcc)
	sk.ak.SetAccount(ctx, moduleAcc)

	return nil
}

// GetModuleAccount for dummy supply keeper
func (sk DummySupplyKeeper) GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI {
	addr := sk.GetModuleAddress(moduleName)

	acc := sk.ak.GetAccount(ctx, addr)
	if acc != nil {
		macc, ok := acc.(types.ModuleAccountI)
		if ok {
			return macc
		}
	}

	moduleAddress := sk.GetModuleAddress(moduleName)
	baseAcc := types.NewBaseAccountWithAddress(moduleAddress)

	// create a new module account
	macc := &moduleAccount{
		BaseAccount: baseAcc,
		name:        moduleName,
		permissions: []string{"basic"},
	}

	maccI := (sk.ak.NewAccount(ctx, macc)).(types.ModuleAccountI)
	sk.ak.SetAccount(ctx, maccI)
	return maccI
}

// GetModuleAddress for dummy supply keeper
func (sk DummySupplyKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(moduleName)))
}
