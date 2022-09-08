package kafka

import (
	"context"
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

// QueueWriter could use to produce message into kafka
type QueueWriter struct {
	*kafka.Writer
}

// NewQueueWriter start a Kafka producder
func NewQueueWriter(brokers string, topic string) *QueueWriter {
	return &QueueWriter{
		&kafka.Writer{
			Addr:                   kafka.TCP(strings.Split(brokers, ",")...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

// WriteMessage write a message
func (q *QueueWriter) WriteMessage(ctx context.Context, message []byte) error {
	return q.Writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}
