package mock

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc-transfer/types"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	ibckeeper "github.com/cosmos/cosmos-sdk/x/ibc/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"math/rand"
	"os"
	"sort"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeepr "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
)

const chainID = ""

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*bam.BaseApp
	Cdc        *codec.Codec // Cdc is public since the codec is passed into the module anyways
	AppCodec   codec.Marshaler
	KeyMain    *sdk.KVStoreKey
	KeyAccount *sdk.KVStoreKey
	KeyStaking *sdk.KVStoreKey
	KeyParams  *sdk.KVStoreKey
	TKeyParams *sdk.TransientStoreKey

	// TODO: Abstract this out from not needing to be auth specifically
	AccountKeeper    authkeepr.AccountKeeper
	BankKeeper       bankkeeper.BaseKeeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper

	GenesisAccounts  []authtypes.BaseAccount
	GenesisBalances  []banktypes.Balance
	TotalCoinsSupply sdk.Coins
}

// NewApp partially constructs a new app on the memstore for module and genesis
// testing.
func NewApp() *App {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	}

	// Create the cdc with some standard codecs
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	banktypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	appCodec, cdc := encodingConfig.Marshaler, encodingConfig.Amino
	authtypes.RegisterCodec(cdc)

	// Create your application object
	app := &App{
		BaseApp:          bam.NewBaseApp("mock", logger, db, authtypes.DefaultTxDecoder(cdc)),
		Cdc:              cdc,
		AppCodec:         appCodec,
		KeyMain:          sdk.NewKVStoreKey(banktypes.StoreKey),
		KeyAccount:       sdk.NewKVStoreKey(authtypes.StoreKey),
		KeyStaking:       sdk.NewKVStoreKey(stakingtypes.StoreKey),
		KeyParams:        sdk.NewKVStoreKey("params"),
		TKeyParams:       sdk.NewTransientStoreKey("transient_params"),
		TotalCoinsSupply: sdk.NewCoins(),
	}

	blockAddrs := map[string]bool{
		authtypes.FeeCollectorName:     true,
		distrtypes.ModuleName:          true,
		minttypes.ModuleName:           true,
		stakingtypes.BondedPoolName:    true,
		stakingtypes.NotBondedPoolName: true,
		govtypes.ModuleName:            true,
		ibctransfertypes.ModuleName:    true,
	}

	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// define keepers
	app.ParamsKeeper = paramskeeper.NewKeeper(app.AppCodec, app.KeyParams, app.TKeyParams)

	app.AccountKeeper = authkeepr.NewAccountKeeper(
		app.AppCodec,
		app.KeyAccount,
		app.ParamsKeeper.Subspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(app.AppCodec, app.KeyMain, app.AccountKeeper,
		app.ParamsKeeper.Subspace(banktypes.StoreKey), blockAddrs)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, app.KeyStaking, app.AccountKeeper, app.BankKeeper, app.ParamsKeeper.Subspace(stakingtypes.ModuleName),
	)

	app.StakingKeeper = stakingKeeper

	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, sdk.NewKVStoreKey(capabilitytypes.StoreKey), memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)

	app.IBCKeeper = ibckeeper.NewKeeper(
		app.Cdc, appCodec, sdk.NewKVStoreKey(ibchost.StoreKey), app.StakingKeeper, scopedIBCKeeper,
	)

	// Initialize the app. The chainers and blockers can be overwritten before
	// calling complete setup.
	app.SetInitChainer(app.InitChainer)
	app.SetAnteHandler(ante.NewAnteHandler(app.AccountKeeper, app.BankKeeper, *app.IBCKeeper, ante.DefaultSigVerificationGasConsumer, authtypes.LegacyAminoJSONHandler{}))

	// Not sealing for custom extension

	return app
}

