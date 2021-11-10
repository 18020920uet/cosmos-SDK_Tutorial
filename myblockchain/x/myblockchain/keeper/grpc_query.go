package keeper

import (
	"github.com/username/myblockchain/x/myblockchain/types"
)

var _ types.QueryServer = Keeper{}
