package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgreSQLConfig struct {
	User		string
	Password	string
	DBName		string
	Host		string
	Port		string
}

// connect to Postgres with creds from config, sslmode=disabled
func NewPostgresDB(cfg *PostgreSQLConfig) (*sql.DB, error) {
	connString := "postgres://" + cfg.User + ":" + cfg.Password +
		"@" + cfg.Host + ":" + cfg.Port + "/" + cfg.DBName + "?sslmode=disable&connect_timeout=2"

	// returns (conn, nil) in most cases, even if conn is not valid
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// timeout connection after 2 sec
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
