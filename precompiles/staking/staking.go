// Copyright 2022 Evmos Foundation
// This file is part of the Evmos Network packages.
//
// Evmos is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Evmos packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Evmos packages. If not, see https://github.com/evmos/evmos/blob/main/LICENSE

package staking

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/evmos/ethermint/x/evm/statedb"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var _ vm.PrecompiledContract = (*StakingPrecompile)(nil)

func init() {
	addressType, _ := abi.NewType("address", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	uint64Type, _ := abi.NewType("uint64", "", nil)

	DelegateMethod = abi.NewMethod(
		"delegate", // name
		"delegate", // raw name
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{
				Name: "delegatorAddress",
				Type: addressType,
			},
			{
				Name: "validatorAddress",
				Type: stringType,
			},
			{
				Name: "amount",
				Type: uint256Type,
			},
		},
		abi.Arguments{},
	)

	UndelegateMethod = abi.NewMethod(
		"undelegate", // name
		"undelegate", // raw name
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{
				Name: "delegatorAddress",
				Type: addressType,
			},
			{
				Name: "validatorAddress",
				Type: stringType,
			},
			{
				Name: "amount",
				Type: uint256Type,
			},
		},
		abi.Arguments{
			{
				Name: "completionTime",
				Type: uint64Type,
			},
		},
	)

	RedelegateMethod = abi.NewMethod(
		"redelegate", // name
		"redelegate", // raw name
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{
				Name: "delegatorAddress",
				Type: addressType,
			},
			{
				Name: "validatorSrcAddress",
				Type: stringType,
			},
			{
				Name: "validatorDstAddress",
				Type: stringType,
			},
			{
				Name: "amount",
				Type: uint256Type,
			},
		},
		abi.Arguments{
			{
				Name: "completionTime",
				Type: uint64Type,
			},
		},
	)

	CancelUnbondingDelegationMethod = abi.NewMethod(
		"cancelUnbondingDelegation", // name
		"cancelUnbondingDelegation", // raw name
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{
				Name: "delegatorAddress",
				Type: addressType,
			},
			{
				Name: "validatorSrcAddress",
				Type: stringType,
			},
			{
				Name: "validatorDstAddress",
				Type: stringType,
			},
			{
				Name: "amount",
				Type: uint256Type,
			},
		},
		abi.Arguments{},
	)
}

type StakingPrecompile struct {
	stakingKeeper stakingkeeper.Keeper
}

func NewStakingPrecompile(
	stakingKeeper stakingkeeper.Keeper,
) vm.PrecompiledContract {
	return &StakingPrecompile{
		stakingKeeper: stakingKeeper,
	}
}

// RequiredGas calculates the contract gas use
func (sp *StakingPrecompile) RequiredGas(input []byte) uint64 {
	// TODO: estimate required gas since this is stateful
	return 0
}

func (sp *StakingPrecompile) Run(_ []byte) ([]byte, error) {
	return nil, errors.New("should run with RunStateful")
}

func (sp *StakingPrecompile) RunStateful(evm vm.EVM, caller common.Address, input []byte, value *big.Int) ([]byte, error) {
	stateDB, ok := evm.StateDB.(statedb.ExtStateDB)
	if !ok {
		return nil, errors.New("not run in ethermint")
	}

	ctx := stateDB.Context()

	methodID := string(input[:4])
	argsBz := input[4:]

	switch methodID {
	// Staking transactions
	case string(DelegateMethod.ID):
		return sp.Delegate(ctx, argsBz, stateDB)
	case string(UndelegateMethod.ID):
		return sp.Undelegate(ctx, argsBz, stateDB)
	case string(RedelegateMethod.ID):
		return sp.Redelegate(ctx, argsBz, stateDB)
	case string(CancelUnbondingDelegationMethod.ID):
		return sp.CancelUnbondingDelegation(ctx, argsBz, stateDB)
		// Staking queries
	case string(DelegationMethod.ID):
		return sp.Delegation(ctx, argsBz, stateDB)
	case string(UnbondingDelegationMethod.ID):
		return sp.UnbondingDelegation(ctx, argsBz, stateDB)
	case string(ValidatorMethod.ID):
		return sp.Validator(ctx, argsBz, stateDB)

	// TODO: get delegation
	default:
		return nil, fmt.Errorf("unknown method '%s'", methodID)
	}
}
