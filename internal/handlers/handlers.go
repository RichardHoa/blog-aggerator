package handlers

import (
	"net/http"
)

func HandleHealthz(w http.ResponseWriter, r *http.Request) {

	respondText := map[string]string{
		"status": "OK"}

	RespondWithJSON(w, http.StatusOK, respondText)

}

func HandleReturnError(w http.ResponseWriter, r *http.Request) {

	RespondWithError(w, 500, "Internal Server Error")
}
