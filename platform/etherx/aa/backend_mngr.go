package aa

import (
	"context"
	"errors"
	"sync"

	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx"
	"github.com/go-chujang/demo-aa/platform/mongox"
)

type managerBackend struct {
	*accountAbstract
	txr *etherx.Transactor
	mu  sync.Mutex
}

func NewMngrBackend(db *mongox.Client, rpcUri string) (*managerBackend, error) {
	aa, err := newAA(db, rpcUri)
	if err != nil {
		return nil, err
	}
	mngdTxr := &managerBackend{}
	mngdTxr.accountAbstract = aa
	mngdTxr.txr = aa.accountFactory.Transactor
	return mngdTxr, nil
}

func (m *managerBackend) PrepareCreateAccount(msg *model.CreateAccount) BundlePayload {
	packed, _ := m.packCreateAccount(msg.Owner)
	return BundlePayload{
		ca:     m.accountFactory.ContractAddress,
		packed: packed,
	}
}

func (m *managerBackend) PrepareFaucet(msg *model.Faucet) BundlePayload {
	packed, _ := m.packTransfer(msg.Receiver, msg.Value)
	return BundlePayload{
		ca:     m.tokenPaymaster.ContractAddress,
		packed: packed,
	}
}

func (m *managerBackend) BundleExec(ctx context.Context, payloadList []BundlePayload, rawDataHelpers []model.RawDataHelper) ([]error, error) {
	if payloadList == nil {
		return nil, nil
	}
	if len(payloadList) != len(rawDataHelpers) {
		return nil, errors.New("mismatched slice lengths")
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	hashes, rpcErrs, err := m.bundleExec(m.txr, payloadList)
	if err != nil {
		return nil, err
	}
	return m.bunleCommit(ctx, hashes, rpcErrs, rawDataHelpers)
}
