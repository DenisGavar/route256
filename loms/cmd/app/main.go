package main

import (
	"log"
	"net/http"
	serverWrapper "route256/libs/server-wrapper"
	"route256/loms/internal/config"
	cancelOrder "route256/loms/internal/handlers/cancel-order"
	createOrder "route256/loms/internal/handlers/create-order"
	listOrder "route256/loms/internal/handlers/list-order"
	orderPayed "route256/loms/internal/handlers/order-payed"
	"route256/loms/internal/handlers/stocks"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	stocksHandler := stocks.New()
	createOrderHandler := createOrder.New()
	listOrderHandler := listOrder.New()
	orderPayedHandler := orderPayed.New()
	cancelOrderHandler := cancelOrder.New()

	http.Handle("/stocks", serverWrapper.New(stocksHandler.Handle))
	http.Handle("/createOrder", serverWrapper.New(createOrderHandler.Handle))
	http.Handle("/listOrder", serverWrapper.New(listOrderHandler.Handle))
	http.Handle("/orderPayed", serverWrapper.New(orderPayedHandler.Handle))
	http.Handle("/cancelOrder", serverWrapper.New(cancelOrderHandler.Handle))

	log.Println("listening http at", config.ConfigData.Services.Loms.Port)
	err = http.ListenAndServe(config.ConfigData.Services.Loms.Port, nil)
	log.Fatal("cannot listen http", err)
}
