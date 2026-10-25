package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	databases "stpNew/backend/database"

	"github.com/joho/godotenv"
)

func main() {
	upFlag := flag.Bool("up", false, "Apply pending database migrations")
	downFlag := flag.Bool("down", false, "Roll back the latest database migration")
	statusFlag := flag.Bool("status", false, "Print the current database migration status")
	flag.Parse()

	if !*upFlag && !*downFlag && !*statusFlag {
		fmt.Println("Usage: go run ./cmd/migrate [-up | -down | -status]")
		os.Exit(1)
	}

	_ = godotenv.Load(".env", "../.env")

	if err := databases.ConnectDB(); err != nil || databases.DB == nil {
		log.Fatalf("Fatal: Database connection failed: %v", err)
	}

	sqlDB, err := databases.DB.DB()
	if err != nil {
		log.Fatalf("Fatal: Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	if *statusFlag {
		version, err := databases.GetMigrationStatus(sqlDB)
		if err != nil {
			log.Fatalf("Error retrieving migration status: %v", err)
		}
		fmt.Printf("Current database migration version: %d\n", version)
		return
	}

	if *upFlag {
		if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("Migration up completed successfully.")
		return
	}

	if *downFlag {
		if err := databases.ApplyMigrationsDown(sqlDB); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("Migration down completed successfully.")
		return
	}
}
