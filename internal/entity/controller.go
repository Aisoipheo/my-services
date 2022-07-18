package entity

import (
	"database/sql"
)

// Base class for any API
type Controller struct {
	cfg		*Config
	db		*sql.DB
	version	string
}
