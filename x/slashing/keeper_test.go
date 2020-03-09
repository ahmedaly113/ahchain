package slashing

import (
	"testing"

	"github.com/ahmedaly113/ahchain/x/staking"

	app "github.com/ahmedaly113/ahchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	staker := keeper.GetParams(ctx).SlashAdmins[1]
	arg, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeChallenge)
	assert.NoError(t, err)

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[1]
	slash, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NoError(t, err)

	assert.NotZero(t, slash.ID)
	assert.Equal(t, uint64(2), arg.ID)
	assert.Equal(t, slash.Creator, creator)
}

func TestNewSlash_InvalidArgument(t *testing.T) {
	ctx, keeper := mockDB()

	invalidArgumentID := uint64(404)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, _, err := keeper.CreateSlash(ctx, invalidArgumentID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidArgument(invalidArgumentID).Code(), err.Code())
}

func TestNewSlash_InvalidDetailedReason(t *testing.T) {
	ctx, keeper := mockDB()
	_, publicKey, creator, coins := getFakeAppAccountParams()
	_, err := keeper.accountKeeper.CreateAppAccount(ctx, creator, coins, publicKey)
	assert.NoError(t, err)
	stakeID := uint64(1)
	longDetailedReason := "This is a very very very descriptive reason to slash an argument. I am writing it in this detail to make the validation fail. I hope it works!"
	_, _, err = keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonOther, longDetailedReason, creator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidSlashReason("").Code(), err.Code())
}

func TestNewSlash_ErrNotEnoughEarnedStake(t *testing.T) {
	ctx, keeper := mockDB()
	_, publicKey, creator, coins := getFakeAppAccountParams()
	_, err := keeper.accountKeeper.CreateAppAccount(ctx, creator, coins, publicKey)
	assert.NoError(t, err)
	stakeID := uint64(1)
	_, _, err = keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotEnoughEarnedStake(creator).Code(), err.Code())
}

func TestNewSlash_ErrAlreadyUnhelpful(t *testing.T) {
	ctx, keeper := mockDB()
	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.Nil(t, err)
	_, _, err = keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAlreadyUnhelpful().Code(), err.Code())
}

func TestSlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	createdSlash, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.Nil(t, err)

	returnedSlash, err := keeper.Slash(ctx, createdSlash.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdSlash, returnedSlash)
}

func TestSlash_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.Nil(t, err)

	_, err = keeper.Slash(ctx, uint64(404))

	assert.NotNil(t, err)
	assert.Equal(t, ErrSlashNotFound(uint64(404)).Code(), err.Code())
}

func TestSlashes_Success(t *testing.T) {
	ctx, keeper := mockDB()
	_, publicKey1, addr1, coins1 := getFakeAppAccountParams()
	_, publicKey2, addr2, coins2 := getFakeAppAccountParams()

	_, err := keeper.accountKeeper.CreateAppAccount(ctx, addr1, coins1, publicKey1)
	assert.NoError(t, err)
	_, err = keeper.accountKeeper.CreateAppAccount(ctx, addr2, coins2, publicKey2)
	assert.NoError(t, err)
	earned := sdk.NewCoins(sdk.NewInt64Coin("general", 70*app.Shanev))
	usersEarnings := []staking.UserEarnedCoins{
		staking.UserEarnedCoins{Address: addr1, Coins: earned},
		staking.UserEarnedCoins{Address: addr2, Coins: earned},
	}
	genesis := staking.DefaultGenesisState()
	genesis.UsersEarnings = usersEarnings
	staking.InitGenesis(ctx, keeper.stakingKeeper, genesis)

	p := keeper.GetParams(ctx)
	p.MinSlashCount = 2
	keeper.SetParams(ctx, p)

	assert.Equal(t, keeper.GetParams(ctx).MinSlashCount, 2)
	staker := keeper.GetParams(ctx).SlashAdmins[1]
	_, err = keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeBacking)
	assert.NoError(t, err)
	stakeID := uint64(1)

	first, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", addr1)
	assert.NoError(t, err)

	another, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", addr2)
	assert.NoError(t, err)

	all := keeper.Slashes(ctx)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0], first)
	assert.Equal(t, all[1], another)

	a, _ := keeper.stakingKeeper.Argument(ctx, 1)
	assert.True(t, a.IsUnhelpful)
}

