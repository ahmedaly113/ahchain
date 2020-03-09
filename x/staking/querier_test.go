package staking

import (
	"fmt"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	app "github.com/ahmedaly113/ahchain/types"
)

func TestQuerier_EmptyTopArgument(t *testing.T) {
	ctx, k, _ := mockDB()

	querier := NewQuerier(k)
	queryParams := QueryClaimTopArgumentParams{
		ClaimID: 1,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryClaimTopArgument}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryParams)
	bz, err := querier(ctx, []string{QueryClaimTopArgument}, query)
	assert.NoError(t, err)
	argument := Argument{}
	jsonErr := k.codec.UnmarshalJSON(bz, &argument)
	assert.NoError(t, jsonErr)
	assert.Equal(t, uint64(0), argument.ID)

}

func TestQuerier_EarnedCoins(t *testing.T) {
	ctx, k, _ := mockDB()
	_, _, address := keyPubAddr()
	usersEarnings := make([]UserEarnedCoins, 0)
	coins := sdk.NewCoins(sdk.NewInt64Coin("crypto", app.Shanev*10),
		sdk.NewInt64Coin("random", app.Shanev*30))
	userEarnings := UserEarnedCoins{
		Address: address,
		Coins:   coins,
	}
	usersEarnings = append(usersEarnings, userEarnings)
	genesisState := NewGenesisState(nil, nil, usersEarnings, DefaultParams())
	InitGenesis(ctx, k, genesisState)

	querier := NewQuerier(k)
	queryParams := QueryEarnedCoinsParams{
		Address: address,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryEarnedCoins}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryParams)
	bz, err := querier(ctx, []string{QueryEarnedCoins}, query)
	assert.NoError(t, err)
	earnedCoins := sdk.Coins{}
	jsonErr := k.codec.UnmarshalJSON(bz, &earnedCoins)
	assert.NoError(t, jsonErr)
	assert.True(t, coins.IsEqual(earnedCoins))
	// total

	queryTotalParams := QueryTotalEarnedCoinsParams{
		Address: address,
	}

	query = abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryTotalEarnedCoins}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryTotalParams)
	bz, err = querier(ctx, []string{QueryTotalEarnedCoins}, query)
	assert.NoError(t, err)
	totalEarned := sdk.Coin{}
	jsonErr = k.codec.UnmarshalJSON(bz, &totalEarned)
	assert.NoError(t, jsonErr)
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*40), totalEarned)

}

func TestQuerier_ArgumentsByIDs(t *testing.T) {
	ctx, k, mdb := mockDB()
	ctx = ctx.WithBlockTime(time.Now())
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})

	argument1, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)

	argument2, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)

	querier := NewQuerier(k)
	queryParams := QueryArgumentsByIDsParams{
		ArgumentIDs: []uint64{argument1.ID, argument2.ID},
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryArgumentsByIDs}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryParams)
	bz, err := querier(ctx, []string{QueryArgumentsByIDs}, query)
	assert.NoError(t, err)

	var arguments []Argument
	k.codec.UnmarshalJSON(bz, &arguments)
	assert.Len(t, arguments, 2)
}

func TestQuerier_CommunityStakes(t *testing.T) {
	ctx, k, mdb := mockDB()
	ctx = ctx.WithBlockTime(time.Now())
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})

	argument1, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)

	claim1, _ := k.claimKeeper.Claim(ctx, argument1.ClaimID)

	querier := NewQuerier(k)
	queryParams := QueryCommunityStakesParams{
		CommunityID: claim1.CommunityID,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryCommunityStakes}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryParams)
	bz, err := querier(ctx, []string{QueryCommunityStakes}, query)
	assert.NoError(t, err)

	var stakes []Stake
	k.codec.UnmarshalJSON(bz, &stakes)
	assert.Len(t, stakes, 2)
}

func TestQuerier_UserCommunityStakes(t *testing.T) {
	ctx, k, mdb := mockDB()
	ctx = ctx.WithBlockTime(time.Now())
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})

	argument1, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)

	claim1, _ := k.claimKeeper.Claim(ctx, argument1.ClaimID)

	querier := NewQuerier(k)
	queryParams := QueryUserCommunityStakesParams{
		Address:     argument1.Creator,
		CommunityID: claim1.CommunityID,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryUserCommunityStakes}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryParams)
	bz, err := querier(ctx, []string{QueryUserCommunityStakes}, query)
	assert.NoError(t, err)

	var stakes []Stake
	k.codec.UnmarshalJSON(bz, &stakes)
	assert.Len(t, stakes, 2)
}

func TestQueryParams_Success(t *testing.T) {
	ctx, keeper, _ := mockDB()

	onChainParams := keeper.GetParams(ctx)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryParams),
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryParams}, query)
	assert.Nil(t, err)

	var returnedParams Params
	sdkErr := ModuleCodec.UnmarshalJSON(resBytes, &returnedParams)
	assert.Nil(t, sdkErr)
	assert.Equal(t, returnedParams, onChainParams)
}
