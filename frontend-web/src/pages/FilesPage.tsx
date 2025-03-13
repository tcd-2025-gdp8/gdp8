// import { useState, useEffect, DragEvent } from 'react';
// import { getStorage, ref, uploadBytesResumable, getDownloadURL, listAll } from 'firebase/storage';
// import { Loader } from 'lucide-react';
// import { app } from '../auth/firebase'; // ✅ Import initialized Firebase app

// // ✅ Use existing `app` instance for Firebase Storage
// const storage = getStorage(app);

// export default function SharedFiles() {
//   const [files, setFiles] = useState<string[]>([]);
//   const [uploading, setUploading] = useState(false);
//   const [dragging, setDragging] = useState(false);

//   useEffect(() => {
//     void fetchFiles();
//   }, []);

//   const fetchFiles = async () => {
//     try {
//       const filesRef = ref(storage, 'shared-files/');
//       const fileList = await listAll(filesRef);
//       const urls = await Promise.all(fileList.items.map(item => getDownloadURL(item)));
//       setFiles(urls);
//     } catch (error) {
//       console.error('Error fetching files:', error);
//     }
//   };

//   const handleFileUpload = async (selectedFiles: FileList | null) => {
//     if (!selectedFiles || selectedFiles.length === 0) return;

//     setUploading(true);
//     try {
//       for (const file of Array.from(selectedFiles)) {
//         const fileRef = ref(storage, `shared-files/${file.name}`);
//         const uploadTask = uploadBytesResumable(fileRef, file);
//         await new Promise<void>((resolve, reject) => {
//           uploadTask.on('state_changed', null, (err) => reject(err), () => resolve());
//         });
//       }
//       await fetchFiles();
//     } catch (error) {
//       console.error('Error uploading files:', error);
//     } finally {
//       setUploading(false);
//     }
//   };

//   // Drag-and-Drop Handlers
//   const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
//     e.preventDefault();
//     setDragging(true);
//   };

//   const handleDragLeave = () => {
//     setDragging(false);
//   };

//   const handleDrop = (e: DragEvent<HTMLDivElement>) => {
//     e.preventDefault();
//     setDragging(false);
//     void handleFileUpload(e.dataTransfer.files);
//   };

//   return (
//     <div className="p-6">
//       <h1 className="text-2xl font-semibold mb-4">Shared Files</h1>

//       {/* Drag-and-Drop Area */}
//       <div
//         className={`border-2 border-dashed p-6 text-center cursor-pointer rounded-lg ${
//           dragging ? 'border-blue-500 bg-blue-100' : 'border-gray-300'
//         }`}
//         onDragOver={handleDragOver}
//         onDragLeave={handleDragLeave}
//         onDrop={handleDrop}
//         onClick={() => document.getElementById('fileInput')?.click()}
//       >
//         <input
//           id="fileInput"
//           type="file"
//           multiple
//           className="hidden"
//           onChange={(e) => {
//             void handleFileUpload(e.target.files);
//           }}
//         />
//         <p>Drag & drop files here, or click to select files</p>
//       </div>

//       {/* Uploading Indicator */}
//       {uploading && (
//         <p className="mt-2 flex items-center">
//           <Loader className="animate-spin mr-2" /> Uploading...
//         </p>
//       )}

//       {/* Uploaded Files List */}
//       <h2 className="text-xl font-medium mt-6">Uploaded Files</h2>
//       <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
//         {files.map((file, index) => (
//           <div key={index} className="p-4 flex justify-between items-center border rounded-lg shadow">
//             <a
//               href={file}
//               target="_blank"
//               rel="noopener noreferrer"
//               className="text-blue-600 hover:underline"
//             >
//               {decodeURIComponent(file.split('/').pop() ?? '')}
//             </a>
//           </div>
//         ))}
//       </div>
//     </div>
//   );
// }

// import { useState, useEffect, DragEvent } from "react";
// import { getStorage, ref, uploadBytesResumable, getDownloadURL, listAll } from "firebase/storage";
// import { app } from "../auth/firebase";
// import { Box, Typography, CircularProgress, Paper, List, ListItem, ListItemText } from "@mui/material";
// import { UploadFile as UploadIcon, CloudUpload as CloudUploadIcon } from "@mui/icons-material";

// // ✅ Firebase Storage Instance
// const storage = getStorage(app);

// export default function FilesPage() {
//     const [files, setFiles] = useState<string[]>([]);
//     const [uploading, setUploading] = useState(false);
//     const [dragging, setDragging] = useState(false);

//     useEffect(() => {
//         void fetchFiles();
//     }, []);

