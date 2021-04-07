// +build simulator

package simulator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dfinance/dnode/cmd/config/genesis/defaults"
)

type SimProfile struct {
	ID          string
	SimDuration time.Duration
	//
	BlockTimeMin time.Duration
	BlockTimeMax time.Duration
	//
	MainTokensBalanceWODec    int64
	BondingTokensBalanceWODec int64
	LPTokensBalanceWODec      int64
	//
	Accounts           uint
	POAValidators      uint
	TMValidatorsTotal  uint
	TMValidatorsActive uint
	//
	OpCreateValidator time.Duration
	//
	OpDelegateBonding               time.Duration
	OpDelegateBondingAmountRatio    sdk.Dec
	OpDelegateBondingMaxSupplyRatio sdk.Dec
	//
	OpDelegateLP               time.Duration
	OpDelegateLPAmountRatio    sdk.Dec
	OpDelegateLPMaxSupplyRatio sdk.Dec
	//
	OpRedelegateBonding            time.Duration
	OpRedelegateBondingAmountRatio sdk.Dec
	//
	OpRedelegateLP            time.Duration
	OpRedelegateLPAmountRatio sdk.Dec
	//
	OpUndelegateBonding            time.Duration
	OpUndelegateBondingAmountRatio sdk.Dec
	//
	OpUndelegateLP            time.Duration
	OpUndelegateLPAmountRatio sdk.Dec
	//
	OpGetValidatorRewards time.Duration
	OpGetDelegatorRewards time.Duration
	//
	OpLockRewards      time.Duration
	OpLockRewardsRatio sdk.Dec
}

func (p SimProfile) String() string {
	str := strings.Builder{}
	str.WriteString("Simulation:\n")
	str.WriteString(fmt.Sprintf("  - ID: %s\n", p.ID))
	str.WriteString(fmt.Sprintf("  - SimDuration:  %v\n", p.SimDuration))
	str.WriteString(fmt.Sprintf("  - BlockTimeMin: %v\n", p.BlockTimeMin))
	str.WriteString(fmt.Sprintf("  - BlockTimeMax: %v\n", p.BlockTimeMax))
	str.WriteString("Initial balances:\n")
	str.WriteString(fmt.Sprintf("  - MainTokens:    %d.0%s\n", p.MainTokensBalanceWODec, defaults.MainDenom))
	str.WriteString(fmt.Sprintf("  - StakingTokens: %d.0%s\n", p.BondingTokensBalanceWODec, defaults.StakingDenom))
	str.WriteString(fmt.Sprintf("  - LPTokens:      %d.0%s\n", p.LPTokensBalanceWODec, defaults.LiquidityProviderDenom))
	str.WriteString("Total number of:\n")
	str.WriteString(fmt.Sprintf("  - Accounts:               %d\n", p.Accounts))
	str.WriteString(fmt.Sprintf("  - PoA validators:         %d\n", p.POAValidators))
	str.WriteString(fmt.Sprintf("  - TM validators (total):  %d\n", p.TMValidatorsTotal))
	str.WriteString(fmt.Sprintf("  - TM validators (active): %d\n", p.TMValidatorsActive))
	str.WriteString("Operations:\n")
	str.WriteString(fmt.Sprintf("  - Create validator every:                    %s\n", FormatDuration(p.OpCreateValidator)))
	str.WriteString(fmt.Sprintf("  - Delegate bonding tokens every:             %s\n", FormatDuration(p.OpDelegateBonding)))
	str.WriteString(fmt.Sprintf("      Delegate amount ratio (of acc balance):  %s\n", p.OpDelegateBondingAmountRatio))
	str.WriteString(fmt.Sprintf("      Max limit ratio (staked ratio):          %s\n", p.OpDelegateBondingMaxSupplyRatio))
	str.WriteString(fmt.Sprintf("  - Delegate LP tokens every:                  %s\n", FormatDuration(p.OpDelegateLP)))
	str.WriteString(fmt.Sprintf("      Delegate amount ratio (of acc balance):  %s\n", p.OpDelegateLPAmountRatio))
	str.WriteString(fmt.Sprintf("      Max limit ratio (staked ratio):          %s\n", p.OpDelegateLPMaxSupplyRatio))
	str.WriteString(fmt.Sprintf("  - Redelegate bonding tokens every:           %s\n", FormatDuration(p.OpRedelegateBonding)))
	str.WriteString(fmt.Sprintf("      Redelegate amount ratio (of del shares): %s\n", p.OpRedelegateBondingAmountRatio))
	str.WriteString(fmt.Sprintf("  - Redelegate LP tokens every:                %s\n", FormatDuration(p.OpRedelegateLP)))
	str.WriteString(fmt.Sprintf("      Redelegate amount ratio (of del shares): %s\n", p.OpRedelegateLPAmountRatio))
	str.WriteString(fmt.Sprintf("  - Undelegate bonding tokens every:           %s\n", FormatDuration(p.OpUndelegateBonding)))
	str.WriteString(fmt.Sprintf("      Undelegate amount ratio (of del shares): %s\n", p.OpUndelegateBondingAmountRatio))
	str.WriteString(fmt.Sprintf("  - Undelegate LP tokens every:                %s\n", FormatDuration(p.OpUndelegateLP)))
	str.WriteString(fmt.Sprintf("      Undelegate amount ratio (of del shares): %s\n", p.OpUndelegateLPAmountRatio))
	str.WriteString(fmt.Sprintf("  - Withdraw all validators comissions every:  %s\n", FormatDuration(p.OpGetValidatorRewards)))
	str.WriteString(fmt.Sprintf("  - Withdraw all delegators rewards every:     %s\n", FormatDuration(p.OpGetDelegatorRewards)))
	str.WriteString(fmt.Sprintf("  - Lock rewards every:                        %s\n", FormatDuration(p.OpLockRewards)))
	str.WriteString(fmt.Sprintf("      Ratio of all validators:                 %s\n", p.OpLockRewardsRatio))

	return str.String()
}

