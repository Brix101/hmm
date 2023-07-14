package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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
		AllowedOrigins:   []string{"http://localhost:5173"}, // Add your allowed origins here
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
	filesPath := filepath.Join(workDir, "files")
	spaPath := filepath.Join(workDir, "client", "dist")

	filesDir := http.Dir(filesPath)
	FileServer(r, "/static", filesDir)

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

	r.Route("/api", func(r chi.Router) {
		r.Post("/files", func(w http.ResponseWriter, r *http.Request) {
			// Parse the JSON request body to retrieve the folder name
			type FolderRequest struct {
				Name string `json:"name"`
			}
			var folderReq FolderRequest
			err := json.NewDecoder(r.Body).Decode(&folderReq)
			if err != nil {
				// Log the error
				log.Println("Error parsing request body:", err)

				// Return an appropriate response
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Create the new folder
			newFolderPath := filepath.Join(filesPath, folderReq.Name)
			err = os.Mkdir(newFolderPath, 0755)
			if err != nil {
				// Log the error
				log.Println("Error creating folder:", err)

				// Split the error message
				errorMessage := strings.SplitN(err.Error(), ":", 2)
				errorDetail := ""
				if len(errorMessage) > 1 {
					errorDetail = strings.TrimSpace(errorMessage[1])
				} else {
					errorDetail = strings.TrimSpace(errorMessage[0])
				}

				// Create an error response
				errorResponse := struct {
					Detail string `json:"detail"`
				}{
					Detail: errorDetail,
				}
				jsonResponse, err := json.Marshal(errorResponse)
				if err != nil {
					// Log the error
					log.Println("Error encoding JSON:", err)

					// Return an appropriate response
					http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
					return
				}

				// Set response headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonResponse)
				return
			}

			// Return a success response
			w.WriteHeader(http.StatusCreated)
		})

		r.Get("/files", func(w http.ResponseWriter, r *http.Request) {
			// Build the folder structure
			folderStructure, err := FolderStructure(filesPath)
			if err != nil {
				// Log the error
				log.Println("Error building folder structure:", err)

				// Return an appropriate response
				http.Error(w, "Error building folder structure", http.StatusInternalServerError)
				return
			}

			// Convert folder structure to JSON
			jsonData, err := json.Marshal(folderStructure)
			if err != nil {
				// Log the error
				log.Println("Error encoding JSON:", err)

				// Return an appropriate response
				http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
				return
			}

			// Set response headers
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Write the JSON response
			w.Write(jsonData)
		})

		r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
			// Simulates some hard work.
			//
			// We want this handler to complete successfully during a shutdown signal,
			// so consider the work here as some background routine to fetch a long running
			// search query to find as many results as possible, but, instead we cut it short
			// and respond with what we have so far. How a shutdown is handled is entirely
			// up to the developer, as some code blocks are preemptible, and others are not.
			time.Sleep(5 * time.Second)

			w.Write([]byte(fmt.Sprintf("all done.\n")))
		})
	})
	return r
}

func FolderStructure(folderPath string) (FileInfo, error) {
	// Get folder information
	folderInfo, err := os.Stat(folderPath)
	if err != nil {
		return FileInfo{}, err
	}

	// Create FileInfo for the folder
	folder := FileInfo{
		Name:  folderInfo.Name(),
		Size:  folderInfo.Size(),
		Files: []FileInfo{},
	}

	// Read folder contents
	fileInfos, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return FileInfo{}, err
	}

	// Traverse through folder contents
	for _, fileInfo := range fileInfos {
		// Skip hidden files and folders
		if fileInfo.Name()[0] == '.' {
			continue
		}

		// Build file/folder path
		itemPath := filepath.Join(folderPath, fileInfo.Name())

		if fileInfo.IsDir() {
			// Recursively build folder structure for subfolders
			subfolder, err := FolderStructure(itemPath)
			if err != nil {
				return FileInfo{}, err
			}

			// Add subfolder to the folder's items
			folder.Files = append(folder.Files, subfolder)
		} else {
			// Get the file type and MIME type
			fileType := filepath.Ext(fileInfo.Name())
			mimeType := mime.TypeByExtension(fileType)

			// Create FileInfo for the file
			file := FileInfo{
				Name:     fileInfo.Name(),
				Size:     fileInfo.Size(),
				FileType: fileType,
			}

			// Set the MIME type if available
			if mimeType != "" {
				file.FileType = mimeType
			}

			// Add file to the folder's items
			folder.Files = append(folder.Files, file)
		}
	}

	return folder, nil
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
