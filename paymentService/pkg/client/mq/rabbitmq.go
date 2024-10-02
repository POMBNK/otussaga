package mq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQ ...
type RabbitMQ struct {
	Cfg        *Cfg
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewRabbitMQ instantiates the RabbitMQ instances using configuration defined in environment variables.
func NewRabbitMQ(ctx context.Context, cfg Cfg) (*RabbitMQ, error) {
	if cfg.Vhost == "/" {
		cfg.Vhost = ""
	}
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Vhost)
	log.Println(rabbitURL)
	var conn *amqp.Connection
	err := again(func() error {
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var err error
		conn, err = amqp.Dial(rabbitURL)
		if err != nil {
			return fmt.Errorf("amqp.Dial %w", err)
		}

		return nil
	}, cfg.MaxAttempts, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("amqp tries limit exceeded")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("conn.Channel %w", err)
	}

	err = ch.ExchangeDeclare(
		cfg.ExchangeSettigns.Name, // name
		cfg.ExchangeSettigns.Kind, // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("ch.ExchangeDeclare %w", err)
	}

	if err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, fmt.Errorf("ch.Qos %w", err)
	}

	// XXX: Dead Letter Exchange will be implemented in future episodes

	return &RabbitMQ{
		Cfg:        &cfg,
		Connection: conn,
		Channel:    ch,
	}, nil
}

// Close ...
func (r *RabbitMQ) Close() {
	r.Connection.Close()
}
