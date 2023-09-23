package minipool

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils"
)

const (
	eventScanInterval uint64 = 10000
)

// ===============
// === Structs ===
// ===============

// Basic binding for version-agnostic RocketMinipool contracts
type minipoolCommon struct {
	*MinipoolCommonDetails
	contract *core.Contract
	rp       *rocketpool.RocketPool
	mpMgr    *core.Contract
	mpQueue  *core.Contract
	mpStatus *core.Contract
}

// Basic details about a minipool, version-agnostic
type MinipoolCommonDetails struct {
	// Core parameters
	Address                    common.Address                            `json:"address"`
	Version                    uint8                                     `json:"version"`
	NodeAddress                common.Address                            `json:"nodeAddress"`
	Status                     core.Uint8Parameter[types.MinipoolStatus] `json:"status"`
	StatusBlock                core.Uint256Parameter[uint64]             `json:"statusBlock"`
	StatusTime                 core.Uint256Parameter[time.Time]          `json:"statusTime"`
	IsFinalised                bool                                      `json:"isFinalized"`
	NodeFee                    core.Uint256Parameter[float64]            `json:"nodeFee"`
	NodeDepositBalance         *big.Int                                  `json:"nodeDepositBalance"`
	NodeRefundBalance          *big.Int                                  `json:"nodeRefundBalance"`
	NodeDepositAssigned        bool                                      `json:"nodeDepositAssigned"`
	UserDepositBalance         *big.Int                                  `json:"userDepositBalance"`
	UserDepositAssigned        bool                                      `json:"userDepositAssigned"`
	UserDepositAssignedTime    core.Uint256Parameter[time.Time]          `json:"userDepositAssignedTime"`
	IsUseLatestDelegateEnabled bool                                      `json:"IsUseLatestDelegateEnabled"`
	DelegateAddress            common.Address                            `json:"delegateAddress"`
	PreviousDelegateAddress    common.Address                            `json:"previousDelegateAddress"`
	EffectiveDelegateAddress   common.Address                            `json:"effectiveDelegateAddress"`
	PenaltyCount               core.Uint256Parameter[uint64]             `json:"penaltyCount"`

	// MinipoolManager
	Exists                bool                                       `json:"exists"`
	Pubkey                types.ValidatorPubkey                      `json:"pubkey"`
	WithdrawalCredentials common.Hash                                `json:"withdrawalCredentials"`
	RplSlashed            bool                                       `json:"rplSlashed"`
	DepositType           core.Uint8Parameter[types.MinipoolDeposit] `json:"depositType"`

	// BondReducer
	IsBondReduceCancelled        bool                             `json:"isBondReduceCancelled"`
	ReduceBondTime               core.Uint256Parameter[time.Time] `json:"reduceBondTime"`
	ReduceBondValue              *big.Int                         `json:"reduceBondValue"`
	LastBondReductionTime        core.Uint256Parameter[time.Time] `json:"lastBondReductionTime"`
	LastBondReductionPrevValue   *big.Int                         `json:"lastBondReductionPrevValue"`
	LastBondReductionPrevNodeFee core.Uint256Parameter[float64]   `json:"lastBondReductionPrevNodeFee"`

	// MinipoolQueue
	QueuePosition core.Uint256Parameter[int64] `json:"queuePosition"`
}

// The data from a minipool's MinipoolPrestaked event
type MinipoolPrestakeEvent struct {
	Pubkey                []byte   `abi:"validatorPubkey"`
	Signature             []byte   `abi:"validatorSignature"`
	DepositDataRoot       [32]byte `abi:"depositDataRoot"`
	Amount                *big.Int `abi:"amount"`
	WithdrawalCredentials []byte   `abi:"withdrawalCredentials"`
	Time                  *big.Int `abi:"time"`
}

// Formatted MinipoolPrestaked event data
type PrestakeData struct {
	Pubkey                types.ValidatorPubkey    `json:"pubkey"`
	WithdrawalCredentials common.Hash              `json:"withdrawalCredentials"`
	Amount                *big.Int                 `json:"amount"`
	Signature             types.ValidatorSignature `json:"signature"`
	DepositDataRoot       common.Hash              `json:"depositDataRoot"`
	Time                  time.Time                `json:"time"`
}

