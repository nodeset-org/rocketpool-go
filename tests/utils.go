package tests

import (
	"fmt"
	"math/big"

	batchquery "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/tokens"
)

// Mint old RPL for unit testing
func MintLegacyRpl(rp *rocketpool.RocketPool, ownerAccount *Account, toAccount *Account, amount *big.Int) (*core.TransactionInfo, error) {
	fsrpl, err := rp.GetContract(rocketpool.ContractName_RocketTokenRPLFixedSupply)
	if err != nil {
		return nil, fmt.Errorf("error creating legacy RPL contract: %w", err)
	}

	return core.NewTransactionInfo(fsrpl, "mint", ownerAccount.Transactor, toAccount.Address, amount)
}

// Registers a new Rocket Pool node
func RegisterNode(rp *rocketpool.RocketPool, account *Account, timezone string) (*node.Node, error) {
	// Create the node
	node, err := node.NewNode(rp, account.Address)
	if err != nil {
		return nil, fmt.Errorf("error creating node %s: %w", account.Address.Hex(), err)
	}

	// Register the node
	err = rp.CreateAndWaitForTransaction(func() (*core.TransactionInfo, error) {
		return node.Register(timezone, account.Transactor)
	}, true, account.Transactor)
	if err != nil {
		return nil, fmt.Errorf("error registering node %s: %w", account.Address.Hex(), err)
	}

	return node, nil
}

// Bootstraps a node into the Oracle DAO, taking care of all of the details involved
func BootstrapNodeToOdao(rp *rocketpool.RocketPool, owner *Account, nodeAccount *Account, timezone string, id string, url string) (*node.Node, error) {
	// Get some contract bindings
	odaoMgr, err := oracle.NewOracleDaoManager(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting oDAO manager binding: %w", err)
	}
	oma, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrustedActions)
	if err != nil {
		return nil, fmt.Errorf("error getting OMA contract: %w", err)
	}
	fsrpl, err := tokens.NewTokenRplFixedSupply(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting FSRPL binding: %w", err)
	}
	rpl, err := tokens.NewTokenRpl(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting RPL binding: %w", err)
	}
	rplContract, err := rp.GetContract(rocketpool.ContractName_RocketTokenRPL)
	if err != nil {
		return nil, fmt.Errorf("error getting RPL contract: %w", err)
	}

	// Register the node
	node, err := RegisterNode(rp, nodeAccount, timezone)
	if err != nil {
		return nil, fmt.Errorf("error registering node: %w", err)
	}

	// Get the amount of RPL to mint
	oSettings := odaoMgr.Settings
	err = rp.Query(func(mc *batchquery.MultiCaller) error {
		odaoMgr.GetMemberCount(mc)
		oSettings.GetRplBond(mc)
		return nil
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting network info: %w", err)
	}
	rplBond := oSettings.Members.RplBond

	// Bootstrap it and mint RPL for it
	err = rp.BatchCreateAndWaitForTransactions([]func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return odaoMgr.BootstrapMember(id, url, nodeAccount.Address, owner.Transactor)
		},
		func() (*core.TransactionInfo, error) {
			return MintLegacyRpl(rp, owner, nodeAccount, rplBond)
		},
	}, true, owner.Transactor)
	if err != nil {
		return nil, fmt.Errorf("error bootstrapping node and minting RPL: %w", err)
	}

	// Swap RPL and Join the oDAO
	err = rp.BatchCreateAndWaitForTransactions([]func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return fsrpl.Approve(*rplContract.Address, rplBond, nodeAccount.Transactor)
		},
		func() (*core.TransactionInfo, error) {
			return rpl.SwapFixedSupplyRplForRpl(rplBond, nodeAccount.Transactor)
		},
		func() (*core.TransactionInfo, error) {
			return rpl.Approve(*oma.Address, rplBond, nodeAccount.Transactor)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Join(nodeAccount.Transactor)
		},
	}, false, nodeAccount.Transactor)
	if err != nil {
		return nil, fmt.Errorf("error joining oDAO: %w", err)
	}

	return node, nil
}
