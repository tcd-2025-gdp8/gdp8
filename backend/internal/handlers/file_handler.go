package handlers

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
)

const fileSizeLimitMb = 20

type FileHandler struct {
	studyGroupService services.StudyGroupService
}

func NewFileHandler(studyGroupService services.StudyGroupService) *FileHandler {
	return &FileHandler{
		studyGroupService: studyGroupService,
	}
}

func (h *FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing chatID", http.StatusBadRequest)
		return
	}
	chatID := parts[3]

	files, err := h.getFilesByID(chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"files.zip\"")
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	for _, file := range files {
		f, err := zipWriter.Create(file.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = f.Write(file.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(fileSizeLimitMb << 20); err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	chatID, userID, err := h.extractFormValues(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.isUserMember(chatID, userID) {
		http.Error(w, "User not a member of the study group", http.StatusForbidden)
		return
	}

	if err := h.saveFile(file, header.Filename, chatID); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, map[string]string{"message": "File uploaded successfully"})
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	chatID, userID, err := h.extractFormValues(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	filename := r.FormValue("filename")
	canDelete := h.hasDeletionRights(filename, chatID, userID)

	if !canDelete {
		http.Error(w, "Cannot delete the file", http.StatusUnauthorized)
	}

	deleted, err := h.deleteFileByName(filename)
	if err != nil {
		http.Error(w, "Couldn't delete the file", http.StatusInternalServerError)

	}
	status := "File uploaded successfully"
	if !deleted {
		status = fmt.Sprintf("The file %s does not exist", filename)

	}
	sendJSONResponse(w, map[string]string{"message": status})
}

func (h *FileHandler) extractFormValues(r *http.Request) (string, string, error) {
	chatID := r.FormValue("chatID")
	if chatID == "" {
		return "", "", errors.New("Missing chatID")
	}

	userID := r.FormValue("userID")
	if userID == "" {
		return "", "", errors.New("Missing userID")
	}

	return chatID, userID, nil
}

func (h *FileHandler) isUserMember(chatID string, userID string) bool {
	studyGroupID, err := strconv.Atoi(chatID)
	if err != nil {
		return false
	}

	studyGroup, err := h.studyGroupService.GetStudyGroupByID(models.StudyGroupID(studyGroupID))
	if err != nil {
		return false
	}

	return isUserMemberOfStudyGroup(models.UserID(userID), studyGroup)
}

func (h *FileHandler) saveFile(file io.Reader, filename string, chatID string) error {
	filename = filepath.Base(filename)
	if filename == "." || filename == "" {
		return errors.New("invalid filename")
	}

	dirPath := filepath.Join("uploads", chatID)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	dstPath := filepath.Join(dirPath, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return err
}

func isUserMemberOfStudyGroup(userID models.UserID, studyGroup *models.StudyGroupView) bool {
	for _, member := range studyGroup.Members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}
func (h *FileHandler) getFilesByID(chatID string) ([]File, error) {
	dirPath := filepath.Join("uploads", chatID)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var files []File
	for _, entry := range entries {
		if !entry.IsDir() {
			content, err := os.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, File{Name: entry.Name(), Content: content})
		}
	}

	return files, nil
}

func (h *FileHandler) deleteFileByName(filename string) (bool, error) {
	_ = h
	_ = filename
	return true, nil
}

func (h *FileHandler) hasDeletionRights(filename string, chatID string, userID string) bool {
	_ = h
	_ = filename
	_ = chatID
	_ = userID
	return true
}

type File struct {
	Name    string
	Content []byte
}
