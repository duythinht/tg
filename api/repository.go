package api

import (
	"fmt"
	"time"

	"github.com/duythinht/tg/model"
	"github.com/labstack/echo/v4"
)

type createRepositoryRequest struct {
	URL string `json:"url"`
}

type createRepositoryResponse struct {
	RepositoryID int64 `json:"repositoryId"`
}

func (api *API) createRepository(c echo.Context) error {

	req := &createRepositoryRequest{}

	err := c.Bind(req)

	if err != nil {
		return fmt.Errorf("binding request: %w", err)
	}

	repository, err := model.RepositoryFromURL(req.URL)

	if err != nil {
		return c.JSON(400, errorResponse{
			Error: err.Error(),
		})
	}

	var lastInsertId int64

	now := time.Now()

	err = api.db.QueryRow(c.Request().Context(), "INSERT INTO repositories(host, owner, repository, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		repository.Host,
		repository.Owner,
		repository.Repository,
		now,
		now,
	).Scan(&lastInsertId)

	if err != nil {
		return c.JSON(400, errorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(201, createRepositoryResponse{
		RepositoryID: lastInsertId,
	})
}

type repositoriesItem struct {
	model.Repository
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type listRepositoriesResponse struct {
	Repositories []repositoriesItem `json:"repositories"`
}

func (api *API) listRepositories(c echo.Context) error {
	rows, err := api.db.Query(c.Request().Context(), "SELECT id, host, owner, repository, created_at, updated_at FROM repositories")
	if err != nil {
		return c.JSON(400, errorResponse{
			Error: err.Error(),
		})
	}

	result := make([]repositoriesItem, 0)

	for rows.Next() {
		var item model.Repository
		err := rows.Scan(&item.ID, &item.Host, &item.Owner, &item.Repository, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return c.JSON(400, errorResponse{
				Error: err.Error(),
			})
		}
		//item.URL = item.Repository.URL()
		result = append(result, repositoriesItem{
			Repository: item,
			URL:        item.URL(),
			CreatedAt:  item.CreatedAt.Time,
			UpdatedAt:  item.UpdatedAt.Time,
		})
	}

	return c.JSON(200, &listRepositoriesResponse{
		Repositories: result,
	})
}
