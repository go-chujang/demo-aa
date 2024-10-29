package aa

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	defaultCallGasLimit         = big.NewInt(350_000)
	defaultVerificationGasLimit = big.NewInt(150_000)
	defaultPreVerificationGas   = big.NewInt(21_000)
	defaultMaxPriorityFeePerGas = big.NewInt(1_000000000)
	defaultMaxFeePerGas         = big.NewInt(20_000000000)
	defaultValidationGasLimit   = big.NewInt(150_000)
	defaultPostOpGasLimit       = big.NewInt(100_000)
)

func userOperationDefault(paymaster common.Address) *PackedUserOperation {
	return new(PackedUserOperation).
		SetInitCode(nil).
		SetAccountGasLimits(defaultCallGasLimit, defaultVerificationGasLimit).
		SetPreVerificationGas(defaultPreVerificationGas).
		SetGasFees(defaultMaxPriorityFeePerGas, defaultMaxFeePerGas).
		SetPaymasterAndData(paymaster, defaultValidationGasLimit, defaultPostOpGasLimit)
}
