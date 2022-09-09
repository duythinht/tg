package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"strings"
	"time"

	"github.com/duythinht/tg/lib/repository/github"
	"github.com/duythinht/tg/lib/rule"
	"github.com/duythinht/tg/lib/rule/secrets"
	"github.com/duythinht/tg/lib/store"
)

type ScanRequest struct {
	ScanID       int64  `json:"scan_id"`
	RepositoryID int64  `json:"repository_id"`
	Owner        string `json:"owner"`
	Repository   string `json:"repository"`
}

type Worker struct {
	queue store.QueueReader
	db    store.DB
	rules []rule.Rule
}

func New(db store.DB, reader store.QueueReader) *Worker {
	return &Worker{
		db:    db,
		queue: reader,

		//currenly, only use 1 rule is secrets rules
		rules: []rule.Rule{
			&secrets.Rule{},
		},
	}
}

func (w *Worker) Run(ctx context.Context) error {
	return w.queue.Consume(ctx, func(ctx context.Context, message []byte) error {
		log.Printf("start scan for request: %s", message)

		var req ScanRequest
		err := json.Unmarshal(message, &req)

		if err != nil {
			return fmt.Errorf("unmarshal scan request %w", err)
		}

		_, err = w.db.Exec(ctx, "UPDATE scans SET status=$1, scanning_at=$2 WHERE id=$3", "In Progress", time.Now(), req.ScanID)

		if err != nil {
			return fmt.Errorf("update scan status %w", err)
		}

		repo := github.OpenRepository(req.Owner, req.Repository)

		rfs, err := repo.OpenFS(ctx)

		if err != nil {
			return fmt.Errorf("open fs %w", err)
		}

		result := make([]rule.Report, 0)

		err = fs.WalkDir(rfs, "", func(path string, d fs.DirEntry, err1 error) error {

			// skip if error before this path
			if err1 != nil {
				return err1
			}

			// skip dir
			if d.IsDir() {
				return nil
			}

			f, err := rfs.Open(path)
			if err != nil {
				return fmt.Errorf("worker open repo FS %w", err)
			}

			for _, r := range w.rules {

				log.Printf("rule %s is scanning for %s...", r.Name(), path)

				reports, err := r.Check(f)
				if err != nil {
					return fmt.Errorf("rule check error %w", err)
				}

				for _, report := range reports {
					// due to zipfs start as pattern <repo-name>-<hash>/files...
					// we need to trim first part of folder to make sure location path is align with source tree
					report.Location.Path = strings.SplitN(path, "/", 2)[1]
					result = append(result, report)
				}

			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("walk github zipfs %w", err)
		}

		if len(result) < 1 {
			_, err = w.db.Exec(ctx, "UPDATE scans SET status=$1, finished_at=$2 WHERE id=$3", "Success", time.Now(), req.ScanID)
			return err
		}

		_, err = w.db.Exec(ctx, "UPDATE scans SET status=$1, finished_at=$2, findings=$3 WHERE id=$4", "Failure", time.Now(), result, req.ScanID)

		if err != nil {
			return fmt.Errorf("update scan failure status %w", err)
		}
		return nil
	})
}
