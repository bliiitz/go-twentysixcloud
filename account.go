package client

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type TwentySixAccount struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

func NewTwentySixAccountFromPrivateKey(privateKey string) (TwentySixAccount, error) {
	var privateKeyBytes []byte
	if privateKey[0:2] == "0x" {
		pk, err := hexutil.Decode(privateKey)
		if err != nil {
			return TwentySixAccount{}, errors.New("error casting private key to bytes")
		}

		privateKeyBytes = pk
	} else {
		pk, err := hex.DecodeString(privateKey)
		if err != nil {
			return TwentySixAccount{}, errors.New("error casting private key to bytes")
		}

		privateKeyBytes = pk
	}

	privateKeyEcdsa, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return TwentySixAccount{}, errors.New("error casting private key to ECDSA")
	}

	publicKey := privateKeyEcdsa.Public()

	publicKeyEcdsa, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return TwentySixAccount{}, errors.New("error casting public key to ECDSA")
	}

	return TwentySixAccount{
		PrivateKey: privateKeyEcdsa,
		PublicKey:  publicKeyEcdsa,
		Address:    crypto.PubkeyToAddress(*publicKeyEcdsa).Hex(),
	}, nil
}

func NewTwentySixAccountFromMnemonic(mnemonic string, derivationPath string) (TwentySixAccount, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	var dpath string
	if len(derivationPath) == 0 {
		dpath = "m/44'/60'/0'/0/0"
	} else {
		dpath = derivationPath
	}

	path := hdwallet.MustParseDerivationPath(dpath)
	account, err := wallet.Derive(path, true)
	if err != nil {
		return TwentySixAccount{}, err
	}

	publicKey, err := wallet.PublicKey(account)
	if err != nil {
		return TwentySixAccount{}, err
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return TwentySixAccount{}, err
	}

	return TwentySixAccount{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    crypto.PubkeyToAddress(*publicKey).Hex(),
	}, nil
}