func Test_punishment(t *testing.T) {
	ctx, keeper := mockDB()
	staker := keeper.GetParams(ctx).SlashAdmins[0]
	slasher := keeper.GetParams(ctx).SlashAdmins[1]
	slashMagnitude := keeper.GetParams(ctx).SlashMagnitude
	stakerStartingBalance := keeper.bankKeeper.GetCoins(ctx, staker)
	slasherStartingBalance := keeper.bankKeeper.GetCoins(ctx, slasher)
	assert.Equal(t, "300000000utru", stakerStartingBalance.String())
	assert.Equal(t, "300000000utru", slasherStartingBalance.String())

	claim, _ := keeper.claimKeeper.Claim(ctx, 1)
	assert.Equal(t, "0utru", claim.TotalChallenged.String())

	argument, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg2", "summary2", staker, claim.ID, staking.StakeChallenge)
	assert.NoError(t, err)

	stake, _ := keeper.stakingKeeper.Stake(ctx, 2)
	assert.Equal(t, argument.ID, stake.ArgumentID)
	assert.Equal(t, "50000000utru", stake.Amount.String())

	claim, _ = keeper.claimKeeper.Claim(ctx, 1)
	assert.Equal(t, stake.Amount.String(), claim.TotalChallenged.String())

	// staker should have = starting balance - stake amount
	stakerEndingBalance := keeper.bankKeeper.GetCoins(ctx, staker)
	expectedBalance := stakerStartingBalance.Sub(sdk.Coins{stake.Amount})
	assert.Equal(t, expectedBalance.String(), stakerEndingBalance.String())

	// this also does a punish because slasher is an admin
	_, _, err = keeper.CreateSlash(ctx, argument.ID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", slasher)
	assert.NoError(t, err)

	// staker should have = starting balance - (stake amount * slashMagnitude)
	slashPenalty := sdk.NewCoin(stake.Amount.Denom, stake.Amount.Amount.MulRaw(int64(slashMagnitude)))
	stakerEndingBalance = keeper.bankKeeper.GetCoins(ctx, staker)
	expectedBalance = stakerStartingBalance.Sub(sdk.Coins{slashPenalty})
	assert.Equal(t, expectedBalance.String(), stakerEndingBalance.String())

	// slasher should have = starting balance + reward (25% stake)
	slasherEndingBalance := keeper.bankKeeper.GetCoins(ctx, slasher)
	reward := stake.Amount.Amount.ToDec().Mul(sdk.NewDecWithPrec(25, 2)).TruncateInt()
	rewardCoin := sdk.NewCoin(stake.Amount.Denom, reward)
	expectedBalance = slasherStartingBalance.Add(sdk.Coins{rewardCoin})
	assert.Equal(t, expectedBalance.String(), slasherEndingBalance.String())

	claim, _ = keeper.claimKeeper.Claim(ctx, 1)
	assert.Equal(t, "0utru", claim.TotalChallenged.String())
}

func TestAddAdmin_Success(t *testing.T) {
	ctx, keeper := mockDB()

	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, _, newAdmin, _ := getFakeAppAccountParams()

	err := keeper.AddAdmin(ctx, newAdmin, creator)
	assert.Nil(t, err)

	newAdmins := keeper.GetParams(ctx).SlashAdmins
	assert.Subset(t, newAdmins, []sdk.AccAddress{newAdmin})
}

func TestAddAdmin_CreatorNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	invalidCreator := sdk.AccAddress([]byte{1, 2})
	_, _, newAdmin, _ := getFakeAppAccountParams()

	err := keeper.AddAdmin(ctx, newAdmin, invalidCreator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}

func TestRemoveAdmin_Success(t *testing.T) {
	ctx, keeper := mockDB()

	currentAdmins := keeper.GetParams(ctx).SlashAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.RemoveAdmin(ctx, adminToRemove, adminToRemove) // removing self
	assert.Nil(t, err)
	newAdmins := keeper.GetParams(ctx).SlashAdmins
	assert.Equal(t, len(currentAdmins)-1, len(newAdmins))
}

func TestRemoveAdmin_RemoverNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	invalidRemover := sdk.AccAddress([]byte{1, 2})
	currentAdmins := keeper.GetParams(ctx).SlashAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.AddAdmin(ctx, adminToRemove, invalidRemover)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}
