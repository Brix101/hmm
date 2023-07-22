package utils

import "encoding/json"

type ErrorData struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

type ErrorMap map[string]ErrorData

type ErrorResponse struct {
	Errors map[string]ErrorData `json:"errors"`
}

// NewErrorMap creates a new instance of the DynamicMap
func NewErrorMap() ErrorMap {
	return make(ErrorMap)
}

// Add adds a new entry to the DynamicMap
func (em ErrorMap) Add(key string, value ErrorData) {
	em[key] = value
}

// Set updates an existing entry in the ErrorMap
func (em ErrorMap) Set(key, innerKey, innerValue string) {
	if _, ok := em[key]; !ok {
		em[key] = ErrorData{}
	}

	// Create a new ErrorData instance with updated fields
	data := em[key]
	switch innerKey {
	case "message":
		data.Message = innerValue
	case "type":
		data.Type = innerValue
	}
	em[key] = data
}

// Get retrieves a value from the ErrorMap
func (em ErrorMap) Get(key, innerKey string) string {
	if data, ok := em[key]; ok {
		switch innerKey {
		case "message":
			return data.Message
		case "type":
			return data.Type
		}
	}
	return ""
}

func (em ErrorMap) Json() ([]byte, error) {
	result := ErrorResponse{
		Errors: make(ErrorMap),
	}
	for key, data := range em {
		result.Errors[key] = data
	}
	return json.Marshal(result)
}
