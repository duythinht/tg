package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/duythinht/tg/api"
	"github.com/duythinht/tg/config"
	"github.com/duythinht/tg/lib/store/kafka"
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

	queue := kafka.NewQueueWriter(cfg.Queue.Brokers, cfg.Queue.Topic)

	handler := api.New(db, queue)

	if err != nil {
		log.Fatal(err)
	}

	s := http.Server{
		Addr:        ":8080",
		Handler:     handler,
		ReadTimeout: 60 * time.Second, // customize http.Server timeouts
	}

	log.Fatal(s.ListenAndServe())
}
