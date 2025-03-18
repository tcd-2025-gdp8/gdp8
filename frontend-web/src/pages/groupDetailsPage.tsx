import { useState, useEffect } from "react";
import { Box, Typography, Button, CircularProgress, Alert } from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";

import Sidebar from "../components/Sidebar";
import CustomAppBar from "../components/CustomAppBar";
import DialogEditStudyGroup, { EditStudyGroupData } from "../components/DialogEditStudyGroup";
import { fetchStudyGroupById, updateStudyGroup, deleteStudyGroup } from "../api/studyGroups";
import { Module } from "../api/modules";
import { useAuth } from "../auth/useAuth";
import { getUserModules } from "../api/users";

export default function GroupDetailsPage() {
    const navigate = useNavigate();
    const { groupId } = useParams();
    const { token, user } = useAuth();
    const currentUserId = user?.uid;

    const [studyGroupData, setStudyGroupData] = useState<EditStudyGroupData | null>(null);
    const [modulesList, setModulesList] = useState<Module[]>([]);
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

        const fetchModules = async () => {
            if (!currentUserId || !token) return;

            try {
                const data = await getUserModules(token, currentUserId);
                setModulesList(data || []);
            } catch (err) {
                console.error("Error fetching modules:", err);
            }
        };

        void fetchStudyGroup();
        void fetchModules();
    }, [groupId, token, currentUserId]);

    const handleDelete = async () => {
        if (!groupId || !token) return;
        const confirmDelete = window.confirm("Are you sure you want to delete this study group?");
        if (!confirmDelete) return;

        try {
            await deleteStudyGroup(token, Number(groupId));
            navigate("/study-groups");
        } catch (err) {
            console.error("Error deleting study group:", err);
            setError("Failed to delete study group.");
        }
    };

    const handleDeleteClick = () => {
        handleDelete().catch((error) => console.error("Error handling delete:", error));
    };

    const handleUpdateStudyGroup = async (updatedData: EditStudyGroupData) => {
        if (!token || !groupId) return;
        
        try {
            const updatedGroup = await updateStudyGroup(token, Number(groupId), updatedData);
            setStudyGroupData(updatedGroup);
            setOpenEditDialog(false);
        } catch (err) {
            console.error("Error updating study group:", err);
            setError("Failed to update study group.");
        }
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
                                modules={modulesList}
                                onSave={handleUpdateStudyGroup}
                            />
                        </>
                    )}
                </Box>
            </Box>
        </Box>
    );
}
