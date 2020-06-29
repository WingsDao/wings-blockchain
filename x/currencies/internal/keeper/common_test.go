// +build unit

package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/dfinance/dnode/x/currencies/internal/types"
	"github.com/dfinance/dnode/x/multisig"
	"github.com/dfinance/dnode/x/poa"
)

const (
	defDenom    = "testcoin"
	defDecimals = 0
	defIssueID1 = "issue1"
	defIssueID2 = "issue2"
)

var (
	maccPerms map[string][]string = map[string][]string{
		auth.FeeCollectorName: nil,
	}
	defAmount = sdk.NewInt(10)
)

type TestInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	//
	keyMain    *sdk.KVStoreKey
	keyAccount *sdk.KVStoreKey
	keyCC      *sdk.KVStoreKey
	keySupply  *sdk.KVStoreKey
	keyParams  *sdk.KVStoreKey
	tkeyParams *sdk.TransientStoreKey
	keyPoa     *sdk.KVStoreKey
	keyMS      *sdk.KVStoreKey
	//
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	supplyKeeper  supply.Keeper
	paramsKeeper  params.Keeper
	keeper        Keeper
}

func (input *TestInput) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (input *TestInput) CreateAccount(t *testing.T, accName string, coins sdk.Coins) (accAddress sdk.AccAddress) {
	if coins == nil {
		coins = sdk.NewCoins()
	}

	addr := sdk.AccAddress(accName)
	acc := input.accountKeeper.NewAccountWithAddress(input.ctx, addr)
	require.NoError(t, acc.SetCoins(coins), "setting coins for accName: %s", accName)

	input.accountKeeper.SetAccount(input.ctx, acc)
	require.True(t, input.bankKeeper.GetCoins(input.ctx, addr).IsEqual(coins), "checking accName created: %s", accName)

	return addr
}

func NewTestInput(t *testing.T) TestInput {
	input := TestInput{
		cdc:        codec.New(),
		keyParams:  sdk.NewKVStoreKey(params.StoreKey),
		keyAccount: sdk.NewKVStoreKey(auth.StoreKey),
		keySupply:  sdk.NewKVStoreKey(supply.StoreKey),
		keyPoa:     sdk.NewKVStoreKey(poa.StoreKey),
		keyMS:      sdk.NewKVStoreKey(multisig.StoreKey),
		keyCC:      sdk.NewKVStoreKey(types.StoreKey),
		tkeyParams: sdk.NewTransientStoreKey(params.TStoreKey),
	}

	// register codec
	sdk.RegisterCodec(input.cdc)
	auth.RegisterCodec(input.cdc)
	bank.RegisterCodec(input.cdc)
	staking.RegisterCodec(input.cdc)
	codec.RegisterCrypto(input.cdc)
	multisig.RegisterCodec(input.cdc)
	types.RegisterCodec(input.cdc)

	// init in-memory DB
	db := dbm.NewMemDB()
	mstore := store.NewCommitMultiStore(db)
	mstore.MountStoreWithDB(input.keyParams, sdk.StoreTypeIAVL, db)
	mstore.MountStoreWithDB(input.keyAccount, sdk.StoreTypeIAVL, db)
	mstore.MountStoreWithDB(input.keySupply, sdk.StoreTypeIAVL, db)
	mstore.MountStoreWithDB(input.keyPoa, sdk.StoreTypeIAVL, db)
	mstore.MountStoreWithDB(input.keyMS, sdk.StoreTypeIAVL, db)
	mstore.MountStoreWithDB(input.keyCC, sdk.StoreTypeIAVL, db)
	mstore.MountStoreWithDB(input.tkeyParams, sdk.StoreTypeTransient, db)
	err := mstore.LoadLatestVersion()
	if err != nil {
		t.Fatal(err)
	}

	// create target and dependant keepers
	input.paramsKeeper = params.NewKeeper(input.cdc, input.keyParams, input.tkeyParams)
	input.accountKeeper = auth.NewAccountKeeper(input.cdc, input.keyAccount, input.paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	input.bankKeeper = bank.NewBaseKeeper(input.accountKeeper, input.paramsKeeper.Subspace(bank.DefaultParamspace), input.ModuleAccountAddrs())
	input.supplyKeeper = supply.NewKeeper(input.cdc, input.keySupply, input.accountKeeper, input.bankKeeper, maccPerms)
	input.keeper = NewKeeper(input.cdc, input.keyCC, input.bankKeeper)

	// create context
	input.ctx = sdk.NewContext(mstore, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return input
}
