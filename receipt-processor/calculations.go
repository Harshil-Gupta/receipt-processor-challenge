package main

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

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

func roundTotalPoints(totalStr string) int {
	points := 0
	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return points
	}
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
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			continue
		}
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int(math.Ceil(price * 0.2))
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
