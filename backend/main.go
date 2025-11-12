package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"stocks/domain"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("Did not find .env file")
	} else {
		log.Println("Found .env file!")
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	log.Println("Db connection successfull")

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting db object")
	}
	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("Db not reachable %v", err)
	}
	log.Println("Db reachable")
	return db
}

func getFirstStocks(db *gorm.DB) ([]domain.Stock, error) {
	var stocks []domain.Stock
	result := db.Limit(5).Find(&stocks)
	return stocks, result.Error
}

func main() {
	db := initDb()
	stocks, err := getFirstStocks(db)
	if err != nil {
		log.Fatalf("Error getting stocks %v", err)
	} else {
		fmt.Printf("First 5 stocks: %v\n", stocks)
	}
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(stocks)
	})

	http.ListenAndServe(":3000", cors(mux))
}
