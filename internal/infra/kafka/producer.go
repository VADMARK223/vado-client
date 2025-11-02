package kafka

import (
	"context"
	"time"
	"vado-client/internal/config/port"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	log    *zap.SugaredLogger
}

func NewProducer(topic string, log *zap.SugaredLogger) *Producer {
	brokers := []string{"localhost:" + port.Kafka}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{}, // Балансировщик для распределения сообщений по партициям (можно использовать другие: Hash, RoundRobin)
		AllowAutoTopicCreation: true,                // Авто создание топика
	}

	return &Producer{
		writer: writer,
		log:    log,
	}
}

// SendMessage — безопасная отправка сообщения с ретраями и логами
func (p *Producer) SendMessage(key, value []byte) error {
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

	topic := p.writer.Topic
	p.log.Debugw("Kafka message sent", "topic", topic, "key", string(key), "size", len(value))
	return nil
}

// Close — корректно закрывает соединение
func (p *Producer) Close() error {
	p.log.Info("Closing Kafka producer...")
	return p.writer.Close()
}
