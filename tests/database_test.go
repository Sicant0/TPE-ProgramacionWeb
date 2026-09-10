package tests

import (
    "context"
    "database/sql"
    "fmt"
    "os"
    "testing"
    "time"

    db "proyectoweb/db/sqlc"

    _ "github.com/jackc/pgx/v5/stdlib"
)

var (
    database *sql.DB
    queries  *db.Queries
)

func TestMain(m *testing.M) {
    connectionString := fmt.Sprintf(
        "postgres://%s:%s@%s:5432/%s?sslmode=disable",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_NAME"),
    )

    var err error
    database, err = sql.Open("pgx", connectionString)
    if err != nil {
        panic(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := database.PingContext(ctx); err != nil {
        panic(err)
    }

    queries = db.New(database)

    exitCode := m.Run()

    database.Close()
    os.Exit(exitCode)
}