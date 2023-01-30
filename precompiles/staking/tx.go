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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

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
func (p *Precompile) Delegate(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	argsBz []byte,
) ([]byte, error) {
	args, err := method.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, err
	}

	msg, err := NewMsgDelegate(args)
	if err != nil {
		return nil, err
	}

	// TODO: verify that the delegator is the contract caller
	// if msg.DelegatorAddress != contract.Caller() {
	// 	return nil, sdkerrors.ErrUnauthorized
	// }

	// calculate gas used in the Cosmos transaction

	initialGas := ctx.GasMeter().GasConsumed()

	// set the default SDK gas configuration to track gas usage
	ctx = ctx.WithKVGasConfig(storetypes.KVGasConfig()).
		WithKVGasConfig(storetypes.TransientGasConfig())

	msgSrv := stakingkeeper.NewMsgServerImpl(p.stakingKeeper)

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

	return []byte{}, nil
}

func (p *Precompile) Undelegate(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	argsBz []byte,
) ([]byte, error) {
	args, err := method.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, err
	}

	msg, err := NewMsgUndelegate(args)
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(p.stakingKeeper)

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

func (p *Precompile) Redelegate(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	argsBz []byte,
) ([]byte, error) {
	var redelegateInput RedelegateInput
	err := precompiles.UnpackIntoInterface(&redelegateInput, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	msg, err := redelegateInput.ToMessage()
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(p.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	res, err := msgSrv.BeginRedelegate(sdk.WrapSDKContext(cacheCtx), msg)
	if err != nil {
		return nil, err
	}

	output := new(RedelegateOutput).FromMessage(res)
	bz, err := output.Pack(method.Outputs)
	if err != nil {
		return nil, err
	}

	writeFn()

	return bz, nil
}

func (p *Precompile) CancelUnbondingDelegation(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	argsBz []byte,
) ([]byte, error) {
	var cancelUnbondingDelegationInput CancelUnbondingDelegationInput
	err := precompiles.UnpackIntoInterface(&cancelUnbondingDelegationInput, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	msg, err := cancelUnbondingDelegationInput.ToMessage()
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(p.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	if _, err := msgSrv.CancelUnbondingDelegation(sdk.WrapSDKContext(cacheCtx), msg); err != nil {
		return nil, err
	}

	writeFn()

	return nil, nil
}
