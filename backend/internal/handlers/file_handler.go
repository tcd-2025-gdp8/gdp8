package handlers

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"

	"gdp8-backend/internal/chatbot"
	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
)

const fileSizeLimitMb = 20

type File struct {
	Name    string
	Content []byte
}

type FileHandler struct {
	studyGroupService services.StudyGroupService
	fileService       services.FileService
	chatbotService    services.ChatBotService
}

var ErrFiletypeNotAllowed = errors.New("filetype not allowed")

func NewFileHandler(studyGroupService services.StudyGroupService,
	fileService services.FileService,
	chatbotService services.ChatBotService) *FileHandler {
	return &FileHandler{
		studyGroupService: studyGroupService,
		fileService:       fileService,
		chatbotService:    chatbotService,
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
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unabe to retrieve userID", http.StatusForbidden)
		return
	}
	groupID, err := parseGroupID(chatID)
	if err != nil {
		http.Error(w, "Unabe to retrieve groupID", http.StatusForbidden)
		return
	}
	if !h.isUserMember(groupID, userID) {
		http.Error(w, "User not a member of the study group", http.StatusForbidden)
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
	if err := r.ParseMultipartForm(fileSizeLimitMb << fileSizeLimitMb); err != nil {
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
	if err = h.saveFile(file, header.Filename, chatID); err != nil {
		if errors.Is(err, ErrFiletypeNotAllowed) {
			http.Error(w, "File type not allowed", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		log.Printf("Error saving the file: %s", err)
		return
	}
	if err = h.createFile(header.Filename, chatID, userID); err != nil {
		http.Error(w, "Error storing the file", http.StatusInternalServerError)
		log.Printf("Error storing the file: %s", err)
		return
	}
	if strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		filePath := filepath.Join("uploads", fmt.Sprintf("%d", chatID), header.Filename)
		content, err := chatbot.ReadPDF(filePath)
		if err == nil {
			groupID, err := parseGroupID(fmt.Sprintf("%v", chatID))
			if err == nil {
				_ = h.chatbotService.StoreMemory(groupID, &models.FileContext{
					Name: header.Filename,
					Data: content,
				})
			}
		}
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
		http.Error(w, "Cannot delete the file", http.StatusForbidden)
		return
	}
	err = h.deleteFileByName(filename, chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	groupID, err := parseGroupID(fmt.Sprintf("%v", chatID))
	if err == nil {
		_ = h.chatbotService.DeleteMemory(groupID, filename)
	}
	sendJSONResponse(w, map[string]string{"message": fmt.Sprintf("the file '%s' has been deleted", filename)})
}

func (h *FileHandler) extractFormValues(r *http.Request) (models.StudyGroupID, models.UserID, error) {
	chatID := r.FormValue("chatID")
	groupID, err := parseGroupID(chatID)
	if err != nil {
		return models.StudyGroupID(0), models.UserID(""), errors.New("Invalid chatID")
	}
	userID, err := getUserID(r)
	if err != nil {
		return models.StudyGroupID(0), models.UserID(""), errors.New("Invalid user")
	}
	return groupID, userID, nil
}

func (h *FileHandler) isUserMember(studyGroupID models.StudyGroupID, userID models.UserID) bool {
	studyGroup, err := h.studyGroupService.GetStudyGroupByID(studyGroupID)
	if err != nil {
		return false
	}
	return isUserMemberOfStudyGroup(userID, studyGroup)
}

func (h *FileHandler) saveFile(file io.Reader, filename string, chatID models.StudyGroupID) error {
	filename = filepath.Base(filename)
	if filename == "." || filename == "" {
		return errors.New("invalid filename")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	kind, err := filetype.Match(data)
	if err != nil {
		return errors.New("failed to detect filetype")
	}
	if !isFiletypeAllowed(kind) {
		return ErrFiletypeNotAllowed
	}

	dirPath := fmt.Sprintf("uploads/%d", chatID)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}
	dstPath := filepath.Join(dirPath, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.Write(data)
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

func (h *FileHandler) deleteFileByName(filename string, groupID models.StudyGroupID) error {
	file := fmt.Sprintf("uploads/%d/%s", groupID, filename)
	return os.Remove(file)
}

func (h *FileHandler) hasDeletionRights(filename string, groupID models.StudyGroupID, userID models.UserID) bool {
	rights, err := h.fileService.HasDeletionRights(filename, groupID, userID)
	if err != nil {
		return false
	}
	return rights
}

func (h *FileHandler) createFile(filename string, groupID models.StudyGroupID, userID models.UserID) error {
	newFile := models.File{
		Name:    filename,
		UserID:  userID,
		GroupID: groupID,
	}
	_, err := h.fileService.CreateFile(newFile)
	return err
}

func parseGroupID(chatID string) (models.StudyGroupID, error) {
	id, err := strconv.Atoi(chatID)
	if err != nil {
		return models.StudyGroupID(0), err
	}
	return models.StudyGroupID(id), nil
}

func isFiletypeAllowed(kind types.Type) bool {
	if kind == filetype.Unknown {
		return false
	}

	allowed := map[string]bool{
		// Documents
		"application/pdf":    true,
		"application/msword": true, // .doc
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
		"application/vnd.oasis.opendocument.text":                                 true, // .odt

		// Spreadsheets
		"application/vnd.ms-excel": true, // .xls
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true, // .xlsx
		"application/vnd.oasis.opendocument.spreadsheet":                    true, // .ods
		"text/csv": true, // .csv

		// Presentations
		"application/vnd.ms-powerpoint":                                             true, // .ppt
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // .pptx
		"application/vnd.oasis.opendocument.presentation":                           true, // .odp

		// Images
		"image/png":     true,
		"image/jpeg":    true,
		"image/gif":     true,
		"image/svg+xml": true,

		// Plain Text & Markdown
		"text/plain":    true, // .txt
		"text/markdown": true, // .md
	}

	return allowed[kind.MIME.Value]
}
