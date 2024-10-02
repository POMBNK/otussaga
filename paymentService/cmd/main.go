package main

import (
	"context"
	paymentrepo "github.com/POMBNK/paymentService/internal/repository/payment"
	"github.com/POMBNK/paymentService/internal/service/fail"
	"github.com/POMBNK/paymentService/internal/service/order"
	"github.com/POMBNK/paymentService/internal/service/payment"
	"github.com/POMBNK/paymentService/pkg/client/mq"
	"github.com/POMBNK/paymentService/pkg/client/postgres"
)

func main() {
	ctx := context.Background()
	dbConn, err := postgres.NewClient(ctx, postgres.Cfg{
		Login:       "pombnk",
		Password:    "postgres",
		Host:        "localhost",
		Port:        "5432",
		Database:    "paymentdb",
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
		ExchangeSettigns: mq.DefaultExchangeSettings("paymentExchange", "topic"),
	})
	if err != nil {
		panic(err)
	}
	//payment block
	paymentRepo := paymentrepo.NewRepo(dbConn)
	paymentEventSender := payment.NewSender(rabbitConn)
	paymentService := payment.NewService(paymentRepo, paymentEventSender)

	// orderconsumer block
	go order.NewConsumer(rabbitConn, paymentService).ProcessIncomePayment(ctx)
	fail.NewConsumer(rabbitConn, paymentService).ProcessFailedPayments(ctx)

}
