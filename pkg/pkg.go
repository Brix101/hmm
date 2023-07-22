package pkg

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Brix101/network-file-manager/pkg/auth"
	"github.com/Brix101/network-file-manager/pkg/features"
	"github.com/Brix101/network-file-manager/pkg/files"
	"github.com/Brix101/network-file-manager/pkg/handlers"
	"github.com/Brix101/network-file-manager/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Test struct {
	Name  string
	Age   int
	Email string
}

func Initialize(r chi.Router, conn *sqlx.DB) http.Handler {
	workDir, _ := os.Getwd()

	// Specify the relative path to the files folder
	filesPath := filepath.Join(workDir, "data")
	templatesPath := filepath.Join(workDir, "templates")
	filesDir := http.Dir(filesPath)

	clientRouter := handlers.ClientRouter()
	userServices := users.NewUserServices(conn)
	featureServices := features.NewFeatureServices(conn)

	filesResource := files.FilesResource{FilesPath: filesPath}
	userResource := users.UsersResource{UserServices: userServices}
	authResource := auth.AuthResource{UserServices: userServices}
	featureResource := features.FeaturesResource{FeatureServices: featureServices}

	clientTest(r, templatesPath)

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

type Film struct {
	Title    string
	Director string
}

//go:embed all:templates/*
var tmplFS embed.FS

func clientTest(r chi.Router, templatesPath string) {
	staticFiles := http.FileServer(http.Dir(templatesPath + "/static"))

	fmt.Println("++++++++++++++++++++", tmplFS)
	films := map[string][]Film{
		"Films": {
			{Title: "The Godfather", Director: "Francis Ford Coppola"},
			{Title: "Blade Runner", Director: "Ridley Scott"},
			{Title: "The Thing", Director: "John Carpenter"},
		},
	}
	// handler function #1 - returns the index.html template, with film data
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(templatesPath + "/index.html"))

		tmpl.Execute(w, films)
	}

	// handler function #2 - returns the template block with the newly added film, as an HTMX response
	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")
		// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
		// tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl := template.Must(template.ParseFiles(templatesPath + "/index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})
	}
	r.Handle("/static/*", http.StripPrefix("/static/", staticFiles))
	// define handlers
	r.HandleFunc("/asd", h1)
	r.HandleFunc("/asd/add-film/", h2)
}
