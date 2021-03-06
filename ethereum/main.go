package ethereum

import (
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"

	"bytes"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"math/big"

	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/scmo/apayment-backend/smart-contracts/apayment-token"
	"github.com/scmo/apayment-backend/smart-contracts/rbac"
)

type EthereumController struct {
	Auth     *bind.TransactOpts
	Client   *ethclient.Client
	Keystore *keystore.KeyStore
}

var ethereumController EthereumController

func Init() {
	pathToEthereum := beego.AppConfig.String("ethereumRootPath")

	pathToIPCEndpoint := pathToEthereum + "geth.ipc"

	// Create an IPC based RPC connection to a remote node
	client, err := ethclient.Dial(pathToIPCEndpoint)

	if err != nil {
		beego.Critical("Failed to connect to the Ethereum client: ", err)
	}
	beego.Info("Path to IPC endpoint: ", pathToIPCEndpoint)

	// Keystore to administrate accounts
	ks := keystore.NewKeyStore(pathToEthereum+"keystore", keystore.LightScryptN, keystore.LightScryptP)
	ethereumController = EthereumController{Auth: nil, Client: client, Keystore: ks}
	auth := GetAuth(beego.AppConfig.String("systemAccountAddress"))
	ethereumController.Auth = auth

	deployRoleBasedAccessControlContract()
	deployAPaymentTokenContract()
}

func deployRoleBasedAccessControlContract() {
	if beego.AppConfig.String("accessControlContract") == "" {
		beego.Info("Deploy new RBAC Contract")
		address, _, _, err := rbac.DeployRBACContract(ethereumController.Auth, ethereumController.Client)
		if err != nil {
			beego.Critical("Error while creating RBAC: ", err)
		}
		beego.Debug(address.String())
		beego.AppConfig.Set("accessControlContract", address.String())
	}
}

func deployAPaymentTokenContract() {
	if beego.AppConfig.String("apaymentTokenContract") == "" {
		beego.Info("Deploy new aPayment Token Contract")
		tokenSupply, err := beego.AppConfig.Int("tokenSupply")
		if err != nil {
			beego.Critical("tokenSupply not found. ", err)
		}

		address, _, _, err := apaymenttoken.DeployAPaymentTokenContract(ethereumController.Auth, ethereumController.Client, big.NewInt(int64(tokenSupply)))
		if err != nil {
			beego.Critical("Error while deploying APaymentTokenContract: ", err)
		}
		beego.Debug(address.Str())
		beego.AppConfig.Set("tokenSupply", address.String())
	}
}

func SendWei(from string, to string, amount *big.Int) {
	beego.Info("Send ether from: ", from, "to: ", to, "amount:", amount)
	ctx := context.Background()
	fromAccount, err := ethereumController.Keystore.Find(accounts.Account{Address: common.HexToAddress(from)})
	if fromAccount.Address.String() != from {
		beego.Debug(fromAccount.Address.String(), from)
		beego.Error("something wrong")
		return
	}
	nonce, err := ethereumController.Client.PendingNonceAt(ctx, fromAccount.Address)
	if err != nil {
		beego.Critical("Failed to get nounce: ", err)
	}
	estimateGas, err := ethereumController.Client.EstimateGas(ctx, ethereum.CallMsg{From: fromAccount.Address, Value: amount, Data: nil})
	if err != nil {
		beego.Critical("Failed to estimate gas: ", err)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(to), amount, estimateGas, big.NewInt(50000000000), nil)
	chainId, err := beego.AppConfig.Int64("chainId")
	if err != nil {
		beego.Critical("Failed to get chainID: ", err)
	}
	tx, err = ethereumController.Keystore.SignTxWithPassphrase(fromAccount, beego.AppConfig.String("systemAccountPassword"), tx, big.NewInt(chainId))
	if err != nil {
		beego.Critical("Failed to Sign Transaction: ", err)
	}
	err = ethereumController.Client.SendTransaction(ctx, tx)
	if err != nil {
		beego.Critical("Failed to Send Transaction: ", err)
	}
}

func GetAuth(address string) *bind.TransactOpts {
	account, err := ethereumController.Keystore.Find(accounts.Account{Address: common.HexToAddress(address)})
	if err != nil {
		beego.Error("Error find: ", err)
	}
	dat, err := ioutil.ReadFile(account.URL.Path)

	password := beego.AppConfig.String("userAccountPassword")
	if address == beego.AppConfig.String("systemAccountAddress") {
		password = beego.AppConfig.String("systemAccountPassword")
	}
	auth, err := bind.NewTransactor(bytes.NewReader(dat), password)
	if err != nil {
		beego.Critical("Failed  to create authorized transactor: ", err)
	}
	return auth
}

func GetEthereumController() EthereumController {
	return ethereumController
}
