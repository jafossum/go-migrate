package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

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
		// Run DOWN Migration
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
	var mig migrate.SqlMigrations
	for i := 0; i < 4; i++ {
		// Define migrations from fileReader
		fu, err := os.Open(fmt.Sprintf("0%d-up.sql", i+1))
		if err != nil {
			return nil, err
		}
		fd, err := os.Open(fmt.Sprintf("0%d-down.sql", i+1))
		if err != nil {
			return nil, err
		}
		m, err := migrate.CreateSqlMigration(fmt.Sprintf("1%d", i+1), fu, fd)
		if err != nil {
			return nil, err
		}
		mig = append(mig, m)
	}
	return mig, nil
}
