package fail

import (
	"context"
	"github.com/POMBNK/paymentService/internal/service/payment"
	"github.com/POMBNK/paymentService/pkg/client/mq"
	"github.com/streadway/amqp"
	"log/slog"
	"strconv"
)

type Consumer struct {
	orderService *payment.Service
	rabbitMQ     *mq.RabbitMQ
}

func NewConsumer(rabbitMQ *mq.RabbitMQ, orderService *payment.Service) *Consumer {
	return &Consumer{
		rabbitMQ:     rabbitMQ,
		orderService: orderService,
	}
}

func (c *Consumer) ProcessFailedPayments(ctx context.Context) {
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

			if err := c.orderService.RollbackPayment(ctx, orderID); err != nil {
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
		"order_failed_to_payments",
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
