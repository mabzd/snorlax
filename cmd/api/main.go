package main

import (
	"log"

	"github.com/mabzd/snorlax/internal/config"
	"github.com/mabzd/snorlax/pkg/rest"
)

// Snorlax REST service (api).
func main() {
	log.SetPrefix("[api] ")
	log.Println("Running snorlax api")
	cfg := config.LoadConfig()
	server := rest.NewServer(cfg)
	log.Fatal(server.ListenAndServe())
}
