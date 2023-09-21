package protocol

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rocket-pool/rocketpool-go/dao"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
)

// Estimate the gas of a proposal submission
func EstimateProposalGas(rp *rocketpool.RocketPool, message string, payload []byte, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	return rocketDAOProtocolProposals.GetTransactionGasInfo(opts, "propose", message, payload, blockNumber, treeNodes)
}

// Submit a trusted node DAO proposal
// Returns the ID of the new proposal
func SubmitProposal(rp *rocketpool.RocketPool, message string, payload []byte, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (uint64, common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return 0, common.Hash{}, err
	}
	proposalCount, err := dao.GetProposalCount(rp, nil)
	if err != nil {
		return 0, common.Hash{}, err
	}
	tx, err := rocketDAOProtocolProposals.Transact(opts, "propose", message, payload, blockNumber, treeNodes)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("error submitting Protocol DAO proposal: %w", err)
	}
	return proposalCount + 1, tx.Hash(), nil
}

// Estimate the gas of ProposeSetMulti
func EstimateProposeSetMultiGas(rp *rocketpool.RocketPool, message string, contractNames []string, settingPaths []string, settingTypes []types.ProposalSettingType, values []any, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	encodedValues, err := abiEncodeMultiValues(settingTypes, values)
	if err != nil {
		return rocketpool.GasInfo{}, fmt.Errorf("error ABI encoding values: %w", err)
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingMulti", contractNames, settingPaths, settingTypes, encodedValues)
	if err != nil {
		return rocketpool.GasInfo{}, fmt.Errorf("error setting multi-set proposal payload: %w", err)
	}
	return EstimateProposalGas(rp, message, payload, blockNumber, treeNodes, opts)
}

// Submit a proposal to update multiple Protocol DAO settings at once
func ProposeSetMulti(rp *rocketpool.RocketPool, message string, contractNames []string, settingPaths []string, settingTypes []types.ProposalSettingType, values []any, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (uint64, common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return 0, common.Hash{}, err
	}
	encodedValues, err := abiEncodeMultiValues(settingTypes, values)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("error ABI encoding values: %w", err)
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingMulti", contractNames, settingPaths, settingTypes, encodedValues)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("error setting multi-set proposal payload: %w", err)
	}
	return SubmitProposal(rp, message, payload, blockNumber, treeNodes, opts)
}

// Estimate the gas of ProposeSetBool
func EstimateProposeSetBoolGas(rp *rocketpool.RocketPool, message, contractName, settingPath string, value bool, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingBool", contractName, settingPath, value)
	if err != nil {
		return rocketpool.GasInfo{}, fmt.Errorf("error setting bool setting proposal payload: %w", err)
	}
	return EstimateProposalGas(rp, message, payload, blockNumber, treeNodes, opts)
}

