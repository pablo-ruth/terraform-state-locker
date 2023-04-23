package store

import (
	"fmt"
	"sync"
)

var (
	ErrTableNotFound      = fmt.Errorf("table not found")
	ErrEntryNotFound      = fmt.Errorf("entry not found")
	ErrEntryAlreadyExists = fmt.Errorf("entry already exists")
	ErrInvalidPrimaryKey  = fmt.Errorf("invalid primary key")
)

type Store interface {
	Get(table, id string) (map[string]string, error)
	Put(table, id string, notExists bool, values map[string]string) error
	Delete(table, id string) error
}

type InMemoryStore struct {
	mu     sync.Mutex
	tables map[string]InMemoryStoreTable
}

type InMemoryStoreTable struct {
	entries map[string]InMemoryStoreEntry
}

type InMemoryStoreEntry struct {
	attributes []struct {
		key   string
		value string
	}
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		tables: make(map[string]InMemoryStoreTable),
	}
}

func (s *InMemoryStore) Get(table, id string) (map[string]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	storeTable, ok := s.tables[table]
	if !ok {
		return nil, ErrTableNotFound
	}

	storeEntry, ok := storeTable.entries[id]
	if !ok {
		return nil, ErrEntryNotFound
	}

	result := make(map[string]string)
	for _, entry := range storeEntry.attributes {
		result[entry.key] = entry.value
	}

	return result, nil
}

func (s *InMemoryStore) Put(table, id string, notExists bool, attributes map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	storeTable, ok := s.tables[table]
	if !ok {
		storeTable = InMemoryStoreTable{
			entries: make(map[string]InMemoryStoreEntry),
		}
	}

	if notExists {
		_, ok := storeTable.entries[id]
		if ok {
			return ErrEntryAlreadyExists
		}
	}

	storeEntry := InMemoryStoreEntry{}
	for key, value := range attributes {
		storeEntry.attributes = append(storeEntry.attributes, struct {
			key   string
			value string
		}{
			key:   key,
			value: value,
		})
	}

	storeTable.entries[id] = storeEntry
	s.tables[table] = storeTable

	return nil
}

func (s *InMemoryStore) Delete(table, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	storeTable, ok := s.tables[table]
	if !ok {
		return ErrTableNotFound
	}

	_, ok = storeTable.entries[id]
	if ok {
		delete(storeTable.entries, id)
		return nil
	}

	return ErrEntryNotFound
}
