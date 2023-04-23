package store

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInMemoryStorePut(t *testing.T) {

	cases := []struct {
		name        string
		store       *InMemoryStore
		table       string
		id          string
		notExists   bool
		attributes  map[string]string
		expectedErr error
		expected    InMemoryStore
	}{
		{
			name:  "put one item",
			store: NewInMemoryStore(),
			table: "terraform-lock-table",
			id:    "tfstates/dynamodbtest",
			attributes: map[string]string{
				"Info": "Test",
			},
			expectedErr: nil,
			expected: InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "put one item with notExists",
			store: &InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
						},
					},
				},
			},
			table: "terraform-lock-table",
			id:    "tfstates/dynamodbtest",
			attributes: map[string]string{
				"Info": "Test2",
			},
			notExists:   true,
			expectedErr: ErrEntryAlreadyExists,
			expected:    InMemoryStore{},
		},
		{
			name: "put one item with notExists at false",
			store: &InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
						},
					},
				},
			},
			table: "terraform-lock-table",
			id:    "tfstates/dynamodbtest",
			attributes: map[string]string{
				"Info": "Test2",
			},
			notExists:   false,
			expectedErr: nil,
			expected: InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test2",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.name == "put one item with notExists" {
				fmt.Println(c.store)
			}
			err := c.store.Put(c.table, c.id, c.notExists, c.attributes)
			if err != c.expectedErr {
				t.Errorf("Expected error %v, got %v", c.expectedErr, err)
			} else if err != nil {
				return
			}

			if !reflect.DeepEqual(*c.store, c.expected) {
				t.Errorf("Expected %v, got %v", c.expected, c.store)
			}
		})
	}
}

func TestLocakStoreGet(t *testing.T) {

	cases := []struct {
		name        string
		store       *InMemoryStore
		table       string
		id          string
		expectedErr error
		expected    map[string]string
	}{
		{
			name: "get one item",
			store: &InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
						},
					},
				},
			},
			table:       "terraform-lock-table",
			id:          "tfstates/dynamodbtest",
			expectedErr: nil,
			expected: map[string]string{
				"Info": "Test",
			},
		},
		{
			name: "get one item that does not exist",
			store: &InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
						},
					},
				},
			},
			table:       "terraform-lock-table",
			id:          "tfstates/dynamodbtest2",
			expectedErr: ErrEntryNotFound,
			expected: map[string]string{
				"Info": "Test",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			attributes, err := c.store.Get(c.table, c.id)
			if err != c.expectedErr {
				t.Errorf("Expected error %v, got %v", c.expectedErr, err)
			} else if err != nil {
				return
			}

			if !reflect.DeepEqual(attributes, c.expected) {
				t.Errorf("Expected %v, got %v", c.expected, attributes)
			}
		})
	}
}

func TestInMemoryStoreDelete(t *testing.T) {

	cases := []struct {
		name        string
		store       *InMemoryStore
		table       string
		id          string
		expectedErr error
		expected    InMemoryStore
	}{
		{
			name: "delete one item",
			store: &InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
							"tfstates/dynamodbtest2": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
						},
					},
				},
			},
			table:       "terraform-lock-table",
			id:          "tfstates/dynamodbtest",
			expectedErr: nil,
			expected: InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest2": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "delete one item that does not exist",
			store: &InMemoryStore{
				tables: map[string]InMemoryStoreTable{
					"terraform-lock-table": {
						entries: map[string]InMemoryStoreEntry{
							"tfstates/dynamodbtest": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
							"tfstates/dynamodbtest2": {
								attributes: []struct {
									key   string
									value string
								}{
									{
										key:   "Info",
										value: "Test",
									},
								},
							},
						},
					},
				},
			},
			table:       "terraform-lock-table",
			id:          "tfstates/dynamodbtest3",
			expectedErr: ErrEntryNotFound,
			expected:    InMemoryStore{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.store.Delete(c.table, c.id)
			if err != c.expectedErr {
				t.Errorf("Expected error %v, got %v", c.expectedErr, err)
			} else if err != nil {
				return
			}

			if !reflect.DeepEqual(*c.store, c.expected) {
				t.Errorf("Expected %v, got %v", c.expected, c.store)
			}
		})
	}
}

// func TestInMemoryStoreDelete(t *testing.T) {

// 	cases := []struct {
// 		name        string
// 		table       string
// 		key         string
// 		value       string
// 		InMemoryStore  InMemoryStore
// 		expectedErr error
// 		expected    InMemoryStore
// 	}{
// 		{
// 			name:  "delete one item",
// 			table: "terraform-lock-table",
// 			key:   "LockID",
// 			value: "tfstates/dynamodbtest2",
// 			InMemoryStore: InMemoryStore{
// 				tables: map[string]InMemoryStoreTable{
// 					"terraform-lock-table": {
// 						entries: []InMemoryStoreEntry{
// 							{
// 								attributes: map[string]struct {
// 									key   string
// 									value string
// 								}{
// 									"LockID": {
// 										key:   "LockID",
// 										value: "tfstates/dynamodbtest",
// 									},
// 									"Info": {
// 										key:   "Info",
// 										value: "Test",
// 									},
// 								},
// 							},
// 							{
// 								attributes: map[string]struct {
// 									key   string
// 									value string
// 								}{
// 									"LockID": {
// 										key:   "LockID",
// 										value: "tfstates/dynamodbtest2",
// 									},
// 									"Info": {
// 										key:   "Info",
// 										value: "Test2",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			expectedErr: nil,
// 			expected: InMemoryStore{
// 				tables: map[string]InMemoryStoreTable{
// 					"terraform-lock-table": {
// 						entries: []InMemoryStoreEntry{
// 							{
// 								attributes: map[string]struct {
// 									key   string
// 									value string
// 								}{
// 									"LockID": {
// 										key:   "LockID",
// 										value: "tfstates/dynamodbtest",
// 									},
// 									"Info": {
// 										key:   "Info",
// 										value: "Test",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "delete one item with wrong key",
// 			table: "terraform-lock-table",
// 			key:   "LockID",
// 			value: "tfstates/dynamodbtest3",
// 			InMemoryStore: InMemoryStore{
// 				tables: map[string]InMemoryStoreTable{
// 					"terraform-lock-table": {
// 						entries: []InMemoryStoreEntry{},
// 					},
// 				},
// 			},
// 			expectedErr: ErrEntryNotFound,
// 			expected: InMemoryStore{
// 				tables: map[string]InMemoryStoreTable{
// 					"terraform-lock-table": {
// 						entries: []InMemoryStoreEntry{},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			store := tc.InMemoryStore

// 			err := store.Delete(tc.table, tc.key, tc.value)
// 			if err != nil {
// 				if err != tc.expectedErr {
// 					t.Fatalf("unexpected error: %v, expected: %v", err, tc.expectedErr)
// 				}
// 			}

// 			if !reflect.DeepEqual(store, tc.expected) {
// 				t.Fatalf("expected: %v, got: %v", tc.expected, store)
// 			}
// 		})
// 	}
// }
