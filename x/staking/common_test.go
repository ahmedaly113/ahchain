package staking

import (
	"time"

	"github.com/ahmedaly113/ahchain/x/account"

	app "github.com/ahmedaly113/ahchain/types"
	trubank "github.com/ahmedaly113/ahchain/x/bank"
	"github.com/ahmedaly113/ahchain/x/claim"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type mockedAccountKeeper struct {
	jailStatus   map[string]bool
	forceFailure bool
}

func newAccountKeeper() *mockedAccountKeeper {
	return &mockedAccountKeeper{
		jailStatus: make(map[string]bool),
	}
}

func (m *mockedAccountKeeper) jail(address sdk.AccAddress) {
	m.jailStatus[address.String()] = true
}

func (m *mockedAccountKeeper) fail() {
	m.forceFailure = true
}
func (m *mockedAccountKeeper) IsJailed(ctx sdk.Context, address sdk.AccAddress) (bool, sdk.Error) {
	if m.forceFailure {
		m.forceFailure = false
		return false, sdk.ErrInternal("error")
	}
	j, _ := m.jailStatus[address.String()]
	if j {
		return true, nil
	}
	return false, nil
}

func (m *mockedAccountKeeper) UnJail(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	if m.forceFailure {
		m.forceFailure = false
		return sdk.ErrInternal("error")
	}
	m.jailStatus[address.String()] = false
	return nil
}

func (m *mockedAccountKeeper) IterateAppAccounts(ctx sdk.Context, cb func(acc account.AppAccount) (stop bool)) {

}

type mockClaimKeeper struct {
	claims           map[uint64]claim.Claim
	enableTrackStake bool
}

func newMockedClaimKeeper() *mockClaimKeeper {
	return &mockClaimKeeper{
		claims:           make(map[uint64]claim.Claim),
		enableTrackStake: false,
	}
}

func (m *mockClaimKeeper) SetClaims(claims map[uint64]claim.Claim) {
	m.claims = claims
}

func (m *mockClaimKeeper) Claim(ctx sdk.Context, id uint64) (claim.Claim, bool) {
	if len(m.claims) == 0 {
		return claim.Claim{CommunityID: "testunit",
			TotalBacked:     sdk.NewInt64Coin(app.StakeDenom, 0),
			TotalChallenged: sdk.NewInt64Coin(app.StakeDenom, 0),
		}, true
	}
	c, ok := m.claims[id]
	return c, ok
}
func (m *mockClaimKeeper) AddBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	if !m.enableTrackStake {
		return nil
	}
	c, ok := m.Claim(ctx, id)
	if !ok {
		return sdk.ErrInternal("unknown claim")
	}
	c.TotalBacked = c.TotalBacked.Add(stake)
	m.claims[id] = c
	return nil
}

func (m *mockClaimKeeper) SetFirstArgumentTime(ctx sdk.Context, id uint64, firstArgumentTime time.Time) sdk.Error {
	if !m.enableTrackStake {
		return nil
	}
	c, ok := m.Claim(ctx, id)
	if !ok {
		return sdk.ErrInternal("unknown claim")
	}
	c.FirstArgumentTime = firstArgumentTime
	m.claims[id] = c
	return nil
}

func (m *mockClaimKeeper) AddChallengeStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	if !m.enableTrackStake {
		return nil
	}
	c, ok := m.Claim(ctx, id)
	if !ok {
		return sdk.ErrInternal("unknown claim")
	}
	c.TotalChallenged = c.TotalChallenged.Add(stake)
	m.claims[id] = c
	return nil
}

// SubtractBackingStake adds a stake amount to the total backing amount
func (m *mockClaimKeeper) SubtractBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	c, ok := m.Claim(ctx, id)
	if !ok {
		return sdk.ErrInternal("unknown claim")
	}
	c.TotalBacked = c.TotalBacked.Sub(stake)
	m.claims[id] = c

	return nil
}

// SubtractChallengeStake adds a stake amount to the total challenge amount
func (m *mockClaimKeeper) SubtractChallengeStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error {
	c, ok := m.Claim(ctx, id)
	if !ok {
		return sdk.ErrInternal("unknown claim")
	}
	c.TotalChallenged = c.TotalChallenged.Sub(stake)
	m.claims[id] = c

	return nil
}