//     const fetchFiles = async () => {
//         try {
//             const filesRef = ref(storage, "shared-files/");
//             const fileList = await listAll(filesRef);
//             const urls = await Promise.all(fileList.items.map((item) => getDownloadURL(item)));
//             setFiles(urls);
//         } catch (error) {
//             console.error("Error fetching files:", error);
//         }
//     };

//     const handleFileUpload = async (selectedFiles: FileList | null) => {
//         if (!selectedFiles || selectedFiles.length === 0) return;

//         setUploading(true);
//         try {
//             for (const file of Array.from(selectedFiles)) {
//                 const fileRef = ref(storage, `shared-files/${file.name}`);
//                 const uploadTask = uploadBytesResumable(fileRef, file);
//                 await new Promise<void>((resolve, reject) => {
//                     uploadTask.on("state_changed", null, reject, resolve);
//                 });
//             }
//             await fetchFiles();
//         } catch (error) {
//             console.error("Error uploading files:", error);
//         } finally {
//             setUploading(false);
//         }
//     };

//     const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
//         e.preventDefault();
//         setDragging(true);
//     };

//     const handleDragLeave = () => {
//         setDragging(false);
//     };

//     const handleDrop = (e: DragEvent<HTMLDivElement>) => {
//         e.preventDefault();
//         setDragging(false);
//         void handleFileUpload(e.dataTransfer.files);
//     };

//     return (
//         <Box
//             sx={{
//                 display: "flex",
//                 flexDirection: "column",
//                 alignItems: "center",
//                 justifyContent: "center",
//                 width: "100vw",
//                 height: "100vh",
//                 bgcolor: "#f0f8ff",
//                 padding: 4,
//             }}
//         >
//             {/* Title */}
//             <Typography variant="h4" fontWeight="bold" color="primary" gutterBottom>
//                 File Upload
//             </Typography>

//             {/* Upload Box (Centered) */}
//             <Paper
//                 sx={{
//                     width: "80%",
//                     height: "250px",
//                     border: "3px dashed",
//                     borderColor: dragging ? "#1976d2" : "#90caf9",
//                     borderRadius: "12px",
//                     display: "flex",
//                     alignItems: "center",
//                     justifyContent: "center",
//                     flexDirection: "column",
//                     backgroundColor: "#ffffff",
//                     transition: "0.3s",
//                     cursor: "pointer",
//                     "&:hover": {
//                         borderColor: "#1565c0",
//                         backgroundColor: "#e3f2fd",
//                     },
//                 }}
//                 onDragOver={handleDragOver}
//                 onDragLeave={handleDragLeave}
//                 onDrop={handleDrop}
//                 onClick={() => document.getElementById("fileInput")?.click()}
//             >
//                 <CloudUploadIcon sx={{ fontSize: 80, color: "#1976d2", mb: 1 }} />
//                 <Typography variant="h6" color="textSecondary">
//                     Drag & Drop files here or Click to Upload
//                 </Typography>
//                 <input
//                     id="fileInput"
//                     type="file"
//                     multiple
//                     style={{ display: "none" }}
//                     onChange={(e) => {
//                         void handleFileUpload(e.target.files);
//                     }}
//                 />
//             </Paper>

//             {/* Uploading Indicator */}
//             {uploading && (
//                 <Box mt={2} display="flex" alignItems="center">
//                     <CircularProgress size={30} sx={{ mr: 1, color: "#1976d2" }} />
//                     <Typography>Uploading...</Typography>
//                 </Box>
//             )}

//             {/* Uploaded Files Section (Full Width) */}
//             <Box sx={{ width: "80%", mt: 4 }}>
//                 <Typography variant="h5" fontWeight="bold" color="primary">
//                     Uploaded Files
//                 </Typography>
//                 <Paper sx={{ mt: 2, p: 2, bgcolor: "#ffffff" }}>
//                     <List>
//                         {files.length === 0 ? (
//                             <Typography color="textSecondary" textAlign="center">
//                                 No files uploaded yet.
//                             </Typography>
//                         ) : (
//                             files.map((file, index) => (
//                                 <ListItem
//                                     key={index}
//                                     sx={{
//                                         borderBottom: "1px solid #e0e0e0",
//                                         "&:last-child": { borderBottom: "none" },
//                                     }}
//                                 >
//                                     <UploadIcon color="primary" sx={{ mr: 1 }} />
//                                     <ListItemText>
//                                         <a href={file} target="_blank" rel="noopener noreferrer" style={{ color: "#1976d2", textDecoration: "none" }}>
//                                             {decodeURIComponent(file.split("/").pop() ?? "")}
//                                         </a>
//                                     </ListItemText>
//                                 </ListItem>
//                             ))
//                         )}
//                     </List>
//                 </Paper>
//             </Box>
//         </Box>
//     );
// }

