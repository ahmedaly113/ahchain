package account

import (
	"github.com/ahmedaly113/ahchain/x/bank/exported"
	"github.com/ahmedaly113/ahchain/x/distribution"
)

const (
	TransactionGift    = exported.TransactionGift
	TransactionBacking = exported.TransactionBacking

	UserGrowthPoolName = distribution.UserGrowthPoolName
)

var (
	FromModuleAccount = exported.FromModuleAccount
)
