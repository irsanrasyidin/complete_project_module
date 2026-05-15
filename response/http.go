package response

import (
	"encoding/json"
	"net/http"
	"time"
)

func WriteJSON[T any](w http.ResponseWriter, res ApiResponse[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	_ = json.NewEncoder(w).Encode(res)
}

func WriteSuccess[T any](w http.ResponseWriter, correlationID string, tin time.Time, statusCode int, data T) {
	WriteJSON(w, Success(correlationID, tin, statusCode, data))
}

func WriteError(w http.ResponseWriter, correlationID string, tin time.Time, statusCode int, message string) {
	WriteJSON(w, Error[any](correlationID, tin, statusCode, message))
}

func WriteRawError(w http.ResponseWriter, errRes *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errRes.StatusCode)
	_ = json.NewEncoder(w).Encode(errRes)
}
