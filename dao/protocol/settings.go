package protocol

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Protocol DAO settings
type ProtocolDaoSettings struct {
	Auction struct {
		IsCreateLotEnabled    *ProtocolDaoBoolSetting              `json:"isCreateLotEnabled"`
		IsBidOnLotEnabled     *ProtocolDaoBoolSetting              `json:"isBidOnLotEnabled"`
		LotMinimumEthValue    *ProtocolDaoUintSetting              `json:"lotMinimumEthValue"`
		LotMaximumEthValue    *ProtocolDaoUintSetting              `json:"lotMaximumEthValue"`
		LotDuration           *ProtocolDaoCompoundSetting[uint64]  `json:"lotDuration"`
		LotStartingPriceRatio *ProtocolDaoCompoundSetting[float64] `json:"lotStartingPriceRatio"`
		LotReservePriceRatio  *ProtocolDaoCompoundSetting[float64] `json:"lotReservePriceRatio"`
	} `json:"auction"`

	Deposit struct {
		IsDepositingEnabled                    *ProtocolDaoBoolSetting              `json:"isDepositingEnabled"`
		AreDepositAssignmentsEnabled           *ProtocolDaoBoolSetting              `json:"areDepositAssignmentsEnabled"`
		MinimumDeposit                         *ProtocolDaoUintSetting              `json:"minimumDeposit"`
		MaximumDepositPoolSize                 *ProtocolDaoUintSetting              `json:"maximumDepositPoolSize"`
		MaximumAssignmentsPerDeposit           *ProtocolDaoCompoundSetting[uint64]  `json:"maximumAssignmentsPerDeposit"`
		MaximumSocialisedAssignmentsPerDeposit *ProtocolDaoCompoundSetting[uint64]  `json:"maximumSocialisedAssignmentsPerDeposit"`
		DepositFee                             *ProtocolDaoCompoundSetting[float64] `json:"depositFee"`
	} `json:"deposit"`

	Inflation struct {
		IntervalRate *ProtocolDaoCompoundSetting[float64]   `json:"intervalRate"`
		StartTime    *ProtocolDaoCompoundSetting[time.Time] `json:"startTime"`
	} `json:"inflation"`

	Minipool struct {
		IsSubmitWithdrawableEnabled *ProtocolDaoBoolSetting                    `json:"isSubmitWithdrawableEnabled"`
		LaunchTimeout               *ProtocolDaoCompoundSetting[time.Duration] `json:"launchTimeout"`
		IsBondReductionEnabled      *ProtocolDaoBoolSetting                    `json:"isBondReductionEnabled"`
		MaximumCount                *ProtocolDaoCompoundSetting[uint64]        `json:"maximumCount"`
		UserDistributeWindowStart   *ProtocolDaoCompoundSetting[time.Duration] `json:"userDistributeWindowStart"`
		UserDistributeWindowLength  *ProtocolDaoCompoundSetting[time.Duration] `json:"userDistributeWindowLength"`
	} `json:"minipool"`

	Network struct {
		OracleDaoConsensusThreshold *ProtocolDaoCompoundSetting[float64] `json:"oracleDaoConsensusThreshold"`
		NodePenaltyThreshold        *ProtocolDaoCompoundSetting[float64] `json:"nodePenaltyThreshold"`
		PerPenaltyRate              *ProtocolDaoCompoundSetting[float64] `json:"perPenaltyRate"`
		IsSubmitBalancesEnabled     *ProtocolDaoBoolSetting              `json:"isSubmitBalancesEnabled"`
		SubmitBalancesFrequency     *ProtocolDaoCompoundSetting[uint64]  `json:"submitBalancesFrequency"`
		IsSubmitPricesEnabled       *ProtocolDaoBoolSetting              `json:"isSubmitPricesEnabled"`
		SubmitPricesFrequency       *ProtocolDaoCompoundSetting[uint64]  `json:"submitPricesFrequency"`
		MinimumNodeFee              *ProtocolDaoCompoundSetting[float64] `json:"minimumNodeFee"`
		TargetNodeFee               *ProtocolDaoCompoundSetting[float64] `json:"targetNodeFee"`
		MaximumNodeFee              *ProtocolDaoCompoundSetting[float64] `json:"maximumNodeFee"`
		NodeFeeDemandRange          *ProtocolDaoUintSetting              `json:"nodeFeeDemandRange"`
		TargetRethCollateralRate    *ProtocolDaoCompoundSetting[float64] `json:"targetRethCollateralRate"`
		NetworkPenaltyThreshold     *ProtocolDaoCompoundSetting[float64] `json:"networkPenaltyThreshold"`
		NetworkPenaltyPerRate       *ProtocolDaoCompoundSetting[float64] `json:"networkPenaltyPerRate"`
		IsSubmitRewardsEnabled      *ProtocolDaoBoolSetting              `json:"isSubmitRewardsEnabled"`
	} `json:"network"`

	Proposals struct {
		VoteTime            *ProtocolDaoCompoundSetting[time.Duration] `json:"voteTime"`
		VoteDelayTime       *ProtocolDaoCompoundSetting[time.Duration] `json:"voteDelayTime"`
		ExecuteTime         *ProtocolDaoCompoundSetting[time.Duration] `json:"executeTime"`
		ProposalBond        *ProtocolDaoUintSetting                    `json:"proposalBond"`
		ChallengeBond       *ProtocolDaoUintSetting                    `json:"challengeBond"`
		ChallengePeriod     *ProtocolDaoCompoundSetting[time.Duration] `json:"challengePeriod"`
		ProposalQuorum      *ProtocolDaoCompoundSetting[float64]       `json:"proposalQuorum"`
		ProposalVetoQuorum  *ProtocolDaoCompoundSetting[float64]       `json:"proposalVetoQuorum"`
		PropoaslMaxBlockAge *ProtocolDaoCompoundSetting[uint64]        `json:"propoaslMaxBlockAge"`
	} `json:"proposals"`

	Node struct {
		IsRegistrationEnabled              *ProtocolDaoBoolSetting              `json:"isRegistrationEnabled"`
		IsSmoothingPoolRegistrationEnabled *ProtocolDaoBoolSetting              `json:"isSmoothingPoolRegistrationEnabled"`
		IsDepositingEnabled                *ProtocolDaoBoolSetting              `json:"isDepositingEnabled"`
		AreVacantMinipoolsEnabled          *ProtocolDaoBoolSetting              `json:"areVacantMinipoolsEnabled"`
		MinimumPerMinipoolStake            *ProtocolDaoCompoundSetting[float64] `json:"minimumPerMinipoolStake"`
		MaximumPerMinipoolStake            *ProtocolDaoCompoundSetting[float64] `json:"maximumPerMinipoolStake"`
	} `json:"node"`

	Rewards struct {
		IntervalTime *ProtocolDaoCompoundSetting[time.Duration] `json:"intervalTime"`
	} `json:"rewards"`

	// === Internal fields ===
	rp            *rocketpool.RocketPool
	pdaoMgr       *ProtocolDaoManager
	dps_auction   *core.Contract
	dps_deposit   *core.Contract
	dps_inflation *core.Contract
	dps_minipool  *core.Contract
	dps_network   *core.Contract
	dps_node      *core.Contract
	dps_proposals *core.Contract
	dps_rewards   *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProtocolDaoSettings binding
func newProtocolDaoSettings(pdaoMgr *ProtocolDaoManager) (*ProtocolDaoSettings, error) {
	// Get the contracts
	contracts, err := pdaoMgr.rp.GetContracts([]rocketpool.ContractName{
		rocketpool.ContractName_RocketDAOProtocolSettingsAuction,
		rocketpool.ContractName_RocketDAOProtocolSettingsDeposit,
		rocketpool.ContractName_RocketDAOProtocolSettingsInflation,
		rocketpool.ContractName_RocketDAOProtocolSettingsMinipool,
		rocketpool.ContractName_RocketDAOProtocolSettingsNetwork,
		rocketpool.ContractName_RocketDAOProtocolSettingsNode,
		rocketpool.ContractName_RocketDAOProtocolSettingsProposals,
		rocketpool.ContractName_RocketDAOProtocolSettingsRewards,
	}...)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO settings contracts: %w", err)
	}

	s := &ProtocolDaoSettings{
		rp:      pdaoMgr.rp,
		pdaoMgr: pdaoMgr,

		dps_auction:   contracts[0],
		dps_deposit:   contracts[1],
		dps_inflation: contracts[2],
		dps_minipool:  contracts[3],
		dps_network:   contracts[4],
		dps_node:      contracts[5],
		dps_proposals: contracts[6],
		dps_rewards:   contracts[7],
	}

	// Auction
	s.Auction.IsCreateLotEnabled = newBoolSetting(s.dps_auction, pdaoMgr, "auction.lot.create.enabled")
	s.Auction.IsBidOnLotEnabled = newBoolSetting(s.dps_auction, pdaoMgr, "auction.lot.bidding.enabled")
	s.Auction.LotMinimumEthValue = newUintSetting(s.dps_auction, pdaoMgr, "auction.lot.value.minimum")
	s.Auction.LotMaximumEthValue = newUintSetting(s.dps_auction, pdaoMgr, "auction.lot.value.maximum")
	s.Auction.LotDuration = newCompoundSetting[uint64](s.dps_auction, pdaoMgr, "auction.lot.duration")
	s.Auction.LotStartingPriceRatio = newCompoundSetting[float64](s.dps_auction, pdaoMgr, "auction.price.start")
	s.Auction.LotReservePriceRatio = newCompoundSetting[float64](s.dps_auction, pdaoMgr, "auction.price.reserve")

	// Deposit
	s.Deposit.IsDepositingEnabled = newBoolSetting(s.dps_deposit, pdaoMgr, "deposit.enabled")
	s.Deposit.AreDepositAssignmentsEnabled = newBoolSetting(s.dps_deposit, pdaoMgr, "deposit.assign.enabled")
	s.Deposit.MinimumDeposit = newUintSetting(s.dps_deposit, pdaoMgr, "deposit.minimum")
	s.Deposit.MaximumDepositPoolSize = newUintSetting(s.dps_deposit, pdaoMgr, "deposit.pool.maximum")
	s.Deposit.MaximumAssignmentsPerDeposit = newCompoundSetting[uint64](s.dps_deposit, pdaoMgr, "deposit.assign.maximum")
	s.Deposit.MaximumSocialisedAssignmentsPerDeposit = newCompoundSetting[uint64](s.dps_deposit, pdaoMgr, "deposit.assign.socialised.maximum")
	s.Deposit.DepositFee = newCompoundSetting[float64](s.dps_deposit, pdaoMgr, "deposit.fee")

	// Inflation
	s.Inflation.IntervalRate = newCompoundSetting[float64](s.dps_inflation, pdaoMgr, "rpl.inflation.interval.rate")
	s.Inflation.StartTime = newCompoundSetting[time.Time](s.dps_inflation, pdaoMgr, "rpl.inflation.interval.start")

	// Minipool
	s.Minipool.IsSubmitWithdrawableEnabled = newBoolSetting(s.dps_minipool, pdaoMgr, "minipool.submit.withdrawable.enabled")
	s.Minipool.LaunchTimeout = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, "minipool.launch.timeout")
	s.Minipool.IsBondReductionEnabled = newBoolSetting(s.dps_minipool, pdaoMgr, "minipool.bond.reduction.enabled")
	s.Minipool.MaximumCount = newCompoundSetting[uint64](s.dps_minipool, pdaoMgr, "minipool.maximum.count")
	s.Minipool.UserDistributeWindowStart = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, "minipool.user.distribute.window.start")
	s.Minipool.UserDistributeWindowLength = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, "minipool.user.distribute.window.length")

	// Network
	s.Network.OracleDaoConsensusThreshold = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.consensus.threshold")
	s.Network.NodePenaltyThreshold = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.penalty.threshold")
	s.Network.PerPenaltyRate = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.penalty.per.rate")
	s.Network.IsSubmitBalancesEnabled = newBoolSetting(s.dps_network, pdaoMgr, "network.submit.balances.enabled")
	s.Network.SubmitBalancesFrequency = newCompoundSetting[uint64](s.dps_network, pdaoMgr, "network.submit.balances.frequency")
	s.Network.IsSubmitPricesEnabled = newBoolSetting(s.dps_network, pdaoMgr, "network.submit.prices.enabled")
	s.Network.SubmitPricesFrequency = newCompoundSetting[uint64](s.dps_network, pdaoMgr, "network.submit.prices.frequency")
	s.Network.MinimumNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.node.fee.minimum")
	s.Network.TargetNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.node.fee.target")
	s.Network.MaximumNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.node.fee.maximum")
	s.Network.NodeFeeDemandRange = newUintSetting(s.dps_network, pdaoMgr, "network.node.fee.demand.range")
	s.Network.TargetRethCollateralRate = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.reth.collateral.target")
	s.Network.NetworkPenaltyThreshold = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.penalty.threshold")
	s.Network.NetworkPenaltyPerRate = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.penalty.per.rate")
	s.Network.IsSubmitRewardsEnabled = newBoolSetting(s.dps_network, pdaoMgr, "network.submit.rewards.enabled")

	// Node
	s.Node.IsRegistrationEnabled = newBoolSetting(s.dps_node, pdaoMgr, "node.registration.enabled")
	s.Node.IsSmoothingPoolRegistrationEnabled = newBoolSetting(s.dps_node, pdaoMgr, "node.smoothing.pool.registration.enabled")
	s.Node.IsDepositingEnabled = newBoolSetting(s.dps_node, pdaoMgr, "node.deposit.enabled")
	s.Node.AreVacantMinipoolsEnabled = newBoolSetting(s.dps_node, pdaoMgr, "node.vacant.minipools.enabled")
	s.Node.MinimumPerMinipoolStake = newCompoundSetting[float64](s.dps_node, pdaoMgr, "node.per.minipool.stake.minimum")
	s.Node.MaximumPerMinipoolStake = newCompoundSetting[float64](s.dps_node, pdaoMgr, "node.per.minipool.stake.maximum")

	// Proposals
	s.Proposals.VoteTime = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, "proposal.vote.time")
	s.Proposals.VoteDelayTime = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, "proposal.vote.delay.time")
	s.Proposals.ExecuteTime = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, "proposal.execute.time")
	s.Proposals.ProposalBond = newUintSetting(s.dps_proposals, pdaoMgr, "proposal.bond")
	s.Proposals.ChallengeBond = newUintSetting(s.dps_proposals, pdaoMgr, "proposal.challenge.bond")
	s.Proposals.ChallengePeriod = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, "proposal.challenge.period")
	s.Proposals.ProposalQuorum = newCompoundSetting[float64](s.dps_proposals, pdaoMgr, "proposal.quorum")
	s.Proposals.ProposalVetoQuorum = newCompoundSetting[float64](s.dps_proposals, pdaoMgr, "proposal.veto.quorum")
	s.Proposals.PropoaslMaxBlockAge = newCompoundSetting[uint64](s.dps_proposals, pdaoMgr, "proposal.max.block.age")

	// Rewards
	s.Rewards.IntervalTime = newCompoundSetting[time.Duration](s.dps_rewards, pdaoMgr, "rpl.rewards.claim.period.time")

	return s, nil
}

