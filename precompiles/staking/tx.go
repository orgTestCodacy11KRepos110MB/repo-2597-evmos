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

	"github.com/ethereum/go-ethereum/core/vm"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/evmos/ethermint/x/evm/statedb"

	"github.com/evmos/evmos/v11/precompiles"
)

const (
	// DelegateMethod defines the ABI method name for the staking Delegate
	// transaction.
	DelegateMethod = "delegate"
	// UndelegateMethod defines the ABI method name for the staking Undelegate
	// transaction.
	UndelegateMethod = "undelegate"
	// RedelegateMethod defines the ABI method name for the staking Redelegate
	// transaction.
	RedelegateMethod = "redelegate"
	// CancelUnbondingDelegationMethod defines the ABI method name for the staking
	// CancelUnbondingDelegation transaction.
	CancelUnbondingDelegationMethod = "cancelUnbondingDelegation"
)

// Delegate performs the staking delegation.
func (sp *StakingPrecompile) Delegate(
	ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	if readOnly {
		return nil, vm.ErrWriteProtection
	}

	method, ok := sp.ABI.Methods[DelegateMethod]
	if !ok {
		return nil, fmt.Errorf("no method with id: %s", DelegationMethod)
	}

	var delegateInput DelegateInput
	err := precompiles.UnpackIntoInterface(&delegateInput, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	// verify that the delegator is the contract caller
	if delegateInput.DelegatorAddress != contract.Caller() {
		return nil, sdkerrors.ErrUnauthorized
	}

	msg, err := delegateInput.ToMessage()
	if err != nil {
		return nil, err
	}

	// calculate gas used in the Cosmos transaction

	initialGas := ctx.GasMeter().GasConsumed()

	// set the default SDK gas configuration to track gas usage
	ctx = ctx.WithKVGasConfig(storetypes.KVGasConfig()).
		WithKVGasConfig(storetypes.TransientGasConfig())

	msgSrv := stakingkeeper.NewMsgServerImpl(sp.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	if _, err := msgSrv.Delegate(sdk.WrapSDKContext(cacheCtx), msg); err != nil {
		return nil, err
	}

	cost := cacheCtx.GasMeter().GasConsumed() - initialGas

	if !contract.UseGas(cost) {
		return nil, vm.ErrOutOfGas
	}

	// commit the changes to state
	writeFn()

	return nil, nil
}

func (sp *StakingPrecompile) Undelegate(
	ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	if readOnly {
		return nil, vm.ErrWriteProtection
	}

	method, ok := sp.ABI.Methods[UndelegateMethod]
	if !ok {
		return nil, fmt.Errorf("no method with id: %s", DelegationMethod)
	}

	var undelegateInput UndelegateInput
	err := precompiles.UnpackIntoInterface(&undelegateInput, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	msg, err := undelegateInput.ToMessage()
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(sp.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	res, err := msgSrv.Undelegate(sdk.WrapSDKContext(cacheCtx), msg)
	if err != nil {
		return nil, err
	}

	output := new(UndelegateOutput).FromMessage(res)
	bz, err := output.Pack(method.Outputs)
	if err != nil {
		return nil, err
	}

	writeFn()

	return bz, nil
}

func (sp *StakingPrecompile) Redelegate(
	ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	if readOnly {
		return nil, vm.ErrWriteProtection
	}

	var redelegateInput RedelegateInput
	err := precompiles.UnpackIntoInterface(&redelegateInput, RedelegateMethod.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	msg, err := redelegateInput.ToMessage()
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(sp.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	res, err := msgSrv.BeginRedelegate(sdk.WrapSDKContext(cacheCtx), msg)
	if err != nil {
		return nil, err
	}

	output := new(RedelegateOutput).FromMessage(res)
	bz, err := output.Pack(RedelegateMethod.Outputs)
	if err != nil {
		return nil, err
	}

	writeFn()

	return bz, nil
}

func (sp *StakingPrecompile) CancelUnbondingDelegation(
	ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	if readOnly {
		return nil, vm.ErrWriteProtection
	}

	var cancelUnbondingDelegationInput CancelUnbondingDelegationInput
	err := precompiles.UnpackIntoInterface(&cancelUnbondingDelegationInput, CancelUnbondingDelegationMethod.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	msg, err := cancelUnbondingDelegationInput.ToMessage()
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(sp.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	if _, err := msgSrv.CancelUnbondingDelegation(sdk.WrapSDKContext(cacheCtx), msg); err != nil {
		return nil, err
	}

	writeFn()

	return nil, nil
}
