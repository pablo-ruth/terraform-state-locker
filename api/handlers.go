package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/pablo-ruth/terraform-state-locker/store"
)

func handlePutItem(w http.ResponseWriter, r *http.Request, s store.Store) {

	putItemRequest, err := ParsePutItemRequest(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		log.Printf("Parsing request: %v", err)
		return
	}

	var notExists bool
	if putItemRequest.ConditionExpression != "" {
		re := regexp.MustCompile(`attribute_not_exists\((.*)\)`)
		match := re.FindStringSubmatch(putItemRequest.ConditionExpression)
		if len(match) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid condition expression"))
			log.Printf("Parsing condition expression: %v", err)
			return
		}

		if match[1] != "LockID" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Condition expression must be on LockID"))
			log.Printf("Parsing condition expression, must be on LockID: %v", err)
			return
		}

		notExists = true
	}

	lockID, ok := putItemRequest.Item["LockID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("LockID is missing"))
		log.Printf("LockID is missing")
		return
	}

	attributes := map[string]string{}
	for k, v := range putItemRequest.Item {
		if k == "LockID" {
			continue
		}

		attributes[k] = v.S
	}

	err = s.Put(putItemRequest.TableName, lockID.S, notExists, attributes)
	if err != nil {
		if err == store.ErrEntryAlreadyExists {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Conflict"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		log.Printf("Storing lock: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func handleGetItem(w http.ResponseWriter, r *http.Request, s store.Store) {

	getItemRequest, err := ParseGetItemRequest(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		log.Printf("Parsing request: %v", err)
		return
	}

	lockID, ok := getItemRequest.Key["LockID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("LockID is missing"))
		return
	}

	values, err := s.Get(getItemRequest.TableName, lockID.S)
	if err != nil {
		if err == store.ErrEntryNotFound || err == store.ErrTableNotFound {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
			log.Printf("Entry not found: %v", err)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		log.Printf("Storing lock: %v", err)
		return
	}

	var resp GetItemResponse
	resp.Item = map[string]Item{}
	for k, v := range values {
		if k == "LockID" {
			continue
		}

		resp.Item[k] = Item{S: v}
	}

	itemJSON, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(itemJSON))
}

func handleDeleteItem(w http.ResponseWriter, r *http.Request, s store.Store) {

	deleteItemRequest, err := ParseDeleteItemRequest(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		log.Printf("Parsing request: %v", err)
		return
	}

	lockID, ok := deleteItemRequest.Key["LockID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("LockID is missing"))
		return
	}

	err = s.Delete(deleteItemRequest.TableName, lockID.S)
	if err != nil {

		if err == store.ErrEntryNotFound || err == store.ErrTableNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		log.Printf("Deleting lock: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
