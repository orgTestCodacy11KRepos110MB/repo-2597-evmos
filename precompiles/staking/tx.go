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

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/evmos/ethermint/x/evm/statedb"
)

var (
	// DelegateMethod defines the ABI method signature for the staking Delegate
	// function.
	DelegateMethod abi.Method
	// UndelegateMethod defines the ABI method signature for the staking Undelegate
	// function.
	UndelegateMethod abi.Method

	// RedelegateMethod
	RedelegateMethod abi.Method
	// CancelUnbondingDelegationMethod
	CancelUnbondingDelegationMethod abi.Method
)

func (sp *StakingPrecompile) Delegate(ctx sdk.Context, argsBz []byte, stateDB statedb.ExtStateDB) ([]byte, error) {
	args, err := DelegateMethod.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, errors.New("fail to unpack input arguments")
	}

	denom := sp.stakingKeeper.BondDenom(ctx)

	msg, err := checkDelegateArgs(denom, args)
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(sp.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	if _, err := msgSrv.Delegate(sdk.WrapSDKContext(cacheCtx), msg); err != nil {
		return nil, err
	}

	writeFn()

	return nil, nil
}

func (sp *StakingPrecompile) Undelegate(ctx sdk.Context, argsBz []byte, stateDB statedb.ExtStateDB) ([]byte, error) {
	args, err := UndelegateMethod.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, errors.New("fail to unpack input arguments")
	}

	denom := sp.stakingKeeper.BondDenom(ctx)

	msg, err := checkUndelegateArgs(denom, args)
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(sp.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	res, err := msgSrv.Undelegate(sdk.WrapSDKContext(cacheCtx), msg)
	if err != nil {
		return nil, err
	}

	completionTimestamp := res.CompletionTime.UTC().Unix()
	bz, err := UndelegateMethod.Outputs.Pack(completionTimestamp)
	if err != nil {
		return nil, err
	}

	writeFn()

	return bz, nil
}

func (sp *StakingPrecompile) Redelegate(ctx sdk.Context, argsBz []byte, stateDB statedb.ExtStateDB) ([]byte, error) {
	args, err := RedelegateMethod.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, errors.New("fail to unpack input arguments")
	}

	denom := sp.stakingKeeper.BondDenom(ctx)

	msg, err := checkRedelegateArgs(denom, args)
	if err != nil {
		return nil, err
	}

	msgSrv := stakingkeeper.NewMsgServerImpl(sp.stakingKeeper)

	cacheCtx, writeFn := ctx.CacheContext()

	res, err := msgSrv.BeginRedelegate(sdk.WrapSDKContext(cacheCtx), msg)
	if err != nil {
		return nil, err
	}

	completionTimestamp := res.CompletionTime.UTC().Unix()
	bz, err := RedelegateMethod.Outputs.Pack(completionTimestamp)
	if err != nil {
		return nil, err
	}

	writeFn()

	return bz, nil
}

func (sp *StakingPrecompile) CancelUnbondingDelegation(ctx sdk.Context, argsBz []byte, stateDB statedb.ExtStateDB) ([]byte, error) {
	args, err := CancelUnbondingDelegationMethod.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, errors.New("fail to unpack input arguments")
	}

	denom := sp.stakingKeeper.BondDenom(ctx)

	msg, err := checkCancelUnbondingDelegationArgs(denom, args)
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
