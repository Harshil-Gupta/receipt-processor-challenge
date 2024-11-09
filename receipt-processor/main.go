package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
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
	router := setupRouter()
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.HandleFunc("/", serveHomePage).Methods("GET")
	router.HandleFunc("/receipts/process", processReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

	return router
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String()
	receiptStore[receiptID] = receipt
	jsonResponse(w, map[string]string{"id": receiptID})
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	receipt, exists := receiptStore[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := calculatePoints(receipt)
	jsonResponse(w, map[string]int{"points": points})
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func calculatePoints(receipt Receipt) int {
	points := alphanumericPoints(receipt.Retailer)
	points += roundTotalPoints(receipt.Total)
	points += itemPoints(receipt.Items)
	points += datePoints(receipt.PurchaseDate)
	points += timePoints(receipt.PurchaseTime)
	return points
}

func alphanumericPoints(name string) int {
	points := 0
	for _, c := range name {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			points++
		}
	}
	return points
}

func roundTotalPoints(total float64) int {
	points := 0
	if total == math.Floor(total) {
		points += 50
	}
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}
	return points
}

func itemPoints(items []Item) int {
	points := 0
	points += (len(items) / 2) * 5
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}
	return points
}

func datePoints(dateStr string) int {
	points := 0
	date, err := time.Parse("2006-01-02", dateStr)
	if err == nil && date.Day()%2 != 0 {
		points += 6
	}
	return points
}

func timePoints(timeStr string) int {
	points := 0
	timeVal, err := time.Parse("15:04", timeStr)
	if err == nil && timeVal.Hour() >= 14 && timeVal.Hour() < 16 {
		points += 10
	}
	return points
}
