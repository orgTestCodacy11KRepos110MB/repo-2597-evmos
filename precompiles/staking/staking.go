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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

const abiPath = "./abi.json"

// FIXME: fix on fork
// var _ vm.PrecompiledContract = &Precompile{}

// Precompile defines the precompiled contract for staking.
type Precompile struct {
	abi.ABI
	stakingKeeper stakingkeeper.Keeper
}

// NewPrecompile creates a new staking Precompile instance as a
// PrecompiledContract interface.
func NewPrecompile(
	stakingKeeper stakingkeeper.Keeper,
) (*Precompile, error) {
	abiJSON, err := ioutil.ReadFile(filepath.Clean(abiPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open abi.json file: %w", err)
	}

	abi, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		return nil, fmt.Errorf("invalid abi.json file: %w", err)
	}

	return &Precompile{
		ABI:           abi,
		stakingKeeper: stakingKeeper,
	}, nil
}

// Address defines the address of the staking compile contract.
// address: 0x0000000000000000000000000000000000000100
func (Precompile) Address() common.Address {
	return common.BytesToAddress([]byte{100})
}

// IsStateful returns true since the precompile contract has access to the
// staking state.
func (Precompile) IsStateful() bool {
	return true
}

// RequiredGas calculates the contract gas use
func (*Precompile) RequiredGas(input []byte) uint64 {
	return 0
}

// Run executes the precompile contract staking methods defined in the ABI.
func (p *Precompile) Run(evm *vm.EVM, contract *vm.Contract, input []byte, readOnly bool) ([]byte, error) {
	// TODO:
	// stateDB, ok := evm.StateDB.(*statedb.StateDB)
	// if !ok {
	// 	return nil, errors.New("not run in ethermint")
	// }

	// ctx := stateDB.GetContext()
	ctx := sdk.Context{}

	methodID := input[:4]

	// NOTE: this function iterates over the method map and returns
	// the method with the given ID
	method, err := p.ABI.MethodById(methodID)
	if err != nil {
		return nil, err
	}

	// return error if trying to write to state during a read-only call
	if readOnly && p.IsTransaction(method.Name) {
		return nil, vm.ErrWriteProtection
	}

	argsBz := input[4:]

	switch method.Name {
	// Staking transactions
	case DelegateMethod:
		return p.Delegate(ctx, contract, method, argsBz)
	case UndelegateMethod:
		return p.Undelegate(ctx, contract, method, argsBz)
	case RedelegateMethod:
		return p.Redelegate(ctx, contract, method, argsBz)
	case CancelUnbondingDelegationMethod:
		return p.CancelUnbondingDelegation(ctx, contract, method, argsBz)
		// Staking queries
	case DelegationMethod:
		return p.Delegation(ctx, contract, method, argsBz)
	case UnbondingDelegationMethod:
		return p.UnbondingDelegation(ctx, contract, method, argsBz)
	// case ValidatorMethod:
	// 	return p.Validator(ctx, contract, argsBz, stateDB, readOnly)
	// case RedelegationsMethod:
	// 	return p.Redelegations(ctx, contract, argsBz, stateDB, readOnly)

	// TODO: Add other queries
	// TODO: how do we handle paginations?
	default:
		return nil, fmt.Errorf("no method with id: %#x", methodID)
	}
}

func (Precompile) IsTransaction(methodID string) bool {
	switch methodID {
	case DelegateMethod,
		UndelegateMethod,
		RedelegateMethod,
		CancelUnbondingDelegationMethod:
		return true
	default:
		return false
	}
}
