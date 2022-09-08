package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/duythinht/tg/api"
	"github.com/duythinht/tg/lib/store/kafka"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:x@localhost:5432/postgres")

	if err != nil {
		log.Fatal(err)
	}

	queue := kafka.NewQueueWriter("localhost:9092", "scan")

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

/*
func main() {
	repo := github.OpenRepository("duythinht", "zhttp")

	gfs, err := repo.OpenFS(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fs.WalkDir(gfs, "", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fmt.Printf("path: `%s`\n", path)
		f, err := gfs.Open(path)
		if err != nil {
			return err
		}

		stat, _ := f.Stat()

		fmt.Printf("name: %#v\n", stat.Name())

		// reports, err := r.Check(f)

		// if err != nil {
		// 	return err
		// }

		// //fmt.Println(reports)
		return nil
	})
}

*/
