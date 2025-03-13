import { useState, DragEvent } from "react";
import { Box, Typography, CircularProgress, Paper, List, ListItem, ListItemText, IconButton } from "@mui/material";
import { UploadFile as UploadIcon, CloudUpload as CloudUploadIcon, Delete as DeleteIcon } from "@mui/icons-material";

interface LocalFile {
    name: string;
    url: string;
}

export default function FilesPage() {
    const [files, setFiles] = useState<LocalFile[]>([]);
    const [uploading, setUploading] = useState(false);
    const [dragging, setDragging] = useState(false);

    const handleFileUpload = (selectedFiles: FileList | null): void => {
        if (!selectedFiles || selectedFiles.length === 0) return;

        setUploading(true);
        setTimeout(() => {
            const newFiles = Array.from(selectedFiles).map(file => ({
                name: file.name,
                url: URL.createObjectURL(file),
            }));

            setFiles(prevFiles => [...prevFiles, ...newFiles]);
            setUploading(false);
        }, 1000); 
    };

    const handleFileDelete = (fileUrl: string): void => {
        setFiles(prevFiles => prevFiles.filter(file => file.url !== fileUrl));
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
        handleFileUpload(e.dataTransfer.files);
    };

    return (
        <Box
            sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                width: "100vw",
                height: "100vh",
                bgcolor: "#f0f8ff",
                padding: 4,
            }}
        >
            <Typography variant="h4" fontWeight="bold" color="primary" gutterBottom>
                File Upload Center
            </Typography>

            <Paper
                sx={{
                    width: "80%",
                    height: "250px",
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
                    onChange={(e) => handleFileUpload(e.target.files)}
                />
            </Paper>

            {uploading && (
                <Box mt={2} display="flex" alignItems="center">
                    <CircularProgress size={30} sx={{ mr: 1, color: "#1976d2" }} />
                    <Typography>Uploading...</Typography>
                </Box>
            )}

            <Box sx={{ width: "80%", mt: 4 }}>
                <Typography variant="h5" fontWeight="bold" color="primary">
                    Uploaded Files
                </Typography>
                <Paper sx={{ mt: 2, p: 2, bgcolor: "#ffffff" }}>
                    <List>
                        {files.map((file, index) => (
                            <ListItem key={index}>
                                <UploadIcon color="primary" sx={{ mr: 1 }} />
                                <ListItemText>{file.name}</ListItemText>
                                <IconButton onClick={() => handleFileDelete(file.url)} sx={{ color: "red" }}>
                                    <DeleteIcon />
                                </IconButton>
                            </ListItem>
                        ))}
                    </List>
                </Paper>
            </Box>
        </Box>
    );
}