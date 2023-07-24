package files

import (
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Reader struct {
	basePath   string
	activePath string
	hidden     bool
}

func NewReader(activePath string, hidden bool) Reader {
	workDir, _ := os.Getwd()

	// Specify the relative path to the files folder
	filesPath := filepath.Join(workDir, "data")

	itemPath := filepath.Join(filesPath, activePath)

	return Reader{
		basePath:   filesPath,
		activePath: itemPath,
		hidden:     hidden,
	}
}

func (r *Reader) Set(activePath string, hidden bool) {
	workDir, _ := os.Getwd()

	// Specify the relative path to the files folder
	filesPath := filepath.Join(workDir, "data")

	itemPath := filepath.Join(filesPath, activePath)

	r.activePath = itemPath
	r.hidden = hidden
}

func (r *Reader) GetContent(activePath string, hidden bool) (FileInfo, error) {
	r.Set(activePath, hidden)

	// Get folder information
	folderStat, err := os.Stat(r.activePath)
	if err != nil {
		return FileInfo{}, err
	}

	// Create FileInfo for the folder
	trimmedFolderPath := strings.TrimPrefix(r.activePath, r.basePath)

	folder := FileInfo{
		PATH:    trimmedFolderPath,
		Name:    folderStat.Name(),
		Size:    folderStat.Size(),
		ModTime: folderStat.ModTime(),
		IsDir:   true,
		Files:   []FileInfo{},
	}

	// Read folder contents
	fileStats, err := ioutil.ReadDir(r.activePath)
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

		// if fileStat.IsDir() {
		// 	// Recursively build folder structure for subfolders
		// 	// subContent := NewReader(fileStat.Name(), r.hidden)
		// 	//
		// 	// subfolder, err := subContent.GetContent()
		// 	// if err != nil {
		// 	// 	return FileInfo{}, err
		// 	// }
		// 	//
		// 	// // Add subfolder to the folder's items
		// 	folder.Files = append(folder.Files)
		// } else {
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
			IsDir:    fileStat.IsDir(),
		}

		// Set the MIME type if available
		if mimeType != "" {
			file.FileType = mimeType
		}

		// Add file to the folder's items
		folder.Files = append(folder.Files, file)
		// }
	}

	return folder, nil
}
