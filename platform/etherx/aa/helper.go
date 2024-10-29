package aa

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/internal/storedquery"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx"
	"github.com/go-chujang/demo-aa/platform/mongox"
)

func getMngdAccounts(db *mongox.Client, rpcUri string, readOnly bool) (
	wallet *model.MngdWallet,
	entrypoint *etherx.Contractor,
	accountImpl *etherx.Contractor,
	accountFactory *etherx.Contractor,
	tokenPaymaster *etherx.Contractor,
	gamble *etherx.Contractor,
	err error) {

	var manager *etherx.Transactor
	if !readOnly {
		wallet, err = storedquery.GetMngdWallet(db)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
		manager, err = etherx.NewTransactor(rpcUri, wallet.Manager.PrivateKey)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
	}

	contracts, err := storedquery.GetMngdContracts(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	for _, contract := range contracts {
		switch contract.Id {
		case EntrypointName.String():
			entrypoint, _ = etherx.NewContractor(rpcUri, EntrypointAbi, contract.Address)
		case AccountFactoryName.String():
			accountFactory, _ = etherx.NewContractor(rpcUri, AccountFactoryAbi, contract.Address, manager)
		case PaymasterName.String():
			tokenPaymaster, _ = etherx.NewContractor(rpcUri, PaymasterAbi, contract.Address, manager)
		case GambleName.String():
			gamble, _ = etherx.NewContractor(rpcUri, GambleAbi, contract.Address, manager)
		}
	}
	accountImpl, _ = etherx.NewContractor(rpcUri, AccountAbi, common.Address{})
	return wallet, entrypoint, accountImpl, accountFactory, tokenPaymaster, gamble, nil
}
