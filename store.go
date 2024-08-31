package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (client *TwentySixClient) StoreFile(filePath string) (Message, string, error) {
	now := float64(time.Now().UnixMilli()) / 1000
	file, err := os.Open(filePath)
	if err != nil {
		return Message{}, "", err
	}

	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return Message{}, "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//Generate metadata
	metadatapart, err := writer.CreateFormField("metadata")
	if err != nil {
		return Message{}, "", err
	}

	itemContent := StoreMessageContent{
		Address:  client.account.Address,
		Time:     now,
		ItemHash: hex.EncodeToString(hash.Sum(nil)),
		ItemType: StorageMessageItem,
	}

	jsonItem, err := json.Marshal(itemContent)
	if err != nil {
		return Message{}, "", err
	}

	contentHash := sha256.Sum256(jsonItem)

	message := Message{
		Chain:       EthereumChain,
		Sender:      client.account.Address,
		Channel:     client.channel,
		Time:        now,
		Type:        StoreMessageType,
		ItemType:    InlineMessageItem,
		ItemHash:    hex.EncodeToString(contentHash[:]),
		ItemContent: string(jsonItem),
	}

	message.SignMessage(client.account)

	req := BroadcastRequest{
		Message: message,
		Sync:    false,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return Message{}, "", err
	}

	metadata := bytes.NewReader(jsonReq)
	io.Copy(metadatapart, metadata)

	//Upload file
	filepart, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return Message{}, "", err
	}

	file, err = os.Open(filePath)
	if err != nil {
		return Message{}, "", err
	}

	defer file.Close()

	io.Copy(filepart, file)
	writer.Close()

	storeEndpoint := client.apiUrl + "/api/v0/storage/add_file"
	request, err := http.NewRequest("POST", storeEndpoint, body)
	if err != nil {
		return Message{}, "", err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Accept", "application/json")

	response, err := client.http.Do(request)
	if err != nil {
		return Message{}, "", err
	}

	resultBody, err := io.ReadAll(response.Body)
	if err != nil {
		return Message{}, "", err
	}

	var storeFileResponse StoreIPFSFileResponse
	if err := json.Unmarshal(resultBody, &storeFileResponse); err != nil {
		return Message{}, "", err
	}

	defer response.Body.Close()

	time.Sleep(5 * time.Second)

	createdMessage, err := client.GetStoreMessageByItemHash(storeFileResponse.Hash)
	if err != nil {
		return Message{}, "", err
	}

	return createdMessage, storeFileResponse.Hash, nil
}

func (client *TwentySixClient) GetStoreMessages(size uint64, page uint64) ([]Message, uint64, error) {
	return client.GetMessages(size, page, []string{}, []string{client.account.Address}, []string{client.channel}, []MessageType{StoreMessageType})
}

func (client *TwentySixClient) GetStoreMessageByItemHash(hash string) (Message, error) {
	var page uint64 = 1
	var parsingEnded = false

	for !parsingEnded {
		volumes, remainingItems, err := client.GetStoreMessages(50, page)
		if err != nil {
			return Message{}, err
		}

		for i := 0; i < len(volumes); i++ {
			var itemContent StoreMessageContent
			json.Unmarshal([]byte(volumes[i].ItemContent), &itemContent)

			if itemContent.ItemHash == hash {
				return volumes[i], nil
			}
		}

		if remainingItems > 0 {
			page += 1
		} else {
			parsingEnded = true
		}
	}

	return Message{}, errors.New("store message not found")
}
