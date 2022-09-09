package main

import (
	"context"
	"log"

	"github.com/duythinht/tg/config"
	"github.com/duythinht/tg/lib/store/kafka"
	"github.com/duythinht/tg/worker"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.Connect(context.Background(), cfg.DB.DSN)

	if err != nil {
		log.Fatal(err)
	}

	queue := kafka.NewQueueReader(cfg.Queue.Brokers, cfg.Queue.Topic, cfg.Queue.GroupID)

	w := worker.New(db, queue)

	log.Printf("Start scanning worker")

	log.Fatal(w.Run(context.Background()))
}
