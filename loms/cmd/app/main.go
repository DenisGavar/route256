package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

func main() {
	stocksHandler := stocks.New()

	http.Handle("/stocks", srvwrapper.New(stocksHandler.Handle))

	log.Println("listening http at", port)
	err := http.ListenAndServe(port, nil)
	log.Fatal("cannot listen http", err)
}
