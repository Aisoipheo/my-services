package postgres

import (
	"database/sql"

	"github.com/lib/pq"
)

type PostgreSQLConfig {
	User		string
	Password	string
	DBName		string
	Host		string
	Port		string
}

func NewPostgresDB(cfg *PostgreSQLConfig) *sql.DB, error {
	connString := "postgresql://" + cfg.User + ":" + cfg.Password +
		"@" + cfg.Host + ":" + cfg.Port + "/" + cfg.DBName + "?sslmode=disabled"

	conn, err := sql.Open("postgres", connString)
	if err != nil {
		// TODO log error
		return nil, err
	}
	return conn, nil
}
