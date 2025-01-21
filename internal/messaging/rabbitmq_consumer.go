package messaging

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitMQConsumer(conn *amqp091.Connection) (*RabbitMQConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *RabbitMQConsumer) StartConsuming(queueName string, handler func([]byte) error) error {
	q, err := c.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			err := handler(msg.Body)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack(false, true) // negative acknowledgement, requeue
			} else {
				msg.Ack(false) // positive acknowledgement
			}
		}
	}()

	return nil
}

func (c *RabbitMQConsumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}
