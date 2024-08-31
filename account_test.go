package client

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestAccountCreationFromPrivateKey(t *testing.T) {

	pkey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	acc, err := NewTwentySixAccountFromPrivateKey(pkey)
	if err != nil {
		t.Fatalf(`NewTwentySixAccountFromMnemonic failed to be instanciated: %v`, err)
	}

	if hexutil.Encode(crypto.FromECDSA(acc.PrivateKey)) != "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" {
		t.Fatalf(`Bad private key generated`)
	}
	if hexutil.Encode(crypto.FromECDSAPub(acc.PublicKey)) != "0x048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5" {
		t.Fatalf(`Bad public key generated`)
	}
	if crypto.PubkeyToAddress(*acc.PublicKey).Hex() != "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266" {
		t.Fatalf(`Bad address key generated`)
	}
}

func TestAccountCreationFrom0xPrivateKey(t *testing.T) {

	pkey := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	acc, err := NewTwentySixAccountFromPrivateKey(pkey)
	if err != nil {
		t.Fatalf(`NewTwentySixAccountFromMnemonic failed to be instanciated: %v`, err)
	}

	if hexutil.Encode(crypto.FromECDSA(acc.PrivateKey)) != "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" {
		t.Fatalf(`Bad private key generated`)
	}
	if hexutil.Encode(crypto.FromECDSAPub(acc.PublicKey)) != "0x048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5" {
		t.Fatalf(`Bad public key generated`)
	}
	if crypto.PubkeyToAddress(*acc.PublicKey).Hex() != "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266" {
		t.Fatalf(`Bad address key generated`)
	}
}

func TestAccountCreationFromMnemonic(t *testing.T) {

	mnemonic := "test test test test test test test test test test test junk"

	acc, err := NewTwentySixAccountFromMnemonic(mnemonic, "m/44'/60'/0'/0/0")
	if err != nil {
		t.Fatalf(`NewTwentySixAccountFromMnemonic failed to be instanciated: %v`, err)
	}

	if hexutil.Encode(crypto.FromECDSA(acc.PrivateKey)) != "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" {
		t.Fatalf(`Bad private key generated`)
	}
	if hexutil.Encode(crypto.FromECDSAPub(acc.PublicKey)) != "0x048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5" {
		t.Fatalf(`Bad public key generated`)
	}
	if crypto.PubkeyToAddress(*acc.PublicKey).Hex() != "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266" {
		t.Fatalf(`Bad address key generated`)
	}
}
