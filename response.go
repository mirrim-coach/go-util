package goutils

import (
	"encoding/json"
	"net/http"
)

// GenericResponse is the structure for the error response messages
type GenericResponse struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}

// ResponseWithJSON is for returning json
func ResponseWithJSON(w http.ResponseWriter, response interface{}) {
	w.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.Encode(response)
}

// ReturnGeneric is for returning a generic response
func ReturnGeneric(w http.ResponseWriter, message interface{}, status int) {
	response := GenericResponse{
		Status:  status,
		Message: message,
	}
	w.WriteHeader(status)
	ResponseWithJSON(w, response)
}
