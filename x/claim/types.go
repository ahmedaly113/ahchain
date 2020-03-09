package claim

import (
	"fmt"
	"net/url"
	"time"

	app "github.com/ahmedaly113/ahchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	StoreKey          = ModuleName
	DefaultParamspace = ModuleName
)

// Claim stores data about a claim
type Claim struct {
	ID                uint64         `json:"id"`
	CommunityID       string         `json:"community_id"`
	Body              string         `json:"body"`
	Creator           sdk.AccAddress `json:"creator"`
	Source            url.URL        `json:"source,omitempty"`
	TotalStakers      uint64         `json:"total_stakers,omitempty"`
	TotalBacked       sdk.Coin       `json:"total_backed,omitempty"`
	TotalChallenged   sdk.Coin       `json:"total_challenged,omitempty"`
	CreatedTime       time.Time      `json:"created_time"`
	FirstArgumentTime time.Time      `json:"first_argument_time"`
}

// Claims is an array of claims
type Claims []Claim

// NewClaim creates a new claim object
func NewClaim(id uint64, communityID string, body string, creator sdk.AccAddress, source url.URL, createdTime time.Time) Claim {
	return Claim{
		ID:              id,
		CommunityID:     communityID,
		Body:            body,
		Creator:         creator,
		Source:          source,
		TotalStakers:    0,
		TotalBacked:     sdk.NewCoin(app.StakeDenom, sdk.ZeroInt()),
		TotalChallenged: sdk.NewCoin(app.StakeDenom, sdk.ZeroInt()),
		CreatedTime:     createdTime,
	}
}

func (c Claim) String() string {
	return fmt.Sprintf(`Claim %d:
  CommunityID: %s
  Body:		   %s
  Creator:     %s
  Source:      %s
  CreatedTime  %s`,
		c.ID, c.CommunityID, c.Body, c.Creator.String(), c.Source.String(), c.CreatedTime.String())
}
