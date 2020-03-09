package staking

import (
	"fmt"
	"reflect"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	app "github.com/ahmedaly113/ahchain/types"
)

var (
	ParamKeyPeriod                   = []byte("period")
	ParamKeyArgumentCreationStake    = []byte("argumentCreationStake")
	ParamKeyArgumentBodyMaxLength    = []byte("argumentBodyMaxLength")
	ParamKeyArgumentBodyMinLength    = []byte("argumentBodyMinLength")
	ParamKeyArgumentSummaryMaxLength = []byte("argumentSummaryMaxLength")
	ParamKeyArgumentSummaryMinLength = []byte("argumentSummaryMinLength")
	ParamKeyUpvoteStake              = []byte("upvoteStake")
	ParamKeyCreatorShare             = []byte("creatorShare")
	ParamKeyInterestRate             = []byte("interestRate")
	ParamKeyStakingAdmins            = []byte("stakingAdmins")
	ParamKeyStakeLimitPercent        = []byte("stakeLimitPercent")
	ParamKeyStakeLimitDays           = []byte("stakeLimitDays")
	ParamKeyUnjailUpvotes            = []byte("unjailUpvotes")
	ParamKeyMaxArgumentsPerClaim     = []byte("maxArgumentsPerClaim")
)

type Params struct {
	Period                   time.Duration    `json:"period"`
	ArgumentCreationStake    sdk.Coin         `json:"argument_creation_stake"`
	ArgumentBodyMaxLength    int              `json:"argument_body_max_length"`
	ArgumentBodyMinLength    int              `json:"argument_body_min_length"`
	ArgumentSummaryMaxLength int              `json:"argument_summary_max_length"`
	ArgumentSummaryMinLength int              `json:"argument_summary_min_length"`
	UpvoteStake              sdk.Coin         `json:"upvote_stake"`
	CreatorShare             sdk.Dec          `json:"creator_share"`
	InterestRate             sdk.Dec          `json:"interest_rate"`
	StakingAdmins            []sdk.AccAddress `json:"staking_admins"`
	// deprecated
	StakeLimitPercent sdk.Dec `json:"stake_limit_percent"`
	// deprecated
	StakeLimitDays       time.Duration `json:"stake_limit_days"`
	UnjailUpvotes        int           `json:"unjail_upvotes"`
	MaxArgumentsPerClaim int           `json:"max_arguments_per_claim"`
}

func DefaultParams() Params {
	return Params{
		Period:                   time.Hour * 24 * 7,
		ArgumentCreationStake:    sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
		ArgumentBodyMaxLength:    1250,
		ArgumentBodyMinLength:    25,
		ArgumentSummaryMaxLength: 140,
		ArgumentSummaryMinLength: 25,
		UpvoteStake:              sdk.NewInt64Coin(app.StakeDenom, app.Shanev*10),
		CreatorShare:             sdk.NewDecWithPrec(50, 2),
		InterestRate:             sdk.NewDecWithPrec(105, 2),
		StakingAdmins:            []sdk.AccAddress{},
		StakeLimitPercent:        sdk.NewDecWithPrec(667, 3),
		StakeLimitDays:           time.Hour * 24 * 7,
		UnjailUpvotes:            1,
		MaxArgumentsPerClaim:     5,
	}
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: ParamKeyPeriod, Value: &p.Period},
		{Key: ParamKeyArgumentCreationStake, Value: &p.ArgumentCreationStake},
		{Key: ParamKeyArgumentBodyMaxLength, Value: &p.ArgumentBodyMaxLength},
		{Key: ParamKeyArgumentBodyMinLength, Value: &p.ArgumentBodyMinLength},
		{Key: ParamKeyArgumentSummaryMaxLength, Value: &p.ArgumentSummaryMaxLength},
		{Key: ParamKeyArgumentSummaryMinLength, Value: &p.ArgumentSummaryMinLength},
		{Key: ParamKeyUpvoteStake, Value: &p.UpvoteStake},
		{Key: ParamKeyCreatorShare, Value: &p.CreatorShare},
		{Key: ParamKeyInterestRate, Value: &p.InterestRate},
		{Key: ParamKeyStakingAdmins, Value: &p.StakingAdmins},
		{Key: ParamKeyStakeLimitPercent, Value: &p.StakeLimitPercent},
		{Key: ParamKeyStakeLimitDays, Value: &p.StakeLimitDays},
		{Key: ParamKeyUnjailUpvotes, Value: &p.UnjailUpvotes},
		{Key: ParamKeyMaxArgumentsPerClaim, Value: &p.MaxArgumentsPerClaim},
	}
}

// ParamKeyTable for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the staking module
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for staking module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", ModuleName)
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("loaded staking params: %+v", params))
}

// UpdateParams updates the required params
func (k Keeper) UpdateParams(ctx sdk.Context, updater sdk.AccAddress, updates Params, updatedFields []string) sdk.Error {
	if !k.isAdmin(ctx, updater) {
		err := ErrAddressNotAuthorised()
		return err
	}

	current := k.GetParams(ctx)
	updated := k.getUpdatedParams(current, updates, updatedFields)
	k.SetParams(ctx, updated)

	return nil
}

func (k Keeper) getUpdatedParams(current Params, updates Params, updatedFields []string) Params {
	updated := current
	mapParams(updates, func(param string, index int, field reflect.StructField) {
		if isIn(param, updatedFields) {
			reflect.ValueOf(&updated).Elem().FieldByName(field.Name).Set(
				reflect.ValueOf(
					reflect.ValueOf(updates).FieldByName(field.Name).Interface(),
				),
			)
		}
	})

	return updated
}

func isIn(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}

// mapParams walks over each param, and ignores the *_admins param because they are out of scope for this CLI command
func mapParams(params interface{}, fn func(param string, index int, field reflect.StructField)) {
	rParams := reflect.TypeOf(params)
	for i := 0; i < rParams.NumField(); i++ {
		field := rParams.Field(i)
		param := field.Tag.Get("json")
		fn(param, i, field)
	}
}
