package ports

type MessageConsumer interface {
	StartConsuming(queueName string, handler func([]byte) error) error
	Close() error
}

type MessageHandler interface {
	HandleMessage([]byte) error
}
