package main

import (
	"context"
	eventrepo "github.com/POMBNK/orderservice/internal/repository/event"
	orderrepo "github.com/POMBNK/orderservice/internal/repository/order"
	"github.com/POMBNK/orderservice/internal/repository/tx"
	"github.com/POMBNK/orderservice/internal/service/event"
	"github.com/POMBNK/orderservice/internal/service/fail"
	orderservice "github.com/POMBNK/orderservice/internal/service/order"
	orderhandler "github.com/POMBNK/orderservice/internal/transport/http/order"
	"github.com/POMBNK/orderservice/pkg/client/mq"
	"github.com/POMBNK/orderservice/pkg/client/postgres"
	_ "github.com/POMBNK/orderservice/statik"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	dbConn, err := postgres.NewClient(ctx, postgres.Cfg{
		Login:       "pombnk",
		Password:    "postgres",
		Host:        "postgres",
		Port:        "5432",
		Database:    "orderdb",
		MaxAttempts: 3,
	})
	if err != nil {
		panic(err)
	}

	rabbitConn, err := mq.NewRabbitMQ(ctx, mq.Cfg{
		MaxAttempts:      3,
		User:             "rmuser",
		Password:         "rmpassword",
		Host:             "definition",
		Port:             "5672",
		Vhost:            "/",
		ExchangeSettigns: mq.DefaultExchangeSettings("orders", "fanout"),
	})
	if err != nil {
		panic(err)
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/swagger/", http.FileServer(statikFS)).ServeHTTP(w, r)
	})
	txManager := tx.New(dbConn)
	//event
	eventRepo := eventrepo.NewRepo(dbConn)
	eventSender := event.NewSender(eventRepo, rabbitConn)
	eventSvc := event.NewService(eventRepo)

	//order
	storage := orderrepo.NewRepo(dbConn)
	svc := orderservice.NewService(storage, eventSvc, eventSender, txManager)
	orderHandler := orderhandler.NewServer(svc)

	//worker
	//eventSender.StartProcessingEvents(ctx, 1*time.Second)
	go fail.NewConsumer(rabbitConn, svc).ProcessFailedOrders(ctx)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", orderHandler.Register(mux, "/api/v1")))
}
