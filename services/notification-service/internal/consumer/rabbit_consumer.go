package consumer

import (
	"context"
	"fmt"
	"log"

	"notification-service/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	ch         *amqp.Channel
	queue      string
	exchange   string
	routingKey string
	svc        service.NotificationService
}

func NewRabbitConsumer(ch *amqp.Channel, queue, exchange, routingKey string, svc service.NotificationService) RabbitConsumer {
	return RabbitConsumer{
		ch:         ch,
		queue:      queue,
		exchange:   exchange,
		routingKey: routingKey,
		svc:        svc,
	}
}

func (c RabbitConsumer) Start(ctx context.Context) error {
	if err := c.setup(); err != nil {
		return err
	}
	msgs, err := c.ch.Consume(c.queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				if err := c.svc.HandleTransactionCompleted(ctx, msg.Body); err != nil {
					log.Printf("handle event failed: %v", err)
					_ = msg.Nack(false, true)
					continue
				}
				_ = msg.Ack(false)
			}
		}
	}()
	return nil
}

func (c RabbitConsumer) setup() error {
	if err := c.ch.ExchangeDeclare(c.exchange, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}
	_, err := c.ch.QueueDeclare(c.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}
	if err := c.ch.QueueBind(c.queue, c.routingKey, c.exchange, false, nil); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}
	return nil
}
