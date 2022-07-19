package entity

import (
	"database/sql"
)

// Base class for any API
type Controller struct {
	Cfg		*Config
	DB		*sql.DB
	Version	string
}