// =============
// === Calls ===
// =============

// Get all of the settings, organized by the type used in proposals and boostraps
func (c *ProtocolDaoSettings) GetSettings() ([]IProtocolDaoSetting[bool], []IProtocolDaoSetting[*big.Int]) {
	boolSettings := []IProtocolDaoSetting[bool]{}
	uintSettings := []IProtocolDaoSetting[*big.Int]{}

	settingsType := reflect.TypeOf(c)
	settingsVal := reflect.ValueOf(c)
	fieldCount := settingsType.NumField()
	for i := 0; i < fieldCount; i++ {
		categoryFieldType := settingsType.Field(i).Type

		// A container struct for settings by category
		if categoryFieldType.Kind() == reflect.Struct {
			// Get all of the settings in this cateogry
			categoryFieldVal := settingsVal.Field(i)
			settingCount := categoryFieldType.NumField()
			for j := 0; j < settingCount; j++ {
				setting := categoryFieldVal.Field(i).Interface()

				// Try bool settings
				boolSetting, isBoolSetting := setting.(IProtocolDaoSetting[bool])
				if isBoolSetting {
					boolSettings = append(boolSettings, boolSetting)
					continue
				}

				// Try uint settings
				uintSetting, isUintSetting := setting.(IProtocolDaoSetting[*big.Int])
				if isUintSetting {
					uintSettings = append(uintSettings, uintSetting)
				}
			}

		}
	}

	return boolSettings, uintSettings
}