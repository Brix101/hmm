package main

import (
	"context"
	"home-server/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// The HTTP Server
	server := &http.Server{Addr: "0.0.0.0:5000", Handler: service()}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

type FileInfo struct {
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	FileType string     `json:"fileType,omitempty"`
	Files    []FileInfo `json:"files,omitempty"`
}

func service() http.Handler {
	r := chi.NewRouter()
	// Set up CORS middleware
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://190.160.15.44:5173"}, // Add your allowed origins here
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value in seconds for the Access-Control-Max-Age header
	})
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	// Create a route along /files that will serve contents from
	// the ./data/ folder.
	workDir, _ := os.Getwd()

	// Specify the relative path to the files folder
	spaPath := filepath.Join(workDir, "client", "dist")

	// Serve static files from the specified directory
	staticDir := "/assets/"
	staticFiles := http.FileServer(http.Dir(filepath.Join(spaPath, staticDir)))
	r.Handle(staticDir+"*", http.StripPrefix(staticDir, staticFiles))

	mediaDir := "/"
	mediaFiles := http.FileServer(http.Dir(filepath.Join(spaPath, mediaDir)))
	r.Handle(mediaDir+"*", http.StripPrefix(mediaDir, mediaFiles))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(spaPath, "index.html")
		http.ServeFile(w, r, indexPath)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(spaPath, "index.html")
		http.ServeFile(w, r, indexPath)
	})

	fileRouter := handlers.FileRouter()

	r.Route("/api", func(r chi.Router) {
		r.Mount("/", fileRouter)
	})

	return r
}
