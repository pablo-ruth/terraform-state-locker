package api

import (
	"encoding/json"
	"io"
)

type Item struct {
	S string `json:"S"`
}

type PutItemRequest struct {
	Item                map[string]Item `json:"Item"`
	TableName           string          `json:"TableName"`
	ConditionExpression string          `json:"ConditionExpression"`
}

type GetItemRequest struct {
	Key       map[string]Item `json:"Key"`
	TableName string          `json:"TableName"`
}

type DeleteItemRequest struct {
	Key       map[string]Item `json:"Key"`
	TableName string          `json:"TableName"`
}

type GetItemResponse struct {
	Item map[string]Item `json:"Item"`
}

func ParsePutItemRequest(body io.Reader) (PutItemRequest, error) {

	var putItemRequest PutItemRequest
	err := json.NewDecoder(body).Decode(&putItemRequest)

	return putItemRequest, err
}

func ParseGetItemRequest(body io.Reader) (GetItemRequest, error) {

	var getItemRequest GetItemRequest
	err := json.NewDecoder(body).Decode(&getItemRequest)

	return getItemRequest, err
}

func ParseDeleteItemRequest(body io.Reader) (DeleteItemRequest, error) {

	var deleteItemRequest DeleteItemRequest
	err := json.NewDecoder(body).Decode(&deleteItemRequest)

	return deleteItemRequest, err
}
