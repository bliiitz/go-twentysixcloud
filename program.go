package client

import (
	"encoding/json"
	"errors"
	"time"
)

func (client *TwentySixClient) CreateProgram(function ProgramMessageContent) (Message, MessageResponse, error) {
	now := float64(time.Now().UnixMilli()) / 1000

	functionMessage := function
	functionMessage.Time = now
	functionMessage.Address = client.account.Address

	message, res, err := client.SendMessage(InstanceMessageType, functionMessage, now)
	if err != nil {
		return Message{}, MessageResponse{}, err
	}

	var createfunctionResponse MessageResponse
	if err := json.Unmarshal(res, &createfunctionResponse); err != nil {
		return Message{}, MessageResponse{}, err
	}

	return message, createfunctionResponse, nil
}

func (client *TwentySixClient) GetProgramMessages(size uint64, page uint64) ([]Message, uint64, error) {
	return client.GetMessages(size, page, []string{}, []string{client.account.Address}, []string{client.channel}, []MessageType{ProgramMessageType})
}

func (client *TwentySixClient) GetProgramMessageByItemHash(hash string) (Message, error) {
	var page uint64 = 1
	var parsingEnded = false

	for !parsingEnded {
		volumes, remainingItems, err := client.GetProgramMessages(50, page)
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

	return Message{}, errors.New("instance message not found")
}
