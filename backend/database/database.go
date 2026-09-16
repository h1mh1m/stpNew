package databases

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Try loading .env from local or parent directory without failing if system env vars exist
	_ = godotenv.Load(".env", "../.env")

	dbName := os.Getenv("DBNAME")
	host := os.Getenv("HOSTNAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("PASSWORD")
	port := os.Getenv("PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Gagal menyambung ke database:", err)
		return
	}
	log.Println("Berhasil menyambung ke database")

	DB = database
}
