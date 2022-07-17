package healthz

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Healthz(c *gin.Context) error {
	return c.String(http.StatusOK, "Service is ready")
}
