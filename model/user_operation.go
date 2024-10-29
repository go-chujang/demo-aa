package model

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/IBM/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/kafka"
)

var (
	_ Message          = (*PackedUserOperation)(nil)
	_ RawDataHelper    = (*PackedUserOperation)(nil)
	_ json.Marshaler   = (*PackedUserOperation)(nil)
	_ json.Unmarshaler = (*PackedUserOperation)(nil)
)

type PackedUserOperation struct {
	Sender             common.Address `json:"sender"`
	Nonce              *big.Int       `json:"nonce"`
	InitCode           []byte         `json:"initCode,omitempty"`
	CallData           []byte         `json:"callData"`
	AccountGasLimits   [32]byte       `json:"accountGasLimits"`
	PreVerificationGas *big.Int       `json:"preVerificationGas"`
	GasFees            [32]byte       `json:"gasFees"`
	PaymasterAndData   []byte         `json:"paymasterAndData,omitempty"`
	Signature          []byte         `json:"signature"`
}

func (m PackedUserOperation) Topic() kafka.Topic { return kafka.TopicOperation }
func (m PackedUserOperation) KeyValue() (key sarama.Encoder, value sarama.Encoder, err error) {
	if m.Sender.Cmp(ethutil.ZeroAddress) == 0 {
		return nil, nil, errors.New("address must not be zero")
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, nil, err
	}
	return sarama.ByteEncoder(m.Sender.Bytes()), sarama.ByteEncoder(bytes), nil
}

func (m *PackedUserOperation) Parse(msg *sarama.ConsumerMessage) error {
	return json.Unmarshal(msg.Value, &m)
}

func (m PackedUserOperation) Hint() *string { return toHintString(&m) }
func (m PackedUserOperation) RawData() (hint *string, data map[string]interface{}) {
	return m.Hint(), map[string]interface{}{}
}

func (o PackedUserOperation) MarshalJSON() ([]byte, error) {
	type PackedUserOperation struct {
		Sender             common.Address `json:"sender"`
		Nonce              *hexutil.Big   `json:"nonce"`
		InitCode           hexutil.Bytes  `json:"initCode,omitempty"`
		CallData           hexutil.Bytes  `json:"callData"`
		AccountGasLimits   string         `json:"accountGasLimits"`
		PreVerificationGas *hexutil.Big   `json:"preVerificationGas"`
		GasFees            string         `json:"gasFees"`
		PaymasterAndData   hexutil.Bytes  `json:"paymasterAndData,omitempty"`
		Signature          hexutil.Bytes  `json:"signature"`
	}
	return json.Marshal(&PackedUserOperation{
		Sender:             o.Sender,
		Nonce:              (*hexutil.Big)(o.Nonce),
		InitCode:           hexutil.Bytes(o.InitCode),
		CallData:           hexutil.Bytes(o.CallData),
		AccountGasLimits:   hexutil.Encode(o.AccountGasLimits[:]),
		PreVerificationGas: (*hexutil.Big)(o.PreVerificationGas),
		GasFees:            hexutil.Encode(o.GasFees[:]),
		PaymasterAndData:   hexutil.Bytes(o.PaymasterAndData),
		Signature:          hexutil.Bytes(o.Signature),
	})
}

func (o *PackedUserOperation) UnmarshalJSON(input []byte) error {
	type PackedUserOperation struct {
		Sender             common.Address `json:"sender"`
		Nonce              *hexutil.Big   `json:"nonce"`
		InitCode           hexutil.Bytes  `json:"initCode,omitempty"`
		CallData           hexutil.Bytes  `json:"callData"`
		AccountGasLimits   string         `json:"accountGasLimits"`
		PreVerificationGas *hexutil.Big   `json:"preVerificationGas"`
		GasFees            string         `json:"gasFees"`
		PaymasterAndData   hexutil.Bytes  `json:"paymasterAndData,omitempty"`
		Signature          hexutil.Bytes  `json:"signature"`
	}
	var dec PackedUserOperation
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	o.Sender = dec.Sender
	o.Nonce = (*big.Int)(dec.Nonce)
	o.InitCode = []byte(dec.InitCode)
	o.CallData = []byte(dec.CallData)

	accountGasLimits, err := hexutil.Decode(dec.AccountGasLimits)
	if err != nil {
		return err
	}
	copy(o.AccountGasLimits[:], accountGasLimits)

	o.PreVerificationGas = (*big.Int)(dec.PreVerificationGas)

	gasFees, err := hexutil.Decode(dec.GasFees)
	if err != nil {
		return err
	}
	copy(o.GasFees[:], gasFees)

	o.PaymasterAndData = []byte(dec.PaymasterAndData)
	o.Signature = []byte(dec.Signature)
	return nil
}
