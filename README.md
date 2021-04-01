# go-migrate

Golang simple PostgreSQL migration tool

# Using Migration Tool

## Create Migrations from String

```go
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
m := migrate.SqlMigrations{m1, m2}
```

## Create Migrations from Files

```go
// Define migrations from fileReader
var mig migrate.SqlMigrations
fu, err := os.Open("01d-up.sql")
if err != nil {
    return nil, err
}
fd, err := os.Open("01-down.sql")
if err != nil {
    return nil, err
}
m, err := migrate.CreateSqlMigration("01", fu, fd)
if err != nil {
    return nil, err
}
mig = append(mig, m)
```

## Run Migration

```go
// Setup migration object
mig := migrate.New(db,
    migrate.WithDebugLog(),
    migrate.WithSqlMigrations(m),
)

// Run UP Migration
if err := mig.MigrateUp(); err != nil {
    panic(err)
}

// Run DOWN Migration
if err := mig.MigrateDown(); err != nil {
    panic(err)
}
```

## Examples

See [Examples folder](./cmd/examples)