// ====================
// === Constructors ===
// ====================

// Create a minipool common binding from an explicit version number
func newMinipoolCommonFromVersion(rp *rocketpool.RocketPool, contract *core.Contract, version uint8) (*minipoolCommon, error) {
	mpMgr, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolManager)
	if err != nil {
		return nil, fmt.Errorf("error creating minipool manager: %w", err)
	}

	mpQueue, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolQueue)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool queue contract: %w", err)
	}

	mpStatus, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolStatus)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool status contract: %w", err)
	}

	return &minipoolCommon{
		MinipoolCommonDetails: &MinipoolCommonDetails{
			Address: *contract.Address,
			Version: version,
		},
		rp:       rp,
		contract: contract,
		mpMgr:    mpMgr,
		mpQueue:  mpQueue,
		mpStatus: mpStatus,
	}, nil
}

// =============
// === Calls ===
// =============

// Gets the underlying minipool's contract
func (c *minipoolCommon) GetContract() *core.Contract {
	return c.contract
}

// Gets the common details for all minipool types
func (c *minipoolCommon) GetCommonDetails() *MinipoolCommonDetails {
	return c.MinipoolCommonDetails
}

// === Minipool ===

// Get the minipool's penalty count
func (c *minipoolCommon) GetPenaltyCount(mc *batch.MultiCaller) {
	// This isn't in the manager, it's in RocketStorage
	key := crypto.Keccak256Hash([]byte("network.penalties.penalty"), c.Address.Bytes())
	c.rp.Storage.GetUint(mc, &c.PenaltyCount.RawValue, key)
}

// Get the minipool's status
func (c *minipoolCommon) GetStatus(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Status.RawValue, "getStatus")
}

// Get the block that the minipool's status last changed
func (c *minipoolCommon) GetStatusBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.StatusBlock.RawValue, "getStatusBlock")
}

// Get the time that the minipool's status last changed
func (c *minipoolCommon) GetStatusTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.StatusTime.RawValue, "getStatusTime")
}

// Check if the minipool has been finalised
func (c *minipoolCommon) GetFinalised(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.IsFinalised, "getFinalised")
}

// Get the minipool's node address
func (c *minipoolCommon) GetNodeAddress(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.NodeAddress, "getNodeAddress")
}

// Get the minipool's commission rate
func (c *minipoolCommon) GetNodeFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.NodeFee.RawValue, "getNodeFee")
}

// Get the balance the node has deposited to the minipool
func (c *minipoolCommon) GetNodeDepositBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.NodeDepositBalance, "getNodeDepositBalance")
}

// Get the amount of ETH ready to be refunded to the node
func (c *minipoolCommon) GetNodeRefundBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.NodeRefundBalance, "getNodeRefundBalance")
}

// Check if the node deposit has been assigned to the minipool
func (c *minipoolCommon) GetNodeDepositAssigned(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.NodeDepositAssigned, "getNodeDepositAssigned")
}

// Get the balance the pool stakers have deposited to the minipool
func (c *minipoolCommon) GetUserDepositBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.UserDepositBalance, "getUserDepositBalance")
}

// Check if the pool staker deposits has been assigned to the minipool
func (c *minipoolCommon) GetUserDepositAssigned(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.UserDepositAssigned, "getUserDepositAssigned")
}

// Get the time at which the pool stakers were assigned to the minipool
func (c *minipoolCommon) GetUserDepositAssignedTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.UserDepositAssignedTime.RawValue, "getUserDepositAssignedTime")
}

// Check if the "use latest delegate" flag is enabled
func (c *minipoolCommon) GetUseLatestDelegate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.IsUseLatestDelegateEnabled, "getUseLatestDelegate")
}

// Get the address of the current delegate the minipool has recorded
func (c *minipoolCommon) GetDelegate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.DelegateAddress, "getDelegate")
}

