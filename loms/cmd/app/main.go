package main

import (
	"log"
	"net/http"
	srvwrapper "route256/libs/server-wrapper"
	"route256/loms/internal/config"
	cancelorder "route256/loms/internal/handlers/cancel-order"
	createorder "route256/loms/internal/handlers/create-order"
	listorder "route256/loms/internal/handlers/list-order"
	orderpayed "route256/loms/internal/handlers/order-payed"
	"route256/loms/internal/handlers/stocks"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	stocksHandler := stocks.New()
	createOrderHandler := createorder.New()
	listOrderHandler := listorder.New()
	orderPayedHandler := orderpayed.New()
	cancelOrderHandler := cancelorder.New()

	http.Handle("/stocks", srvwrapper.New(stocksHandler.Handle))
	http.Handle("/createOrder", srvwrapper.New(createOrderHandler.Handle))
	http.Handle("/listOrder", srvwrapper.New(listOrderHandler.Handle))
	http.Handle("/orderPayed", srvwrapper.New(orderPayedHandler.Handle))
	http.Handle("/cancelOrder", srvwrapper.New(cancelOrderHandler.Handle))

	log.Println("listening http at", config.ConfigData.Services.Loms.Port)
	err = http.ListenAndServe(config.ConfigData.Services.Loms.Port, nil)
	log.Fatal("cannot listen http", err)
}
