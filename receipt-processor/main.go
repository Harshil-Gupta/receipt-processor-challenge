package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode"
)

var receiptStore = make(map[string]Receipt)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
}

func main() {
	// Initialize a new router
	router := mux.NewRouter()

	// Serve static files (e.g., CSS, JS)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Serve the basic webpage on the root path
	router.HandleFunc("/", homePageHandler).Methods("GET")

	// Receipt processing routes
	router.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

	// Start the HTTP server on port 8080
	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

// HomePageHandler serves the index.html file
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file
	http.ServeFile(w, r, "index.html")
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	receiptID := uuid.New().String()
	receiptStore[receiptID] = receipt

	if err := json.NewEncoder(w).Encode(map[string]string{"id": receiptID}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	receipt, exists := receiptStore[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := calculatePoints(receipt)
	if err := json.NewEncoder(w).Encode(map[string]int{"points": points}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// calculatePoints applies all the defined rules to calculate total points for a receipt
func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: 1 point for each alphanumeric character in the retailer name
	for _, c := range receipt.Retailer {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			points++
		}
	}

	// Rule 2: 50 points if the total is a round dollar amount (no cents)
	if receipt.Total == math.Floor(receipt.Total) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: Points based on item description length
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			// Multiply the price by 0.2, round up to the nearest integer, and add to points
			itemPoints := int(math.Ceil(item.Price * 0.2))
			points += itemPoints
		}
	}

	// Rule 6: 6 points if the purchase day is odd
	if date, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		day := date.Day()
		if day%2 != 0 {
			points += 6
		}
	}

	// Rule 7: 10 points if the purchase time is between 2:00pm and 4:00pm
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		if purchaseTime.Hour() == 14 || (purchaseTime.Hour() == 15 && purchaseTime.Minute() < 60) {
			points += 10
		}
	}

	return points
}
