package api

import (
	"net/http"

	"github.com/duythinht/tg/lib/store"
	"github.com/labstack/echo/v4"
)

type API struct {
	*echo.Echo
	db    store.DB
	queue store.QueueWriter
}

func New(db store.DB, queue store.QueueWriter) *API {
	e := echo.New()

	api := &API{
		Echo:  e,
		db:    db,
		queue: queue,
	}

	api.setupRoutes()

	return api
}

func (api *API) setupRoutes() {
	api.GET("/-/health", api.health)
	api.POST("/api/v1/repositories", api.createRepository)
	api.GET("/api/v1/repositories", api.listRepositories)
	api.GET("/api/v1/scans/:repositoryId", api.listScans)
	api.POST("/api/v1/scans/:repositoryId", api.triggerScan)
}

func (api *API) health(c echo.Context) error {
	return c.String(200, "ok")
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.Echo.ServeHTTP(w, r)
}

type errorResponse struct {
	Error string `json:"error"`
}
