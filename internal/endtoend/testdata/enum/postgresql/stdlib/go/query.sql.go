// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const deleteBySize = `-- name: DeleteBySize :exec
DELETE FROM users
WHERE shoe_size = $1 AND shirt_size = $2
`

type DeleteBySizeParams struct {
	ShoeSize  Size
	ShirtSize NullSize
}

func (q *Queries) DeleteBySize(ctx context.Context, db DBTX, arg DeleteBySizeParams) error {
	_, err := db.ExecContext(ctx, deleteBySize, arg.ShoeSize, arg.ShirtSize)
	return err
}

const getAll = `-- name: GetAll :many
SELECT id, first_name, last_name, age, shoe_size, shirt_size FROM users
`

func (q *Queries) GetAll(ctx context.Context, db DBTX) ([]User, error) {
	rows, err := db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Age,
			&i.ShoeSize,
			&i.ShirtSize,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newUser = `-- name: NewUser :exec
INSERT INTO users (
    id,
    first_name,
    last_name,
    age,
    shoe_size,
    shirt_size
) VALUES
($1, $2, $3, $4, $5, $6)
`

type NewUserParams struct {
	ID        int32
	FirstName string
	LastName  sql.NullString
	Age       int32
	ShoeSize  Size
	ShirtSize NullSize
}

func (q *Queries) NewUser(ctx context.Context, db DBTX, arg NewUserParams) error {
	_, err := db.ExecContext(ctx, newUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Age,
		arg.ShoeSize,
		arg.ShirtSize,
	)
	return err
}

const updateSizes = `-- name: UpdateSizes :exec
UPDATE users
SET shoe_size = $2, shirt_size = $3
WHERE id = $1
`

type UpdateSizesParams struct {
	ID        int32
	ShoeSize  Size
	ShirtSize NullSize
}

func (q *Queries) UpdateSizes(ctx context.Context, db DBTX, arg UpdateSizesParams) error {
	_, err := db.ExecContext(ctx, updateSizes, arg.ID, arg.ShoeSize, arg.ShirtSize)
	return err
}
