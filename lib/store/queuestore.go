package store

import "context"

// QueueWriter present a message sender
type QueueWriter interface {
	WriteMessage(ctx context.Context, message []byte) error
}

type QueueReaderFunc func(ctx context.Context, message []byte) error

// QueueReader present a message consumer
type QueueReader interface {
	Consume(ctx context.Context, fn QueueReaderFunc) error
}
