import { useState, DragEvent, useEffect, useCallback } from "react";
import {
    Box,
    Typography,
    CircularProgress,
    Paper,
    List,
    ListItem,
    ListItemText,
    IconButton,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button
} from "@mui/material";
import {
    UploadFile as UploadIcon,
    CloudUpload as CloudUploadIcon,
    Delete as DeleteIcon
} from "@mui/icons-material";
import ChatbotChat from "../components/ChatbotChat";
import { BackendFile, apiFetchFiles, apiUploadFiles, apiDeleteFiles} from "../utils/apiFile"
import { useAuth } from "../auth/useAuth";

const getChatID = () => {
  const parts = window.location.toString().split("/");
  for (let i = 0; i < parts.length; i++) {
    if (parts[i] === "study-groups" && i + 1 < parts.length) {
      const candidate = parts[i + 1];
      if (!isNaN(Number(candidate))) {
        return candidate;
      }
    }
  }
  return "invalid";
}


export default function FilesPage() {
    const [files, setFiles] = useState<BackendFile[]>([]);
    const [uploading, setUploading] = useState(false);
    const [dragging, setDragging] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [fileToDelete, setFileToDelete] = useState<string | null>(null);
    const [chatbotOpen, setChatbotOpen] = useState(false);

    const chatID = getChatID()
    const { token } = useAuth();

    const fetchFiles = useCallback(async () => {
        if (!chatID) return;
        try {
            const fetchedFiles = await apiFetchFiles(chatID, token);
            setFiles(fetchedFiles);
        } catch (error) {
            console.error("Error fetching files:", error);
        }
    }, [chatID, token]);

    useEffect(() => {
        void fetchFiles();
    }, [chatID, token, fetchFiles]);


    const handleFileUpload = async (selectedFiles: FileList | null): Promise<void> => {
        if (!selectedFiles || selectedFiles.length === 0 || !chatID) return;
        setUploading(true);
        try {
            for (const file of Array.from(selectedFiles)) {
                await apiUploadFiles(file, chatID, token);
            }
            await fetchFiles();

            const length = Array.from(selectedFiles).length;
            const popupMessage = length == 1 ?
                                "Your file has been uploaded successfully!" :
                                "Your files have been uploaded successfully!";
            alert(popupMessage);
        } catch (error) {
            const errorMessage = (error as Error)?.message;
            if (errorMessage.includes("File type not allowed")) {
                alert("File type not allowed.")
            }
            console.error("Error uploading file(s):", error);
        } finally {
            setUploading(false);
        }
    };


    const handleDeleteConfirmed = async (): Promise<void> => {
        if (fileToDelete && chatID) {
            try {
                await apiDeleteFiles(fileToDelete, chatID, token);
                await fetchFiles();
            } catch (error) {
                console.error("Error deleting file:", error);
            }
        }
        setDeleteDialogOpen(false);
        setFileToDelete(null);
    };



    const confirmDeleteFile = (filename: string): void => {
        setFileToDelete(filename);
        setDeleteDialogOpen(true);
    };

    const handleDragOver = (e: DragEvent<HTMLDivElement>): void => {
        e.preventDefault();
        setDragging(true);
    };

    const handleDragLeave = (): void => {
        setDragging(false);
    };

    const handleDrop = (e: DragEvent<HTMLDivElement>): void => {
        e.preventDefault();
        setDragging(false);
        void handleFileUpload(e.dataTransfer.files);
    };

    return (
        <>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    justifyContent: "flex-start",
                    width: "calc(87% - 260px)", //do not change
                    bgcolor: "#f0f8ff",
                    padding: 3,
                    gap: 1,
                    marginLeft: "240px",
                    boxSizing: "border-box",
                }}
            >
                <Typography variant="h4" fontWeight="bold" color="primary" gutterBottom>
                    File Upload
                </Typography>

                <Paper
                    sx={{
                        width: "100%",
                        maxWidth: "780px",
                        height: "120px", //do not change
                        border: "3px dashed",
                        borderColor: dragging ? "#1976d2" : "#90caf9",
                        borderRadius: "12px",
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        flexDirection: "column",
                        backgroundColor: "#ffffff",
                        transition: "0.3s",
                        cursor: "pointer",
                        "&:hover": {
                            borderColor: "#1565c0",
                            backgroundColor: "#e3f2fd",
                        },
                    }}
                    onDragOver={handleDragOver}
                    onDragLeave={handleDragLeave}
                    onDrop={handleDrop}
                    onClick={() => document.getElementById("fileInput")?.click()}
                >
                    <CloudUploadIcon sx={{ fontSize: 80, color: "#1976d2", mb: 1 }} />
                    <Typography variant="h6" color="textSecondary">
                        Drag & Drop files here or Click to Upload
                    </Typography>
                    <input
                        id="fileInput"
                        type="file"
                        multiple
                        style={{ display: "none" }}
                        onChange={(e) => void handleFileUpload(e.target.files)}
                    />
                </Paper>

                {uploading && (
                    <Box mt={2} display="flex" alignItems="center">
                        <CircularProgress size={30} sx={{ mr: 1, color: "#1976d2" }} />
                        <Typography>Uploading...</Typography>
                    </Box>
                )}

                <Box sx={{ width: "80%", mt: 2 }}>
                    <Typography variant="h5" fontWeight="bold" color="primary">
                        Uploaded Files
                    </Typography>
                    <Paper
                        sx={{
                            mt: 2,
                            p: 2,
                            bgcolor: "#3b5998",
                            borderRadius: "8px",
                            maxHeight: "220px",
                            overflowY: "auto",
                        }}
                    >
                        <List>
                            {files.map((file, index) => (
                                <ListItem
                                    key={index}
                                    sx={{
                                        backgroundColor: "white",
                                        borderRadius: "4px",
                                        mb: 1,
                                        padding: 1.5,
                                        "&:hover": {
                                            backgroundColor: "#f5f5f5",
                                        },
                                    }}
                                    onClick={() => file.download()}
                                >
                                    <UploadIcon color="primary" sx={{ mr: 1 }} />
                                    <ListItemText
                                        primary={file.name}
                                        primaryTypographyProps={{ color: "text.primary" }}
                                    />
                                    <IconButton
                                        onClick={(e) => {
                                        e.stopPropagation();
                                        confirmDeleteFile(file.name);
                                        }}
                                        sx={{ color: "error.main" }}
                                    >
                                        <DeleteIcon />
                                    </IconButton>
                                </ListItem>
                            ))}
                        </List>
                    </Paper>
                </Box>

                {/* Delete Confirmation Dialog */}
                <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
                    <DialogTitle>Confirm Deletion</DialogTitle>
                    <DialogContent>Are you sure you want to delete this file?</DialogContent>
                    <DialogActions>
                        <Button onClick={() => setDeleteDialogOpen(false)} color="primary">
                            No
                        </Button>
                        <Button onClick={() => void handleDeleteConfirmed()} color="error">
                            Yes
                        </Button>
                    </DialogActions>
                </Dialog>
            </Box>

            {/* Floating Chatbot Button */}
            <Button
                onClick={() => setChatbotOpen(true)}
                variant="contained"
                sx={{
                    position: "fixed",
                    bottom: "20px",
                    right: "20px",
                    zIndex: 1000,
                }}
            >
                Chatbot
            </Button>

            {chatbotOpen && <ChatbotChat onClose={() => setChatbotOpen(false)} />}
        </>
    );
}
