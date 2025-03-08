package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const fileSizeLimitMb = 20

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
	fmt.Println("GetFiles called with chatID:", chatID)

	// TODO: integrate with SCRUM 110
	response := map[string]string{
		"message": "Fetched files for chatID " + chatID,
	}

	sendJSONResponse(w, response)
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(fileSizeLimitMb << 20)
	if err != nil {
		fmt.Println("Error parsing multipart form:", err)
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving file:", err)
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	chatID := r.FormValue("chatID")
	if chatID == "" {
		fmt.Println("chatID is missing in form data")
		http.Error(w, "Missing chatID", http.StatusBadRequest)
		return
	}

	userID := r.FormValue("userID")
	if userID == "" {
		fmt.Println("userID is missing in form data")
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Println("Error creating uploads directory:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	dstPath := filepath.Join(uploadDir, header.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error writing file:", err)
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}
	response := map[string]string{
		"message": "File uploaded successfully",
		"chatID":  chatID,
		"file":    header.Filename,
		"userID":  userID,
	}

	sendJSONResponse(w, response)
}
