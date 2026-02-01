package queue

import (
	_ "NotificationService/internal/app"
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitQueue struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
	msg   amqp.Delivery
}

func NewRabbitQueue() *RabbitQueue {
	conn := newConn()
	chMq, err := conn.Channel()
	failOnError(err, "Failed to create Channel")
	queue, err := chMq.QueueDeclare("EventKeyQueue", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	return &RabbitQueue{
		conn:  conn,
		ch:    chMq,
		queue: queue.Name,
	}
}

func newConn() *amqp.Connection {
	conn, err := amqp.Dial("")
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func (q *RabbitQueue) CloseConnection() {
	err := q.conn.Close()
	failOnError(err, "Failed to close connection")
}

func (q *RabbitQueue) CloseChannel() {
	err := q.ch.Close()
	failOnError(err, "Failed to close channel")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (q *RabbitQueue) PublishEvent(ctx context.Context, key string) error {	
	err := q.ch.PublishWithContext(ctx, "", q.queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(key),
	})
	failOnError(err, "Failed to publish a message")
	return nil
}

func (q *RabbitQueue) ConsumeEvent(ctx context.Context) (string, error) {
	msg, err := q.ch.ConsumeWithContext(ctx, q.queue, "", false, false, false, false, nil)		
	failOnError(err, "Failed to consume event")
	q.msg = <-msg
	return string(q.msg.Body), nil
}
