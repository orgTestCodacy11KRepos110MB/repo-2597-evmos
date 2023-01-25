package common

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Coin struct {
	Denom  string
	Amount *big.Int
}

func (c Coin) ToSDKType() sdk.Coin {
	return sdk.NewCoin(c.Denom, sdk.NewIntFromBigInt(c.Amount))
}
