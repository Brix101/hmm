package handlers

import (
	"encoding/json"
	"fmt"
	"home-server/internal/services"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type AuthResource struct {
	UserServices *services.UserServices
}

// Routes creates a REST router for the todos resource
func (rs AuthResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Post("/signin", rs.SignIn)    // POST /users - create a new user and persist it
	r.Get("/me", rs.GetCurrentUser) // GET /users - read a list of users

	return r
}

type signInRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (rs AuthResource) SignIn(w http.ResponseWriter, r *http.Request) {
	// get the info from the request body
	var body signInRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusUnprocessableEntity)
		return
	}

	// Validate the user struct
	validate := validator.New()
	err = validate.Struct(body)
	if err != nil {
		// Initialize the map to hold the validation errors
		validationErrors := make(map[string][]string)

		// Get the validation errors from the validator and store them in the map
		for _, err := range err.(validator.ValidationErrors) {
			// Make the field name lowercase
			fieldName := strings.ToLower(err.Field())

			// Append the error message to the field's slice in the map
			errorMessage := fmt.Sprintf("%s is %s", fieldName, err.Tag())
			validationErrors[fieldName] = append(validationErrors[fieldName], errorMessage)
		}

		// Create the error response JSON
		errorResponse := ErrorResponse{
			Errors: validationErrors,
		}

		// Marshal the error response JSON
		errorJSON, _ := json.Marshal(errorResponse)

		// Return the validation errors in the HTTP response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorJSON)
		return
	}

	existingUser, err := rs.UserServices.GetUserByEmail(body.Email)
	if err != nil {
		// Handle the error, such as "user not found"
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Compare the password from the database with the provided password
	if existingUser.Password != body.Password {
		// Passwords do not match
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}
	// Store the user ID in a cookie
	cookie := http.Cookie{
		Name:     "session",
		Value:    existingUser.Id.String(),            // Use the user ID as the cookie value
		Expires:  time.Now().Add(30 * 24 * time.Hour), // Set the expiration time to 30 days from now
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// Serialize the created user to JSON
	userJSON, err := json.Marshal(existingUser)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	// Set the appropriate content type in the response header
	w.Header().Set("Content-Type", "application/json")
	// Write the user JSON to the response
	w.Write(userJSON)
}

func (rs AuthResource) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Read the cookie from the request
	cookie, err := r.Cookie("session")
	if err != nil {
		// No session cookie found, the user is not authenticated
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract the user ID from the cookie value
	userID, err := uuid.Parse(cookie.Value)
	if err != nil {
		// Invalid or corrupted cookie value
		http.Error(w, "Invalid cookie", http.StatusBadRequest)
		return
	}

	// Fetch the user from the database using the user ID
	existingUser, err := rs.UserServices.GetUserByID(userID)
	if err != nil {
		// Handle the error, such as "user not found"
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Serialize the user to JSON
	userJSON, err := json.Marshal(existingUser)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	// Set the appropriate content type in the response header
	w.Header().Set("Content-Type", "application/json")
	// Write the user JSON to the response
	w.Write(userJSON)
}
