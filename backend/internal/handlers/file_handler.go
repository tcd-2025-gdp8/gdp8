package handlers

import (
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
	fmt.Println("GetFiles called with chatID:", chatID)

	// TODO: integrate with SCRUM 110

	response := map[string]string{
		"message": "Fetched files for chatID " + chatID,
	}

	sendJSONResponse(w, response)
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
	dirPath := "uploads/" + chatID
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