// CompleteSetup completes the application setup after the routes have been
// registered.
func (app *App) CompleteSetup(newKeys ...sdk.StoreKey) error {
	newKeys = append(
		newKeys,
		app.KeyMain, app.KeyAccount, app.KeyParams, app.TKeyParams,
	)

	for _, key := range newKeys {
		switch key.(type) {
		case *sdk.KVStoreKey:
			app.MountStore(key, sdk.StoreTypeIAVL)
		case *sdk.TransientStoreKey:
			app.MountStore(key, sdk.StoreTypeTransient)
		default:
			return fmt.Errorf("unsupported StoreKey: %+v", key)
		}
	}

	err := app.LoadLatestVersion()

	return err
}

// InitChainer performs custom logic for initialization.
// nolint: errcheck
func (app *App) InitChainer(ctx sdk.Context, _ abci.RequestInitChain) abci.ResponseInitChain {

	// Load the genesis accounts
	for _, genBal := range app.GenesisBalances {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, genBal.GetAddress())
		app.BankKeeper.SetBalances(ctx, genBal.GetAddress(), genBal.Coins)
		app.AccountKeeper.SetAccount(ctx, acc)
	}

	auth.InitGenesis(ctx, app.AccountKeeper, authtypes.DefaultGenesisState())

	return abci.ResponseInitChain{}
}

// Type that combines an Address with the privKey and pubKey to that address
type AddrKeys struct {
	Address sdk.AccAddress
	PubKey  crypto.PubKey
	PrivKey crypto.PrivKey
}

func NewAddrKeys(address sdk.AccAddress, pubKey crypto.PubKey,
	privKey crypto.PrivKey) AddrKeys {

	return AddrKeys{
		Address: address,
		PubKey:  pubKey,
		PrivKey: privKey,
	}
}

// implement `Interface` in sort package.
type AddrKeysSlice []AddrKeys

func (b AddrKeysSlice) Len() int {
	return len(b)
}

// Sorts lexographically by Address
func (b AddrKeysSlice) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i].Address.Bytes(), b[j].Address.Bytes()) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
	}
}

func (b AddrKeysSlice) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// CreateGenAccounts generates genesis accounts loaded with coins, and returns
// their addresses, pubkeys, and privkeys.
func CreateGenAccounts(numAccs int, genCoins sdk.Coins) (genAccs []authtypes.BaseAccount, genBals []banktypes.Balance,
	addrs []sdk.AccAddress, pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {

	addrKeysSlice := AddrKeysSlice{}

	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())

		addrKeysSlice = append(addrKeysSlice, NewAddrKeys(addr, pubKey, privKey))
	}

	sort.Sort(addrKeysSlice)

	for i := range addrKeysSlice {
		addrs = append(addrs, addrKeysSlice[i].Address)
		pubKeys = append(pubKeys, addrKeysSlice[i].PubKey)
		privKeys = append(privKeys, addrKeysSlice[i].PrivKey)
		genAccs = append(genAccs, authtypes.BaseAccount{
			Address: addrKeysSlice[i].Address,
		})
		genBals = append(genBals, banktypes.Balance{
			Address: addrKeysSlice[i].Address,
			Coins:   genCoins,
		})
	}

	return
}

// SetGenesis sets the mock app genesis accounts.
func SetGenesis(app *App, accs []authtypes.BaseAccount, bals []banktypes.Balance) {
	// Pass the accounts in via the application (lazy) instead of through
	// RequestInitChain.
	app.GenesisAccounts = accs
	app.GenesisBalances = bals

	app.InitChain(abci.RequestInitChain{})
	app.Commit()
}

