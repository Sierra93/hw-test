package internalmessagebroker

import (
	"encoding/json"
	"fmt"

	"github.com/Sierra93/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Sierra93/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/rabbitmq/amqp091-go" //nolint:depguard
)

type Producer struct {
	logger app.Logger
	urls   []string // список URL брокеров RabbitMQ
	queue  string
	admin  *ClusterAdmin
}

// NewProducer создаёт продюсера для RabbitMQ.
// urls — список URL брокеров (например: "amqp://user:pass@host:5672/").
// queue — имя очереди, в которую будет публиковаться сообщение.
func NewProducer(logger app.Logger, urls []string, queue string) *Producer {
	return &Producer{
		logger: logger,
		urls:   urls,
		queue:  queue,
		admin:  NewClusterAdmin(logger, urls),
	}
}

func (s *Producer) Send(event *storage.Event) error {
	// Убедимся, что админ и соединение готовы
	if s.admin == nil {
		s.admin = NewClusterAdmin(s.logger, s.urls)
	}
	// Подключение (ClusterAdmin имеет неэкспортированное connect; т.к. Producer в той же папке,
	// мы можем вызывать его напрямую)
	if err := s.admin.connect(); err != nil {
		return err
	}

	// Убедимся, что очередь существует
	if err := s.admin.CreateQueue(s.queue); err != nil {
		return err
	}

	// Сериализация события
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("sender failed to marshal event: %w", err)
	}

	// Публикация в очередь через дефолтный обменник
	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        eventJSON,
	}

	if err := s.admin.channel.Publish(
		"",      // exchange
		s.queue, // routing key = имя очереди
		false,   // mandatory
		false,   // immediate
		msg,
	); err != nil {
		return fmt.Errorf("sender failed to publish event: %w", err)
	}

	s.logger.Debug("Event [%+v] sent to RabbitMQ queue %s", event, s.queue)
	return nil
}
