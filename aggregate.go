package client

import (
	"encoding/json"
	"errors"
	"time"
)

func (client *TwentySixClient) CreateAggregate(aggregate AggregateMessageContent) (Message, MessageResponse, error) {
	now := float64(time.Now().UnixMilli()) / 1000

	aggregateMessage := aggregate
	aggregateMessage.Time = now
	aggregateMessage.Address = client.account.Address

	message, res, err := client.SendMessage(InstanceMessageType, aggregateMessage, now)
	if err != nil {
		return Message{}, MessageResponse{}, err
	}

	var createAggregateResponse MessageResponse
	if err := json.Unmarshal(res, &createAggregateResponse); err != nil {
		return Message{}, MessageResponse{}, err
	}

	return message, createAggregateResponse, nil
}

func (client *TwentySixClient) GetAggregateMessages(size uint64, page uint64) ([]Message, uint64, error) {
	return client.GetMessages(size, page, []string{}, []string{client.account.Address}, []string{client.channel}, []MessageType{AggregateMessageType})
}

func (client *TwentySixClient) GetAggregateMessageByItemHash(hash string) (Message, error) {
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

	return Message{}, errors.New("aggregate message not found")
}
