package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./receipt-processor/static"))))
	router.HandleFunc("/", serveHomePage).Methods("GET")
	router.HandleFunc("/receipts/process", processReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")
	router.HandleFunc("/receipts/{id}", updateById).Methods("PUT")

	return router
}
