package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pablo-ruth/terraform-state-locker/store"
)

func Serve(addr, cert, key string, store store.Store) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Header.Get("X-Amz-Target") {
		case "":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("X-Amz-Target header is missing"))
			log.Printf("X-Amz-Target header is missing")
		case "DynamoDB_20120810.PutItem":
			handlePutItem(w, r, store)
		case "DynamoDB_20120810.GetItem":
			handleGetItem(w, r, store)
		case "DynamoDB_20120810.DeleteItem":
			handleDeleteItem(w, r, store)
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown X-Amz-Target header"))
			log.Printf("Unknown X-Amz-Target header")
		}
	})

	fmt.Printf("Server is running on port https://%s\n", addr)
	err := http.ListenAndServeTLS(addr, cert, key, r)

	return err
}
