package testhelper

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"api-go/internal/model"
)

type TestDB struct {
	DB        *bun.DB
	container testcontainers.Container
}

func NewTestDB(t *testing.T) *TestDB {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	db.RegisterModel((*model.PostTag)(nil))

	tdb := &TestDB{DB: db, container: container}
	tdb.migrate(t, ctx)

	t.Cleanup(func() {
		db.Close()
		container.Terminate(ctx)
	})

	return tdb
}

func (t *TestDB) migrate(tb testing.TB, ctx context.Context) {
	tb.Helper()

	// Use the underlying sql.DB to avoid bun rewriting DDL placeholders
	sqlDB := t.DB.DB

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			name       VARCHAR(255) NOT NULL DEFAULT '',
			email      VARCHAR(255) NOT NULL UNIQUE,
			password   VARCHAR(255) NOT NULL DEFAULT '',
			is_admin   BOOLEAN NOT NULL DEFAULT FALSE,
			status     INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id           SERIAL PRIMARY KEY,
			title        VARCHAR(255) NOT NULL,
			slug         VARCHAR(255) NOT NULL UNIQUE,
			preview_text TEXT,
			detail_text  TEXT,
			created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id          SERIAL PRIMARY KEY,
			title       VARCHAR(255) NOT NULL,
			slug        VARCHAR(255) NOT NULL UNIQUE,
			content     TEXT,
			description TEXT,
			image       VARCHAR(255),
			status      BOOLEAN NOT NULL DEFAULT FALSE,
			is_featured BOOLEAN NOT NULL DEFAULT FALSE,
			views_count INT NOT NULL DEFAULT 0,
			user_id     INT NOT NULL REFERENCES users(id),
			category_id INT REFERENCES categories(id),
			created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id         SERIAL PRIMARY KEY,
			title      VARCHAR(255) NOT NULL,
			slug       VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS post_tags (
			post_id INT NOT NULL REFERENCES posts(id),
			tag_id  INT NOT NULL REFERENCES tags(id),
			PRIMARY KEY (post_id, tag_id)
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id          SERIAL PRIMARY KEY,
			author_name VARCHAR(255) NOT NULL DEFAULT 'anon',
			text        TEXT NOT NULL,
			post_id     INT NOT NULL REFERENCES posts(id),
			user_id     INT REFERENCES users(id),
			status      INT NOT NULL DEFAULT 0,
			created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS posts_likes (
			id            SERIAL PRIMARY KEY,
			post_id       INT NOT NULL REFERENCES posts(id),
			ip            VARCHAR(255),
			device_memory INT,
			created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS aphorism (
			id          SERIAL PRIMARY KEY,
			detail_text VARCHAR(255) NOT NULL,
			created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS incomings (
			id         SERIAL PRIMARY KEY,
			name       TEXT NOT NULL,
			email      TEXT,
			phone      TEXT,
			message    TEXT NOT NULL,
			status     INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS subscriptions (
			id         SERIAL PRIMARY KEY,
			email      VARCHAR(255) NOT NULL UNIQUE,
			token      VARCHAR(255),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
	}

	for _, q := range queries {
		if _, err := sqlDB.ExecContext(ctx, q); err != nil {
			tb.Fatalf("migration failed: %v\nquery: %s", err, q)
		}
	}
}

func (t *TestDB) Truncate(tb testing.TB, tables ...string) {
	tb.Helper()
	ctx := context.Background()
	sqlDB := t.DB.DB
	for _, table := range tables {
		if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)); err != nil {
			tb.Fatalf("truncate %s failed: %v", table, err)
		}
	}
}

// SQLDB returns the underlying *sql.DB for raw queries in tests.
func (t *TestDB) SQLDB() *sql.DB {
	return t.DB.DB
}