// GenTx generates a signed mock transaction.
func GenTx(msgs []sdk.Msg, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) authtypes.StdTx {
	// Make the transaction free
	fee := authtypes.StdFee{
		Amount: sdk.NewCoins(sdk.NewInt64Coin("foocoin", 0)),
		Gas:    100000,
	}

	sigs := make([]authtypes.StdSignature, len(priv))
	memo := "testmemotestmemo"

	for i, p := range priv {
		sig, err := p.Sign(authtypes.StdSignBytes(chainID, accnums[i], seq[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}

		sigs[i] = authtypes.StdSignature{
			PubKey:    p.PubKey().Bytes(),
			Signature: sig,
		}
	}

	return authtypes.NewStdTx(msgs, fee, sigs, memo)
}

// GeneratePrivKeys generates a total n secp256k1 private keys.
func GeneratePrivKeys(n int) (keys []crypto.PrivKey) {
	// TODO: Randomize this between ed25519 and secp256k1
	keys = make([]crypto.PrivKey, n)
	for i := 0; i < n; i++ {
		keys[i] = secp256k1.GenPrivKey()
	}

	return
}

// GeneratePrivKeyAddressPairs generates a total of n private key, address
// pairs.
func GeneratePrivKeyAddressPairs(n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	keys = make([]crypto.PrivKey, n)
	addrs = make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		if rand.Int63()%2 == 0 {
			keys[i] = secp256k1.GenPrivKey()
		} else {
			keys[i] = ed25519.GenPrivKey()
		}
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}

// GeneratePrivKeyAddressPairsFromRand generates a total of n private key, address
// pairs using the provided randomness source.
func GeneratePrivKeyAddressPairsFromRand(rand *rand.Rand, n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	keys = make([]crypto.PrivKey, n)
	addrs = make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		secret := make([]byte, 32)
		_, err := rand.Read(secret)
		if err != nil {
			panic("Could not read randomness")
		}
		if rand.Int63()%2 == 0 {
			keys[i] = secp256k1.GenPrivKeySecp256k1(secret)
		} else {
			keys[i] = ed25519.GenPrivKeyFromSecret(secret)
		}
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}

// RandomSetGenesis set genesis accounts with random coin values using the
// provided addresses and coin denominations.
// nolint: errcheck
func RandomSetGenesis(r *rand.Rand, app *App, addrs []sdk.AccAddress, denoms []string) {
	accts := make([]authtypes.BaseAccount, len(addrs))
	randCoinIntervals := []BigInterval{
		{sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 1)},
		{sdk.NewIntWithDecimal(1, 2), sdk.NewIntWithDecimal(1, 3)},
		{sdk.NewIntWithDecimal(1, 40), sdk.NewIntWithDecimal(1, 50)},
	}

	for i := 0; i < len(accts); i++ {
		coins := make([]sdk.Coin, len(denoms))

		// generate a random coin for each denomination
		for j := 0; j < len(denoms); j++ {
			coins[j] = sdk.Coin{Denom: denoms[j],
				Amount: RandFromBigInterval(r, randCoinIntervals).ToDec(),
			}
		}

		app.TotalCoinsSupply = app.TotalCoinsSupply.Add(coins...)
		baseAcc := authtypes.NewBaseAccountWithAddress(addrs[i])

		app.BankKeeper.SetBalances(app.NewContext(false, abci.Header{}), baseAcc.GetAddress(), coins)
		accts[i] = *baseAcc
	}
	app.GenesisAccounts = accts
}

// GenSequenceOfTxs generates a set of signed transactions of messages, such
// that they differ only by having the sequence numbers incremented between
// every transaction.
func GenSequenceOfTxs(msgs []sdk.Msg, accnums []uint64, initSeqNums []uint64, numToGenerate int, priv ...crypto.PrivKey) []authtypes.StdTx {
	txs := make([]authtypes.StdTx, numToGenerate)
	for i := 0; i < numToGenerate; i++ {
		txs[i] = GenTx(msgs, accnums, initSeqNums, priv...)
		incrementAllSequenceNumbers(initSeqNums)
	}

	return txs
}

func incrementAllSequenceNumbers(initSeqNums []uint64) {
	for i := 0; i < len(initSeqNums); i++ {
		initSeqNums[i]++
	}
}
