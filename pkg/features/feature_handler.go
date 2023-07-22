package features

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FeaturesResource struct {
	FeatureServices *FeatureServices
}

// Routes creates a REST router for the features resource
func (rs FeaturesResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List) // GET /features - read a list of features

	return r
}

func (rs FeaturesResource) List(w http.ResponseWriter, r *http.Request) {
	features, err := rs.FeatureServices.GetAll()
	if err != nil {
		panic(err)
	}

	featuresJSON := features.ToJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(featuresJSON)
}
