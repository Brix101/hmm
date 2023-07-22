package main

import (
	"context"
	"home-server/internal/handlers"
	"home-server/internal/middlewares"
	"home-server/internal/services"
	"home-server/pkg/db"
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
	"github.com/jmoiron/sqlx"
)

func main() {
	// The HTTP Server
	conn := db.CreateConnectionPool()

	server := &http.Server{Addr: "0.0.0.0:5000", Handler: service(conn)}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel() // Call the cancel function when the shutdown function finishes

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
	log.Println("🚀🚀🚀 Server at http://" + server.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func service(conn *sqlx.DB) http.Handler {
	r := chi.NewRouter()
	workDir, _ := os.Getwd()

	// Specify the relative path to the files folder
	filesPath := filepath.Join(workDir, "data")
	filesDir := http.Dir(filesPath)
	// Set up CORS middleware
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://192.168.254.180:5173"}, // Add your allowed origins here
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value in seconds for the Access-Control-Max-Age header
	})

	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middlewares.ErrorHandler)

	clientRouter := handlers.ClientRouter()
	userServices := services.NewUserServices(conn)
	featureServices := services.NewFeatureServices(conn)

	filesResource := handlers.FilesResource{FilesPath: filesPath}

	userResource := handlers.UsersResource{UserServices: userServices}
	authResource := handlers.AuthResource{UserServices: userServices}
	featureResource := handlers.FeaturesResource{FeatureServices: featureServices}

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
