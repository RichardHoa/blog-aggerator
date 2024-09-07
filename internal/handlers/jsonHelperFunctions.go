package handlers

import (
	"net/http"
	"encoding/json"
)


func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	payload := map[string]string{"error": msg}

	response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
	
}