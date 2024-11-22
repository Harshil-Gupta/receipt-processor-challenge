package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := setupRouter()
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
