package main

import (
	"context"
	reservationRepo "github.com/POMBNK/warehouseService/internal/repository/reservation"
	"github.com/POMBNK/warehouseService/internal/repository/tx"
	"github.com/POMBNK/warehouseService/internal/service/fail"
	"github.com/POMBNK/warehouseService/internal/service/order"
	"github.com/POMBNK/warehouseService/internal/service/reservation"
	"github.com/POMBNK/warehouseService/pkg/client/mq"
	"github.com/POMBNK/warehouseService/pkg/client/postgres"
)

func main() {
	ctx := context.Background()
	dbConn, err := postgres.NewClient(ctx, postgres.Cfg{
		Login:       "pombnk",
		Password:    "postgres",
		Host:        "postgres",
		Port:        "5432",
		Database:    "warehousedb",
		MaxAttempts: 3,
	})
	if err != nil {
		panic(err)
	}

	rabbitConn, err := mq.NewRabbitMQ(ctx, mq.Cfg{
		MaxAttempts:      3,
		User:             "rmuser",
		Password:         "rmpassword",
		Host:             "localhost",
		Port:             "5672",
		Vhost:            "/",
		ExchangeSettigns: mq.DefaultExchangeSettings("warehouseExchange", "topic"),
	})
	if err != nil {
		panic(err)
	}

	txManager := tx.New(dbConn)
	//event
	reservRepo := reservationRepo.NewRepo(dbConn)
	reservSender := reservation.NewSender(rabbitConn)
	reservSvc := reservation.NewService(reservRepo, txManager, reservSender)

	//consumers block
	go order.NewConsumer(rabbitConn, reservSvc).ProcessIncomeReservation(ctx)
	fail.NewConsumer(rabbitConn, reservSvc).ProcessFailedReservations(ctx)

}
