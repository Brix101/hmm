package utils

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

type Errors struct {
	StatusCode int              `json:"-"`
	Errors     map[string]Error `json:"errors"`
}

type ErrorsResponse struct {
	Errors map[string]Error `json:"errors"`
}

func NewErrors() Errors {
	return Errors{
		StatusCode: 400,
		Errors:     make(map[string]Error),
	}
}

func (err *Errors) HttpStatus(code int) {
	err.StatusCode = code
}

func (err *Errors) Add(key string, value Error) {
	err.Errors[key] = value
}

func (err Errors) ToResponse(w http.ResponseWriter) {
	errorJSON, _ := json.Marshal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	w.Write(errorJSON)
}
