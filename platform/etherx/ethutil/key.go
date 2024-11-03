package ethutil

import (
	"crypto/ecdsa"
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-chujang/demo-aa/common/utils/ternary"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

func ParseMnemonic(mnemonic string, count uint32, derivationPath ...accounts.DerivationPath) ([]*ecdsa.PrivateKey, error) {
	seed := bip39.NewSeed(mnemonic, "")
	master, err := hdkeychain.NewMaster(seed, &chaincfg.Params{})
	if err != nil {
		return nil, err
	}

	parsed := make([]*ecdsa.PrivateKey, 0, count)
	path := ternary.VArgs(nil, accounts.DefaultBaseDerivationPath, derivationPath...)
	for idx := range count {
		path[len(path)-1] = idx

		key := master
		for _, v := range path {
			key, _ = key.DeriveNonStandard(v)
		}
		pv, err := key.ECPrivKey()
		if err != nil {
			return nil, err
		}
		// parsed = append(parsed, pv.ToECDSA())
		ecdsaKey, err := crypto.ToECDSA(crypto.FromECDSA(pv.ToECDSA()))
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, ecdsaKey)
	}
	return parsed, nil
}

func ParseMnemonic2(mnemonic string, count uint32, derivationPath ...accounts.DerivationPath) ([]*ecdsa.PrivateKey, error) {
	seed := bip39.NewSeed(mnemonic, "")
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return nil, err
	}

	parsed := make([]*ecdsa.PrivateKey, 0, count)
	path := ternary.VArgs(nil, accounts.DefaultBaseDerivationPath, derivationPath...)
	for idx := range count {
		path[len(path)-1] = idx
		account, err := wallet.Derive(path, false)
		if err != nil {
			return nil, err
		}
		pv, err := wallet.PrivateKey(account)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, pv)
	}
	return parsed, nil
}

func PvKey2Address(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	if privateKey == nil {
		return common.Address{}, errors.New("privateKey must not be nil")
	}
	return crypto.PubkeyToAddress(privateKey.PublicKey), nil
}

func PvKey2Hex(privateKey *ecdsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", errors.New("privateKey must not be nil")
	}
	return hexutil.Encode(crypto.FromECDSA(privateKey)), nil
}

func PvHex2Key(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	privateKeyHex = Trim0x(privateKeyHex)
	return crypto.HexToECDSA(privateKeyHex)
}
