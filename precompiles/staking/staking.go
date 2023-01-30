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

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
func (p *Precompile) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) ([]byte, error) {
	// TODO:
	// stateDB, ok := evm.StateDB.(*statedb.StateDB)
	// if !ok {
	// 	return nil, errors.New("not run in ethermint")
	// }

	// ctx := stateDB.GetContext()
	ctx := sdk.Context{}

	methodID := contract.Input[:4]

	// NOTE: this function iterates over the method map and returns
	// the method with the given ID
	method, err := p.MethodById(methodID)
	if err != nil {
		return nil, err
	}

	// return error if trying to write to state during a read-only call
	if readOnly && p.IsTransaction(method.Name) {
		return nil, vm.ErrWriteProtection
	}

	argsBz := contract.Input[4:]

	initialGas := ctx.GasMeter().GasConsumed()

	// set the default SDK gas configuration to track gas usage
	ctx = ctx.WithKVGasConfig(storetypes.KVGasConfig()).
		WithKVGasConfig(storetypes.TransientGasConfig())

	// reset the gas configuration after state transition
	defer func() {
		ctx = ctx.WithKVGasConfig(storetypes.GasConfig{}).
			WithTransientKVGasConfig(storetypes.GasConfig{})
	}()

	// cache the context to avoid writing to state in case of failure or
	// out of gas
	cacheCtx, writeFn := ctx.CacheContext()

	var bz []byte
	switch method.Name {
	// Staking transactions
	case DelegateMethod:
		bz, err = p.Delegate(cacheCtx, contract, method, argsBz)
	case UndelegateMethod:
		bz, err = p.Undelegate(cacheCtx, contract, method, argsBz)
	case RedelegateMethod:
		bz, err = p.Redelegate(cacheCtx, contract, method, argsBz)
	case CancelUnbondingDelegationMethod:
		bz, err = p.CancelUnbondingDelegation(cacheCtx, contract, method, argsBz)
		// Staking queries
	case DelegationMethod:
		bz, err = p.Delegation(cacheCtx, contract, method, argsBz)
	case UnbondingDelegationMethod:
		bz, err = p.UnbondingDelegation(cacheCtx, contract, method, argsBz)
	// case ValidatorMethod:
	// 	bz, err = p.Validator(cacheCtx, method, argsBz)
	// case RedelegationsMethod:
	// 	bz, err = p.Redelegations(cacheCtx, method, argsBz)

	// TODO: Add other queries
	// TODO: how do we handle paginations?
	default:
		return nil, fmt.Errorf("no method with id: %#x", methodID)
	}

	if err != nil {
		return nil, err
	}

	cost := ctx.GasMeter().GasConsumed() - initialGas

	if !contract.UseGas(cost) {
		return nil, vm.ErrOutOfGas
	}

	// commit the changes to state
	writeFn()

	return bz, nil
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