// Get the address of the previous delegate the minipool will use after a rollback
func (c *minipoolCommon) GetPreviousDelegate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.PreviousDelegateAddress, "getPreviousDelegate")
}

// Get the address of the delegate the minipool will use (may be different than DelegateAddress if UseLatestDelegate is enabled)
func (c *minipoolCommon) GetEffectiveDelegate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.EffectiveDelegateAddress, "getEffectiveDelegate")
}

// === MinipoolManager ===

// Check if a minipool exists
func (c *minipoolCommon) GetExists(mc *batch.MultiCaller) {
	// TODO: Is this really necessary?
	core.AddCall(mc, c.mpMgr, &c.Exists, "getMinipoolExists", c.Address)
}

// Get the minipool's pubkey
func (c *minipoolCommon) GetPubkey(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.Pubkey, "getMinipoolPubkey", c.Address)
}

// Get the minipool's 0x01-based withdrawal credentials
func (c *minipoolCommon) GetWithdrawalCredentials(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.WithdrawalCredentials, "getMinipoolWithdrawalCredentials", c.Address)
}

// Check if the minipool's RPL has been slashed
func (c *minipoolCommon) GetRplSlashed(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.WithdrawalCredentials, "getMinipoolRPLSlashed", c.Address)
}

// Get the minipool's deposit type
func (c *minipoolCommon) GetDepositType(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.DepositType.RawValue, "getMinipoolDepositType", c.Address)
}

// === MinipoolQueue ===

// Get queue position of the minipool (-1 means not in the queue, otherwise 0-indexed).
func (c *minipoolCommon) GetQueuePosition(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpQueue, &c.QueuePosition.RawValue, "getMinipoolPosition", c.Address)
}

// Get the basic details
func (c *minipoolCommon) QueryAllDetails(mc *batch.MultiCaller) {
	c.GetPenaltyCount(mc)
	c.GetStatus(mc)
	c.GetStatusBlock(mc)
	c.GetStatusTime(mc)
	c.GetFinalised(mc)
	c.GetNodeAddress(mc)
	c.GetNodeFee(mc)
	c.GetNodeDepositBalance(mc)
	c.GetNodeRefundBalance(mc)
	c.GetNodeDepositAssigned(mc)
	c.GetUserDepositBalance(mc)
	c.GetUserDepositAssigned(mc)
	c.GetUserDepositAssignedTime(mc)
	c.GetUseLatestDelegate(mc)
	c.GetDelegate(mc)
	c.GetPreviousDelegate(mc)
	c.GetEffectiveDelegate(mc)
	c.GetExists(mc)
	c.GetPubkey(mc)
	c.GetWithdrawalCredentials(mc)
	c.GetRplSlashed(mc)
	c.GetDepositType(mc)
}

// ====================
// === Transactions ===
// ====================

// === Minipool ===

// Get info for refunding node ETH from the minipool
func (c *minipoolCommon) Refund(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "refund", opts)
}

// Get info for progressing the prelaunch minipool to staking
func (c *minipoolCommon) Stake(validatorSignature types.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "stake", opts, validatorSignature[:], depositDataRoot)
}

// Get info for dissolving the initialized or prelaunch minipool
func (c *minipoolCommon) Dissolve(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "dissolve", opts)
}

// Get info for withdrawing node balances from the dissolved minipool and closing it
func (c *minipoolCommon) Close(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "close", opts)
}

// Get info for finalising a minipool to get the RPL stake back
func (c *minipoolCommon) Finalise(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "finalise", opts)
}

// Get info for upgrading this minipool to the latest network delegate contract
func (c *minipoolCommon) DelegateUpgrade(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "delegateUpgrade", opts)
}

// Get info for rolling back to the previous delegate contract
func (c *minipoolCommon) DelegateRollback(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "delegateRollback", opts)
}

// Get info for setting the UseLatestDelegate flag (if set to true, will automatically use the latest delegate contract)
func (c *minipoolCommon) SetUseLatestDelegate(setting bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "setUseLatestDelegate", opts, setting)
}

