package migrate

import "database/sql"

type migrate struct {
	db         *sql.DB
	migrations SqlMigrations
	debug      bool
}

type SqlMigrations []SqlMigration

type SqlMigration struct {
	ID       string
	Migrate  func(tx *sql.Tx) error
	Rollback func(tx *sql.Tx) error
}

type migrationType int

const (
	unknown migrationType = iota
	UP
	DOWN
)
