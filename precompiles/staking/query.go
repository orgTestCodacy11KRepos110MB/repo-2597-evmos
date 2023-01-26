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

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/evmos/ethermint/x/evm/statedb"
	"github.com/evmos/evmos/v11/precompiles"
)

const (
	// DelegationMethod defines the ABI method name for the staking Delegation
	// query.
	DelegationMethod = "delegation"
	// UnbondingDelegationMethod defines the ABI method name for the staking
	// UnbondingDelegationMethod query.
	UnbondingDelegationMethod = "unbonding"
	// ValidatorMethod defines the ABI method name for the staking
	// Validator query.
	ValidatorMethod = "validator"
	// ValidatorsMethod defines the ABI method name for the staking
	// Validators query.
	ValidatorsMethod = "validators"
)

func (sp *StakingPrecompile) Delegation(ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	method, ok := sp.ABI.Methods[DelegationMethod]
	if !ok {
		return nil, fmt.Errorf("no method with id: %s", DelegationMethod)
	}

	var delegationInput DelegationInput
	err := precompiles.UnpackIntoInterface(&delegationInput, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	req := delegationInput.ToRequest()

	queryServer := stakingkeeper.Querier{Keeper: sp.stakingKeeper}

	res, err := queryServer.Delegation(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	out := new(DelegationOutput).FromResponse(res)

	return out.Pack(method.Outputs)
}

func (sp *StakingPrecompile) UnbondingDelegation(ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	method, ok := sp.ABI.Methods[UnbondingDelegationMethod]
	if !ok {
		return nil, fmt.Errorf("no method with id: %s", DelegationMethod)
	}

	var input UnbondingDelegationInput
	err := precompiles.UnpackIntoInterface(&input, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	req := input.ToRequest()

	queryServer := stakingkeeper.Querier{Keeper: sp.stakingKeeper}

	res, err := queryServer.UnbondingDelegation(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	bz, err := method.Outputs.Pack(res)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func (sp *StakingPrecompile) Validator(ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	method, ok := sp.ABI.Methods[ValidatorMethod]
	if !ok {
		return nil, fmt.Errorf("no method with id: %s", DelegationMethod)
	}

	var input ValidatorInput
	err := precompiles.UnpackIntoInterface(&input, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	req := input.ToRequest()

	queryServer := stakingkeeper.Querier{Keeper: sp.stakingKeeper}

	res, err := queryServer.Validator(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	out := new(ValidatorOutput).FromResponse(res)

	return out.Pack(method.Outputs)
}

func (sp *StakingPrecompile) Validators(ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	method, ok := sp.ABI.Methods[ValidatorsMethod]
	if !ok {
		return nil, fmt.Errorf("no method with id: %s", DelegationMethod)
	}

	var input ValidatorsInput
	err := precompiles.UnpackIntoInterface(&input, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	req := input.ToRequest()

	queryServer := stakingkeeper.Querier{Keeper: sp.stakingKeeper}

	res, err := queryServer.Validators(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	out := new(ValidatorsOutput).FromResponse(res)

	return out.Pack(method.Outputs)
}

func (sp *StakingPrecompile) Redelegations(ctx sdk.Context,
	contract *vm.Contract,
	argsBz []byte,
	stateDB *statedb.StateDB,
	readOnly bool,
) ([]byte, error) {
	method, ok := sp.ABI.Methods[RedelegationsMethod]
	if !ok {
		return nil, fmt.Errorf("no method with id: %s", DelegationMethod)
	}

	var input RedelegationsInput
	err := precompiles.UnpackIntoInterface(&input, method.Inputs, argsBz)
	if err != nil {
		return nil, err
	}

	req := input.ToRequest()

	queryServer := stakingkeeper.Querier{Keeper: sp.stakingKeeper}

	res, err := queryServer.Redelegations(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, err
	}

	out := new(RedelegationsOutput).FromResponse(res)

	return out.Pack(method.Outputs)
}
