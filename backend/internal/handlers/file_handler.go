package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const fileSizeLimitMb = 5

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing chatID", http.StatusBadRequest)
		return
	}
	chatID := parts[3]

	// TODO: integrate with SCRUM 110
	response := map[string]string{
		"message": "Fetched files for chatID " + chatID,
	}

	sendJSONResponse(w, response)
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(fileSizeLimitMb << 20)
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	chatID := r.FormValue("chatID")
	if chatID == "" {
		http.Error(w, "Missing chatID", http.StatusBadRequest)
		return
	}

	dst, err := os.Create(filepath.Join("uploads", header.Filename))
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "File uploaded successfully",
		"chatID":  chatID,
		"file":    header.Filename,
	}

	sendJSONResponse(w, response)
}
