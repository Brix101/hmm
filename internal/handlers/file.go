package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type FileInfo struct {
	Name      string     `json:"name"`
	Size      int64      `json:"size"`
	FileType  string     `json:"fileType,omitempty"`
	PATH      string     `json:"path,omitempty"`
	Files     []FileInfo `json:"files,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type FilesResource struct{ FilesPath string }

// Routes creates a REST router for the todos resource
func (rs FilesResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/*", rs.List)    // GET /todos - read a list of todos
	r.Post("/*", rs.Create) // POST /todos - create a new todo and persist it
	// r.Put("/", rs.Delete)

	// r.Route("/{id}", func(r chi.Router) {
	// 	// r.Use(rs.TodoCtx) // lets have a todos map, and lets actually load/manipulate
	// 	r.Get("/", rs.Get)       // GET /todos/{id} - read a single todo by :id
	// 	r.Put("/", rs.Update)    // PUT /todos/{id} - update a single todo by :id
	// 	r.Delete("/", rs.Delete) // DELETE /todos/{id} - delete a single todo by :id
	// 	r.Get("/sync", rs.Sync)
	// })

	return r
}

func (rs FilesResource) List(w http.ResponseWriter, r *http.Request) {
	// Extract the nested route from the request URL
	path := strings.TrimPrefix(r.URL.Path, "/")

	// Remove "api/files" segment from the path
	path = strings.TrimPrefix(path, "api/files")

	folderStructure, err := folderStructureReader(rs.FilesPath+path, rs.FilesPath)
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
}

func (rs FilesResource) Create(w http.ResponseWriter, r *http.Request) {
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
	newFolderPath := filepath.Join(rs.FilesPath+path, folderReq.Name)
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
		Name:      folderInfo.Name(),
		Size:      folderInfo.Size(),
		PATH:      trimmedFolderPath,
		CreatedAt: folderInfo.ModTime(),
		Files:     []FileInfo{},
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
				Name:      fileInfo.Name(),
				Size:      fileInfo.Size(),
				PATH:      folder.PATH + "/" + fileInfo.Name(),
				FileType:  fileType,
				CreatedAt: fileInfo.ModTime(),
			}

			// Set the MIME type if available
			if mimeType != "" {
				file.FileType = mimeType
			}

			// Add file to the folder's items
			folder.Files = append(folder.Files, file)
		}
	}

	// Sort the folder.Files slice by file type for folders, and by name for files
	sort.Slice(folder.Files, func(i, j int) bool {
		isFolderI := folder.Files[i].Files != nil && len(folder.Files[i].Files) > 0
		isFolderJ := folder.Files[j].Files != nil && len(folder.Files[j].Files) > 0

		if isFolderI && !isFolderJ {
			return true // Folders come before files
		} else if !isFolderI && isFolderJ {
			return false // Files come after folders
		} else if isFolderI && isFolderJ {
			return folder.Files[i].FileType < folder.Files[j].FileType // Sort folders by file type
		}

		// Both are files, sort by name
		return folder.Files[i].Name < folder.Files[j].Name
	})

	return folder, nil
}
