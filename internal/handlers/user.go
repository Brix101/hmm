package handlers

import (
	"encoding/json"
	"fmt"
	"home-server/internal/services"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type UsersResource struct {
	UserServices *services.UserServices
}

// Routes creates a REST router for the todos resource
func (rs UsersResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /users - read a list of users
	r.Post("/", rs.Create) // POST /users - create a new user and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a users map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /users/{id} - read a single user by :id
		r.Put("/", rs.Update)    // PUT /users/{id} - update a single user by :id
		r.Delete("/", rs.Delete) // DELETE /users/{id} - delete a single user by :id
	})

	return r
}

func (rs UsersResource) List(w http.ResponseWriter, r *http.Request) {
	users, err := rs.UserServices.GetAll()
	if err != nil {
		// Log the error
		log.Println("Error getting users:", err)

		// Return an appropriate response
		http.Error(w, "Error getting users", http.StatusInternalServerError)
		return
	}

	// Convert the 'users' data to JSON
	usersJSON, err := json.Marshal(users)
	if err != nil {
		// Log the error
		log.Println("Error marshaling users to JSON:", err)

		// Return an appropriate response
		http.Error(w, "Error marshaling users to JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to indicate JSON response
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(usersJSON)
}

type userRequestBody struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}

func (rs UsersResource) Create(w http.ResponseWriter, r *http.Request) {
	// get the info from the request body
	var user userRequestBody
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusUnprocessableEntity)
		return
	}

	// Validate the user struct
	validate := validator.New()
	err = validate.Struct(user)
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

	// create the user
	createdUser, err := rs.UserServices.CreateUser(services.NewUser{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})

	fmt.Println(createdUser)
	if err != nil {
		// Return the error message as JSON response

		fmt.Println(err)
		errorMsg := struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}

		// Serialize the error message to JSON
		errorJSON, _ := json.Marshal(errorMsg)

		// Set the appropriate content type and HTTP status code in the response header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		// Write the error JSON to the response
		w.Write(errorJSON)
		return
	}
	// Serialize the created user to JSON
	userJSON, err := json.Marshal(createdUser)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	// Set the appropriate content type in the response header
	w.Header().Set("Content-Type", "application/json")
	// Write the user JSON to the response
	w.Write(userJSON)
}

func (rs UsersResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user get"))
}

func (rs UsersResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user update"))
}

func (rs UsersResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user delete"))
}
