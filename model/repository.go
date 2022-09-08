package model

import (
	"fmt"
	"strings"

	"github.com/jackc/pgtype"
)

// Repository model for pg schema
type Repository struct {
	ID         int64            `json:"id"`
	Host       string           `json:"host"`
	Owner      string           `json:"owner"`
	Repository string           `json:"repository"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
}

// URL return repository url
func (r *Repository) URL() string {

	switch r.Host {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s", r.Owner, r.Repository)
	default:
		// currently only support for github (other host is not support on this version), don't be panic here
		panic(fmt.Sprintf("host: %s does not have support", r.Host))
	}
}

func RepositoryFromURL(rawURL string) (*Repository, error) {
	switch {
	// handle for ssh url
	case strings.HasPrefix(rawURL, "git@github.com:"):
		repoPath := strings.Split(rawURL[15:], "/")

		if len(repoPath) < 2 {
			return nil, fmt.Errorf("%s is not match github repo ssh url", rawURL)
		}
		return &Repository{
			Host:       "github",
			Owner:      repoPath[0],
			Repository: strings.TrimSuffix(repoPath[1], ".git"),
		}, nil
	case strings.HasPrefix(rawURL, "https://github.com/"):
		repoPath := strings.Split(rawURL[19:], "/")
		if len(repoPath) < 2 {
			return nil, fmt.Errorf("%s is not match github repo ssh url", rawURL)
		}

		return &Repository{
			Host:       "github",
			Owner:      repoPath[0],
			Repository: strings.TrimSuffix(repoPath[1], ".git"),
		}, nil
	}

	return nil, fmt.Errorf("%s, unknown url pattern", rawURL)
}
