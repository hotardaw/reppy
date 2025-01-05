package handlers

import (
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
)

type WorkoutSetHandler struct {
	queries   *sqlc.Queries
	jwtSecret []byte
}

func NewWorkoutSetHandler(q *sqlc.Queries, jwtSecret []byte) *WorkoutSetHandler {
	return &WorkoutSetHandler{
		queries:   q,
		jwtSecret: jwtSecret,
	}
}

// req, resp structs

func (h *WorkoutSetHandler) HandleWorkoutSets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateWorkoutSet(w, r)
	}
}

func (h *WorkoutSetHandler) CreateWorkoutSet(w http.ResponseWriter, r *http.Request) {}
