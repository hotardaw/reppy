package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Sets "Content-Type" to "application/json" and returns an error response with the error's message and status code.
func SendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   message,
	})
}

// Sets "Content-Type" to "application/json" and returns a successful response, any returned data, and an optional status code. Only to be used by REST APIs with JSON responses.
func SendSuccess(w http.ResponseWriter, data interface{}, code ...int) {
	w.Header().Set("Content-Type", "application/json")
	if len(code) > 0 {
		w.WriteHeader(code[0])
	}
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

// Both of the below are possible with variadic input param "code ...int"

// // WITHOUT status code - will use default 200 OK
// response.SendSuccess(w, data)

// WITH status code
// response.SendSuccess(w, data, http.StatusCreated)  // 201
