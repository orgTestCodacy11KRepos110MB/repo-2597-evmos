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

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func checkDelegationArgs(args []interface{}) (*stakingtypes.QueryDelegationRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid input arguments. Expected 3, got %d", len(args))
	}

	delegatorAddr, _ := args[0].(string)
	validatorAddr, _ := args[1].(string)

	req := &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}

	return req, nil
}

func checkUnbondingDelegationArgs(args []interface{}) (*stakingtypes.QueryUnbondingDelegationRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid input arguments. Expected 3, got %d", len(args))
	}

	delegatorAddr, _ := args[0].(string)
	validatorAddr, _ := args[1].(string)

	req := &stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}

	return req, nil
}

func checkValidatorArgs(args []interface{}) (*stakingtypes.QueryValidatorRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid input arguments. Expected 3, got %d", len(args))
	}

	validatorAddr, _ := args[1].(string)

	req := &stakingtypes.QueryValidatorRequest{
		ValidatorAddr: validatorAddr,
	}

	return req, nil
}
