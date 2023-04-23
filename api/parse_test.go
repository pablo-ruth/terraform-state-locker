package api

import (
	"reflect"
	"strings"
	"testing"
)

func TestParsePutItemRequest(t *testing.T) {

	body := `{"ConditionExpression":"attribute_not_exists(LockID)","Item":{"Info":{"S":"{\"ID\":\"bc4abeab-0f07-e6b0-8b6d-a68460074a8e\",\"Operation\":\"OperationTypePlan\",\"Info\":\"\",\"Who\":\"pablo@APORDSI28\",\"Version\":\"1.3.2\",\"Created\":\"2023-04-17T17:42:37.305265252Z\",\"Path\":\"tfstates/dynamodbtest\"}"},"LockID":{"S":"tfstates/dynamodbtest"}},"TableName":"terraform-lock-table"}`
	expected := PutItemRequest{
		Item: map[string]Item{
			"Info": {
				S: "{\"ID\":\"bc4abeab-0f07-e6b0-8b6d-a68460074a8e\",\"Operation\":\"OperationTypePlan\",\"Info\":\"\",\"Who\":\"pablo@APORDSI28\",\"Version\":\"1.3.2\",\"Created\":\"2023-04-17T17:42:37.305265252Z\",\"Path\":\"tfstates/dynamodbtest\"}",
			},
			"LockID": {
				S: "tfstates/dynamodbtest",
			},
		},
		TableName: "terraform-lock-table",
	}

	putItemRequest, err := ParsePutItemRequest(strings.NewReader(body))
	if err != nil {
		t.Errorf("Error parsing PutItemRequest: %v", err)
	}

	if reflect.DeepEqual(expected, putItemRequest) {
		t.Errorf("Expected: %v\nGot: %v", expected, putItemRequest)
	}
}

func TestParseGetItemRequest(t *testing.T) {

	body := `{"ConsistentRead":true,"Key":{"LockID":{"S":"tfstates/dynamodbtest"}},"ProjectionExpression":"LockID,Info","TableName":"terraform-lock-table"}`
	expected := GetItemRequest{
		Key: map[string]Item{
			"LockID": {
				S: "tfstates/dynamodbtest",
			},
		},
		TableName: "terraform-lock-table",
	}

	getItemRequest, err := ParseGetItemRequest(strings.NewReader(body))
	if err != nil {
		t.Errorf("Error parsing GetItemRequest: %v", err)
	}

	if reflect.DeepEqual(expected, getItemRequest) {
		t.Errorf("Expected: %v\nGot: %v", expected, getItemRequest)
	}
}
