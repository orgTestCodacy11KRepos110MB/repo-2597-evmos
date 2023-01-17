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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/evmos/ethermint/x/evm/statedb"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var _ vm.PrecompiledContract = &StakingPrecompile{}

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

// Address defines the address of the staking compile contract.
// address: 0x0000000000000000000000000000000000000100
func (StakingPrecompile) Address() common.Address {
	return common.BytesToAddress([]byte{100})
}

// IsStateful returns true since the precompile contract has access to the
// staking state.
func (StakingPrecompile) IsStateful() bool {
	return true
}

// RequiredGas calculates the contract gas use
func (sp *StakingPrecompile) RequiredGas(input []byte) uint64 {
	return 0
}

func (sp *StakingPrecompile) Run(evm *vm.EVM, contract *vm.Contract, input []byte, readOnly bool) ([]byte, error) {
	stateDB, ok := evm.StateDB.(*statedb.StateDB)
	if !ok {
		return nil, errors.New("not run in ethermint")
	}

	// ctx := stateDB.GetContext()
	ctx := sdk.Context{}

	methodID := string(input[:4])
	argsBz := input[4:]

	switch methodID {
	// Staking transactions
	case string(DelegateMethod.ID):
		return sp.Delegate(ctx, contract, argsBz, stateDB, readOnly)
	case string(UndelegateMethod.ID):
		return sp.Undelegate(ctx, contract, argsBz, stateDB, readOnly)
	case string(RedelegateMethod.ID):
		return sp.Redelegate(ctx, contract, argsBz, stateDB, readOnly)
	case string(CancelUnbondingDelegationMethod.ID):
		return sp.CancelUnbondingDelegation(ctx, contract, argsBz, stateDB, readOnly)
		// Staking queries
	case string(DelegationMethod.ID):
		return sp.Delegation(ctx, contract, argsBz, stateDB, readOnly)
	case string(UnbondingDelegationMethod.ID):
		return sp.UnbondingDelegation(ctx, contract, argsBz, stateDB, readOnly)
	case string(ValidatorMethod.ID):
		return sp.Validator(ctx, contract, argsBz, stateDB, readOnly)

	// TODO: Add other queries
	// TODO: how do we handle paginations?
	default:
		return nil, fmt.Errorf("unknown method '%s'", methodID)
	}
}
