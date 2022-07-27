package middleware

import (
	"database/sql"

	"feed-service/internal/models"
)

// Base class for any API
type Controller struct {
	Cfg		*models.Config
	DB		*sql.DB
}
