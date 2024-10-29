package etherx

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-chujang/demo-aa/common/utils/ternary"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/etherx/rpcx"
	"golang.org/x/sync/errgroup"
)

const (
	defaultDataGasLimit uint64 = 300_000
	valueGasLimit       uint64 = 21_000
)

type Transactor struct {
	rpcUri     string
	chainId    *big.Int
	signer     types.Signer
	priveteKey *ecdsa.PrivateKey
	Address    common.Address
}

func NewTransactor(uri string, privateKey *ecdsa.PrivateKey, chainIdOps ...*big.Int) (*Transactor, error) {
	addr, err := ethutil.PvKey2Address(privateKey)
	if err != nil {
		return nil, err
	}

	var chainId *big.Int
	if chainIdOps != nil {
		chainId = chainIdOps[0]
	} else {
		chainId, err = rpcx.EasyBig(uri, rpcx.MethodChainId)
		if err != nil {
			return nil, err
		}
	}
	return &Transactor{
		rpcUri:     uri,
		chainId:    chainId,
		signer:     types.NewEIP155Signer(chainId),
		priveteKey: privateKey,
		Address:    addr,
	}, nil
}

func (x *Transactor) prepareTx() (nonce uint64, gasPrice *big.Int, err error) {
	_, res, errs, err := rpcx.EasyBatch(x.rpcUri, []rpcx.Request{
		rpcx.ReqId(rpcx.MethodNonceAt, 0, x.Address, rpcx.BlockParamPending),
		rpcx.ReqId(rpcx.MethodGasPrice, 1),
	})
	if err != nil {
		return
	}
	if err = errors.Join(errs...); err != nil {
		return
	}

	if nonce, err = ethutil.Hex2X[uint64](res[0]); err != nil {
		return
	}
	gasPrice, err = ethutil.Hex2X[*big.Int](res[1])
	return
}

func (x *Transactor) NonceAndGasPrice() (nonce uint64, gasPrice *big.Int, err error) {
	return x.prepareTx()
}

func (x *Transactor) DataTxs(to *common.Address, gas uint64, packedDataList ...[]byte) ([]*types.Transaction, error) {
	if packedDataList == nil {
		return nil, errors.New("packedDataList must not be nil")
	}
	nonce, gasPrice, err := x.prepareTx()
	if err != nil {
		return nil, err
	}

	var (
		txs      = make([]*types.Transaction, 0, len(packedDataList))
		gasLimit = ternary.Cond(gas == 0, defaultDataGasLimit, gas)
	)
	for i, packed := range packedDataList {
		txs = append(txs, types.NewTx(&types.LegacyTx{
			Nonce:    nonce + uint64(i),
			Gas:      gasLimit,
			GasPrice: gasPrice,
			To:       to,
			Data:     packed,
		}))
	}
	return txs, nil
}

func (x *Transactor) ValueTxs(gas *uint64, tos []*common.Address, values []*big.Int, packedDataList ...[]byte) ([]*types.Transaction, error) {
	if len(tos) != len(values) {
		return nil, errors.New("tos and values must have the same length")
	}
	if packedDataList != nil && len(packedDataList) != len(tos) {
		return nil, errors.New("packedDataList and tos must have the same length")
	}
	nonce, gasPrice, err := x.prepareTx()
	if err != nil {
		return nil, err
	}
	gaslimit := valueGasLimit
	if gas != nil {
		gaslimit = *gas
	}
	txs := make([]*types.Transaction, 0, len(tos))
	for i, v := range values {
		to := tos[i]
		if to == nil || v == nil {
			return nil, errors.New("to address and value must not be nil")
		}
		var data []byte
		if packedDataList != nil {
			data = packedDataList[i]
		}
		txs = append(txs, types.NewTx(&types.LegacyTx{
			Nonce:    nonce + uint64(i),
			Gas:      gaslimit,
			GasPrice: gasPrice,
			To:       to,
			Value:    v,
			Data:     data,
		}))
	}
	return txs, nil
}

func (x *Transactor) SendTransactionsWithSign(txs []*types.Transaction) (txHashes []string, rpcErrs []error, err error) {
	var (
		eg     = new(errgroup.Group)
		count  = len(txs)
		rawTxs = make([]string, count)
	)

	txHashes = make([]string, len(txs))
	for idx := range count {
		i := idx
		eg.Go(func() error {
			tx := txs[i]
			hash, rawtx, err := ethutil.SignWithEncode(tx, x.signer, x.priveteKey)
			if err != nil {
				return err
			}
			rawTxs[i] = rawtx
			txHashes[i] = hash
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}
	batchReqs := make([]rpcx.Request, 0, count)
	for id, v := range rawTxs {
		batchReqs = append(batchReqs, rpcx.ReqId(rpcx.MethodSendRawTransaction, id, v))
	}
	_, _, rpcErrs, err = rpcx.EasyBatch(x.rpcUri, batchReqs)
	return txHashes, rpcErrs, err
}
