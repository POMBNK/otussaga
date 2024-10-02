package order

import (
	"context"
	"encoding/json"
	"github.com/POMBNK/paymentService/internal/entity"
	"github.com/POMBNK/paymentService/internal/service/payment"
	"github.com/POMBNK/paymentService/pkg/client/mq"
	"github.com/streadway/amqp"
	"log/slog"
)

type Consumer struct {
	reservationService *payment.Service
	rabbitMQ           *mq.RabbitMQ
}

func NewConsumer(rabbitMQ *mq.RabbitMQ, reservationService *payment.Service) *Consumer {
	return &Consumer{
		rabbitMQ:           rabbitMQ,
		reservationService: reservationService,
	}
}

func (c *Consumer) ProcessIncomePayment(ctx context.Context) {
	msgs, err := c.consume(ctx)
	if err != nil {
		slog.Error("consume error", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopped consuming reservations")
			return
		case msg := <-msgs:
			if len(msg.Body) == 0 {
				continue
			}
			slog.Info("Received message", "payload", string(msg.Body))
			var o entity.OrderEvent
			err := json.Unmarshal(msg.Body, &o.Payload)
			if err != nil {
				slog.Error("unmarshal error", err)
				continue
			}
			reserv := entity.Payment{
				OrderID: o.Payload.ID,
			}

			exist, err := c.reservationService.IsPaymentAlreadyExist(ctx, reserv)
			if err != nil {
				slog.Error("is reservation already exist error", err)
			}

			if exist {
				slog.Info("reservation already exist")
				err = msg.Ack(false)
				if err != nil {
					slog.Error("ack error", err)
				}
				continue
			}

			if _, err := c.reservationService.CreatePayment(ctx, reserv); err != nil {
				//send to failed topic to rollback
				err = c.reservationService.EventSender.SendPaymentFailed(ctx, reserv.OrderID)
				if err != nil {
					slog.Error("send to failed topic error", err)
				}
				slog.Error("create payment error", err)
			}

			err = msg.Ack(false)
			if err != nil {
				slog.Error("ack error", err)
			}
		default:
		}
	}
}

func (c *Consumer) consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	q, err := c.rabbitMQ.Channel.QueueDeclare(
		"orders_to_payment",
		c.rabbitMQ.Cfg.ExchangeSettigns.Durable,
		c.rabbitMQ.Cfg.ExchangeSettigns.AutoDelete,
		c.rabbitMQ.Cfg.ExchangeSettigns.Internal,
		c.rabbitMQ.Cfg.ExchangeSettigns.NoWait,
		c.rabbitMQ.Cfg.ExchangeSettigns.Args,
	)
	if err != nil {
		return nil, err
	}

	if err = c.rabbitMQ.Channel.Qos(1, 0, false); err != nil {
		return nil, err
	}
	err = c.rabbitMQ.Channel.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		"orders", // exchange
		false,
		nil,
	)
	msgs, err := c.rabbitMQ.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
