import { useState, useEffect } from "react";
import { Box, Typography, Button, CircularProgress, Alert } from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";

import Sidebar from "../components/Sidebar";
import CustomAppBar from "../components/CustomAppBar";
import DialogEditStudyGroup, { EditStudyGroupData } from "../components/DialogEditStudyGroup";
import { fetchStudyGroupById } from "../api/studyGroups";
import { useAuth } from "../auth/useAuth";

export default function GroupDetailsPage() {
    const navigate = useNavigate();
    const { groupId } = useParams();
    const { token } = useAuth();

    const [studyGroupData, setStudyGroupData] = useState<EditStudyGroupData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [openEditDialog, setOpenEditDialog] = useState(false);

    useEffect(() => {
        const fetchStudyGroup = async () => {
            if (!groupId || !token) return;
    
            try {
                const data = await fetchStudyGroupById(token, Number(groupId));
                console.log("Fetched Study Group Data:", data);
                setStudyGroupData(data);
            } catch (err) {
                console.error("Error fetching study group:", err);
                setError("Failed to load study group details.");
            } finally {
                setLoading(false);
            }
        };
    
        void fetchStudyGroup();
    }, [groupId, token]);
    

    const handleDelete = async () => {
        if (!groupId || !token) return;
        const confirmDelete = window.confirm("Are you sure you want to delete this study group?");
        if (!confirmDelete) return;

        try {
            const response = await fetch(`/api/study-groups/${groupId}`, {
                method: "DELETE",
                headers: { Authorization: `Bearer ${token}` },
            });

            if (!response.ok) {
                throw new Error("Failed to delete study group.");
            }

            void navigate("/study-groups");
        } catch (err) {
            console.error("Error deleting study group:", err);
            setError("Failed to delete study group.");
        }
    };

    const handleDeleteClick = () => {
        handleDelete().catch((error) => console.error("Error handling delete:", error));
    };

    const handleSaveUpdates = (updatedData: EditStudyGroupData) => {
        setStudyGroupData(updatedData);
    };

    return (
        <Box sx={{ display: "flex" }}>
            <Sidebar />
            <Box sx={{ flexGrow: 1, ml: "300px" }}>
                <CustomAppBar />
                <Box sx={{ p: 2, mt: 10, position: "relative", paddingLeft: "80px" }}>

                    {loading && <CircularProgress />}
                    {error && <Alert severity="error">{error}</Alert>}

                    {studyGroupData && (
                        <>
                            <Typography variant="h4" sx={{ mb: 2 }}>{studyGroupData.name}</Typography>
                            <Typography variant="body1" sx={{ mb: 1 }}>Description: {studyGroupData.description}</Typography>
                            <Typography variant="body1" sx={{ mb: 1 }}>Type: {studyGroupData.type}</Typography>
                            <Typography variant="body1" sx={{ mb: 1 }}>Module ID: {studyGroupData.moduleId}</Typography>
                            <Typography variant="body1" sx={{ mb: 3 }}>Max Members: {studyGroupData.maxMembers}</Typography>

                            <Button variant="contained" color="primary" onClick={() => setOpenEditDialog(true)} sx={{ mr: 2 }}>
                                Edit Study Group
                            </Button>
                            <Button variant="contained" color="error" onClick={handleDeleteClick}>
                                Delete Study Group
                            </Button>

                            <DialogEditStudyGroup
                                open={openEditDialog}
                                onClose={() => setOpenEditDialog(false)}
                                existingData={studyGroupData}
                                modules={[
                                    { id: 101, code: "MATH101", name: "Algebra" },
                                    { id: 102, code: "BIO102", name: "Biology" },
                                ]}
                                onSave={handleSaveUpdates}
                            />
                        </>
                    )}
                </Box>
            </Box>
        </Box>
    );
}
