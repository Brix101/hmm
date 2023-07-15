package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

type FileInfo struct {
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	FileType string     `json:"fileType,omitempty"`
	PATH     string     `json:"path,omitempty"`
	Files    []FileInfo `json:"files,omitempty"`
}

func FileRouter() http.Handler {
	r := chi.NewRouter()
	workDir, _ := os.Getwd()

	// Specify the relative path to the files folder
	filesPath := filepath.Join(workDir, "data")
	filesDir := http.Dir(filesPath)
	fileServer(r, "/data", filesDir)

	r.Post("/files/*", func(w http.ResponseWriter, r *http.Request) {
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

		// Extract the nested route from the request URL
		path := strings.TrimPrefix(r.URL.Path, "/")

		// Remove "api/files" segment from the path
		path = strings.TrimPrefix(path, "api/files")

		// Create the new folder
		newFolderPath := filepath.Join(filesPath+path, folderReq.Name)
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

	r.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
		// Extract the nested route from the request URL
		path := strings.TrimPrefix(r.URL.Path, "/")

		// Remove "api/files" segment from the path
		path = strings.TrimPrefix(path, "api/files")

		folderStructure, err := folderStructureReader(filesPath+path, filesPath)
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

	/*
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
	*/
	return r
}

// static files from a http.FileSystem.
func fileRouteReader(r chi.Router, path string, root http.FileSystem) {
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

func folderStructureReader(folderPath string, basePath string) (FileInfo, error) {
	// Get folder information
	folderInfo, err := os.Stat(folderPath)
	if err != nil {
		return FileInfo{}, err
	}

	// Create FileInfo for the folder
	trimmedFolderPath := strings.TrimPrefix(folderPath, basePath)

	folder := FileInfo{
		Name:  folderInfo.Name(),
		Size:  folderInfo.Size(),
		PATH:  trimmedFolderPath,
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
			subfolder, err := folderStructureReader(itemPath, basePath)
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
				PATH:     folder.PATH + "/" + fileInfo.Name(),
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

// fileServer conveniently sets up a http.fileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
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