// Get info for voting to scrub a minipool
func (c *minipoolCommon) VoteScrub(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "voteScrub", opts)
}

// === MinipoolStatus ===

// Get info for submitting a minipool withdrawable event
func (c *minipoolCommon) SubmitMinipoolWithdrawable(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mpStatus, "submitMinipoolWithdrawable", opts, c.Address)
}

// =============
// === Utils ===
// =============

// Given a validator balance, calculates how much belongs to the node (taking into consideration rewards and penalties)
func (c *minipoolCommon) CalculateNodeShare(mc *batch.MultiCaller, share_Out **big.Int, balance *big.Int) {
	core.AddCall(mc, c.contract, share_Out, "calculateNodeShare", balance)
}

// Given a validator balance, calculates how much belongs to rETH pool stakers (taking into consideration rewards and penalties)
func (c *minipoolCommon) CalculateUserShare(mc *batch.MultiCaller, share_Out **big.Int, balance *big.Int) {
	core.AddCall(mc, c.contract, share_Out, "calculateUserShare", balance)
}

// Get the data from this minipool's MinipoolPrestaked event
func (c *minipoolCommon) GetPrestakeEvent(intervalSize *big.Int, opts *bind.CallOpts) (PrestakeData, error) {

	addressFilter := []common.Address{c.Address}
	topicFilter := [][]common.Hash{{c.contract.ABI.Events["MinipoolPrestaked"].ID}}

	// Grab the latest block number
	currentBlock, err := c.rp.Client.BlockNumber(context.Background())
	if err != nil {
		return PrestakeData{}, fmt.Errorf("error getting current block %s: %w", c.Address.Hex(), err)
	}

	// Grab the lowest block number worth querying from (should never have to go back this far in practice)
	deployBlockHash := crypto.Keccak256Hash([]byte("deploy.block"))
	var fromBlockBig *big.Int
	err = c.rp.Query(func(mc *batch.MultiCaller) error {
		c.rp.Storage.GetUint(mc, &fromBlockBig, deployBlockHash)
		return nil
	}, opts)
	if err != nil {
		return PrestakeData{}, fmt.Errorf("error getting deploy block %s: %w", c.Address.Hex(), err)
	}

	fromBlock := fromBlockBig.Uint64()
	var log gethtypes.Log
	found := false

	// Backwards scan through blocks to find the event
	for i := currentBlock; i >= fromBlock; i -= eventScanInterval {
		from := i - eventScanInterval + 1
		if from < fromBlock {
			from = fromBlock
		}

		fromBig := big.NewInt(0).SetUint64(from)
		toBig := big.NewInt(0).SetUint64(i)

		logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, fromBig, toBig, nil)
		if err != nil {
			return PrestakeData{}, fmt.Errorf("error getting prestake logs for minipool %s: %w", c.Address.Hex(), err)
		}

		if len(logs) > 0 {
			log = logs[0]
			found = true
			break
		}
	}

	if !found {
		// This should never happen
		return PrestakeData{}, fmt.Errorf("error finding prestake log for minipool %s", c.Address.Hex())
	}

	// Decode the event
	prestakeEvent := new(MinipoolPrestakeEvent)
	c.contract.Contract.UnpackLog(prestakeEvent, "MinipoolPrestaked", log)
	if err != nil {
		return PrestakeData{}, fmt.Errorf("error unpacking prestake data: %w", err)
	}

	// Convert the event to a more useable struct
	prestakeData := PrestakeData{
		Pubkey:                types.BytesToValidatorPubkey(prestakeEvent.Pubkey),
		WithdrawalCredentials: common.BytesToHash(prestakeEvent.WithdrawalCredentials),
		Amount:                prestakeEvent.Amount,
		Signature:             types.BytesToValidatorSignature(prestakeEvent.Signature),
		DepositDataRoot:       prestakeEvent.DepositDataRoot,
		Time:                  time.Unix(prestakeEvent.Time.Int64(), 0),
	}
	return prestakeData, nil
}