import { useState, useEffect, DragEvent } from "react";
import { getStorage, ref, uploadBytesResumable, getDownloadURL, listAll, deleteObject } from "firebase/storage";
import { app } from "../auth/firebase";
import { Box, Typography, CircularProgress, Paper, List, ListItem, ListItemText, IconButton } from "@mui/material";
import { UploadFile as UploadIcon, CloudUpload as CloudUploadIcon, Delete as DeleteIcon } from "@mui/icons-material";

// ✅ Firebase Storage Instance
const storage = getStorage(app);

export default function FilesPage() {
    const [files, setFiles] = useState<string[]>([]);
    const [uploading, setUploading] = useState(false);
    const [dragging, setDragging] = useState(false);

    useEffect(() => {
        void fetchFiles();
    }, []);

    const fetchFiles = async () => {
        try {
            const filesRef = ref(storage, "shared-files/");
            const fileList = await listAll(filesRef);
            const urls = await Promise.all(fileList.items.map((item) => getDownloadURL(item)));
            setFiles(urls);
        } catch (error) {
            console.error("Error fetching files:", error);
        }
    };

    const handleFileUpload = async (selectedFiles: FileList | null) => {
        if (!selectedFiles || selectedFiles.length === 0) return;

        setUploading(true);
        try {
            for (const file of Array.from(selectedFiles)) {
                const fileRef = ref(storage, `shared-files/${file.name}`);
                const uploadTask = uploadBytesResumable(fileRef, file);
                await new Promise<void>((resolve, reject) => {
                    uploadTask.on("state_changed", null, reject, resolve);
                });
            }
            await fetchFiles();
        } catch (error) {
            console.error("Error uploading files:", error);
        } finally {
            setUploading(false);
        }
    };

    const handleFileDelete = async (fileUrl: string) => {
        try {
            const fileName = decodeURIComponent(fileUrl.split("/").pop() ?? "");
            const fileRef = ref(storage, `shared-files/${fileName}`);

            // ✅ Delete file from Firebase Storage
            await deleteObject(fileRef);

            // ✅ Refresh file list after deletion
            setFiles((prevFiles) => prevFiles.filter((f) => f !== fileUrl));
            console.log(`Deleted file: ${fileName}`);
        } catch (error) {
            console.error("Error deleting file:", error);
        }
    };

    const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        setDragging(true);
    };

    const handleDragLeave = () => {
        setDragging(false);
    };

    const handleDrop = (e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        setDragging(false);
        void handleFileUpload(e.dataTransfer.files);
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
            {/* Title */}
            <Typography variant="h4" fontWeight="bold" color="primary" gutterBottom>
                File Upload Center
            </Typography>

            {/* Upload Box (Centered) */}
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
                    onChange={(e) => {
                        void handleFileUpload(e.target.files);
                    }}
                />
            </Paper>

            {/* Uploading Indicator */}
            {uploading && (
                <Box mt={2} display="flex" alignItems="center">
                    <CircularProgress size={30} sx={{ mr: 1, color: "#1976d2" }} />
                    <Typography>Uploading...</Typography>
                </Box>
            )}

            {/* Uploaded Files Section (Full Width) */}
            <Box sx={{ width: "80%", mt: 4 }}>
                <Typography variant="h5" fontWeight="bold" color="primary">
                    Uploaded Files
                </Typography>
                <Paper sx={{ mt: 2, p: 2, bgcolor: "#ffffff" }}>
                    <List>
                        {files.length === 0 ? (
                            <Typography color="textSecondary" textAlign="center">
                                No files uploaded yet.
                            </Typography>
                        ) : (
                            files.map((file, index) => (
                                <ListItem
                                    key={index}
                                    sx={{
                                        borderBottom: "1px solid #e0e0e0",
                                        "&:last-child": { borderBottom: "none" },
                                    }}
                                >
                                    <UploadIcon color="primary" sx={{ mr: 1 }} />
                                    <ListItemText>
                                        <a href={file} target="_blank" rel="noopener noreferrer" style={{ color: "#1976d2", textDecoration: "none" }}>
                                            {decodeURIComponent(file.split("/").pop() ?? "")}
                                        </a>
                                    </ListItemText>
                                    {/* ✅ Delete Button */}
                                    <IconButton onClick={() => void handleFileDelete(file)} sx={{ color: "red" }}>
                                        <DeleteIcon />
                                    </IconButton>
                                </ListItem>
                            ))
                        )}
                    </List>
                </Paper>
            </Box>
        </Box>
    );
}