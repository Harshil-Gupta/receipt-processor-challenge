package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || receipt.Total == "" || len(receipt.Items) == 0 {
		http.Error(w, `{"error": "Missing required fields in the receipt."}`, http.StatusBadRequest)
		return
	}

	for _, item := range receipt.Items {
		if item.Price == "" || item.ShortDescription == "" {
			http.Error(w, `{"error": "Missing required fields in items"}`, http.StatusBadRequest)
			return
		}
	}

	_, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		http.Error(w, `{"error": "The purchase time field is incorrect in the receipt."}`, http.StatusBadRequest)
		return
	}

	var totalCalculated float64
	for _, item := range receipt.Items {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			http.Error(w, "Invalid price format", http.StatusBadRequest)
			return
		}
		totalCalculated += price
	}

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		http.Error(w, "Invalid total price format", http.StatusBadRequest)
		return
	}

	if math.Abs(totalCalculated-total) > 0.01 {
		http.Error(w, "Incorrect JSON: Prices do not match the total price", http.StatusBadRequest)
		return
	}

	receipt.Points = calculatePoints(receipt)
	receiptID := generateReceiptID()
	receiptStore[receiptID] = receipt

	jsonResponse(w, map[string]string{"id": receiptID})
}

func updateById(w http.ResponseWriter, r *http.Request) {
	var updatedReceipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&updatedReceipt); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	receiptId := mux.Vars(r)["id"]
	for receipts := range receiptStore {
		if receiptId == receipts {
			receiptStore[receiptId] = updatedReceipt
		}
	}
	jsonResponse(w, updatedReceipt)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	receipt, exists := receiptStore[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"receipt": receipt,
		"points":  receipt.Points,
	}
	jsonResponse(w, response)
}
