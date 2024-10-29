package etherx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-chujang/demo-aa/common/utils/ternary"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/etherx/rpcx"
	"github.com/go-chujang/packany"
)

const (
	defaultMaxGasLimit    uint64 = 700_000
	defaultDeployGasLimit uint64 = 5_000_000
)

type Contractor struct {
	*Transactor
	ContractAddress common.Address
	abi             abi.ABI
	maxGasLimit     uint64
}

func NewContractor(uri string, abi abi.ABI, addr common.Address, txr ...*Transactor) (*Contractor, error) {
	c := &Contractor{
		ContractAddress: addr,
		maxGasLimit:     defaultMaxGasLimit,
		abi:             abi,
	}
	if txr != nil && txr[0] != nil {
		c.Transactor = txr[0]
	} else {
		c.Transactor = &Transactor{rpcUri: uri}
	}
	return c, nil
}

func (c Contractor) trySign() error {
	if c.Transactor.priveteKey == nil {
		return errors.New("nil private key, cannot sign")
	}
	return nil
}

// Warning: SetTxr
func (c *Contractor) SetTxr(txr *Transactor)         { c.Transactor = txr }
func (c *Contractor) SetMaxGasLimit(gasLimit uint64) { c.maxGasLimit = gasLimit }
func (c *Contractor) GetMaxGasLimit() uint64         { return c.maxGasLimit }

func (c *Contractor) Pack(method string, args ...interface{}) ([]byte, error) {
	return packany.PackArgs(c.abi, method, args...)
}

func (c *Contractor) PackAny(method string, args interface{}) ([]byte, error) {
	return packany.PackAny(c.abi, method, args)
}

func (c *Contractor) Unpack(result interface{}, method string, hexData string) error {
	data, err := hexutil.Decode(hexData)
	if err != nil {
		return err
	}
	return c.abi.UnpackIntoInterface(result, method, data)
}

func (c *Contractor) Deploy(bytecode []byte, args ...interface{}) (string, common.Address, error) {
	if err := c.trySign(); err != nil {
		return "", c.ContractAddress, err
	}
	packed, err := c.Pack("", args...)
	if err != nil {
		return "", c.ContractAddress, err
	}
	nonce, gasPrice, err := c.Transactor.prepareTx()
	if err != nil {
		return "", c.ContractAddress, err
	}
	hashes, rpcErrs, err := c.Transactor.SendTransactionsWithSign([]*types.Transaction{types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Gas:      defaultDeployGasLimit,
		GasPrice: gasPrice,
		To:       nil,
		Data:     append(bytecode, packed...),
	})})
	if err != nil {
		return "", c.ContractAddress, err
	}
	if err = errors.Join(rpcErrs...); err != nil {
		return "", c.ContractAddress, c.Error(err)
	}
	c.ContractAddress = crypto.CreateAddress(c.Transactor.Address, nonce)
	return hashes[0], c.ContractAddress, nil
}

type contractCallMsg struct {
	From common.Address `json:"from"`
	To   common.Address `json:"to"`
	Data string         `json:"data"`
}

func (c *Contractor) CallReq(packed []byte, id int, toOps ...common.Address) rpcx.Request {
	return rpcx.ReqId(rpcx.MethodCall, id, contractCallMsg{
		From: ethutil.ZeroAddress,
		To:   ternary.VArgs(nil, c.ContractAddress, toOps...),
		Data: hexutil.Encode(packed),
	}, rpcx.BlockParamLatest)
}

func (c *Contractor) call(result interface{}, method string, packed []byte) error {
	if result == nil || reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}
	res, err := rpcx.Do[string](c.rpcUri, rpcx.MethodCall, contractCallMsg{
		From: ethutil.ZeroAddress,
		To:   c.ContractAddress,
		Data: hexutil.Encode(packed),
	}, rpcx.BlockParamLatest)
	if err != nil {
		return err
	}
	data, err := hexutil.Decode(res)
	if err != nil {
		return err
	}
	return c.abi.UnpackIntoInterface(result, method, data)
}

func (c *Contractor) Call(result interface{}, method string, args ...interface{}) error {
	packed, err := c.Pack(method, args...)
	if err != nil {
		return err
	}
	return c.call(result, method, packed)
}

func (c *Contractor) CallAny(result interface{}, method string, args any) error {
	packed, err := c.PackAny(method, args)
	if err != nil {
		return err
	}
	return c.call(result, method, packed)
}

func (c *Contractor) Execute(packedDataList ...[]byte) (hashes []string, rpcErrs []error, err error) {
	if err := c.trySign(); err != nil {
		return nil, nil, err
	}
	txs, err := c.DataTxs(&c.ContractAddress, c.maxGasLimit, packedDataList...)
	if err != nil {
		return nil, nil, err
	}
	return c.Transactor.SendTransactionsWithSign(txs)
}

func (c *Contractor) ExecuteOne(method string, args ...interface{}) (string, error) {
	if err := c.trySign(); err != nil {
		return "", err
	}
	packed, err := c.Pack(method, args...)
	if err != nil {
		return "", err
	}
	hashes, rpcErrs, err := c.Execute(packed)
	if err != nil || rpcErrs != nil {
		return "", err
	}
	if err = errors.Join(rpcErrs...); err != nil {
		return "", c.Error(err)
	}
	return hashes[0], nil
}

func (c *Contractor) Error(err error) error {
	if err == nil {
		return nil
	}
	rpcerr, ok := err.(*rpcx.Error)
	if !ok {
		return err
	}
	parsed, ok := rpcerr.ErrorData().(map[string]interface{})
	if !ok {
		return rpcerr
	}
	return fmt.Errorf("%v", parsed)
}
