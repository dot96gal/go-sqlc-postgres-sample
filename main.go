package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dot96gal/go-sqlc-postgres-sample/internal/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func run() error {
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := "localhost"
	postgresPort := "5432"

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPass, postgresHost, postgresPort, postgresDB))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := sqlc.New(conn)

	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	author, err := queries.CreateAuthor(ctx, sqlc.CreateAuthorParams{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Name: "Brian Kernighan",
		Bio: pgtype.Text{
			String: "Co-author of The C Programming Language and The Go Programming Language",
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}

	log.Println(author)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
