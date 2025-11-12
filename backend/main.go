package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

func getFirstStocks(db *gorm.DB, offset int, sortby *string, asc bool, queryString string) ([]domain.Stock, error) {
	var stocks []domain.Stock
	orderBy := "ticker"
	if sortby != nil && *sortby != "" {
		orderBy = *sortby
	}
	sortOrder := ""
	if asc {
		sortOrder = " ASC"
	} else {
		sortOrder = " DESC"
	}
	dbQuery := db.Limit(5).Order(orderBy + sortOrder).Offset(offset)
	if len(queryString) > 0 {
		pattern := fmt.Sprintf("%%%s%%", queryString)
		dbQuery = dbQuery.Where("ticker ILIKE ?", pattern)
	}
	result := dbQuery.Find(&stocks)
	return stocks, result.Error
}

func queryStocks(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var StockColumnMap = map[string]string{
		"Id":         "id",
		"Ticker":     "ticker",
		"TargetFrom": "target_from",
		"TargetTo":   "target_to",
		"Company":    "company",
		"Action":     "action",
		"Brokerage":  "brokerage",
		"RatingFrom": "rating_from",
		"RatingTo":   "rating_to",
		"Time":       "time",
	}
	q := r.URL.Query()
	offsetStr := q.Get("offset")
	sortByStr := q.Get("sortby")
	asc := q.Get("asc") == "true"
	query := q.Get("query")
	sortBy := StockColumnMap[sortByStr]
	offset := 0
	if len(offsetStr) > 0 {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		offset = offset * 5
		if err != nil {
			fmt.Println("Invalid offset:", err)
			offset = 0
		}
	}
	stocks, err := getFirstStocks(db, offset, &sortBy, asc, query)
	if err != nil {
		log.Fatalf("Error getting stocks %v", err)
	} else {
		fmt.Printf("First 5 stocks: %v\n", stocks)
	}
	json.NewEncoder(w).Encode(stocks)
}

func recommendStock(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var stocks []domain.Stock
	result := db.Limit(1).Find(&stocks)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if len(stocks) != 1 {
		http.Error(w, fmt.Sprintf("expected 1 stock, got %d", len(stocks)), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(stocks[0])
}

func main() {
	db := initDb()
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	mux := http.NewServeMux()
	mux.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		queryStocks(w, r, db)
	})

	mux.HandleFunc("/recommendation", func(w http.ResponseWriter, r *http.Request) {
		recommendStock(w, r, db)
	})

	http.ListenAndServe(":3000", cors(mux))
}
