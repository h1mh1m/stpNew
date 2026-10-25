package main

import (
	"fmt"
	"log"
	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"

	"github.com/joho/godotenv"
	"gorm.io/gorm/clause"
)

// DevRoomSeeds contains development-only initial facility records aligned with Sirkuit ITS scope.
var DevRoomSeeds = []schemas.Room{
	{
		ID:          "1",
		Name:        "Sirkuit Otomotif ITS",
		Description: "Fasilitas Pengujian & Balap Sirkuit Otomotif Kampus ITS Sukolilo, Surabaya",
		Capacity:    100,
		Price:       150000,
		ImageURL:    "https://images.unsplash.com/photo-1568605117036-5fe5e7bab0b7?auto=format&fit=crop&q=80&w=1000",
		Features:    "Sirkuit Balap & Pengujian",
	},
}

func main() {
	_ = godotenv.Load(".env", "../.env")

	if err := databases.ConnectDB(); err != nil || databases.DB == nil {
		log.Fatalf("Fatal: Database connection failed: %v", err)
	}

	fmt.Println("Seeding development room catalogue...")

	for _, room := range DevRoomSeeds {
		result := databases.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&room)
		if result.Error != nil {
			log.Fatalf("Failed to seed facility '%s' (%s): %v", room.Name, room.ID, result.Error)
		}
		if result.RowsAffected > 0 {
			fmt.Printf("  [+] Inserted facility: %s (%s)\n", room.Name, room.ID)
		} else {
			fmt.Printf("  [=] Facility already exists (preserved): %s (%s)\n", room.Name, room.ID)
		}
	}

	fmt.Println("Room seeding complete. No credentials or customer data were seeded.")
}
