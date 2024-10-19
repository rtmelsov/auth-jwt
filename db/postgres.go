package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

var DB *pgx.Conn

func ConnectDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	DB = conn
	fmt.Println("Successfully connected to database")
}

func CloseDB() {
	err := DB.Close(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
