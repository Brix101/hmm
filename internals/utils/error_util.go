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
	writer     http.ResponseWriter `json:"-"`
	statusCode int                 `json:"-"`
	ErrorMap   map[string]Error    `json:"errors"`
}

func NewErrors(writer http.ResponseWriter) Errors {
	return Errors{
		writer:     writer,
		statusCode: 400,
		ErrorMap:   make(map[string]Error),
	}
}

func (err *Errors) HttpStatus(code int) {
	err.statusCode = code
	switch code {
	case http.StatusUnprocessableEntity:
		err.Add("root", Error{
			Message: "Failed to parse request body",
			Type:    "invalid_type",
		})
	}
}

func (err *Errors) Add(key string, value Error) {
	err.ErrorMap[key] = value
}

func (err Errors) Raise() {
	w := err.writer
	statusCode := err.statusCode
	errorJSON, _ := json.Marshal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(errorJSON)
}
