package order

import (
	"context"
	"encoding/json"
	"github.com/POMBNK/warehouseService/internal/entity"
	"github.com/POMBNK/warehouseService/internal/service/reservation"
	"github.com/POMBNK/warehouseService/pkg/client/mq"
	"github.com/streadway/amqp"
	"log/slog"
)

type Consumer struct {
	reservationService *reservation.Service
	rabbitMQ           *mq.RabbitMQ
}

func NewConsumer(rabbitMQ *mq.RabbitMQ, reservationService *reservation.Service) *Consumer {
	return &Consumer{
		rabbitMQ:           rabbitMQ,
		reservationService: reservationService,
	}
}

func (c *Consumer) ProcessIncomeReservation(ctx context.Context) {
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
			reserv := entity.Reservation{
				OrderID: o.Payload.ID,
				Goods:   o.Payload.Goods,
			}

			exist, err := c.reservationService.IsReservationAlreadyExist(ctx, reserv)
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

			_, err = c.reservationService.CreateReservation(ctx, reserv)

			if err != nil || reserv.OrderID%2 == 0 {
				//send to failed topic to rollback
				err = c.reservationService.EventSender.SendReservationFailed(ctx, reserv.OrderID)
				if err != nil {
					slog.Error("send to failed topic error", err)
				}
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
		"orders_to_ware",
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
