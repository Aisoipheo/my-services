package postgres

import (
	"testing"
)

func TestNewPostgresDB(t *testing.T) {
	errStr := `
This test is supposed to be run with PostgreSQL running
addr		: localhost:5432
Test user	: user
Test pass	: qwerty123
test dbName	: test
`

	//this creds are supposed to be in the test database
	postgresCfg := PostgreSQLConfig {
		User	: "user",
		Password: "qwerty123",
		DBName	: "test",
		Host	: "localhost",
		Port	: "5432",
	}

	db, err := NewPostgresDB(&postgresCfg)
	if err != nil {
		t.Fatalf(errStr, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf(errStr, err)
	}
}
