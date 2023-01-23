// SPDX-License-Identifier: BUSL-only
pragma solidity >=0.8.17;

/// @dev The StakingI contract's address.
address constant STAKING_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000100;

/// @dev The StakingI contract's instance.
StakingI constant STAKING_CONTRACT = StakingI(
    STAKING_PRECOMPILE_ADDRESS
);

struct Validator {
  address operatorAddress;
  string consensusPubkey;
  bool jailed;
  BondStatus status;
  uint256 tokens;
  uint256 delegatorShares;
  string description;
  int64 unbondingHeight;
  uint64 unbondingTime;
  // TODO: Commision
  uint256 minSelfDelegation;
}

struct Coin {
  string denom;
  uint256 amount;
}


/// BondStatus is the status of the validator
enum BondStatus {
  Unspecified,
  Unbonded,
  Unbonding,
  Bonded
}



/// @author Evmos Team
/// @title Staking Precompiled Contract
/// @dev The interface through which solidity contracts will interact with  Staking
/// We follow this same interface including four-byte function selectors, in the precompile that
/// wraps the pallet
/// @custom:address 0x0000000000000000000000000000000000000100
interface StakingI {
    /// @dev Delegates the given amount of the bond denomination to a validator.
    /// @param delegatorAddress the address that we want to confirm is a delegator
    /// @param validatorAddress the address that we want to confirm is a delegator
    /// @param denom the address that we want to confirm is a delegator
    /// @param amount amount to be delegated to the validator
    function delegate(
      address delegatorAddress,
      string memory validatorAddress,
      string memory denom,
      uint256 amount
    ) external;

    /// @dev Undelegate the given amount of the bond denomination to a validator.
    /// @param delegatorAddress the address that we want to confirm is a delegator
    /// @param validatorAddress the address that we want to confirm is a delegator
    /// @param denom the address that we want to confirm is a delegator
    /// @param amount amount to be delegated to the validator
    function undelegate(
      address delegatorAddress,
      string memory validatorAddress,
      string memory denom,
      uint256 amount
    ) external returns (uint256 completionTime);

    /// @dev Redelegates the given amount of the bond denomination to a validator.
    /// @param delegatorAddress the address that we want to confirm is a delegator
    /// @param validatorSrcAddress the address that we want to confirm is a delegator
    /// @param validatorDstAddress the address that we want to confirm is a delegator
    /// @param denom the address that we want to confirm is a delegator
    /// @param amount amount to be delegated to the validator
    function redelegate(
      address delegatorAddress,
      string memory validatorSrcAddress,
      string memory validatorDstAddress,
      string memory denom,
      uint256 amount
    ) external returns (uint256 completionTime);

    /// @dev Delegates the given amount of the bond denomination to a validator.
    /// @param delegatorAddress the address that we want to confirm is a delegator
    /// @param validatorAddress the address that we want to confirm is a delegator
    /// @param denom the address that we want to confirm is a delegator
    /// @param amount amount to be delegated to the validator
    /// @param creationHeight amount to be delegated to the validator
    function cancelUnbondingDelegation(
      address delegatorAddress,
      string memory validatorAddress,
      string memory denom,
      uint256 amount,
      uint256 creationHeight
    ) external returns (uint256 completionTime);

    /// @dev Delegation the given amount of the bond denomination to a validator.
    /// @param delegatorAddress the address that we want to confirm is a delegator
    /// @param validatorAddress the address that we want to confirm is a delegator
    function delegation(
      address delegatorAddress,
      string memory validatorAddress
    ) external view returns (
      uint256 shares,
      string memory denom,
      uint256 amount
    );

    /// @dev Delegation the given amount of the bond denomination to a validator.
    /// @param delegatorAddress the address that we want to confirm is a delegator
    /// @param validatorAddress the address that we want to confirm is a delegator
    function unbondingDelegation(
      address delegatorAddress,
      string memory validatorAddress
    ) external view returns (
      uint256 shares,
      string memory denom,
      uint256 amount
    );

    /// @dev Delegation the given amount of the bond denomination to a validator.
    /// @param validatorAddress the address that we want to confirm is a delegator
    function validator(
      string memory validatorAddress
    ) external view returns (
      string memory operatorAddress,
      string memory consensusPubkey,
      bool jailed,
      BondStatus status,
      uint256 shares,
      string memory denom,
      uint256 amount
    );
}
