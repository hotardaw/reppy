package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"go-fitsync/backend/internal/api/response"
	"go-fitsync/backend/internal/database/sqlc"
)

type MuscleHandler struct {
	queries *sqlc.Queries
}

func NewMuscleHandler(q *sqlc.Queries) *MuscleHandler {
	return &MuscleHandler{
		queries: q,
	}
}

func (h *MuscleHandler) HandleMuscles(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(strings.TrimSuffix(r.URL.Path, "/"))
	parts := strings.Split(cleanPath, "/")

	// Make sure only /muscles endpoint is handled
	if len(parts) != 2 || parts[1] != "muscles" {
		response.SendError(w, "Invalid URL - must be '/muscles'", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetMuscle(w, r)
	case http.MethodPost:
		h.CreateMuscle(w, r)
	case http.MethodDelete:
		h.DeleteMuscle(w, r)
	}
}

func (h *MuscleHandler) GetMuscle(w http.ResponseWriter, r *http.Request) {
	muscleName := r.URL.Query().Get("name") // '/muscles?name=Biceps%20Brachii'
	if muscleName == "" {
		response.SendError(w, "Muscle name is required for GET requests", http.StatusBadRequest)
		return
	}

	muscle, err := h.queries.GetMuscle(r.Context(), muscleName)
	if err != nil {
		response.SendError(w, "Muscle not found", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, muscle)
}

func (h *MuscleHandler) CreateMuscle(w http.ResponseWriter, r *http.Request) {
	var request struct {
		MuscleName  string `json:"muscle_name"`
		MuscleGroup string `json:"muscle_group"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	muscle, err := h.queries.CreateMuscle(r.Context(), sqlc.CreateMuscleParams{
		MuscleName:  request.MuscleName,
		MuscleGroup: request.MuscleGroup,
	})
	if err != nil {
		response.SendError(w, "Failed to create muscle", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, muscle, http.StatusCreated)
}

func (h *MuscleHandler) DeleteMuscle(w http.ResponseWriter, r *http.Request) {
	muscleName := r.URL.Query().Get("name") // '/muscles?name=Biceps%20Brachii'
	if muscleName == "" {
		response.SendError(w, "Muscle name is required for DELETE requests", http.StatusBadRequest)
		return
	}

	muscle, err := h.queries.DeleteMuscle(r.Context(), muscleName)
	if err != nil {
		response.SendError(w, "Failed to delete muscle", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, muscle, http.StatusNoContent)
}
