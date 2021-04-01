package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func New(db *sql.DB, opts ...Option) *migrate {
	m := &migrate{
		db:         db,
		migrations: SqlMigrations{},
	}
	// Loop functional options
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Migrate will run the migrations using the provided db connection.
func (m *migrate) MigrateUp() error {
	m.printf("Migrating UP...\n")
	return m.migrate(UP)
}

// Migrate will run the migrations using the provided db connection.
func (m *migrate) MigrateDown() error {
	m.printf("Migrating DOWN...\n")
	return m.migrate(DOWN)
}

func (m *migrate) migrate(t migrationType) error {
	// Create Migrations Table
	m.printf("Creating/checking migrations table...\n")
	err := m.createMigrationTable()
	if err != nil {
		return err
	}
	// Prepare TX for all migrations
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

	// Create prepared statemnt for finding migration ID
	sid, err := tx.Prepare("SELECT id FROM migrations WHERE id=$1")
	if err != nil {
		return err
	}
	defer sid.Close()

	// Set stamant for Migration table
	switch t {
	case UP:
		m.printf("Running UP Migration...\n")
		// Loop migrations and apply
		for _, mi := range m.migrations {
			if err := m.upMigration(tx, sid, mi); err != nil {
				tx.Rollback()
				m.printf("Migration UP failed", mi.ID, err)
				return err
			}
		}
	case DOWN:
		m.printf("Running DOWN Migration...\n")
		// Loop rollbacks backwards and apply
		for i := len(m.migrations) - 1; i >= 0; i-- {
			mi := m.migrations[i]
			if err := m.downMigration(tx, sid, mi); err != nil {
				tx.Rollback()
				m.printf("Migration DOWN failed", mi.ID, err)
				return err
			}
		}
	default:
		return errors.New("no Migration chosen. Doing nothing")
	}
	// Migration done. Do commit
	m.printf("Migration done...")
	return tx.Commit()
}

func (m *migrate) upMigration(tx *sql.Tx, stmt *sql.Stmt, mig SqlMigration) error {
	ret, err := stmt.Exec(mig.ID)
	switch err {
	case sql.ErrNoRows:
		m.printf("Running migration: %v\n", mig.ID)
		// we need to run  the migration so we continue to code below
	case nil:
		n, _ := ret.RowsAffected()
		if n != 0 {
			m.printf("Skipping migration: %v\n", mig.ID)
			return nil
		}
		m.printf("Running migration: %v\n", mig.ID)
	default:
		return fmt.Errorf("looking up migration by id: %w", err)
	}
	if mig.Migrate == nil {
		m.printf("No migration provided")
		return nil
	}
	ins, err := tx.Prepare("INSERT INTO migrations (id) VALUES ($1)")
	if err != nil {
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(mig.ID)
	if err != nil {
		return fmt.Errorf("migration statement error: %w", err)
	}
	return mig.Migrate(tx)
}

func (m *migrate) downMigration(tx *sql.Tx, stmt *sql.Stmt, mig SqlMigration) error {
	ret, err := stmt.Exec(mig.ID)
	switch err {
	case sql.ErrNoRows:
		m.printf("Skipping migration: %v\n", mig.ID)
		return nil
		// we need to run  the migration so we continue to code below
	case nil:
		n, _ := ret.RowsAffected()
		if n == 0 {
			m.printf("Skipping migration: %v\n", mig.ID)
			return nil
		}
		m.printf("Running migration: %v\n", mig.ID)
	default:
		return fmt.Errorf("looking up migration by id: %w", err)
	}
	if mig.Rollback == nil {
		m.printf("No migration provided")
		return nil
	}
	// Create prepared statemnt for deleting migration
	del, err := tx.Prepare("DELETE FROM migrations WHERE id=$1")
	if err != nil {
		return err
	}
	defer del.Close()
	_, err = del.Exec(mig.ID)
	if err != nil {
		return fmt.Errorf("migration statement error: %w", err)
	}
	return mig.Rollback(tx)
}

// CreateSqlMigration - creates a SqlMigration object
func CreateSqlMigration(id string, up, down io.Reader) (SqlMigration, error) {
	m := SqlMigration{
		ID:       id,
		Migrate:  loadReader(up),
		Rollback: loadReader(down),
	}
	return m, nil
}

func loadReader(r io.Reader) func(tx *sql.Tx) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return func(tx *sql.Tx) error {
			return err
		}
	}
	return func(tx *sql.Tx) error {
		_, err := tx.Exec(string(buf))
		return err
	}
}

func (m *migrate) createMigrationTable() error {
	_, err := m.db.Exec("CREATE TABLE IF NOT EXISTS migrations (id TEXT PRIMARY KEY )")
	if err != nil {
		return fmt.Errorf("creating migrations table: %w", err)
	}
	return nil
}

func (m *migrate) printf(format string, a ...interface{}) {
	if m.debug {
		log.Printf(format, a...)
	}
}
