package mint

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

type MockApp struct {
	*mock.App
	tkeyStaking *sdk.TransientStoreKey
	keyMint     *sdk.KVStoreKey
	mintKeeper  keeper.Keeper
}

func registerCodec(cdc *codec.Codec) {
	banktypes.RegisterCodec(cdc)
}
func getMockApp(t *testing.T, numGenAccs int, balance int64, mintParams types.Params) (mockApp *MockApp, addrKeysSlice mock.AddrKeysSlice) {
	mapp := mock.NewApp()
	registerCodec(mapp.Cdc)
	mockApp = &MockApp{
		App:         mapp,
		tkeyStaking: sdk.NewTransientStoreKey(stakingtypes.StoreKey),
		keyMint:     sdk.NewKVStoreKey(types.StoreKey),
	}

	mockApp.mintKeeper = keeper.NewKeeper(mockApp.AppCodec, mockApp.keyMint,
		mockApp.ParamsKeeper.Subspace(types.DefaultParamspace), &mockApp.StakingKeeper,
		mockApp.AccountKeeper, mockApp.BankKeeper, authtypes.FeeCollectorName)
	//mockApp.Router().AddRoute("", nil)
	mockApp.QueryRouter().AddRoute(types.QuerierRoute, keeper.NewQuerier(mockApp.mintKeeper))
	decCoins, _ := sdk.ParseDecCoins(fmt.Sprintf("%d%s",
		balance, sdk.DefaultBondDenom))
	coins := decCoins
	keysSlice, genAccs, genBals := CreateGenAccounts(numGenAccs, coins)
	addrKeysSlice = keysSlice
	mockApp.SetBeginBlocker(getBeginBlocker(mockApp.mintKeeper))
	mockApp.SetInitChainer(getInitChainer(mockApp.App, mockApp.BankKeeper, mockApp.mintKeeper, mockApp.StakingKeeper,
		genAccs, mintParams, coins))
	// todo: checkTx in mock app
	mockApp.SetAnteHandler(nil)
	require.NoError(t, mockApp.CompleteSetup(
		mockApp.KeyStaking,
		mockApp.keyMint,
	))
	mock.SetGenesis(mockApp.App, genAccs, genBals)
	for i := 0; i < numGenAccs; i++ {
		mock.CheckBalance(t, mockApp.App, keysSlice[i].Address, coins)
	}
	return
}
func getBeginBlocker(keeper keeper.Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		BeginBlocker(ctx, keeper)
		return abci.ResponseBeginBlock{}
	}
}
func getInitChainer(mapp *mock.App, supplyKeeper bankkeeper.BaseKeeper, mintKeeper keeper.Keeper, stakingkeeper stakingkeeper.Keeper,
	genAccs []authtypes.BaseAccount, mintParams types.Params) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		for _, acc := range genAccs {
			mapp.TotalCoinsSupply = mapp.TotalCoinsSupply.Add(mapp.BankKeeper.GetAllBalances(ctx, acc.GetAddress())...)
		}
		supplyKeeper.SetSupply(ctx, banktypes.NewSupply(mapp.TotalCoinsSupply))
		mintKeeper.SetParams(ctx, mintParams)
		mintKeeper.SetMinterCustom(ctx, types.MinterCustom{})
		stakingkeeper.SetParams(ctx, stakingtypes.DefaultParams())
		return abci.ResponseInitChain{}
	}
}
func CreateGenAccounts(numAccs int, coins sdk.Coins) (addrKeysSlice mock.AddrKeysSlice,
	genAccs []authtypes.BaseAccount, genBals []banktypes.Balance) {
	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())
		addrKeys := mock.NewAddrKeys(addr, pubKey, privKey)
		account := authtypes.BaseAccount{
			Address: addr,
		}

		genAccs = append(genAccs, account)

		bal := banktypes.Balance{
			Address: addr,
			Coins: coins,
		}
		genBals = append(genBals, bal)
		addrKeysSlice = append(addrKeysSlice, addrKeys)
	}
	return
}
