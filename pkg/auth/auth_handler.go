package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Brix101/network-file-manager/internals/utils"
	"github.com/Brix101/network-file-manager/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type AuthResource struct {
	UserServices *users.UserServices
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
	errors := utils.NewErrors(w)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		errors.HttpStatus(http.StatusUnprocessableEntity)
		errors.Raise()
		return
	}

	// Validate the user struct
	validate := validator.New()
	err = validate.Struct(body)
	if err != nil {
		// Get the validation errors from the validator and store them in the map
		for _, err := range err.(validator.ValidationErrors) {
			// Make the field name lowercase
			fieldName := strings.ToLower(err.Field())
			errorMessage := fmt.Sprintf("%s is %s", fieldName, err.Tag())

			errors.Add(fieldName, utils.Error{
				Message: errorMessage,
				Type:    "invalid_type",
			})
		}

		errors.HttpStatus(http.StatusBadRequest)
		errors.Raise()
		return
	}

	user, err := rs.UserServices.GetByEmail(body.Email)
	if err != nil {

		errors.Add("email", utils.Error{
			Message: "User not Found",
			Type:    "invalid_type",
		})

		errors.HttpStatus(http.StatusBadRequest)
		errors.Raise()
		return
	}

	// Compare the password from the database with the provided password
	if user.Password != body.Password {

		errors.Add("password", utils.Error{
			Message: "Incorrect password",
			Type:    "invalid_type",
		})

		errors.HttpStatus(http.StatusBadRequest)
		errors.Raise()
		return
	}
	// Store the user ID in a cookie
	cookie := http.Cookie{
		Name:     "session",
		Value:    user.Id.String(),                    // Use the user ID as the cookie value
		Expires:  time.Now().Add(30 * 24 * time.Hour), // Set the expiration time to 30 days from now
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	userJson := user.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)
}

func (rs AuthResource) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	errors := utils.NewErrors(w)
	errors.HttpStatus(http.StatusUnauthorized)
	// Read the cookie from the request
	cookie, err := r.Cookie("session")
	if err != nil {
		errors.Add("root", utils.Error{
			Message: "Unauthorized",
			Type:    "invalid_type",
		})

		errors.Raise()
		return
	}

	// Extract the user ID from the cookie value
	userID, err := uuid.Parse(cookie.Value)
	if err != nil {
		errors.Add("root", utils.Error{
			Message: "Invalid Cookie",
			Type:    "invalid_type",
		})

		errors.Raise()
		return
	}

	// Fetch the user from the database using the user ID
	user, err := rs.UserServices.GetByID(userID)
	if err != nil {
		errors.Add("root", utils.Error{
			Message: "User not found",
			Type:    "invalid_type",
		})

		errors.Raise()
		return
	}

	userJson := user.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)
}
