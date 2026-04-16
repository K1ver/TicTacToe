// database/postgres.go
package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DataBase struct {
	DB *pgxpool.Pool
}

func NewDataBase() *DataBase {
	return &DataBase{
		DB: ConnectPostgresDB(),
	}
}

func ConnectPostgresDB() *pgxpool.Pool {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")

	log.Printf("Connecting to database with URL: postgres://postgres:****@localhost:5432/TicTacToeDB")

	dbpool, _ := pgxpool.New(context.Background(), databaseURL)

	log.Println("Successfully connected to PostgreSQL database")
	return dbpool
}

func CloseDB(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Println("Database connection closed")
	}
}
