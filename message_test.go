package client

import (
	"bytes"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestPrepareMessage(t *testing.T) {

	pkey := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	acc, err := NewTwentySixAccountFromPrivateKey(pkey)
	if err != nil {
		t.Fatalf(`NewTwentySixAccountFromPrivateKey failed to be instanciated: %v`, err)
	}

	now := float64(time.Now().UnixMilli()) / 1000
	msgContent := AggregateMessageContent{
		Key:     "test",
		Address: acc.Address,
		Time:    now,
		Content: map[string]string{
			"Hello": "World",
		},
	}

	message, err := PrepareMessage(acc, "TEST", AggregateMessageType, msgContent, now)
	if err != nil {
		t.Fatal(err)
	}

	sig, err := hexutil.Decode(message.Signature)
	if err != nil {
		t.Fatalf(`Decoding signature failed: %v`, err)
	}

	sig[crypto.RecoveryIDOffset] -= 27

	verificationPayload := accounts.TextHash(message.GetVerificationPayload())

	sigPublicKey, err := crypto.Ecrecover(verificationPayload, sig)
	if err != nil {
		t.Fatal(err)
	}

	matches := bytes.Equal(sigPublicKey, crypto.FromECDSAPub(acc.PublicKey))
	if !matches {
		t.Fatal("ECRecover failed")
	}
}
