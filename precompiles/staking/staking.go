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
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/evmos/ethermint/x/evm/statedb"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

const abiPath = "./abi.json"

var _ vm.PrecompiledContract = &StakingPrecompile{}

// StakingPrecompile defines the precompiled contract for staking.
type StakingPrecompile struct {
	*abi.ABI
	stakingKeeper stakingkeeper.Keeper
}

// NewStakingPrecompile creates a new StakingPrecompile instance as a
// PrecompiledContract interface.
func NewStakingPrecompile(
	stakingKeeper stakingkeeper.Keeper,
) (vm.PrecompiledContract, error) {
	abiJSON, err := ioutil.ReadFile(filepath.Clean(abiPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open abi.json file: %w", err)
	}

	abi, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		return nil, fmt.Errorf("invalid abi.json file: %w", err)
	}

	return &StakingPrecompile{
		ABI:           abi,
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

// Run executes the data
func (sp *StakingPrecompile) Run(evm *vm.EVM, contract *vm.Contract, input []byte, readOnly bool) ([]byte, error) {
	stateDB, ok := evm.StateDB.(*statedb.StateDB)
	if !ok {
		return nil, errors.New("not run in ethermint")
	}

	// ctx := stateDB.GetContext()
	ctx := sdk.Context{}

	methodID := string(input[:4])
	argsBz := input[4:]

	switch string(methodID) {
	// Staking transactions
	case DelegateMethod:
		return sp.Delegate(ctx, contract, argsBz, stateDB, readOnly)
	case UndelegateMethod:
		return sp.Undelegate(ctx, contract, argsBz, stateDB, readOnly)
	case RedelegateMethod:
		return sp.Redelegate(ctx, contract, argsBz, stateDB, readOnly)
	case CancelUnbondingDelegationMethod:
		return sp.CancelUnbondingDelegation(ctx, contract, argsBz, stateDB, readOnly)
		// Staking queries
	case DelegationMethod:
		return sp.Delegation(ctx, contract, argsBz, stateDB, readOnly)
	case UnbondingDelegationMethod:
		return sp.UnbondingDelegation(ctx, contract, argsBz, stateDB, readOnly)
	case ValidatorMethod:
		return sp.Validator(ctx, contract, argsBz, stateDB, readOnly)

	// TODO: Add other queries
	// TODO: how do we handle paginations?
	default:
		return nil, fmt.Errorf("no method with id: %#x", methodID)
	}
}
