package main

import (
	"github.com/google/uuid"
)

var receiptStore = make(map[string]Receipt)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
	Points       int    `json:"-"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func generateReceiptID() string {
	return uuid.New().String()
}