func simulate(t *testing.T, profile SimProfile) {
	go http.ListenAndServe(":8090", nil)

	t.Logf(profile.String())

	// create a tmp directory
	workingDir, err := ioutil.TempDir("/tmp", fmt.Sprintf("dnode-simulator-%s-", profile.ID))
	require.NoError(t, err)

	// genesis accounts balance
	amtDecimals := sdk.NewInt(1000000000000000000)
	genCoins := sdk.NewCoins(
		sdk.NewCoin(defaults.MainDenom, sdk.NewInt(profile.MainTokensBalanceWODec).Mul(amtDecimals)),
		sdk.NewCoin(defaults.StakingDenom, sdk.NewInt(profile.BondingTokensBalanceWODec).Mul(amtDecimals)),
		sdk.NewCoin(defaults.LiquidityProviderDenom, sdk.NewInt(profile.LPTokensBalanceWODec).Mul(amtDecimals)),
	)

	// write profile to file
	{
		f, err := os.Create(path.Join(workingDir, "profile.txt"))
		require.NoError(t, err)
		_, err = f.WriteString(profile.String())
		require.NoError(t, err)
		f.Close()
	}

	// CSV report writer
	reportWriter, writerClose := NewSimReportCSVWriter(t, path.Join(workingDir, "report.csv"))
	defer writerClose()

	// create simulator
	s := NewSimulator(t, workingDir, NewDefferOps(),
		//InMemoryDBOption(),
		BlockTimeOption(profile.BlockTimeMin, profile.BlockTimeMax),
		GenerateWalletAccountsOption(profile.Accounts, profile.POAValidators, genCoins),
		LogOption(log.AllowInfoWith("module", "x/staking")),
		LogOption(log.AllowInfoWith("module", "x/mint")),
		LogOption(log.AllowInfoWith("module", "x/distribution")),
		LogOption(log.AllowInfoWith("module", "x/slashing")),
		LogOption(log.AllowInfoWith("module", "x/evidence")),
		StakingParamsOption(func(state *staking.GenesisState) {
			state.Params.UnbondingTime = 24 * time.Hour
			state.Params.MaxValidators = uint16(profile.TMValidatorsActive)
			state.Params.MaxDelegationsRatio = sdk.NewDecWithPrec(1000, 0)
		}),
		InvariantCheckPeriodOption(1000),
		OperationsOption(
			NewSimInvariantsOp(1*time.Hour),
			NewForceUpdateOp(8*time.Hour),
			//
			NewReportOp(24*time.Hour, false, NewSimReportConsoleWriter(), reportWriter),
			//
			NewCreateValidatorOp(profile.OpCreateValidator, profile.TMValidatorsTotal),
			NewDelegateBondingOp(profile.OpDelegateBonding, profile.OpDelegateBondingAmountRatio, profile.OpDelegateBondingMaxSupplyRatio),
			NewDelegateLPOp(profile.OpDelegateLP, profile.OpDelegateLPAmountRatio, profile.OpDelegateLPMaxSupplyRatio),
			NewRedelegateBondingOp(profile.OpRedelegateBonding, profile.OpRedelegateBondingAmountRatio),
			NewRedelegateLPOp(profile.OpRedelegateLP, profile.OpRedelegateLPAmountRatio),
			NewUndelegateBondingOp(profile.OpUndelegateBonding, profile.OpUndelegateBondingAmountRatio),
			NewUndelegateLPOp(profile.OpUndelegateLP, profile.OpUndelegateLPAmountRatio),
			//
			NewGetValidatorRewardOp(profile.OpGetValidatorRewards),
			NewGetDelegatorRewardOp(profile.OpGetDelegatorRewards),
			NewLockValidatorRewardsOp(profile.OpLockRewards, profile.OpLockRewardsRatio),
		),
	)

	s.Start()

	// work loop
	_, simDur := s.SimulatedDur()
	for simDur < profile.SimDuration {
		s.Next()
		_, simDur = s.SimulatedDur()
	}

	t.Logf("Simulation is done, output dir: %s", s.workingDir)
}

