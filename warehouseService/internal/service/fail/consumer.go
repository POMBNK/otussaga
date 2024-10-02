package fail

import (
	"context"
	"github.com/POMBNK/warehouseService/internal/service/reservation"
	"github.com/POMBNK/warehouseService/pkg/client/mq"
	"github.com/streadway/amqp"
	"log/slog"
	"strconv"
)

type Consumer struct {
	reservation *reservation.Service
	rabbitMQ    *mq.RabbitMQ
}

func NewConsumer(rabbitMQ *mq.RabbitMQ, reservation *reservation.Service) *Consumer {
	return &Consumer{
		rabbitMQ:    rabbitMQ,
		reservation: reservation,
	}
}

func (c *Consumer) ProcessFailedReservations(ctx context.Context) {
	msgs, err := c.consume(ctx)
	if err != nil {
		slog.Error("consume error", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopped consuming failed orders")
			return
		case msg := <-msgs:
			if msg.Body == nil {
				continue
			}
			slog.Info("Received failed order message", "payload", string(msg.Body))
			orderID, err := strconv.Atoi(string(msg.Body))
			if err != nil {
				slog.Error("unmarshal error", err)
				continue
			}

			if err := c.reservation.RollbackReservation(ctx, orderID); err != nil {
				//send to failed topic to rollback
				slog.Error("rollback order error", err)
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
		"order_failed_to_warehouse_reservation",
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
		q.Name,         // queue name
		"",             // routing key
		"order_failed", // exchange
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
