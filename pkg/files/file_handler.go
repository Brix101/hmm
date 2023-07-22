package files

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Brix101/network-file-manager/internals/utils"
	"github.com/go-chi/chi/v5"
)

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

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func (rs FilesResource) Serve(r chi.Router, path string, root http.FileSystem) {
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

func (rs FilesResource) List(w http.ResponseWriter, r *http.Request) {
	// Extract the nested route from the request URL
	path := strings.TrimPrefix(r.URL.Path, "/")

	// Get the "hidden" query parameter, and set default value to true
	hiddenParam := r.URL.Query().Get("hidden")
	hidden := true

	// Check if "hiddenParam" exists and if it's "false", set "hidden" to false
	if hiddenParam == "false" {
		hidden = false
	}
	// Remove "api/files" segment from the path
	path = strings.TrimPrefix(path, "api/files")

	newReader := reader{
		folderPath: rs.FilesPath + path,
		basePath:   rs.FilesPath,
		hidden:     hidden,
	}

	folderStructure, err := folderStructureReader(newReader)
	if err != nil {
		panic("Error building folder structure:" + err.Error())
	}

	jsonData := folderStructure.ToJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

const MAX_UPLOAD_SIZE = 1024 * 1024 * 10 // 1MB

func (rs FilesResource) Create(w http.ResponseWriter, r *http.Request) {
	errors := utils.NewErrors(w)
	clearTempFiles()

	path := strings.TrimPrefix(r.URL.Path, "/")
	path = strings.TrimPrefix(path, "api/files")
	activePath := rs.FilesPath + path
	newReader := reader{
		folderPath: rs.FilesPath + path,
		basePath:   rs.FilesPath,
		hidden:     false,
	}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors.HttpStatus(http.StatusInternalServerError)
		errors.Add("root", utils.Error{
			Message: err.Error(),
			Type:    "invalid_type",
		})
		errors.Raise()
		return
	}
	// Retrieve the name value from the form data
	name := r.FormValue("name")

	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["files"]
	if name == "" && len(files) <= 1 {
		errors.HttpStatus(http.StatusUnprocessableEntity)
		errors.Add("root", utils.Error{
			Message: "Either 'name' or 'files' field is required",
			Type:    "invalid_type",
		})
		errors.Raise()
		return
	}

	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			errors.HttpStatus(http.StatusInternalServerError)
			errors.Add("root", utils.Error{
				Message: "Error open file:" + err.Error(),
				Type:    "invalid_type",
			})
			errors.Raise()
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			errors.HttpStatus(http.StatusInternalServerError)
			errors.Add("root", utils.Error{
				Message: "Error read file:" + err.Error(),
				Type:    "invalid_type",
			})
			errors.Raise()
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			errors.HttpStatus(http.StatusInternalServerError)
			errors.Add("root", utils.Error{
				Message: "Error seek file:" + err.Error(),
				Type:    "invalid_type",
			})
			errors.Raise()
			return
		}

		// Create the destination file in the activePath directory
		destinationFile, err := os.Create(filepath.Join(activePath, fileHeader.Filename))
		if err != nil {
			errors.HttpStatus(http.StatusInternalServerError)
			errors.Add("root", utils.Error{
				Message: "Error creating path:" + err.Error(),
				Type:    "invalid_type",
			})
			errors.Raise()
			return
		}
		defer destinationFile.Close()

		// Copy the contents of the uploaded file to the destination file
		_, err = io.Copy(destinationFile, file)
		if err != nil {
			errors.HttpStatus(http.StatusInternalServerError)
			errors.Add("root", utils.Error{
				Message: "Error copying file:" + err.Error(),
				Type:    "invalid_type",
			})
			errors.Raise()
			return
		}

	}

	if name != "" {
		newFolderPath := filepath.Join(activePath, name)
		err := os.Mkdir(newFolderPath, 0755)
		if err != nil {
			errors.HttpStatus(http.StatusInternalServerError)
			errors.Add("root", utils.Error{
				Message: err.Error(),
				Type:    "invalid_type",
			})
			errors.Raise()
			return
		}
	}

	folderStructure, err := folderStructureReader(newReader)
	if err != nil {
		panic("Error building folder structure:" + err.Error())
	}

	jsonData := folderStructure.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func clearTempFiles() {
	// Get a list of files that start with "multipart" in the /tmp directory
	files, err := filepath.Glob("/tmp/multipart*")
	if err != nil {
		fmt.Println("Error listing files:", err)
		return
	}

	// Remove each file
	for _, file := range files {
		err := os.RemoveAll(file)
		if err != nil {
			fmt.Println("Error removing file:", file, err)
		} else {
			fmt.Println("Removed file:", file)
		}
	}
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

type FileInfo struct {
	PATH     string     `json:"path,omitempty"`
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	FileType string     `json:"fileType,omitempty"`
	Files    []FileInfo `json:"files,omitempty"`
	IsDir    bool       `json:"isDir"`
	ModTime  time.Time  `json:"modTime"`
}

type reader struct {
	folderPath string
	basePath   string
	hidden     bool
}

func (fileInfo FileInfo) ToJSON() []byte {
	jsonData, err := json.Marshal(fileInfo)
	if err != nil {
		panic(err)
	}

	return jsonData
}

func folderStructureReader(r reader) (FileInfo, error) {
	// Get folder information
	folderStat, err := os.Stat(r.folderPath)
	if err != nil {
		return FileInfo{}, err
	}

	// Create FileInfo for the folder
	trimmedFolderPath := strings.TrimPrefix(r.folderPath, r.basePath)

	folder := FileInfo{
		PATH:    trimmedFolderPath,
		Name:    folderStat.Name(),
		Size:    folderStat.Size(),
		ModTime: folderStat.ModTime(),
		IsDir:   true,
		Files:   []FileInfo{},
	}

	// Read folder contents
	fileStats, err := ioutil.ReadDir(r.folderPath)
	if err != nil {
		return FileInfo{}, err
	}

	// Sort fileStats based on IsDir and Name
	sort.SliceStable(fileStats, func(i, j int) bool {
		if fileStats[i].IsDir() && !fileStats[j].IsDir() {
			return true
		}
		if !fileStats[i].IsDir() && fileStats[j].IsDir() {
			return false
		}
		return fileStats[i].Name() < fileStats[j].Name()
	})

	// Traverse through folder contents
	for _, fileStat := range fileStats {
		// Skip hidden files and folders if "hidden" is true and exclude .git if hidden is false
		if (r.hidden && fileStat.Name()[0] == '.') || (!r.hidden && strings.Contains(fileStat.Name(), ".git")) {
			continue
		} // Build file/folder path
		itemPath := filepath.Join(r.folderPath, fileStat.Name())

		if fileStat.IsDir() {
			// Recursively build folder structure for subfolders
			subReader := reader{
				folderPath: itemPath,
				basePath:   r.basePath,
				hidden:     r.hidden,
			}

			subfolder, err := folderStructureReader(subReader)
			if err != nil {
				return FileInfo{}, err
			}

			// Add subfolder to the folder's items
			folder.Files = append(folder.Files, subfolder)
		} else {
			// Get the file type and MIME type
			fileType := filepath.Ext(fileStat.Name())
			mimeType := mime.TypeByExtension(fileType)

			// Create FileInfo for the file
			file := FileInfo{
				Name:     fileStat.Name(),
				Size:     fileStat.Size(),
				PATH:     folder.PATH + "/" + fileStat.Name(),
				FileType: fileType,
				ModTime:  fileStat.ModTime(),
				IsDir:    false,
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