func TestSimInflation(t *testing.T) {
	profile := SimProfile{
		ID:           "low_staking",
		SimDuration:  3 * Month,
		BlockTimeMin: 300 * time.Second,
		BlockTimeMax: 305 * time.Second,
		//
		MainTokensBalanceWODec:    50000,
		BondingTokensBalanceWODec: 500000,
		LPTokensBalanceWODec:      100000,
		//
		Accounts:           150,
		POAValidators:      3,
		TMValidatorsTotal:  110,
		TMValidatorsActive: 100,
		//
		OpCreateValidator: 3 * time.Hour,
		//
		OpDelegateBonding:               2 * time.Hour,
		OpDelegateBondingAmountRatio:    sdk.NewDecWithPrec(50, 2),
		OpDelegateBondingMaxSupplyRatio: sdk.NewDecWithPrec(30, 2),
		//
		OpDelegateLP:               4 * time.Hour,
		OpDelegateLPAmountRatio:    sdk.NewDecWithPrec(50, 2),
		OpDelegateLPMaxSupplyRatio: sdk.NewDecWithPrec(30, 2),
		//
		OpRedelegateBonding:            4 * Day,
		OpRedelegateBondingAmountRatio: sdk.NewDecWithPrec(30, 2),
		//
		OpRedelegateLP:            8 * Day,
		OpRedelegateLPAmountRatio: sdk.NewDecWithPrec(30, 2),
		//
		OpUndelegateBonding:            2 * Day,
		OpUndelegateBondingAmountRatio: sdk.NewDecWithPrec(15, 2),
		//
		OpUndelegateLP:            4 * Day,
		OpUndelegateLPAmountRatio: sdk.NewDecWithPrec(15, 2),
		//
		OpGetValidatorRewards: 1 * Month,
		OpGetDelegatorRewards: 1 * Month,
		//
		OpLockRewards:      4 * Day,
		OpLockRewardsRatio: sdk.NewDecWithPrec(30, 2),
	}

	simulate(t, profile)
}
