package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/duythinht/tg/model"
	"github.com/duythinht/tg/worker"
	"github.com/labstack/echo/v4"
)

type listScansItem struct {
	model.Scan
	QueuedAt   time.Time `json:"queued_at"`
	ScanningAt time.Time `json:"scanning_at"`
	FinishedAt time.Time `json:"finished_at"`
}

type listScansResponse struct {
	Scans []listScansItem `json:"scans"`
}

func (api *API) listScans(c echo.Context) error {
	repositoryId := c.Param("repositoryId")

	rows, err := api.db.Query(c.Request().Context(), "SELECT * FROM scans WHERE repository_id=$1 ORDER BY id DESC", repositoryId)

	if err != nil {
		return c.JSON(400, errorResponse{
			Error: err.Error(),
		})
	}
	result := make([]listScansItem, 0)

	for rows.Next() {
		var item model.Scan
		err = rows.Scan(
			&item.ID,
			&item.RepositoryID,
			&item.Status,
			&item.Findings,
			&item.QueuedAt,
			&item.ScanningAt,
			&item.FinishedAt,
		)
		if err != nil {
			return c.JSON(400, errorResponse{
				Error: err.Error(),
			})
		}
		result = append(result, listScansItem{
			Scan:       item,
			QueuedAt:   item.QueuedAt.Time,
			ScanningAt: item.ScanningAt.Time,
			FinishedAt: item.FinishedAt.Time,
		})
	}

	return c.JSON(200, listScansResponse{
		Scans: result,
	})
}

type triggerScanResponse struct {
	ScanID int64 `json:"scanId"`
}

func (api *API) triggerScan(c echo.Context) error {
	repositoryId, err := strconv.ParseInt(c.Param("repositoryId"), 10, 64)

	if err != nil {
		return fmt.Errorf("manipulate repositoryId %w", err)
	}

	var owner, repo string

	err = api.db.QueryRow(
		c.Request().Context(),
		"SELECT owner, repository FROM repositories WHERE id=$1",
		repositoryId,
	).Scan(&owner, &repo)

	if err != nil {
		return c.JSON(404, errorResponse{
			Error: err.Error(),
		})
	}

	var lastScanInsertedID int64
	err = api.db.QueryRow(
		c.Request().Context(),
		"INSERT INTO scans(repository_id ,status) values($1, $2) RETURNING id",
		repositoryId,
		"Queued",
	).Scan(&lastScanInsertedID)

	if err != nil {
		return c.JSON(400, errorResponse{
			Error: err.Error(),
		})
	}

	message, err := json.Marshal(worker.ScanRequest{
		ScanID:       lastScanInsertedID,
		RepositoryID: repositoryId,
		Owner:        owner,
		Repository:   repo,
	})
	if err != nil {
		return fmt.Errorf("marshal queue message %w", err)
	}

	err = api.queue.WriteMessage(c.Request().Context(), message)

	if err != nil {
		return c.JSON(400, errorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(201, triggerScanResponse{
		ScanID: lastScanInsertedID,
	})
}
