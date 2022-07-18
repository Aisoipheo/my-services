package healthz

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func getHealthz(c *gin.Context) error {
	return c.String(http.StatusOK, "Service is ready")
}
