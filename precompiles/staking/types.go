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
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type DelegateInput struct {
	DelegatorAddress common.Address
	ValidatorAddress string
	Denom            string
	Amount           *big.Int
}

func (di DelegateInput) ToMessage() (*stakingtypes.MsgDelegate, error) {
	msg := &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(di.DelegatorAddress.Bytes()).String(),
		ValidatorAddress: di.ValidatorAddress,
		Amount: sdk.Coin{
			Denom:  di.Denom,
			Amount: sdk.NewIntFromBigInt(di.Amount),
		},
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

type UndelegateInput struct {
	DelegatorAddress common.Address
	ValidatorAddress string
	Denom            string
	Amount           *big.Int
}

func (ui UndelegateInput) ToMessage() (*stakingtypes.MsgUndelegate, error) {
	msg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: sdk.AccAddress(ui.DelegatorAddress.Bytes()).String(),
		ValidatorAddress: ui.ValidatorAddress,
		Amount: sdk.Coin{
			Denom:  ui.Denom,
			Amount: sdk.NewIntFromBigInt(ui.Amount),
		},
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

type UndelegateOutput struct {
	CompletionTime *big.Int
}

func (uo *UndelegateOutput) FromMessage(res *stakingtypes.MsgUndelegateResponse) *UndelegateOutput {
	uo.CompletionTime = big.NewInt(res.CompletionTime.UTC().Unix())
	return uo
}

func (uo UndelegateOutput) Pack(args abi.Arguments) ([]byte, error) {
	return args.Pack(uo.CompletionTime)
}

type RedelegateInput struct {
	DelegatorAddress    common.Address
	ValidatorSrcAddress string
	ValidatorDstAddress string
	Denom               string
	Amount              *big.Int
}

func (ri RedelegateInput) ToMessage() (*stakingtypes.MsgBeginRedelegate, error) {
	msg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    sdk.AccAddress(ri.DelegatorAddress.Bytes()).String(), // bech32 formatted
		ValidatorSrcAddress: ri.ValidatorSrcAddress,
		ValidatorDstAddress: ri.ValidatorDstAddress,
		Amount: sdk.Coin{
			Denom:  ri.Denom,
			Amount: sdk.NewIntFromBigInt(ri.Amount),
		},
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

type RedelegateOutput struct {
	CompletionTime *big.Int
}

func (ro *RedelegateOutput) FromMessage(res *stakingtypes.MsgBeginRedelegateResponse) *RedelegateOutput {
	ro.CompletionTime = big.NewInt(res.CompletionTime.UTC().Unix())
	return ro
}

func (uo RedelegateOutput) Pack(args abi.Arguments) ([]byte, error) {
	return args.Pack(uo.CompletionTime)
}

type CancelUnbondingDelegationInput struct {
	DelegatorAddress common.Address
	ValidatorAddress string
	Denom            string
	Amount           *big.Int
	CreationHeight   *big.Int
}

func (ci CancelUnbondingDelegationInput) ToMessage() (*stakingtypes.MsgCancelUnbondingDelegation, error) {
	msg := &stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: sdk.AccAddress(ci.DelegatorAddress.Bytes()).String(), // bech32 formatted
		ValidatorAddress: ci.ValidatorAddress,
		Amount: sdk.Coin{
			Denom:  ci.Denom,
			Amount: sdk.NewIntFromBigInt(ci.Amount),
		},
		CreationHeight: ci.Amount.Int64(),
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}
