package kafkaTab

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	log    *zap.SugaredLogger
}

func NewProducer(broker, topic string, log *zap.SugaredLogger) *Producer {
	brokers := []string{broker}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{}, // Балансировщик для распределения сообщений по партициям (можно использовать другие: Hash, RoundRobin)
		AllowAutoTopicCreation: true,                // Авто создание топика
	}

	return &Producer{
		writer: writer, /*&kafka.Writer{
			Addr:                   kafka.TCP(broker),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
			// --- важные продакшн-настройки ---
			RequiredAcks: kafka.RequireAll, // дождаться репликации
			MaxAttempts:  5,                // до 5 попыток при ошибке
			BatchTimeout: 200 * time.Millisecond,
			BatchSize:    100,   // пакет до 100 сообщений
			Async:        false, // ждать подтверждения (без потерь)
			Transport: &kafka.Transport{
				//IdleConnTimeout:   30 * time.Second,
				MetadataTTL: time.Minute,
				ClientID:    "vado-producer",
				DialTimeout: 5 * time.Second,
				//WriteTimeout:      5 * time.Second,
				//ReadTimeout:       5 * time.Second,
				//MaxIdleConns:      10,
				//EnableIdleTimeout: true,
			},
		},*/
		log: log,
	}
}

// SendMessage — безопасная отправка сообщения с ретраями и логами
func (p *Producer) SendMessage(ctx context.Context, key, value []byte) error {
	//ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	//defer cancel()

	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	err := p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		p.log.Errorw("Kafka write error", "key", string(key), "err", err)
		return err
	}

	p.log.Debugw("Kafka message sent", "key", string(key), "size", len(value))
	return nil
}

// Close — корректно закрывает соединение
func (p *Producer) Close() error {
	p.log.Info("Closing Kafka producer...")
	return p.writer.Close()
}
