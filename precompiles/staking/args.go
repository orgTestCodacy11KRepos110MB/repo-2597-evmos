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
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
)

func checkRedelegateArgs(denom string, args []interface{}) (*stakingtypes.MsgBeginRedelegate, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("invalid input arguments. Expected 4, got %d", len(args))
	}

	delegatorAddr, _ := args[0].(common.Address)
	validatorSrcAddr, _ := args[1].(string)
	validatorDstAddr, _ := args[2].(string)
	amount, ok := args[3].(*big.Int)
	if !ok || amount == nil {
		amount = big.NewInt(0)
	}

	coin := sdk.Coin{
		Denom:  denom,
		Amount: sdk.NewIntFromBigInt(amount),
	}

	delAddr := sdk.AccAddress(delegatorAddr.Bytes())

	msg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    delAddr.String(), // bech32 formatted
		ValidatorSrcAddress: validatorSrcAddr,
		ValidatorDstAddress: validatorDstAddr,
		Amount:              coin,
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

func checkCancelUnbondingDelegationArgs(denom string, args []interface{}) (*stakingtypes.MsgCancelUnbondingDelegation, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("invalid input arguments. Expected 4, got %d", len(args))
	}

	delegatorAddr, _ := args[0].(common.Address)
	validatorAddr, _ := args[1].(string)
	amount, ok := args[2].(*big.Int)
	if !ok || amount == nil {
		amount = big.NewInt(0)
	}

	creationHeight, _ := args[3].(int64)

	coin := sdk.Coin{
		Denom:  denom,
		Amount: sdk.NewIntFromBigInt(amount),
	}

	delAddr := sdk.AccAddress(delegatorAddr.Bytes())

	msg := &stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: delAddr.String(), // bech32 formatted
		ValidatorAddress: validatorAddr,
		Amount:           coin,
		CreationHeight:   creationHeight,
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

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
