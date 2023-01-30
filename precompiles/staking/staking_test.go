package staking_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/evmos/ethermint/tests"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/evmos/evmos/v11/app"
	"github.com/evmos/evmos/v11/precompiles/staking"
	claimstypes "github.com/evmos/evmos/v11/x/claims/types"
	inflationtypes "github.com/evmos/evmos/v11/x/inflation/types"
)

type PrecompileTestSuite struct {
	suite.Suite

	ctx        sdk.Context
	app        *app.Evmos
	address    common.Address
	validators []stakingtypes.Validator
	ethSigner  ethtypes.Signer
	signer     keyring.Signer
	bondDenom  string

	precompile *staking.Precompile
	contract   *vm.Contract
	stateDB    vm.StateDB
}

var s *PrecompileTestSuite

func TestPrecompileTestSuite(t *testing.T) {
	s = new(PrecompileTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (s *PrecompileTestSuite) SetupTest() {
	s.DoSetupTest(s.T())
}

func (s *PrecompileTestSuite) DoSetupTest(t require.TestingT) {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	s.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	s.signer = tests.NewSigner(priv)

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())

	// init app
	s.app = app.Setup(false, feemarkettypes.DefaultGenesisState())
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{
		Height:          1,
		ChainID:         "evmos_9001-1",
		Time:            time.Now().UTC(),
		ProposerAddress: consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	// bond denom
	params := claimstypes.DefaultParams()
	stakingParams := s.app.StakingKeeper.GetParams(s.ctx)
	stakingParams.BondDenom = params.GetClaimsDenom()
	s.bondDenom = stakingParams.BondDenom
	s.app.StakingKeeper.SetParams(s.ctx, stakingParams)

	err = s.app.BankKeeper.MintCoins(s.ctx, inflationtypes.ModuleName, sdk.Coins{{Denom: s.bondDenom, Amount: sdk.NewInt(1000000)}})
	require.NoError(t, err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, inflationtypes.ModuleName, s.address[:], sdk.Coins{{Denom: s.bondDenom, Amount: sdk.NewInt(10000)}})
	require.NoError(t, err)

	// Set Validator
	valAddr := sdk.ValAddress(s.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, privCons.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	validator = stakingkeeper.TestingUpdateValidator(s.app.StakingKeeper, s.ctx, validator, true)
	err = s.app.StakingKeeper.AfterValidatorCreated(s.ctx, validator.GetOperator())
	require.NoError(t, err)
	err = s.app.StakingKeeper.SetValidatorByConsAddr(s.ctx, validator)
	require.NoError(t, err)

	s.validators = []stakingtypes.Validator{validator}

	s.ethSigner = ethtypes.LatestSignerForChainID(s.app.EvmKeeper.ChainID())

	s.precompile, err = staking.NewPrecompile(s.app.StakingKeeper)
	require.NoError(t, err)
}

func (s *PrecompileTestSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			staking.DelegateMethod,
			s.precompile.Methods[staking.DelegateMethod].Name,
			true,
		},
		{
			staking.UndelegateMethod,
			s.precompile.Methods[staking.UndelegateMethod].Name,
			true,
		},
		{
			staking.RedelegateMethod,
			s.precompile.Methods[staking.RedelegateMethod].Name,
			true,
		},
		{
			staking.CancelUnbondingDelegationMethod,
			s.precompile.Methods[staking.CancelUnbondingDelegationMethod].Name,
			true,
		},
		{
			staking.DelegationMethod,
			s.precompile.Methods[staking.DelegationMethod].Name,
			false,
		},
		{
			"invalid",
			"invalid",
			false,
		},
	}

	for _, tc := range testCases {
		s.Require().Equal(s.precompile.IsTransaction(tc.method), tc.isTx)
	}
}
