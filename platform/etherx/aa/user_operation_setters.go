package aa

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
)

func (p *PackedUserOperation) SetSender(sender common.Address) *PackedUserOperation {
	p.Sender = sender
	return p
}

func (p *PackedUserOperation) SetNonce(nonce *big.Int) *PackedUserOperation {
	p.Nonce = nonce
	return p
}

func (p *PackedUserOperation) SetInitCode(salt []byte) *PackedUserOperation {
	p.InitCode = salt
	return p
}

func (p *PackedUserOperation) SetCallData(calldata []byte) *PackedUserOperation {
	p.CallData = calldata
	return p
}

func (p *PackedUserOperation) SetAccountGasLimits(callGasLimit, verificationGasLimit *big.Int) *PackedUserOperation {
	p.AccountGasLimits = ethutil.PackLowHighByBig(callGasLimit, verificationGasLimit)
	return p
}

func (p *PackedUserOperation) SetPreVerificationGas(preVerificationGas *big.Int) *PackedUserOperation {
	p.PreVerificationGas = preVerificationGas
	return p
}

func (p *PackedUserOperation) SetGasFees(maxPriorityFeePerGas, maxFeePerGas *big.Int) *PackedUserOperation {
	p.GasFees = ethutil.PackLowHighByBig(maxPriorityFeePerGas, maxFeePerGas)
	return p
}

func (p *PackedUserOperation) SetPaymasterAndData(paymaster common.Address, validationGasLimit, postOpGasLimit *big.Int) *PackedUserOperation {
	if paymaster == ethutil.ZeroAddress {
		return p
	}
	var (
		buf  = new(bytes.Buffer)
		data = ethutil.PackLowHighByBig(validationGasLimit, postOpGasLimit)
	)
	buf.Write(paymaster.Bytes())
	buf.Write(data[:])

	p.PaymasterAndData = buf.Bytes()
	return p
}

func (p *PackedUserOperation) SetSignature(sig []byte) *PackedUserOperation {
	p.Signature = sig
	return p
}
