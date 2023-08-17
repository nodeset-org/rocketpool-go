package tests

// Contract addresses and account private keys are based on the following mnemonic:
// test test test test test test test test test test test junk

const (
	evmRpcAddress         string = "http://127.0.0.1:8545"
	rocketStorageAddress  string = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
	multicallAddress      string = "0x8464135c8F25Da09e49BC8782676a84730C318bC"
	balanceBatcherAddress string = "0x71C95911E9a5D330f4D621842EC243EE1343292e"
	evmChainID            int    = 31337
	startTime             int    = 12

	ValidatorPubkey     = "968bcf4081af4a10d054c1cde1dadfd6e85a120a397174173ca869f66bdc72835f9918ea251930778e5ba67a7907e30e"
	ValidatorPubkey2    = "968bcf4081af4a10d054c1cde1dadfd6e85a120a397174173ca869f66bdc72835f9918ea251930778e5ba67a7907e30d"
	ValidatorPubkey3    = "968bcf4081af4a10d054c1cde1dadfd6e85a120a397174173ca869f66bdc72835f9918ea251930778e5ba67a7907e30c"
	ValidatorSignature  = "83757098b3b118c67d993218afb69e80a13eb3b174cd3da9958971f05e6b30b9ff5a55677d644f972b31c24e0544604703e8cf18b109fde1e0d3cde0446147bf2f38f02fefce604e4119a605348dfc8a99935dbd65a64eb773c77508f9150e33"
	ValidatorSignature2 = "83757098b3b118c67d993218afb69e80a13eb3b174cd3da9958971f05e6b30b9ff5a55677d644f972b31c24e0544604703e8cf18b109fde1e0d3cde0446147bf2f38f02fefce604e4119a605348dfc8a99935dbd65a64eb773c77508f9150e34"
	ValidatorSignature3 = "83757098b3b118c67d993218afb69e80a13eb3b174cd3da9958971f05e6b30b9ff5a55677d644f972b31c24e0544604703e8cf18b109fde1e0d3cde0446147bf2f38f02fefce604e4119a605348dfc8a99935dbd65a64eb773c77508f9150e35"
)

// These are generated by default with Hardhat, each account has 10000 ETH on it
var AccountPrivateKeys = []string{
	"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
	"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
	"7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
	"47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
	"8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
	"92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
	"4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
	"dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
	"2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
	"f214f2b2cd398c806f84e317254e0f0b801d0643303237d97a22a48e01628897",
	"701b615bbdfb9de65240bc28bd21bbc0d996645a3dd57e7b12bc2bdf6f192c82",
	"a267530f49f8280200edf313ee7af6b827f2a8bce2897751d06a843f644967b1",
	"47c99abed3324a2707c28affff1267e45918ec8c3f20b8aa892e8b065d2942dd",
	"c526ee95bf44d8fc405a158bb884d9d1238d99f0612e9f33d006bb0789009aaa",
	"8166f546bab6da521a8369cab06c5d2b9e46670292d85c875ee9ec20e84ffb61",
	"ea6c44ac03bff858b476bba40716402b03e41b8e97e276d1baec7c37d42484a0",
	"689af8efa8c651a91ad287602527f3af2fe9f6501a7ac4b061667b5a93e037fd",
	"de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0",
	"df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e",
}
