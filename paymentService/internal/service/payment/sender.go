package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/POMBNK/paymentService/internal/entity"
	"github.com/POMBNK/paymentService/pkg/client/mq"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

type Sender struct {
	rabbit *mq.RabbitMQ
}

func NewSender(rabbit *mq.RabbitMQ) *Sender {
	return &Sender{rabbit: rabbit}
}

func (s *Sender) SendPaymentEvent(ctx context.Context, payment entity.Payment) error {
	payload, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	err = s.publish(ctx, "paymentExchange", "application/json", payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) SendPaymentFailed(ctx context.Context, orderID int) error {
	err := s.publish(ctx, "order_failed", "text/plain", []byte(strconv.Itoa(orderID)))
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) publish(ctx context.Context, exchange string, contentType string, payload []byte) error {
	err := s.initFailChannel()
	if err != nil {
		return err
	}

	err = s.rabbit.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("err due qos: %w", err)
	}

	err = s.rabbit.Channel.Publish(
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			AppId:        "payment-rest-server",
			ContentType:  contentType,
			Body:         payload,
			Timestamp:    time.Now(),
		})

	if err != nil {
		return fmt.Errorf("err due publishing event: %w", err)
	}

	return nil
}

func (s *Sender) initFailChannel() error {
	err := s.rabbit.Channel.ExchangeDeclare(
		"order_failed", // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("err due declaring exchange: %w", err)
	}

	return nil
}
