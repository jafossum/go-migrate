package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	"github.com/jafossum/go-migrate/pkg/migrate"
	_ "github.com/lib/pq"
)

// Contants running on docker-compose test setup
const (
	host     = "localhost"
	port     = 5432
	user     = "migrate-test"
	password = "migrate-test"
	dbname   = "migrate-test"
)

const table = "example_test_table_1"

func main() {

	var mode = flag.String("m", "UP", "Mode UP or DOWN")
	flag.Parse()

	// Get Database SQL connection
	db := getDB()
	defer db.Close()

	// Get migrations
	m, err := getMigrations()
	if err != nil {
		panic(err)
	}

	// Setup migration object
	mig := migrate.New(db,
		migrate.WithDebugLog(),
		migrate.WithSqlMigrations(m),
	)

	switch *mode {
	case "UP":
		// Run UP Migration
		if err := mig.MigrateUp(); err != nil {
			panic(err)
		}
	case "DOWN":
		// Run UP Migration
		if err := mig.MigrateDown(); err != nil {
			panic(err)
		}
	default:
		fmt.Println("Argument must be 'UP' or 'DOWN', not", *mode)
	}

}

func getDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func getMigrations() (migrate.SqlMigrations, error) {
	// Define migrations
	m1, err := migrate.CreateSqlMigration("1",
		strings.NewReader("CREATE TABLE IF NOT EXISTS "+table+"()"),
		strings.NewReader("DROP TABLE "+table))
	if err != nil {
		return nil, err
	}
	m2, err := migrate.CreateSqlMigration("2",
		strings.NewReader("ALTER TABLE "+table+" ADD id serial"),
		strings.NewReader("ALTER TABLE "+table+" DROP id"))
	if err != nil {
		return nil, err
	}
	m3, err := migrate.CreateSqlMigration("3",
		strings.NewReader("INSERT INTO "+table+" VALUES (1)"),
		strings.NewReader("TRUNCATE "+table))
	if err != nil {
		return nil, err
	}
	m4, err := migrate.CreateSqlMigration("4",
		strings.NewReader("INSERT INTO "+table+" VALUES (2)"),
		strings.NewReader("DELETE FROM "+table+" WHERE id = 2"))
	if err != nil {
		return nil, err
	}
	return migrate.SqlMigrations{m1, m2, m3, m4}, nil
}
