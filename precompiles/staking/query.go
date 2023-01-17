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
	// DelegationMethod defines the ABI method signature for the staking Delegation
	// query function.
	DelegationMethod abi.Method
	// UnbondingDelegationMethod defines the ABI method signature for the staking
	// UnbondingDelegationMethod query function.
	UnbondingDelegationMethod abi.Method
	// Validator defines the ABI method signature for the staking
	// Validator query function.
	ValidatorMethod abi.Method
)

func (sp *StakingPrecompile) Delegation(ctx sdk.Context, argsBz []byte, stateDB statedb.ExtStateDB) ([]byte, error) {
	args, err := DelegationMethod.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, errors.New("fail to unpack input arguments")
	}

	req, err := checkDelegationArgs(args)
	if err != nil {
		return nil, err
	}

	queryServer := stakingkeeper.Querier{Keeper: sp.stakingKeeper}

	res, err := queryServer.Delegation(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	bz, err := DelegationMethod.Outputs.Pack(res)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func (sp *StakingPrecompile) UnbondingDelegation(ctx sdk.Context, argsBz []byte, stateDB statedb.ExtStateDB) ([]byte, error) {
	args, err := UnbondingDelegationMethod.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, errors.New("fail to unpack input arguments")
	}

	req, err := checkUnbondingDelegationArgs(args)
	if err != nil {
		return nil, err
	}

	queryServer := stakingkeeper.Querier{Keeper: sp.stakingKeeper}

	res, err := queryServer.UnbondingDelegation(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	bz, err := UnbondingDelegationMethod.Outputs.Pack(res)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func (sp *StakingPrecompile) Validator(ctx sdk.Context, argsBz []byte, stateDB statedb.ExtStateDB) ([]byte, error) {
	args, err := ValidatorMethod.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, errors.New("fail to unpack input arguments")
	}

	req, err := checkValidatorArgs(args)
	if err != nil {
		return nil, err
	}

	queryServer := stakingkeeper.Querier{Keeper: sp.stakingKeeper}

	res, err := queryServer.Validator(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	bz, err := ValidatorMethod.Outputs.Pack(res)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
