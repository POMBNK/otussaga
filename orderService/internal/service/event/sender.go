package event

import (
	"context"
	"fmt"
	"github.com/POMBNK/orderservice/internal/entity"
	"github.com/POMBNK/orderservice/pkg/client/mq"
	"github.com/streadway/amqp"
	"log"
	"log/slog"
	"time"
)

type Sender struct {
	repo   EventRepo
	rabbit *mq.RabbitMQ
}

func NewSender(repo EventRepo, rabbit *mq.RabbitMQ) *Sender {
	return &Sender{repo: repo, rabbit: rabbit}
}

func (s *Sender) StartProcessingEvents(ctx context.Context, handlePeriod time.Duration) {
	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Stopped processing events")
				return
			case <-ticker.C:
			}

			events, err := s.repo.GetNewEvent(ctx)
			if err != nil {
				log.Println(err)
				continue
			}

			if events.Info.ID == 0 {
				continue
			}

			err = s.SendNewEvent(ctx, events)
			if err != nil {

			}

			eventIDs := []int{events.Info.ID}
			if err := s.repo.SetCompleted(ctx, eventIDs); err != nil {
				//log.Println(err)
				continue
			}

		}
	}()
}

func (s *Sender) SendEvents(ctx context.Context, event []entity.Event) error {
	for _, e := range event {
		slog.Info("Sending event", "eventInfo", e.Info, "eventPayload", string(e.Payload))
		err := s.publish(ctx, string(entity.NewOrderEventType), e.Payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Sender) SendNewEvent(ctx context.Context, event entity.Event) error {
	slog.Info("Sending event", "eventInfo", event.Info, "eventPayload", string(event.Payload))
	err := s.publish(ctx, string(entity.NewOrderEventType), event.Payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) SendFailedEvents(ctx context.Context, event []entity.Event) error {
	for _, e := range event {
		slog.Info("Sending event", "eventInfo", e.Info, "eventPayload", string(e.Payload))
		err := s.publish(ctx, string(entity.FailedOrderEventType), e.Payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Sender) publish(ctx context.Context, routingKey string, payload []byte) error {

	err := s.rabbit.Channel.Publish(
		s.rabbit.Cfg.ExchangeSettigns.Name, // exchange
		"",                                 // routing key
		false,                              // mandatory
		false,                              // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			AppId:        "order-rest-server",
			ContentType:  "application/json",
			Body:         payload,
			Timestamp:    time.Now(),
		})

	if err != nil {
		return fmt.Errorf("err due publishing event: %w", err)
	}

	return nil
}
