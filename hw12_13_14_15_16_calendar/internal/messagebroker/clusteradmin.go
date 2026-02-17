package internalmessagebroker //nolint

import (
	"fmt"

	"github.com/Sierra93/hw-test/hw12_13_14_15_calendar/internal/app"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ClusterAdmin struct {
	logger app.Logger
	urls   []string

	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewClusterAdmin создаёт администратора для RabbitMQ.
// urls — список URL-адресов брокеров (например: "amqp://user:pass@host:5672/").
func NewClusterAdmin(logger app.Logger, urls []string) *ClusterAdmin {
	return &ClusterAdmin{
		logger: logger,
		urls:   urls,
	}
}

// connect устанавливает соединение и создаёт канал, если они ещё не созданы.
func (a *ClusterAdmin) connect() error {
	if a.conn != nil && a.channel != nil {
		if a.conn.IsClosed() {
			a.conn = nil
			a.channel = nil
		} else {
			return nil
		}
	}

	var lastErr error
	for _, url := range a.urls {
		conn, err := amqp.Dial(url)
		if err != nil {
			lastErr = err
			a.logger.Error("failed to connect to rabbitmq: %s", err.Error())
			continue
		}
		ch, err := conn.Channel()
		if err != nil {
			_ = conn.Close()
			lastErr = err
			a.logger.Error("failed to open channel: %s", err.Error())
			continue
		}
		a.conn = conn
		a.channel = ch
		a.logger.Info("Successfully connected to RabbitMQ at %s", url)
		return nil
	}
	return lastErr
}

// ensureQueueExists проверяет существование очереди и создает ее, если не существует
func (a *ClusterAdmin) ensureQueueExists(queue string) error {
	if a.channel == nil {
		return fmt.Errorf("channel is nil, cannot declare queue '%s'", queue)
	}
	_, err := a.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		a.logger.Error("failed to declare queue '%s': %s", queue, err.Error())
		return err
	}
	a.logger.Info("queue '%s' ensured", queue)
	return nil
}

// CreateQueue объявляет очередь с указанным именем.
// Заменено на ensureQueueExists для единообразия
func (a *ClusterAdmin) CreateQueue(queue string) error {
	return a.ensureQueueExists(queue)
}

// QueueExists проверяет, существует ли очередь с указанным именем.
func (a *ClusterAdmin) QueueExists(queue string) (bool, error) {
	if err := a.connect(); err != nil {
		return false, err
	}
	_, err := a.channel.QueueDeclarePassive(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		if amqpErr, ok := err.(*amqp.Error); ok && amqpErr.Code == amqp.NotFound {
			return false, nil
		}
		a.logger.Error("failed to check queue '%s' existence: %s", queue, err.Error())
		return false, err
	}
	return true, nil
}

// Close корректно закрывает соединение и канал.
func (a *ClusterAdmin) Close() error {
	var err error
	if a.channel != nil {
		if e := a.channel.Close(); e != nil {
			a.logger.Error("error closing channel: %s", e.Error())
			err = e
		} else {
			a.logger.Info("channel closed")
		}
		a.channel = nil
	}
	if a.conn != nil {
		if e := a.conn.Close(); e != nil {
			a.logger.Error("error closing connection: %s", e.Error())
			err = e
		} else {
			a.logger.Info("connection closed")
		}
		a.conn = nil
	}
	return err
}
