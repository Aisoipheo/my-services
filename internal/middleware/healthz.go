package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetHealthz(c *gin.Context) {
	c.String(http.StatusOK, "Service is ready. " + h.Config.ServiceVersion)
}

