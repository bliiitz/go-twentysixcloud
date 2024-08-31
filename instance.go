package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func (client *TwentySixClient) CreateInstance(instance InstanceMessageContent) (Message, MessageResponse, error) {
	now := float64(time.Now().UnixMilli()) / 1000

	instanceMessage := instance
	instanceMessage.Time = now
	instanceMessage.Address = client.account.Address

	message, res, err := client.SendMessage(InstanceMessageType, instanceMessage, now)
	if err != nil {
		return Message{}, MessageResponse{}, err
	}

	var createInstanceResponse MessageResponse
	if err := json.Unmarshal(res, &createInstanceResponse); err != nil {
		return Message{}, MessageResponse{}, err
	}

	return message, createInstanceResponse, nil
}

func (client *TwentySixClient) GetInstanceState(hash string) (SchedulerAllocation, error) {
	body := &bytes.Buffer{}
	endpoint := "https://scheduler.api.aleph.sh/api/v0/allocation/" + hash

	var res SchedulerAllocation

	request, err := http.NewRequest("GET", endpoint, body)
	if err != nil {
		return res, err
	}

	request.Header.Add("Accept", "application/json")
	response, err := client.http.Do(request)
	if err != nil {
		return res, err
	}

	resultBody, err := io.ReadAll(response.Body)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(resultBody, &res); err != nil {
		return res, err
	}

	return res, nil
}

func (client *TwentySixClient) GetInstanceMessages(size uint64, page uint64) ([]Message, uint64, error) {
	return client.GetMessages(size, page, []string{}, []string{client.account.Address}, []string{client.channel}, []MessageType{InstanceMessageType})
}

func (client *TwentySixClient) GetInstanceMessageByItemHash(hash string) (Message, error) {
	var page uint64 = 1
	var parsingEnded = false

	for !parsingEnded {
		volumes, remainingItems, err := client.GetInstanceMessages(50, page)
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
