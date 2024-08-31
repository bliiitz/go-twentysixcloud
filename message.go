package client

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type MessageConfirmation struct {
	Chain  MessageChain `json:"chain"`
	Hash   string       `json:"hash"`
	Height uint64       `json:"height"`
}

type Message struct {
	Type      MessageType  `json:"type"`
	Chain     MessageChain `json:"chain"`
	Sender    string       `json:"sender"`
	Time      float64      `json:"time"`
	Channel   string       `json:"channel"`
	Signature string       `json:"signature"`

	ItemHash    string          `json:"item_hash"`
	ItemType    MessageItemType `json:"item_type"`
	ItemContent string          `json:"item_content"`

	Confirmations []MessageConfirmation `json:"confirmations,omitempty"`
	Confirmed     bool                  `json:"confirmed,omitempty"`
}

func (msg Message) GetVerificationPayload() []byte {
	//message signing in typescript
	//Buffer.from([this.chain, this.sender, this.type, this.item_hash].join('\n'))

	return []byte(fmt.Sprintf("%s\n%s\n%s\n%s", msg.Chain, msg.Sender, msg.Type, msg.ItemHash))
}

func (msg *Message) SignMessage(account TwentySixAccount) error {
	messageHash := accounts.TextHash(msg.GetVerificationPayload())

	signature, err := crypto.Sign(messageHash, account.PrivateKey)
	if err != nil {
		return err
	}

	signature[crypto.RecoveryIDOffset] += 27

	msg.Signature = hexutil.Encode(signature)
	return nil
}

func (msg *Message) JSON() []byte {
	payload, err := json.Marshal(msg)
	if err != nil {
		return []byte("")
	}

	return payload
}

func PrepareMessage(account TwentySixAccount, channel string, msgType MessageType, content interface{}, at float64) (Message, error) {
	msgContent, err := json.Marshal(content)
	if err != nil {
		return Message{}, err
	}

	contentHash := sha256.Sum256(msgContent)

	message := Message{
		Type:    msgType,
		Chain:   EthereumChain,
		Sender:  account.Address,
		Time:    at,
		Channel: channel,

		ItemHash:    hex.EncodeToString(contentHash[:]),
		ItemType:    InlineMessageItem,
		ItemContent: string(msgContent),
	}

	message.SignMessage(account)

	return message, nil
}
