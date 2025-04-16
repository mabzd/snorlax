package dbm

import (
	"embed"
	"io/fs"
	"log"
	"sort"

	"github.com/mabzd/snorlax/internal/config"
	"github.com/mabzd/snorlax/internal/database"
)

//go:embed resources/*.sql
var embeddedResources embed.FS

type migration struct {
	FileName string
	Sql      string
}

func UpgradeDatabaseIfNeeded(cfg config.Config) {
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	appliedFileNameSet, err := database.GetAppliedMigrationFileNamesSet(db)
	if err != nil {
		log.Fatalf("Failed to get applied file names: %v", err)
	}

	migrations := readMigrations()
	appliedCount := 0

	for _, migration := range migrations {
		if _, ok := appliedFileNameSet[migration.FileName]; !ok {
			log.Printf("Applying migration: %s...\n", migration.FileName)
			database.MustApplyMigration(db, migration.Sql, migration.FileName)
			appliedCount++
		}
	}

	if appliedCount == 0 {
		log.Println("No new migrations to apply")
	} else {
		log.Printf("Applied %d new migrations\n", appliedCount)
	}
}

func readMigrations() []migration {
	entries, err := fs.ReadDir(embeddedResources, "resources")
	if err != nil {
		log.Fatalf("Failed to read resources directory: %v", err)
	}

	sort.Slice(
		entries,
		func(i, j int) bool {
			return entries[i].Name() < entries[j].Name()
		})

	var migrations []migration
	for _, entry := range entries {
		data, err := embeddedResources.ReadFile("resources/" + entry.Name())
		if err != nil {
			log.Fatalf("Failed to read %s: %v\n", entry.Name(), err)
		}

		migrations = append(
			migrations,
			migration{
				FileName: entry.Name(),
				Sql:      string(data),
			})
	}

	return migrations
}
