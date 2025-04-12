package main

import (
	"log"

	"github.com/mabzd/snorlax/pkg/rest"
)

func main() {
	server := rest.NewServer()
	log.Fatal(server.ListenAndServe())
}