// Submit a proposal to update a bool Protocol DAO setting
func ProposeSetBool(rp *rocketpool.RocketPool, message, contractName, settingPath string, value bool, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (uint64, common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return 0, common.Hash{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingBool", contractName, settingPath, value)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("error setting bool setting proposal payload: %w", err)
	}
	return SubmitProposal(rp, message, payload, blockNumber, treeNodes, opts)
}

// Estimate the gas of ProposeSetUint
func EstimateProposeSetUintGas(rp *rocketpool.RocketPool, message, contractName, settingPath string, value *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingUint", contractName, settingPath, value)
	if err != nil {
		return rocketpool.GasInfo{}, fmt.Errorf("error encoding set uint setting proposal payload: %w", err)
	}
	return EstimateProposalGas(rp, message, payload, blockNumber, treeNodes, opts)
}

// Submit a proposal to update a uint Protocol DAO setting
func ProposeSetUint(rp *rocketpool.RocketPool, message, contractName, settingPath string, value *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (uint64, common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return 0, common.Hash{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingUint", contractName, settingPath, value)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("error encoding set uint setting proposal payload: %w", err)
	}
	return SubmitProposal(rp, message, payload, blockNumber, treeNodes, opts)
}

// Estimate the gas of ProposeSetAddress
func EstimateProposeSetAddressGas(rp *rocketpool.RocketPool, message, contractName, settingPath string, value common.Address, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingAddress", contractName, settingPath, value)
	if err != nil {
		return rocketpool.GasInfo{}, fmt.Errorf("error encoding set address setting proposal payload: %w", err)
	}
	return EstimateProposalGas(rp, message, payload, blockNumber, treeNodes, opts)
}

// Submit a proposal to update an address Protocol DAO setting
func ProposeSetAddress(rp *rocketpool.RocketPool, message, contractName, settingPath string, value common.Address, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (uint64, common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return 0, common.Hash{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingAddress", contractName, settingPath, value)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("error encoding set address setting proposal payload: %w", err)
	}
	return SubmitProposal(rp, message, payload, blockNumber, treeNodes, opts)
}

// Estimate the gas of ProposeSetAddress
func EstimateProposeSetRewardsPercentageGas(rp *rocketpool.RocketPool, message, contractName, settingPath string, percentage *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingRewardsClaimer", contractName, percentage)
	if err != nil {
		return rocketpool.GasInfo{}, fmt.Errorf("error encoding set rewards-claimer percent proposal payload: %w", err)
	}
	return EstimateProposalGas(rp, message, payload, blockNumber, treeNodes, opts)
}

// Submit a proposal to update one of percentage allocations of RPL rewards
func ProposeSetRewardsPercentage(rp *rocketpool.RocketPool, message, contractName string, percentage *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (uint64, common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return 0, common.Hash{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSettingRewardsClaimer", contractName, percentage)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("error encoding set rewards-claimer percent proposal payload: %w", err)
	}
	return SubmitProposal(rp, message, payload, blockNumber, treeNodes, opts)
}

// Estimate the gas of ProposeSpendTreasury
func EstimateProposeSpendTreasuryGas(rp *rocketpool.RocketPool, message, invoiceID string, recipient common.Address, amount *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSpendTreasury", invoiceID, recipient, amount)
	if err != nil {
		return rocketpool.GasInfo{}, fmt.Errorf("error encoding set spend-treasury percent proposal payload: %w", err)
	}
	return EstimateProposalGas(rp, message, payload, blockNumber, treeNodes, opts)
}

// Submit a proposal to spend a portion of the Rocket Pool treasury
func ProposeSpendTreasury(rp *rocketpool.RocketPool, message, invoiceID string, recipient common.Address, amount *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (uint64, common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return 0, common.Hash{}, err
	}
	payload, err := rocketDAOProtocolProposals.ABI.Pack("proposalSpendTreasury", invoiceID, recipient, amount)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("error encoding set spend-treasury percent proposal payload: %w", err)
	}
	return SubmitProposal(rp, message, payload, blockNumber, treeNodes, opts)
}

// Estimate the gas of CancelProposal
func EstimateCancelProposalGas(rp *rocketpool.RocketPool, proposalId uint64, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	return rocketDAOProtocolProposals.GetTransactionGasInfo(opts, "cancel", big.NewInt(int64(proposalId)))
}

// Cancel a submitted proposal
func CancelProposal(rp *rocketpool.RocketPool, proposalId uint64, opts *bind.TransactOpts) (common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := rocketDAOProtocolProposals.Transact(opts, "cancel", big.NewInt(int64(proposalId)))
	if err != nil {
		return common.Hash{}, fmt.Errorf("error cancelling Protocol DAO proposal %d: %w", proposalId, err)
	}
	return tx.Hash(), nil
}

// Estimate the gas of VoteOnProposal
func EstimateVoteOnProposalGas(rp *rocketpool.RocketPool, proposalId uint64, support bool, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	return rocketDAOProtocolProposals.GetTransactionGasInfo(opts, "vote", big.NewInt(int64(proposalId)), support)
}

// Vote on a submitted proposal
func VoteOnProposal(rp *rocketpool.RocketPool, proposalId uint64, support bool, opts *bind.TransactOpts) (common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := rocketDAOProtocolProposals.Transact(opts, "vote", big.NewInt(int64(proposalId)), support)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error voting on Protocol DAO proposal %d: %w", proposalId, err)
	}
	return tx.Hash(), nil
}

// Estimate the gas of ExecuteProposal
func EstimateExecuteProposalGas(rp *rocketpool.RocketPool, proposalId uint64, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	return rocketDAOProtocolProposals.GetTransactionGasInfo(opts, "execute", big.NewInt(int64(proposalId)))
}

// Execute a submitted proposal
func ExecuteProposal(rp *rocketpool.RocketPool, proposalId uint64, opts *bind.TransactOpts) (common.Hash, error) {
	rocketDAOProtocolProposals, err := getRocketDAOProtocolProposals(rp, nil)
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := rocketDAOProtocolProposals.Transact(opts, "execute", big.NewInt(int64(proposalId)))
	if err != nil {
		return common.Hash{}, fmt.Errorf("error executing Protocol DAO proposal %d: %w", proposalId, err)
	}
	return tx.Hash(), nil
}

// Get the ABI encoding of multiple values for a ProposeSettingMulti call
func abiEncodeMultiValues(settingTypes []types.ProposalSettingType, values []any) ([][]byte, error) {
	// Sanity check the lengths
	settingCount := len(settingTypes)
	if settingCount != len(values) {
		return nil, fmt.Errorf("settingTypes and values must be the same length")
	}
	if settingCount == 0 {
		return [][]byte{}, nil
	}

	// ABI encode each value
	results := make([][]byte, settingCount)
	for i, settingType := range settingTypes {
		var encodedArg []byte
		switch settingType {
		case types.ProposalSettingType_Uint256:
			arg, success := values[i].(*big.Int)
			if !success {
				return nil, fmt.Errorf("value %d is not a *big.Int, but the setting type is Uint256", i)
			}
			encodedArg = math.U256Bytes(big.NewInt(0).Set(arg))

		case types.ProposalSettingType_Bool:
			arg, success := values[i].(bool)
			if !success {
				return nil, fmt.Errorf("value %d is not a bool, but the setting type is Bool", i)
			}
			if arg {
				encodedArg = math.PaddedBigBytes(common.Big1, 32)
			} else {
				encodedArg = math.PaddedBigBytes(common.Big0, 32)
			}

		case types.ProposalSettingType_Address:
			arg, success := values[i].(common.Address)
			if !success {
				return nil, fmt.Errorf("value %d is not an address, but the setting type is Address", i)
			}
			encodedArg = common.LeftPadBytes(arg.Bytes(), 32)

		default:
			return nil, fmt.Errorf("unknown proposal setting type [%v]", settingType)
		}
		results[i] = encodedArg
	}

	return results, nil
}

// Get contracts
var rocketDAOProtocolProposalsLock sync.Mutex

func getRocketDAOProtocolProposals(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*rocketpool.Contract, error) {
	rocketDAOProtocolProposalsLock.Lock()
	defer rocketDAOProtocolProposalsLock.Unlock()
	return rp.GetContract("rocketDAOProtocolProposals", opts)
}