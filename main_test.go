package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var conn *pgx.Conn

func TestMain(m *testing.M) {
	postgresDB := os.Getenv("TEST_POSTGRES_DB")
	postgresUser := os.Getenv("TEST_POSTGRES_USER")
	postgresPass := os.Getenv("TEST_POSTGRES_PASSWORD")
	// postgresHost := "localhost"
	// postgresPort := "5432"

	// databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPass, postgresHost, postgresPort, postgresDB)

	pwd, _ := os.Getwd()
	schemaFiles := fmt.Sprintf("file://%s/db/migrations", pwd)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	pool.MaxWait = 30 * time.Second

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	runOptions := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.2",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", postgresDB),
			fmt.Sprintf("POSTGRES_USER=%s", postgresUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", postgresPass),
			"listen_addresses='*'",
		},
	}

	resource, err := pool.RunWithOptions(
		runOptions,
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	ctx := context.Background()

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", postgresUser, postgresPass, resource.GetHostPort("5432/tcp"), postgresDB)

	if err := pool.Retry(func() error {
		conn, err = pgx.Connect(ctx, databaseUrl)
		if err != nil {
			return err
		}

		return conn.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// database migration
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Could not connect to database %s", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not instantiate driver: %s", err)
	}
	mig, err := migrate.NewWithDatabaseInstance(schemaFiles, "postgresql", driver)
	if err != nil {
		log.Fatalf("Could not instantiate migrate: %s", err)
	}
	err = mig.Up()
	if err != nil {
		log.Fatalf("Could not migrate database: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
