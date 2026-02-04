package queue

import (
	_ "NotificationService/internal/app"
	"context"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitQueue struct {
	ch         *amqp.Channel
	queue      string
	deliveries <-chan amqp.Delivery
	msg        amqp.Delivery
}

func NewRabbitQueue(ch *amqp.Channel, queue string) (*RabbitQueue, error) {
	if err := ch.Qos(1, 0, false); err != nil {
		return nil, err
	}
	msgs, err := ch.Consume(
		queue,
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
	return &RabbitQueue{
		ch:         ch,
		queue:      queue,
		deliveries: msgs,
	}, nil
}

func newConn(connection string) *amqp.Connection {
	conn, err := amqp.Dial(connection)
	if err != nil {
		return nil
	}
	return conn
}

func (q *RabbitQueue) CloseConnection(conn *amqp.Connection) {
	err := conn.Close()
	if err != nil {
		log.Printf("Failed to close connection")
	}
}

func (q *RabbitQueue) CloseChannel() {
	err := q.ch.Close()
	if err != nil {
		log.Printf("Failed to close channel")
	}
}

func (q *RabbitQueue) PublishEvent(ctx context.Context, key string) error {
	err := q.ch.PublishWithContext(ctx, "", q.queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(key),
	})
	if err != nil {
		return errors.New("Failed to publish event")
	}
	return nil
}

func (q *RabbitQueue) ConsumeEvent(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case msg, ok := <-q.deliveries:
		if !ok {
			return "", errors.New("delivery channel closed")
		}
		q.msg = msg
		return string(msg.Body), nil
	}
}

func (q *RabbitQueue) AckEvent() error {
	if q.msg.DeliveryTag == 0 {
		return errors.New("No message to Ack")
	}
	err := q.msg.Ack(false)
	if err != nil {
		return err
	}
	return nil
}

func (q *RabbitQueue) NackEvent() error {
	if q.msg.DeliveryTag == 0 {
		return errors.New("No message to nack")
	}
	err := q.msg.Nack(false, true)
	if err != nil {
		return err
	}
	return nil
}
