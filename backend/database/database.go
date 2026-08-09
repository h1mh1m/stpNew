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
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln("Error load .env file")
	}
	dbName := os.Getenv("DBNAME")
	host := os.Getenv("HOSTNAME")
	user := os.Getenv("DB_USER")
	Password := os.Getenv("PASSWORD")
	port := os.Getenv("PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, Password, dbName, port)
	fmt.Println(dsn)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Gagal nyambung ke database", err)
	}
	log.Println("Berhasil nyambung ke database")

	DB = database

}