type mockedDB struct {
	authAccKeeper auth.AccountKeeper
	accountKeeper AccountKeeper
	claimKeeper   ClaimKeeper
	bankKeeper    BankKeeper
	supplyKeeper  supply.Keeper
}

func mockDB() (sdk.Context, Keeper, *mockedDB) {
	db := dbm.NewMemDB()
	storeKey := sdk.NewKVStoreKey(ModuleName)
	accKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)
	bankKey := sdk.NewKVStoreKey("bank")
	claimKey := sdk.NewKVStoreKey(claim.StoreKey)
	supplyKey := sdk.NewKVStoreKey(supply.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(bankKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(claimKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(supplyKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	// codec registration
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)

	// Keepers
	pk := params.NewKeeper(cdc, paramsKey, transientParamsKey, params.DefaultCodespace)
	accKeeper := auth.NewAccountKeeper(cdc, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	distAcc := supply.NewEmptyModuleAccount(distribution.ModuleName)
	coins, _ := sdk.ParseCoins("1000000000utru") // 100TRU
	err := distAcc.SetCoins(coins)
	if err != nil {
		panic(err)
	}

	userRewardAcc := supply.NewEmptyModuleAccount(UserRewardPoolName, supply.Burner)
	err = userRewardAcc.SetCoins(coins)
	if err != nil {
		panic(err)
	}

	// so bank cannot use module accounts, only supply keeper
	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[distAcc.String()] = true
	blacklistedAddrs[userRewardAcc.String()] = true

	bankKeeper := bank.NewBaseKeeper(accKeeper,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		distribution.ModuleName: nil,
		UserRewardPoolName:      {supply.Burner},
		UserStakesPoolName:      {supply.Minter, supply.Burner},
	}

	supplyKeeper := supply.NewKeeper(cdc, supplyKey, accKeeper, bankKeeper, maccPerms)
	supplyKeeper.SetModuleAccount(ctx, distAcc)
	supplyKeeper.SetModuleAccount(ctx, userRewardAcc)
	totalSupply := coins
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	trubankKeeper := trubank.NewKeeper(cdc, bankKey, bankKeeper, pk.Subspace(trubank.DefaultParamspace), trubank.DefaultCodespace, supplyKeeper)

	mockedAccountKeeper := newAccountKeeper()
	mockedClaimKeeper := newMockedClaimKeeper()
	mockedClaimKeeper.claims = make(map[uint64]claim.Claim)
	keeper := NewKeeper(cdc, storeKey, mockedAccountKeeper, trubankKeeper, mockedClaimKeeper, supplyKeeper, pk.Subspace(DefaultParamspace), DefaultCodespace)
	_, _, admin1 := keyPubAddr()
	_, _, admin2 := keyPubAddr()
	genesis := DefaultGenesisState()
	genesis.Params.StakingAdmins = append(genesis.Params.StakingAdmins, admin1, admin2)
	InitGenesis(ctx, keeper, genesis)
	trubank.InitGenesis(ctx, trubankKeeper, trubank.DefaultGenesisState())

	mockedDB := &mockedDB{
		claimKeeper:   mockedClaimKeeper,
		accountKeeper: mockedAccountKeeper,
		authAccKeeper: accKeeper,
		bankKeeper:    trubankKeeper,
		supplyKeeper:  supplyKeeper,
	}
	return ctx, keeper, mockedDB
}

func createFakeFundedAccount(ctx sdk.Context, am auth.AccountKeeper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	return addr
}
func setCoins(ctx sdk.Context, am auth.AccountKeeper, coins sdk.Coins, addr sdk.AccAddress) {
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func mustParseTime(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return t
}

func afterCreatedTimeStakes(ctx sdk.Context, k Keeper, addr sdk.AccAddress, after time.Time) (stakes []Stake) {
	k.IterateAfterCreatedTimeUserStakes(ctx, addr, after, func(stake Stake) bool {
		stakes = append(stakes, stake)
		return false
	})
	return
}
