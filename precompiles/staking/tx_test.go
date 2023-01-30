package staking_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/evmos/v11/precompiles/staking"
)

func (s *PrecompileTestSuite) TestDelegate() {
	var argsBz []byte
	method := s.precompile.Methods[staking.DelegateMethod]
	validatorAddr, err := sdk.ValAddressFromBech32(s.validators[0].OperatorAddress)
	s.Require().NoError(err)

	testCases := []struct {
		name       string
		argsFn     func() []byte
		gas        uint64
		delegation *big.Int
		expError   bool
	}{
		{
			"fail - empty input args",
			func() []byte { return []byte{} },
			200000,
			big.NewInt(0),
			true,
		},
		{
			"fail - message validation failed",
			func() []byte {
				argsBz, err = method.Inputs.Pack(
					s.address,
					s.validators[0].OperatorAddress,
					"",
					big.NewInt(1),
				)
				s.Require().NoError(err)
				return argsBz
			},
			200000,
			big.NewInt(1),
			true,
		},
		{
			"fail - delegation failed",
			func() []byte {
				argsBz, err = method.Inputs.Pack(
					s.address,
					s.validators[0].OperatorAddress,
					s.bondDenom,
					big.NewInt(100000000000),
				)
				s.Require().NoError(err)
				return argsBz
			},
			200000,
			big.NewInt(100000000000),
			true,
		},
		{
			"fail - out of gas",
			func() []byte {
				argsBz, err = method.Inputs.Pack(
					s.address,
					s.validators[0].OperatorAddress,
					s.bondDenom,
					big.NewInt(100),
				)
				s.Require().NoError(err)
				return argsBz
			},
			200,
			big.NewInt(100),
			true,
		},
		{
			"success",
			func() []byte {
				argsBz, err = method.Inputs.Pack(
					s.address,
					s.validators[0].OperatorAddress,
					s.bondDenom,
					big.NewInt(100),
				)
				s.Require().NoError(err)
				return argsBz
			},
			20000,
			big.NewInt(100),
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			argsBz := tc.argsFn()
			contract := vm.NewContract(vm.AccountRef(s.address), s.precompile, big.NewInt(0), tc.gas)

			s.ctx = s.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
			initialGas := s.ctx.GasMeter().GasConsumed()
			s.Require().Zero(initialGas)

			bz, err := s.precompile.Delegate(s.ctx, contract, &method, argsBz)
			gasConsumed := s.ctx.GasMeter().GasConsumed()

			delegation := s.app.StakingKeeper.Delegation(s.ctx, s.address[:], validatorAddr)
			if tc.expError {
				s.Require().Error(err)
				s.Require().Empty(bz)
				s.Require().Nil(delegation)
				return
			}

			delegationAmt := sdk.NewIntFromBigInt(tc.delegation)
			s.Require().NoError(err)
			s.Require().Empty(bz)
			s.Require().NotNil(delegation)
			s.Require().Equal(delegationAmt, delegation.GetShares().TruncateInt())
			s.Require().Equal(int64(tc.gas-contract.Gas), int64(gasConsumed), "gas consumed should be equal")
		})
	}
}
