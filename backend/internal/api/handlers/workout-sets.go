package handlers

import "go-fitsync/backend/internal/database/sqlc"

type WorkoutSetHandler struct {
	queries   *sqlc.Queries
	jwtSecret []byte
}
