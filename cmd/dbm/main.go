package main

import (
	"log"

	"github.com/mabzd/snorlax/internal/config"
	"github.com/mabzd/snorlax/pkg/dbm"
)

// Snorlax database migration tool (dbm).
func main() {
	log.SetPrefix("[dbm] ")
	log.Println("Running snorlax database migration")
	cfg := config.LoadConfig()
	dbm.UpgradeDatabaseIfNeeded(cfg)
}
