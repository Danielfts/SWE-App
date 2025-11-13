package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

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

func getCentroids() domain.KMeansModel {
	file, err := os.Open(os.Getenv("CENTROIDS_PATH"))
	if err != nil {
		log.Fatalf("Error getting centroids file")
	}
	defer file.Close()

	var model domain.KMeansModel
	if err := json.NewDecoder(file).Decode(&model); err != nil {
		log.Fatalf("Error decoding JSON")
	}
	fmt.Printf("Obtained centroids %v\n", model)
	return model
}

func normalize(feature, mean, std float64) float64 {
	if std == 0 {
		return 0
	}
	return (feature - mean) / std
}

func getStockFeatures(stock *domain.Stock, model *domain.KMeansModel) (domain.KMeansFeatures, error) {
	// Score maps
	means := model.Means
	stds := model.Stds
	var actionMapping = map[string]int{
		"upgrades":              1,
		"upgraded by":           1,
		"target raised by":      1,
		"downgrades":            -1,
		"downgraded by":         -1,
		"target lowered by":     -1,
		"maintains":             0,
		"initiates coverage on": 0,
		"initiates":             0,
		"initiated by":          0,
		"reiterates":            0,
		"reiterated by":         0,
		"target set by":         0,
		"":                      0,
	}
	var ratingMapping = map[string]float64{
		// Most bullish (1.0)
		"strong-buy": 1.0,
		"strong buy": 1.0,

		// Very bullish (0.75–0.9)
		"buy":               0.75,
		"speculative buy":   0.75,
		"outperform":        0.75,
		"outperformer":      0.75,
		"overweight":        0.75,
		"accumulate":        0.75,
		"market outperform": 0.75,
		"sector outperform": 0.75,

		// Slightly bullish (0.25–0.5)
		"positive": 0.5,

		// Neutral (0)
		"hold":             0,
		"neutral":          0,
		"market perform":   0,
		"equal weight":     0,
		"equal-weight":     0,
		"in-line":          0,
		"sector perform":   0,
		"sector performer": 0,
		"sector weight":    0,
		"peer perform":     0,

		// Slightly bearish (-0.25 to -0.5)
		"negative": -0.5,

		// Very bearish (-0.75 to -0.9)
		"underperform":        -0.75,
		"underweight":         -0.75,
		"reduce":              -0.75,
		"sector underperform": -0.75,

		// Most bearish (-1.0)
		"sell":        -1.0,
		"strong sell": -1.0,

		// Empty/missing
		"": 0,
	}

	// Get target delta
	targetToNum, err := strconv.ParseFloat(stock.TargetTo, 64)
	if err != nil {
		return domain.KMeansFeatures{}, err
	}
	targetFromNum, err := strconv.ParseFloat(stock.TargetFrom, 64)
	if err != nil {
		return domain.KMeansFeatures{}, err
	}
	var targetDelta float64 = ((targetToNum - targetFromNum) / targetFromNum) * 100
	// Get Has brokerage
	var hasBrokerage int
	if len(stock.Brokerage) > 0 {
		hasBrokerage = 1
	} else {
		hasBrokerage = 0
	}

	// Get Action Score
	var actionScore int = actionMapping[stock.Action]

	// Get Rating Score Delta
	ratingFromScore := ratingMapping[stock.RatingFrom]
	ratingToScore := ratingMapping[stock.RatingTo]
	ratingDelta := ratingToScore - ratingFromScore
	// Get Time delta
	timestamp, err := time.Parse(time.RFC3339, stock.Time)
	if err != nil {
		return domain.KMeansFeatures{}, err
	}
	now := time.Now()
	delta := now.Sub(timestamp).Hours() / 24

	features := domain.KMeansFeatures{
		TargetDelta:      normalize(targetDelta, means[0], stds[0]),
		HasBrokerage:     normalize(float64(hasBrokerage), means[1], stds[1]),
		ActionScore:      normalize(float64(actionScore), means[2], stds[2]),
		RatingDeltaScore: normalize(ratingDelta, means[3], stds[3]),
		TimeDelta:        normalize(math.Round(delta), means[4], stds[4]),
	}
	fmt.Printf("Kmeans features: %+v", features)
	return features, nil
}

func recommendStock(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	model := getCentroids()
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
	stock := stocks[0]
	_, err := getStockFeatures(&stock, &model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stock)
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
