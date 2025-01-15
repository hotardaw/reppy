package handlers

import (
	"encoding/json"
	"net/http"

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

type CreateMuscleRequest struct {
	MuscleName  string `json:"muscle_name"`
	MuscleGroup string `json:"muscle_group"`
}

func (h *MuscleHandler) HandleMuscles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateMuscle(w, r)
	case http.MethodGet:
		h.GetMuscle(w, r)
	case http.MethodDelete:
		h.DeleteMuscle(w, r)
	default:
		response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// "/muscles"
func (h *MuscleHandler) CreateMuscle(w http.ResponseWriter, r *http.Request) {
	var request CreateMuscleRequest
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

// "/muscles?name=Biceps%20Brachii"
func (h *MuscleHandler) GetMuscle(w http.ResponseWriter, r *http.Request) {
	muscleName := r.URL.Query().Get("name")
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

// "/muscles?name=Biceps%20Brachii"
func (h *MuscleHandler) DeleteMuscle(w http.ResponseWriter, r *http.Request) {
	muscleName := r.URL.Query().Get("name")
	if muscleName == "" {
		response.SendError(w, "Muscle name is required for DELETE requests", http.StatusBadRequest)
		return
	}

	deletedMuscle, err := h.queries.DeleteMuscle(r.Context(), muscleName)
	if err != nil {
		response.SendError(w, "Failed to delete muscle", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, map[string]interface{}{
		"message": "Muscle deleted successfully",
		"id":      deletedMuscle,
	}, http.StatusOK) // Not StatusNoContent bc this is a soft delete))
}
