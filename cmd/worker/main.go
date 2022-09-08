package main

import (
	"context"
	"log"

	"github.com/duythinht/tg/lib/store/kafka"
	"github.com/duythinht/tg/worker"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:x@localhost:5432/postgres")

	if err != nil {
		log.Fatal(err)
	}

	queue := kafka.NewQueueReader("localhost:9092", "scan", "scan-secrets-worker")

	w := worker.New(db, queue)

	log.Fatal(w.Run(context.Background()))
}
