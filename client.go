package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const AlephApiUrl string = "https://api3.aleph.im"

type TwentySixClient struct {
	account TwentySixAccount
	channel string
	apiUrl  string
	http    http.Client
}

func (client *TwentySixClient) GetMessageByHash(hash string) (Message, error) {

	//https://api2.aleph.im/api/v0/messages.json?hashes=d51f34748974a1e652becd28c28249c2eb5a0cfaf8b718dde7121034d5733981
	messageEndpoint := AlephApiUrl + "/api/v0/messages.json?hashes=" + hash
	request, err := http.NewRequest("GET", messageEndpoint, bytes.NewBuffer([]byte("")))
	if err != nil {
		return Message{}, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := client.http.Do(request)
	if err != nil {
		return Message{}, err
	}

	resultBody, err := io.ReadAll(response.Body)
	if err != nil {
		return Message{}, err
	}

	var result GetMessageResponse
	if err := json.Unmarshal(resultBody, &result); err != nil { // Parse []byte to go struct pointer
		return Message{}, err
	}

	defer response.Body.Close()

	if result.PaginationTotal != 1 {
		return Message{}, errors.New("message not found")
	} else {
		return result.Messages[0], nil
	}
}

func (client *TwentySixClient) WaitMessageConfirmation(hash string, timeout int64, interval int64) error {
	var startAt int64 = time.Now().Unix()
	var message Message

	message, err := client.GetMessageByHash(hash)
	if err != nil {
		return err
	}

	for !message.Confirmed {
		time.Sleep(time.Duration(interval) * time.Second)

		message, err = client.GetMessageByHash(hash)
		if err != nil {
			return err
		}

		now := time.Now().Unix()
		if now > startAt+timeout {
			return errors.New("message confirmation timeout")
		}
	}

	return nil
}

func (client *TwentySixClient) SendMessage(msgType MessageType, content interface{}, at float64) (Message, []byte, error) {

	message, err := PrepareMessage(client.account, client.channel, msgType, content, at)
	if err != nil {
		return Message{}, []byte{}, err
	}

	req := BroadcastRequest{
		Message: message,
		Sync:    false,
	}

	buff, err := json.Marshal(req)
	if err != nil {
		return Message{}, []byte{}, err
	}

	messageEndpoint := AlephApiUrl + "/api/v0/messages"
	request, err := http.NewRequest("POST", messageEndpoint, bytes.NewBuffer(buff))
	if err != nil {
		return Message{}, []byte{}, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := client.http.Do(request)
	if err != nil {
		return Message{}, []byte{}, err
	}

	resultBody, err := io.ReadAll(response.Body)
	if err != nil {
		return Message{}, []byte{}, err
	}

	return message, resultBody, nil
}

func (client *TwentySixClient) GetMessages(size uint64, page uint64, hashes []string, addresses []string, channels []string, msgTypes []MessageType) ([]Message, uint64, error) {
	var messages []Message
	body := &bytes.Buffer{}

	messageEndpoint := AlephApiUrl + "/api/v0/messages.json?"

	params := url.Values{}

	params.Add("page", fmt.Sprint(page))
	params.Add("size", fmt.Sprint(size))

	for i := 0; i < len(hashes); i++ {
		params.Add("hashes", hashes[i])
	}
	for i := 0; i < len(addresses); i++ {
		params.Add("addresses", addresses[i])
	}
	for i := 0; i < len(channels); i++ {
		params.Add("channels", channels[i])
	}
	for i := 0; i < len(msgTypes); i++ {
		params.Add("msgTypes", string(msgTypes[i]))
	}

	filteredEndpoint := messageEndpoint + params.Encode()

	request, err := http.NewRequest("GET", filteredEndpoint, body)
	if err != nil {
		return messages, 0, err
	}

	request.Header.Add("Accept", "application/json")

	response, err := client.http.Do(request)
	if err != nil {
		return messages, 0, err
	}

	resultBody, err := io.ReadAll(response.Body)
	if err != nil {
		return messages, 0, err
	}

	var getMessageResponse GetMessageResponse
	if err := json.Unmarshal(resultBody, &getMessageResponse); err != nil {
		return messages, 0, err
	}

	for i := 0; i < len(getMessageResponse.Messages); i++ {
		messages = append(messages, getMessageResponse.Messages[i])
	}

	var remainingItems uint64
	if getMessageResponse.PaginationPage*getMessageResponse.PaginationPerPage > getMessageResponse.PaginationTotal {
		remainingItems = 0
	} else {
		remainingItems = getMessageResponse.PaginationTotal - (getMessageResponse.PaginationPage * getMessageResponse.PaginationPerPage)
	}

	return messages, remainingItems, nil
}

func (client *TwentySixClient) ForgetMessage(hash string) (MessageResponse, error) {
	now := float64(time.Now().UnixMilli()) / 1000

	itemContent := ForgetMessageContent{
		Address: client.account.Address,
		Time:    now,
		Hashes:  []string{hash},
	}

	msgContent, err := json.Marshal(itemContent)
	if err != nil {
		return MessageResponse{}, err
	}

	contentHash := sha256.Sum256(msgContent)

	message := Message{
		Type:    ForgetMessageType,
		Chain:   EthereumChain,
		Sender:  client.account.Address,
		Time:    now,
		Channel: client.channel,

		ItemHash:    hex.EncodeToString(contentHash[:]),
		ItemType:    InlineMessageItem,
		ItemContent: string(msgContent),
	}

	message.SignMessage(client.account)

	req := BroadcastRequest{
		Message: message,
		Sync:    false,
	}

	buff, err := json.Marshal(req)
	if err != nil {
		return MessageResponse{}, err
	}

	storeEndpoint := AlephApiUrl + "/api/v0/messages"
	request, err := http.NewRequest("POST", storeEndpoint, bytes.NewBuffer(buff))
	if err != nil {
		return MessageResponse{}, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := client.http.Do(request)
	if err != nil {
		return MessageResponse{}, err
	}

	resultBody, err := io.ReadAll(response.Body)
	if err != nil {
		return MessageResponse{}, err
	}

	var parsedRes MessageResponse
	json.Unmarshal(resultBody, &parsedRes)

	return parsedRes, nil
}

func NewTwentySixClient(acc TwentySixAccount, channel string, apiUrl string) TwentySixClient {
	return TwentySixClient{
		account: acc,
		channel: channel,
		apiUrl:  apiUrl,
		http:    http.Client{},
	}
}
