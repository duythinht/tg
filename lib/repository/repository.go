package repository

import (
	"context"
	"io/fs"
)

// Repository presentation of git repository
type Repository interface {
	OpenFS(ctx context.Context, ref string) (fs.FS, error)
}
