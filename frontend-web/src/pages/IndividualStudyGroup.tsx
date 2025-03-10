// pages/IndividualStudyGroup.tsx

import { useState } from "react";
import { Box, Typography, Button } from "@mui/material";
import { useNavigate } from "react-router-dom";

import Sidebar from "../components/Sidebar";
import CustomAppBar from "../components/CustomAppBar";
//import DialogEditStudyGroup from "../components/DialogEditStudyGroup";
import DialogEditStudyGroup, { EditStudyGroupData } from "../components/DialogEditStudyGroup";


export default function IndividualStudyGroup() {
    const navigate = useNavigate();
    const [openEditDialog, setOpenEditDialog] = useState(false);

    // Example existing data for demonstration (no API calls here)
    const [studyGroupData, setStudyGroupData] = useState<EditStudyGroupData>({
        name: "Sample Study Group",
        description: "This is a sample group for demonstration.",
        type: "public", // no 'as const'
        moduleId: 101,
        maxMembers: 5,
    });

    const handleDelete = () => {
        const confirmDelete = window.confirm("Are you sure you want to delete this study group?");
        if (confirmDelete) {
            // Simulate deletion on the frontend
            // You can also navigate back to the list page:
            void navigate("/study-groups");
        }
    };

    const handleSaveUpdates = (updatedData: typeof studyGroupData) => {
        // In a real app, you'd call an API or update state in a global store.
        // Here, we just update local state and show an alert.
        alert("Study group updated (frontend simulation).");
        console.log("Updated data:", updatedData);

        setStudyGroupData(updatedData);
        setOpenEditDialog(false);
    };

    return (
        <Box sx={{ display: "flex" }}>
            {/* Sidebar remains visible */}
            <Sidebar />

            {/* Main Content Area */}
            <Box sx={{ flexGrow: 1, ml: "240px" }}>
                {/* Optional: If you want the custom app bar at the top */}
                <CustomAppBar />

                <Box sx={{ p: 2, mt: 10 }}>
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

                    {/* Buttons */}
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={() => setOpenEditDialog(true)}
                        sx={{ mr: 2 }}
                    >
                        Edit Study Group
                    </Button>
                    <Button variant="contained" color="error" onClick={handleDelete}>
                        Delete Study Group
                    </Button>
                </Box>
            </Box>

            {/* The new "Edit Study Group" dialog */}
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
    );
}
