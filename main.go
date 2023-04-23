package main

import (
	"flag"
	"fmt"

	"github.com/pablo-ruth/terraform-state-locker/api"
	"github.com/pablo-ruth/terraform-state-locker/store"
)

func main() {

	addr := flag.String("l", "0.0.0.0:8000", "Listen address")
	cert := flag.String("c", "cert.pem", "Path to TLS certificate")
	key := flag.String("k", "key.pem", "Path to TLS private key")
	flag.Parse()

	store := store.NewInMemoryStore()

	err := api.Serve(*addr, *cert, *key, store)
	if err != nil {
		fmt.Println(err)
	}
}
