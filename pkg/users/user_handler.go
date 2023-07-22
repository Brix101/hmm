package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Brix101/network-file-manager/internals/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type UsersResource struct {
	UserServices *UserServices
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
		panic(err)
	}

	usersJSON := users.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	w.Write(usersJSON)
}

func (rs UsersResource) Create(w http.ResponseWriter, r *http.Request) {
	errors := utils.NewErrors(w)
	var user UserRequestBody
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errors.HttpStatus(http.StatusUnprocessableEntity)
		errors.Raise()
		return
	}

	// Validate the user struct
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		// Get the validation errors from the validator and store them in the map
		for _, err := range err.(validator.ValidationErrors) {
			// Make the field name lowercase
			fieldName := strings.ToLower(err.Field())

			// Append the error message to the field's slice in the map
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

	// create the user
	createdUser, err := rs.UserServices.CreateUser(NewUser{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		errors.Add("root", utils.Error{
			Message: err.Error(),
			Type:    "invalid_type",
		})

		errors.HttpStatus(http.StatusInternalServerError)
		errors.Raise()
		return
	}

	userJSON := createdUser.ToJSON()

	w.Header().Set("Content-Type", "application/json")
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
