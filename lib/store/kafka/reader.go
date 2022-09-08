package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/duythinht/tg/lib/store"
	kafka "github.com/segmentio/kafka-go"
)

// QueueReader by kafka consumer
type QueueReader struct {
	*kafka.Reader
}

// NewQueueReader return a kafka consumer group
func NewQueueReader(brokers string, topic string, groupId string) *QueueReader {
	return &QueueReader{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: strings.Split(brokers, ","),
			GroupID: groupId,
			Topic:   topic,
		}),
	}
}

// Consume message and process
func (q *QueueReader) Consume(ctx context.Context, fn store.QueueReaderFunc) error {
	for {
		m, err := q.Reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("queue reader consume message %w", err)
		}
		err = fn(ctx, m.Value)
		if err != nil {
			log.Printf("error process message: %s", err.Error())
		}
	}
}
