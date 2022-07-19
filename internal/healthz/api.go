package healthz

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"my-service/internal/entity"
)

func (h *Controller) getHealthz(c *gin.Context) error {
	return c.String(http.StatusOK, "Service is ready. " + h.Version)
}

