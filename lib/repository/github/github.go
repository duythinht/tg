package github

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/duythinht/tg/lib/zipfs"
	"github.com/google/go-github/v47/github"
)

// Repository of github
type Repository struct {
	*github.Client
	Owner string
	Repo  string
}

// OpenRepository return abstraction of github repo
func OpenRepository(owner, repo string) *Repository {
	return &Repository{
		Client: github.NewClient(nil),
		Owner:  owner,
		Repo:   repo,
	}
}

// OpenFS create zipfs for github archive
func (r *Repository) OpenFS(ctx context.Context) (fs.FS, error) {
	url, _, err := r.Repositories.GetArchiveLink(ctx, r.Owner, r.Repo, github.Zipball, &github.RepositoryContentGetOptions{}, true)

	if err != nil {
		return nil, fmt.Errorf("get github archive link %w", err)
	}

	archive, err := http.Get(url.String())

	if err != nil {
		return nil, fmt.Errorf("read github archive %w", err)
	}

	defer archive.Body.Close()

	zipData, err := io.ReadAll(archive.Body)

	if err != nil {
		return nil, fmt.Errorf("open archive %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))

	if err != nil {
		return nil, fmt.Errorf("create zip reader %w", err)
	}

	return zipfs.NewFS(zipReader), nil
}
