// frontend-web/pages/GroupDetailsPage.tsx

import { useState } from "react";
import { Box, Typography, Button } from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";

import Sidebar from "../components/Sidebar";
import CustomAppBar from "../components/CustomAppBar";
import DialogEditStudyGroup, { EditStudyGroupData } from "../components/DialogEditStudyGroup";

export default function GroupDetailsPage() {
    const navigate = useNavigate();
    const { groupId } = useParams();
    console.log("GroupID:", groupId);
    const [openEditDialog, setOpenEditDialog] = useState(false);

    // Example data for demonstration (no real API calls yet)
    const [studyGroupData, setStudyGroupData] = useState<EditStudyGroupData>({
        name: "Sample Study Group",
        description: "This is a sample group for demonstration.",
        type: "public",
        moduleId: 101,
        maxMembers: 5,
    });

    const handleDelete = () => {
        const confirmDelete = window.confirm("Are you sure you want to delete this study group?");
        if (confirmDelete) {
            // Simulate deletion
            void navigate("/study-groups");
        }
    };

    const handleSaveUpdates = (updatedData: typeof studyGroupData) => {
        alert("Study group updated (frontend simulation).");
        console.log("Updated data:", updatedData);

        setStudyGroupData(updatedData);
        setOpenEditDialog(false);
    };

    return (
        <Box sx={{ display: "flex" }}>
            <Sidebar />
            <Box sx={{ flexGrow: 1, ml: "300px" }}>
                <CustomAppBar />
                <Box sx={{ p: 2, mt: 10, position: 'relative', paddingLeft: '80px' }}>
                    <Typography variant="h4" sx={{ mb: 2 }}>
                        {studyGroupData.name}
                    </Typography>

                    <Typography variant="body1" sx={{ mb: 1 }}>
                        Description: {studyGroupData.description}
                    </Typography>
                    <Typography variant="body1" sx={{ mb: 1 }}>
                        Type: {studyGroupData.type}
                    </Typography>
                    <Typography variant="body1" sx={{ mb: 1 }}>
                        Module ID: {studyGroupData.moduleId}
                    </Typography>
                    <Typography variant="body1" sx={{ mb: 3 }}>
                        Max Members: {studyGroupData.maxMembers}
                    </Typography>

                    <Button
                        variant="contained"
                        color="primary"
                        onClick={() => setOpenEditDialog(true)}
                        sx={{ mr: 2 }}
                    >
                        Edit Study Group
                    </Button>
                    <Button
                        variant="contained"
                        color="error"
                        onClick={handleDelete}
                    >
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
                </Box>
            </Box>
        </Box>
    );
}
