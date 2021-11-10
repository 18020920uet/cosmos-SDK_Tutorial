package keeper

import (
	"github.com/username/interchange/x/ibcdex/types"
)

var _ types.QueryServer = Keeper{}
