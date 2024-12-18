// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: seeding.sql

package sqlc

import (
	"context"
)

const checkInitialSeed = `-- name: CheckInitialSeed :one
SELECT EXISTS (
  SELECT 1 FROM app_state WHERE key = 'initial_seed_completed'
)
`

func (q *Queries) CheckInitialSeed(ctx context.Context) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkInitialSeed)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const markInitialSeed = `-- name: MarkInitialSeed :exec
INSERT INTO app_state (key, value) VALUES ('initial_seed_completed', 'true')
`

func (q *Queries) MarkInitialSeed(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, markInitialSeed)
	return err
}
