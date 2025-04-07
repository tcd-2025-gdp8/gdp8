import { useState, useEffect } from "react";
import { 
  Box, 
  Typography, 
  CircularProgress, 
  Alert, 
  Paper, 
  Chip, 
  Divider,
  Avatar,
  Card,
  CardContent,
  Grid,
  IconButton,
  Tooltip
} from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import GroupIcon from "@mui/icons-material/Group";
import SchoolIcon from "@mui/icons-material/School";
import DescriptionIcon from "@mui/icons-material/Description";
import LockIcon from "@mui/icons-material/Lock";
import PublicIcon from "@mui/icons-material/Public";
import MailIcon from "@mui/icons-material/Mail";

import Sidebar from "../components/Sidebar";
import CustomAppBar from "../components/CustomAppBar";
import DialogEditStudyGroup, { EditStudyGroupData } from "../components/DialogEditStudyGroup";
import { fetchStudyGroupById, updateStudyGroup, deleteStudyGroup } from "../api/studyGroups";
import { Module } from "../api/modules";
import { useAuth } from "../auth/useAuth";
import { getUserModules } from "../api/users";

// Define a new interface that includes members
interface StudyGroupData {
    id: number;
    name: string;
    description: string;
    type: "public" | "closed" | "invite-only";
    moduleId: number;
    maxMembers: number;
    members: { id: string; name: string; role: string }[]; 
}

