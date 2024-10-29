package ethutil

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func SignWithEncode(tx *types.Transaction, signer types.Signer, prv *ecdsa.PrivateKey) (string, string, error) {
	if signer == nil {
		signer = types.HomesteadSigner{}
	}
	signed, err := types.SignTx(tx, signer, prv)
	if err != nil {
		return "", "", err
	}
	bin, err := signed.MarshalBinary()
	if err != nil {
		return "", "", err
	}
	return signed.Hash().Hex(), hexutil.Encode(bin), nil
}

// eip-191
func Signature(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if message == nil || privateKey == nil {
		return nil, errors.New("hashbytes or privateKey must not be nil")
	}
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	prefixedMessage := append([]byte(prefix), message...)
	hash := crypto.Keccak256(prefixedMessage)
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}
	if signature[64] < 27 {
		signature[64] += 27
	}
	return signature, nil
}

func VerifySignature(address common.Address, message []byte, signature []byte) (bool, error) {
	if len(signature) != 65 {
		return false, errors.New("invalid signature length")
	}

	v := signature[64]
	if v < 27 {
		v += 27
	}
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	prefixedMessage := append([]byte(prefix), message...)
	hash := crypto.Keccak256(prefixedMessage)
	pubKey, err := crypto.SigToPub(hash, append(signature[:64], v-27))
	if err != nil {
		return false, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr == address, nil
}
