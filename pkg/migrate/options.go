package migrate

type Option func(*migrate)

// Functional Options
func WithSqlMigrations(mig SqlMigrations) Option {
	return func(m *migrate) {
		m.migrations = mig
	}
}

// Functional Options
func WithDebugLog() Option {
	return func(m *migrate) {
		m.debug = true
	}
}

// Functional Options
func WithMigrationTableName(t string) Option {
	return func(m *migrate) {
		m.migTable = t
	}
}
