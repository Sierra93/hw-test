package internalmessagebroker

import (
	"context" // Уже импортируется
	"encoding/json"
	"time"

	"github.com/Sierra93/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Sierra93/hw-test/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go" //nolint:depguard
)

type Consumer struct {
	app        app.Application
	logger     app.Logger
	brokers    []string
	queue      string
	conn       *amqp.Connection // Подключение amqp091-go
	channel    *amqp.Channel    // Канал amqp091-go
	maxRetries int
}

func NewConsumer(app app.Application, logger app.Logger, brokers []string, queue string, maxRetries int) *Consumer {
	return &Consumer{
		app:        app,
		logger:     logger,
		brokers:    brokers,
		queue:      queue,
		maxRetries: maxRetries,
	}
}

// connect устанавливает соединение и канал с RabbitMQ.
func (c *Consumer) connect(ctx context.Context) error {
	var lastErr error
	retryDelay := time.Second

	connectionTimeout := 10 * time.Second

	for i := 0; i < c.maxRetries; i++ {
		if len(c.brokers) == 0 {
			return nil // Или вернуть ошибку, если нет брокеров
		}
		url := c.brokers[0]

		dialerWithTimeout := amqp.DefaultDial(connectionTimeout)

		conn, err := amqp.DialConfig(url, amqp.Config{
			Dial: dialerWithTimeout,
		})
		if err == nil {
			c.conn = conn
			ch, err := c.conn.Channel()
			if err != nil {
				lastErr = err
				c.logger.Error("failed to open a channel: %s", err.Error())
				conn.Close()
				goto retry
			}
			c.channel = ch
			lastErr = nil
			c.logger.Info("Successfully connected to RabbitMQ at %s", url)
			break
		}

		lastErr = err
		c.logger.Warn("failed to connect to RabbitMQ: %s, retrying in %s", err.Error(), retryDelay)

	retry:
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(retryDelay)
			retryDelay *= 2
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return nil
}

// ensureQueueExists проверяет существование очереди и создает ее, если не существует.
func (c *Consumer) ensureQueueExists() error {
	if c.channel == nil {
		return nil // Или ошибка, если канала нет
	}
	_, err := c.channel.QueueDeclare(
		c.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		c.logger.Error("failed to declare queue: %s", err.Error())
		return err
	}
	c.logger.Info("queue '%s' ensured", c.queue)
	return nil
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.connect(ctx); err != nil {
		return err
	}
	defer c.Close()

	if err := c.ensureQueueExists(); err != nil {
		return err
	}

	deliveries, err := c.channel.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("failed to register a consumer: %s", err.Error())
		return err
	}

	c.logger.Info("consumer started for queue '%s'", c.queue)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("context done, stopping consumer")
			return nil
		case d, ok := <-deliveries:
			if !ok {
				c.logger.Warn("delivery channel closed")
				// Используем context.Background().
				// Если нужен исходный контекст, передайте его, но Background надежнее для переподключения
				return c.Start(context.Background())
			}

			var event storage.Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				c.logger.Error("consumer failed to unmarshal message: %s", err.Error())
				_ = d.Ack(false)
				continue
			}

			if _, err := c.app.AddEventNotification(ctx, &event); err != nil {
				c.logger.Error("consumer failed to add notification: %s", err.Error())
				if requeueErr := d.Nack(false, true); requeueErr != nil {
					c.logger.Error("failed to nack message for requeue: %s", requeueErr.Error())
				}
				continue
			}

			if err := d.Ack(false); err != nil {
				c.logger.Error("failed to ack message: %s", err.Error())
			}
		}
	}
}

// Close закрывает соединение и канал с RabbitMQ.
func (c *Consumer) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
		c.channel = nil
		c.logger.Info("channel closed")
	}
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
		c.logger.Info("connection closed")
	}
}
