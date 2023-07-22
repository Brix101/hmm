package handlers

import (
	"encoding/json"
	"home-server/internal/services"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FeaturesResource struct {
	FeatureServices *services.FeatureServices
}

// Routes creates a REST router for the features resource
func (rs FeaturesResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List) // GET /features - read a list of features
	// r.Post("/", rs.Create) // POST /features - create a new todo and persist it
	// r.Put("/", rs.Delete)
	//
	// r.Route("/{id}", func(r chi.Router) {
	// 	// r.Use(rs.TodoCtx) // lets have a features map, and lets actually load/manipulate
	// 	r.Get("/", rs.Get)       // GET /features/{id} - read a single todo by :id
	// 	r.Put("/", rs.Update)    // PUT /features/{id} - update a single todo by :id
	// 	r.Delete("/", rs.Delete) // DELETE /features/{id} - delete a single todo by :id
	// 	r.Get("/sync", rs.Sync)
	// })

	return r
}

func (rs FeaturesResource) List(w http.ResponseWriter, r *http.Request) {
	features, err := rs.FeatureServices.GetAll()
	if err != nil {
		// Log the error
		log.Println("Error getting features:", err)

		// Return an appropriate response
		http.Error(w, "Error getting features", http.StatusInternalServerError)
		return
	}

	// Convert the 'features' data to JSON
	featuresJSON, err := json.Marshal(features)
	if err != nil {
		// Log the error
		log.Println("Error marshaling features to JSON:", err)

		// Return an appropriate response
		http.Error(w, "Error marshaling features to JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to indicate JSON response
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(featuresJSON)
}

// func (rs FeaturesResource) Create(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("features create"))
// }
//
// func (rs FeaturesResource) Get(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("todo get"))
// }
//
// func (rs FeaturesResource) Update(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("todo update"))
// }
//
// func (rs FeaturesResource) Delete(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("todo delete"))
// }
//
// func (rs FeaturesResource) Sync(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("todo sync"))
// }
