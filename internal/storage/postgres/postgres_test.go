//go:build integration
// +build integration

package postgres_test

import (
	"database/sql"
	"flag"
	"testing"

	"github.com/adrianolmedo/go-restapi/config"
	"github.com/adrianolmedo/go-restapi/internal/storage/postgres"
)

// $ go test -v -tags integration -args -dbengine postgres -dbhost 127.0.0.1 -dbport 5432 -dbuser username -dbname foodb -dbpass 12345
var (
	dbhost   = flag.String("dbhost", "", "Database host.")
	dbengine = flag.String("dbengine", "", "Database engine, choose mysql or postgres.")
	dbport   = flag.String("dbport", "", "Database port.")
	dbuser   = flag.String("dbuser", "", "Database user.")
	dbpass   = flag.String("dbpass", "", "Database password.")
	dbname   = flag.String("dbname", "", "Database name.")
)

// TestDB can open & close.
func TestDB(t *testing.T) {
	db := openDB(t)
	closeDB(t, db)
}

func openDB(t *testing.T) *sql.DB {
	dbcfg := config.Database{
		Engine:   *dbengine,
		Host:     *dbhost,
		Port:     *dbport,
		User:     *dbuser,
		Password: *dbpass,
		Name:     *dbname,
	}

	db, err := postgres.New(dbcfg)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func closeDB(t *testing.T, db *sql.DB) {
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}
