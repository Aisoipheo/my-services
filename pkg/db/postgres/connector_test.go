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
		t.Fatal(errStr, err)
	}
	defer db.Close()
}

func TestNewPostgresDBFail(t *testing.T) {
	postgresCfg := PostgreSQLConfig {
		User	: "asd",
		Password: "wwwww",
		DBName	: "abracadabra",
		Host	: "1.1.1.1",
		Port	: "11111",
	}

	_, err := NewPostgresDB(&postgresCfg)
	if err == nil {
		t.Fatal("This is an impossible test. How did this happen?")
	}
}
