package dlq

import (
	"net/http"

	"github.com/mauriciofsnts/hermes/internal/providers/database"
	"github.com/mauriciofsnts/hermes/internal/server/api"
)

// DLQController handles Dead Letter Queue operations.
type DLQController struct {
	dlq *database.DLQService
}

// NewDLQController creates a new DLQ controller.
func NewDLQController(dlq *database.DLQService) *DLQController {
	return &DLQController{dlq: dlq}
}

// GetStats godoc
//
//	@Summary		Get DLQ statistics
//	@Description	Returns statistics about the Dead Letter Queue (total, pending, failed counts)
//	@Tags			DLQ
//	@Produce		json
//	@Param			X-API-Key	header		string					true	"API Key"
//	@Success		200			{object}	map[string]interface{}	"DLQ statistics"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/dlq/stats [get]
func (c *DLQController) GetStats(r *http.Request) api.Response {
	if c.dlq == nil {
		return api.Err(api.InternalServerErr, "DLQ service not available")
	}

	stats, err := c.dlq.GetStats()
	if err != nil {
		return api.Err(api.InternalServerErr, "Failed to get DLQ stats: "+err.Error())
	}

	return api.Ok(stats)
}

// GetPending godoc
//
//	@Summary		Get pending DLQ items
//	@Description	Returns up to 100 pending items from the Dead Letter Queue
//	@Tags			DLQ
//	@Produce		json
//	@Param			X-API-Key	header		string					true	"API Key"
//	@Success		200			{object}	map[string]interface{}	"Pending DLQ items"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/dlq/pending [get]
func (c *DLQController) GetPending(r *http.Request) api.Response {
	if c.dlq == nil {
		return api.Err(api.InternalServerErr, "DLQ service not available")
	}

	letters, err := c.dlq.GetPending(100)
	if err != nil {
		return api.Err(api.InternalServerErr, "Failed to get pending items: "+err.Error())
	}

	return api.Ok(map[string]interface{}{
		"count": len(letters),
		"items": letters,
	})
}

// GetFailed godoc
//
//	@Summary		Get failed DLQ items
//	@Description	Returns up to 100 failed items from the Dead Letter Queue
//	@Tags			DLQ
//	@Produce		json
//	@Param			X-API-Key	header		string					true	"API Key"
//	@Success		200			{object}	map[string]interface{}	"Failed DLQ items"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/dlq/failed [get]
func (c *DLQController) GetFailed(r *http.Request) api.Response {
	if c.dlq == nil {
		return api.Err(api.InternalServerErr, "DLQ service not available")
	}

	letters, err := c.dlq.GetByStatus("failed", 100)
	if err != nil {
		return api.Err(api.InternalServerErr, "Failed to get failed items: "+err.Error())
	}

	return api.Ok(map[string]interface{}{
		"count": len(letters),
		"items": letters,
	})
}
