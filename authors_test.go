package main

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/dot96gal/go-sqlc-postgres-sample/internal/sqlc"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateAuthor(t *testing.T) {
	authorUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams sqlc.CreateAuthorParams
		}
		expected sqlc.Author
	}{
		{
			scenario: "create author",
			input: struct {
				createAuthorParams sqlc.CreateAuthorParams
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					ID:   pgtype.UUID{Bytes: authorUuid, Valid: true},
					Name: "author001",
					Bio:  pgtype.Text{String: "author001", Valid: true},
				},
			},
			expected: sqlc.Author{
				ID:   pgtype.UUID{Bytes: authorUuid, Valid: true},
				Name: "author001",
				Bio:  pgtype.Text{String: "author001", Valid: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()

			fmt.Println(conn)
			queries := sqlc.New(conn)

			// test with transaction
			tx, err := conn.Begin(ctx)
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				err = tx.Rollback(ctx)
				if err != nil {
					t.Error(err)
				}
			})

			queries = queries.WithTx(tx)

			// crete author
			author, err := queries.CreateAuthor(ctx, tt.input.createAuthorParams)
			if err != nil {
				t.Error(err)
			}

			if author != tt.expected {
				t.Errorf("got=%v, want=%v", author, tt.expected)
			}
		})
	}
}

func TestUpdateAuthor(t *testing.T) {
	authorUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams sqlc.CreateAuthorParams
			updateAuthorParams sqlc.UpdateAuthorParams
		}
		expected sqlc.Author
	}{
		{
			scenario: "update author",
			input: struct {
				createAuthorParams sqlc.CreateAuthorParams
				updateAuthorParams sqlc.UpdateAuthorParams
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					ID:   pgtype.UUID{Bytes: authorUuid, Valid: true},
					Name: "author001",
					Bio:  pgtype.Text{String: "author001", Valid: true},
				},
				updateAuthorParams: sqlc.UpdateAuthorParams{
					Name: "Updated: author001",
					Bio:  pgtype.Text{String: "Updated: author001", Valid: true},
					ID:   pgtype.UUID{Bytes: authorUuid, Valid: true},
				},
			},
			expected: sqlc.Author{
				ID:   pgtype.UUID{Bytes: authorUuid, Valid: true},
				Name: "Updated: author001",
				Bio:  pgtype.Text{String: "Updated: author001", Valid: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()

			queries := sqlc.New(conn)

			// test with transaction
			tx, err := conn.Begin(ctx)
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				err = tx.Rollback(ctx)
				if err != nil {
					t.Error(err)
				}
			})

			queries = queries.WithTx(tx)

			// crete author
			_, err = queries.CreateAuthor(ctx, tt.input.createAuthorParams)
			if err != nil {
				t.Error(err)
			}

			// update author
			author, err := queries.UpdateAuthor(ctx, tt.input.updateAuthorParams)
			if err != nil {
				t.Error(err)
			}

			if author != tt.expected {
				t.Errorf("got=%v, want=%v", author, tt.expected)
			}
		})
	}
}

func TestDeleteAuthor(t *testing.T) {
	authorUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams sqlc.CreateAuthorParams
			deleteAuthorUuid   pgtype.UUID
		}
		expected sqlc.Author
	}{
		{
			scenario: "delete author",
			input: struct {
				createAuthorParams sqlc.CreateAuthorParams
				deleteAuthorUuid   pgtype.UUID
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					ID:   pgtype.UUID{Bytes: authorUuid, Valid: true},
					Name: "author001",
					Bio:  pgtype.Text{String: "author001", Valid: true},
				},
				deleteAuthorUuid: pgtype.UUID{Bytes: authorUuid, Valid: true},
			},
			expected: sqlc.Author{
				ID:   pgtype.UUID{Bytes: authorUuid, Valid: true},
				Name: "author001",
				Bio:  pgtype.Text{String: "author001", Valid: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()

			queries := sqlc.New(conn)

			// test with transaction
			tx, err := conn.Begin(ctx)
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				err = tx.Rollback(ctx)
				if err != nil {
					t.Error(err)
				}
			})

			queries = queries.WithTx(tx)

			// crete author
			_, err = queries.CreateAuthor(ctx, tt.input.createAuthorParams)
			if err != nil {
				t.Error(err)
			}

			// delete author
			author, err := queries.DeleteAuthor(ctx, tt.input.deleteAuthorUuid)
			if err != nil {
				t.Error(err)
			}

			if author != tt.expected {
				t.Errorf("got=%v, want=%v", author, tt.expected)
			}
		})
	}
}

func TestListAuthors(t *testing.T) {
	authorUuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParamsList []sqlc.CreateAuthorParams
		}
		expected []sqlc.Author
	}{
		{
			scenario: "list authors",
			input: struct{ createAuthorParamsList []sqlc.CreateAuthorParams }{
				createAuthorParamsList: []sqlc.CreateAuthorParams{
					{
						ID:   pgtype.UUID{Bytes: authorUuids[0], Valid: true},
						Name: "author001",
						Bio:  pgtype.Text{String: "author001", Valid: true},
					},
					{
						ID:   pgtype.UUID{Bytes: authorUuids[1], Valid: true},
						Name: "author002",
						Bio:  pgtype.Text{String: "author002", Valid: true},
					},
				},
			},
			expected: []sqlc.Author{
				{
					ID:   pgtype.UUID{Bytes: authorUuids[0], Valid: true},
					Name: "author001",
					Bio:  pgtype.Text{String: "author001", Valid: true},
				},
				{
					ID:   pgtype.UUID{Bytes: authorUuids[1], Valid: true},
					Name: "author002",
					Bio:  pgtype.Text{String: "author002", Valid: true},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			ctx := context.Background()

			queries := sqlc.New(conn)

			// test with transaction
			tx, err := conn.Begin(ctx)
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				err = tx.Rollback(ctx)
				if err != nil {
					t.Error(err)
				}
			})

			queries = queries.WithTx(tx)

			// crete author
			for _, params := range tt.input.createAuthorParamsList {
				_, err := queries.CreateAuthor(ctx, params)
				if err != nil {
					t.Error(err)
				}
			}

			// list authors
			authors, err := queries.ListAuthors(ctx)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(
				tt.expected,
				func(i, j int) bool {
					return tt.expected[i].Name < tt.expected[j].Name
				},
			)

			for i := range authors {
				if authors[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", authors[i], tt.expected[i])
				}
			}
		})
	}
}
