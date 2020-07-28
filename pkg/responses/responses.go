package responses

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON function to respond with a JSON payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

// RespondWithMessage is a wrapper function for responding on a successful request
func RespondWithMessage(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"message": message})
}

// RespondWithError is a wrapper function for responding on an unsuccessful request
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
