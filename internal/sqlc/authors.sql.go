// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: authors.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO
  authors (id, name, bio)
VALUES
  ($1, $2, $3)
RETURNING
  id, name, bio
`

type CreateAuthorParams struct {
	ID   pgtype.UUID
	Name string
	Bio  pgtype.Text
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, createAuthor, arg.ID, arg.Name, arg.Bio)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :one
DELETE FROM authors
WHERE
  id = $1
RETURNING
  id, name, bio
`

func (q *Queries) DeleteAuthor(ctx context.Context, id pgtype.UUID) (Author, error) {
	row := q.db.QueryRow(ctx, deleteAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const getAuthor = `-- name: GetAuthor :one
SELECT
  id, name, bio
FROM
  authors
WHERE
  id = $1
LIMIT
  1
`

func (q *Queries) GetAuthor(ctx context.Context, id pgtype.UUID) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT
  id, name, bio
FROM
  authors
ORDER BY
  name
`

func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
	rows, err := q.db.Query(ctx, listAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name, &i.Bio); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAuthor = `-- name: UpdateAuthor :one
UPDATE authors
set
  name = $2,
  bio = $3
WHERE
  id = $1
RETURNING
  id, name, bio
`

type UpdateAuthorParams struct {
	ID   pgtype.UUID
	Name string
	Bio  pgtype.Text
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, updateAuthor, arg.ID, arg.Name, arg.Bio)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}