export default function GroupDetailsPage() {
    const navigate = useNavigate();
    const { groupId } = useParams();
    const { token, user } = useAuth();
    const currentUserId = user?.uid;

    const [studyGroupData, setStudyGroupData] = useState<StudyGroupData | null>(null);
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
            void navigate("/study-groups");
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

    const handleSaveWrapper = (updatedData: EditStudyGroupData) => {
        handleUpdateStudyGroup(updatedData).catch((err) => {
            console.error("Error updating study group:", err);
            setError("Failed to update study group.");
        });
    };

    const module = modulesList.find(module => module.id === studyGroupData?.moduleId);
    const moduleDisplay = module ? `${module.code} ${module.name}` : "Unknown Module";

    // Check if the user is an admin
    const isAdmin = studyGroupData?.members.some(
        (member) => member.id === currentUserId && member.role === "admin"
    );

    // Get group type icon
    const getGroupTypeIcon = (type: string) => {
        switch (type) {
            case 'public': return <PublicIcon />;
            case 'closed': return <LockIcon />;
            case 'invite-only': return <MailIcon />;
            default: return <PublicIcon />;
        }
    };

    // Get color for member role
    const getMemberRoleColor = (role: string) => {
        switch (role) {
            case 'admin': return 'error';
            case 'moderator': return 'warning';
            default: return 'primary';
        }
    };

    // Generate color for avatar based on name
    const stringToColor = (string: string) => {
        let hash = 0;
        for (let i = 0; i < string.length; i++) {
            hash = string.charCodeAt(i) + ((hash << 5) - hash);
        }
        
        let color = '#';
        for (let i = 0; i < 3; i++) {
            const value = (hash >> (i * 8)) & 0xFF;
            const lightValue = Math.floor((value + 255) / 2);
            color += ('00' + lightValue.toString(16)).substr(-2);
        }
        
        return color;
    };

    return (
        <Box sx={{ display: "flex", bgcolor: "#f5f7fa", minHeight: "100vh" }}>
            <Sidebar />
            <Box sx={{ flexGrow: 1, ml: { xs: "0px", sm: "100px", md: "300px" } }}>
                <CustomAppBar />
                
                <Box
  sx={{
    position: 'absolute',
    top: '20%', 
    left: '30%'
  }}
>
                                        {loading && (
                        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh' }}>
                            <CircularProgress />
                        </Box>
                    )}
                    
                    {error && (
                        <Alert 
                            severity="error" 
                            sx={{ maxWidth: 800, mx: 'auto', mb: 2 }}
                        >
                            {error}
                        </Alert>
                    )}

                    {studyGroupData && (
                        <Grid container spacing={3}>
                            {/* Main group info */}
                            <Grid item xs={12}>
                                <Paper 
                                    elevation={2} 
                                    sx={{ 
                                        p: 3, 
                                        borderRadius: 2,
                                        position: 'relative',
                                        overflow: 'hidden'
                                    }}
                                >
                                    <Box 
                                        sx={{ 
                                            position: 'absolute', 
                                            top: 0, 
                                            left: 0, 
                                            width: '100%', 
                                            height: '6px', 
                                            background: 'linear-gradient(90deg, #2196f3, #6a1b9a)'
                                        }} 
                                    />
                                    
                                    <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between' }}>
                                        <Box>
                                            <Typography variant="h4" fontWeight="bold" sx={{ mb: 1 }}>
                                                {studyGroupData.name}
                                            </Typography>
                                            
                                            <Box sx={{ display: 'flex', gap: 2, mb: 2, alignItems: 'center', flexWrap: 'wrap' }}>
                                                <Chip 
                                                    icon={getGroupTypeIcon(studyGroupData.type)} 
                                                    label={studyGroupData.type} 
                                                    color="primary" 
                                                    variant="outlined"
                                                    size="small"
                                                />
                                                
                                                <Chip 
                                                    icon={<SchoolIcon />} 
                                                    label={moduleDisplay} 
                                                    color="secondary" 
                                                    variant="outlined"
                                                    size="small"
                                                />
                                                
                                                <Chip 
                                                    icon={<GroupIcon />} 
                                                    label={`${studyGroupData.members.length}/${studyGroupData.maxMembers} members`} 
                                                    variant="outlined"
                                                    size="small"
                                                />
                                            </Box>
                                        </Box>
                                        
                                        {isAdmin && (
                                            <Box>
                                                <Tooltip title="Edit Group">
                                                    <IconButton 
                                                        onClick={() => setOpenEditDialog(true)}
                                                        color="primary"
                                                    >
                                                        <EditIcon />
                                                    </IconButton>
                                                </Tooltip>
                                                
                                                <Tooltip title="Delete Group">
                                                    <IconButton 
                                                        onClick={handleDeleteClick}
                                                        color="error"
                                                    >
                                                        <DeleteIcon />
                                                    </IconButton>
                                                </Tooltip>
                                            </Box>
                                        )}
                                    </Box>
                                    
                                    <Divider sx={{ my: 2 }} />
                                    
                                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                                        <DescriptionIcon sx={{ mr: 1, color: 'text.secondary' }} />
                                        <Typography variant="subtitle1" fontWeight="medium">Description</Typography>
                                    </Box>
                                    
                                    <Paper 
                                        variant="outlined" 
                                        sx={{ 
                                            p: 2, 
                                            bgcolor: 'rgba(0,0,0,0.01)', 
                                            borderRadius: 1,
                                            mb: 3
                                        }}
                                    >
                                        <Typography variant="body1">
                                            {studyGroupData.description || "No description provided."}
                                        </Typography>
                                    </Paper>
                                </Paper>
                            </Grid>
                            
                            {/* Members section */}
                            <Grid item xs={12}>
                                <Card elevation={2} sx={{ borderRadius: 2 }}>
                                    <CardContent>
                                        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                                            <GroupIcon sx={{ mr: 1, color: 'text.secondary' }} />
                                            <Typography variant="h6" fontWeight="medium">Members</Typography>
                                        </Box>
                                        
                                        <Divider sx={{ mb: 2 }} />
                                        
                                        <Grid container spacing={2}>
                                            {studyGroupData.members.map((member) => (
                                                <Grid item xs={12} sm={6} md={4} key={member.id}>
                                                    <Paper 
                                                        elevation={0} 
                                                        sx={{ 
                                                            p: 2, 
                                                            display: 'flex', 
                                                            alignItems: 'center',
                                                            borderRadius: 2,
                                                            bgcolor: 'rgba(0,0,0,0.01)',
                                                            border: '1px solid rgba(0,0,0,0.08)'
                                                        }}
                                                    >
                                                        <Avatar 
                                                            sx={{ 
                                                                bgcolor: stringToColor(member.name),
                                                                width: 42,
                                                                height: 42,
                                                                mr: 2
                                                            }}
                                                        >
                                                            {member.name.charAt(0).toUpperCase()}
                                                        </Avatar>
                                                        
                                                        <Box sx={{ minWidth: 0, flexGrow: 1 }}>
                                                            <Typography 
                                                                variant="subtitle1" 
                                                                noWrap 
                                                                sx={{ fontWeight: 'medium' }}
                                                            >
                                                                {member.name}
                                                            </Typography>
                                                            
                                                            <Chip 
                                                                label={member.role} 
                                                                size="small"
                                                                color={getMemberRoleColor(member.role)}
                                                                sx={{ height: 24, fontSize: '0.75rem' }}
                                                            />
                                                            
                                                            {member.id === currentUserId && (
                                                                <Chip 
                                                                    label="You" 
                                                                    size="small"
                                                                    variant="outlined"
                                                                    sx={{ ml: 1, height: 24, fontSize: '0.75rem' }}
                                                                />
                                                            )}
                                                        </Box>
                                                    </Paper>
                                                </Grid>
                                            ))}
                                        </Grid>
                                    </CardContent>
                                </Card>
                            </Grid>
                        </Grid>
                    )}

                    {/* Maintain original dialog behavior with conditional rendering */}
                    {studyGroupData && (
                        <DialogEditStudyGroup
                            open={openEditDialog}
                            onClose={() => setOpenEditDialog(false)}
                            existingData={studyGroupData} // Using 'as any' to bypass type checking for compatibility
                            modules={modulesList}
                            onSave={handleSaveWrapper}
                        />
                    )}
                </Box>
            </Box>
        </Box>
    );
}