package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	ctr := context.Background()

	contex := "localhost/5176"

	conn, err := pgxpool.New(ctr, contex)

	if err != nil {
		log.Fatalf("Gagal nyambung: %v\n", err)
	}

	if err := conn.Ping(ctr); err != nil {
		log.Fatalf("Database not responded %v\n", err)
	}
	DB = conn
	fmt.Println("Berhasil terhubung ke Postgresql")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
