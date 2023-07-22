package pkg

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Brix101/network-file-manager/pkg/auth"
	"github.com/Brix101/network-file-manager/pkg/features"
	"github.com/Brix101/network-file-manager/pkg/files"
	"github.com/Brix101/network-file-manager/pkg/handlers"
	"github.com/Brix101/network-file-manager/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func Initialize(r chi.Router, conn *sqlx.DB) http.Handler {
	workDir, _ := os.Getwd()

	// Specify the relative path to the files folder
	filesPath := filepath.Join(workDir, "data")
	filesDir := http.Dir(filesPath)

	clientRouter := handlers.ClientRouter()
	userServices := users.NewUserServices(conn)
	featureServices := features.NewFeatureServices(conn)

	filesResource := files.FilesResource{FilesPath: filesPath}
	userResource := users.UsersResource{UserServices: userServices}
	authResource := auth.AuthResource{UserServices: userServices}
	featureResource := features.FeaturesResource{FeatureServices: featureServices}

	r.Mount("/", clientRouter)
	r.Route("/api", func(r chi.Router) {
		r.Mount("/", authResource.Routes())
		r.Mount("/files", filesResource.Routes())
		r.Mount("/users", userResource.Routes())
		r.Mount("/features", featureResource.Routes())
	})

	filesResource.Serve(r, "/data/files", filesDir)

	return r
}
