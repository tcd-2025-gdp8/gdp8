import React, { useState, useEffect, useCallback } from "react";
import {
    Container,
    Typography,
    TextField,
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    IconButton,
    Card,
    CardContent,
    CardActionArea,
    Grid2,
    AppBar,
    Toolbar,
    Box,
    Badge,
    Drawer,
} from "@mui/material";
import { useAuth } from "../auth/useAuth";
import { useNavigate } from "react-router-dom";
import AddIcon from "@mui/icons-material/Add";
import "./ModuleSettings.css";
import NotificationsIcon from "@mui/icons-material/Notifications";
import CloseIcon from "@mui/icons-material/Close";

interface Module {
    id: string;
    name: string;
}

interface Notification {
    id: number;
    message: string;
}

const ModuleSettings: React.FC = () => {
    const navigate = useNavigate();
    const { token, user } = useAuth();
    const userID = user?.uid;
    const [searchQuery, setSearchQuery] = useState("");
    const [selectedModules, setSelectedModules] = useState<string[]>([]);
    const [modulesList, setModulesList] = useState<Module[]>([]);
    const [loading, setLoading] = useState(true);

    const [openDialog, setOpenDialog] = useState(false);
    const [moduleID, setModuleID] = useState("");
    const [moduleName, setModuleName] = useState("");

    const [notifications, setNotifications] = useState<Notification[]>([
            { id: 1, message: 'Your request to join "The Prefects" has been accepted.' },
            { id: 2, message: 'New study group "CS Wizards" has been created for CSU44051: Human Factors.' },
            { id: 3, message: 'New study group "The Elites" has been created for CSU44052: Computer Graphics.' },
        ]);
    const [openNotifications, setOpenNotifications] = useState(false);
    
        // Handler to remove a notification
        const handleDeleteNotification = (notificationId: number) => {
            setNotifications((prev) =>
                prev.filter((notification) => notification.id !== notificationId)
            );
        };

    const fetchModules = useCallback(async () => {
        if (!token) return;
        try {
            const response = await fetch("http://localhost:8080/api/modules", {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
            });
            if (!response.ok) {
                throw new Error("Failed to fetch modules");
            }
            const data = (await response.json()) as Module[];
            setModulesList(data);
        } catch (error) {
            console.error("Error fetching modules:", error);
        } finally {
            setLoading(false);
        }
    }, [token]);

    const fetchUserModules = useCallback(async () => {
        if (!token || !userID) return;
        try {
            const response = await fetch(`http://localhost:8080/api/user/${userID}/modules`, {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
            });
            if (!response.ok) {
                throw new Error("Failed to fetch user modules");
            }
            const data = (await response.json()) as Module[];
            const ids = data.map((mod) => mod.id);
            setSelectedModules(ids);
        } catch (error) {
            console.error("Error fetching user modules:", error);
        }
    }, [token, userID]);

    useEffect(() => {
        void fetchModules();
        void fetchUserModules();
    }, [token, fetchModules, fetchUserModules]);


    const handleCreateModule = async () => {
        if (!token) {
            alert("You must be logged in to create a module.");
            return;
        }
        if (!moduleID.trim() || !moduleName.trim()) {
            alert("Both Module ID and Module Name are required.");
            return;
        }
        try {
            const response = await fetch("http://localhost:8080/api/modules", {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ id: moduleID, name: moduleName }),
            });
            if (response.ok) {
                alert("Module created successfully!");
                void fetchModules();
                setModuleID("");
                setModuleName("");
                setOpenDialog(false);
            } else {
                alert("Error creating module.");
            }
        } catch (error) {
            console.error("Error:", error);
            alert("Failed to create module.");
        }
    };

    // Create a synchronous wrapper to call the async handleCreateModule.
    const handleCreateModuleClick = (): void => {
        void handleCreateModule();
    };

    const filteredModules = modulesList.filter((module) =>
        module.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const handleToggleModule = (moduleId: string) => {
        setSelectedModules((prevSelected) =>
            prevSelected.includes(moduleId)
                ? prevSelected.filter((id) => id !== moduleId)
                : [...prevSelected, moduleId]
        );
    };

    const handleSave = async () => {
        if (!token) {
            alert("You must be logged in to save preferences.");
            return;
        }

        try {
            const response = await fetch(`http://localhost:8080/api/user/${userID}/modules`, {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ selectedModules }),
            });
            if (response.ok) {
                alert("Modules saved successfully!");
            } else {
                alert("Error saving modules.");
            }
        } catch (error) {
            console.error("Error:", error);
            alert("Failed to save modules.");
        }
    };

    return (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh" width="100vw">
            {/* Back to Landing Button */}
       <Button
        variant="contained"
        sx={{
          backgroundColor: "#0056b3",
          "&:hover": {
              backgroundColor: "#004494",
          }}}
        onClick={() => { void navigate("/landing"); }}
        style={{ position: "absolute", top: "10px", left: "10px" }}
      >
        Back to Landing
      </Button>
      <AppBar
            position="fixed"
            sx={{
                backgroundColor: "#ffffff",
                color: "#000000",
                zIndex: 1201,
                width: "calc(100% - 240px)",
                marginLeft: "240px",
            }}
        >
            <Toolbar sx={{ flexDirection: "row", alignItems: "center" }}>
                <Box sx={{ flexGrow: 1, display: "flex", alignItems: "center" }}>
                    <Typography variant="h6">
                        Blackboard + StudyWise
                    </Typography>
                </Box>
                {user && (
                    <Typography variant="subtitle1" sx={{ mr: 2 }}>
                        {user.email}
                    </Typography>
                )}
                <IconButton color="inherit" onClick={() => setOpenNotifications(true)}>
                    <Badge badgeContent={notifications.length} color="error">
                        <NotificationsIcon />
                    </Badge>
                </IconButton>
            </Toolbar>
        </AppBar>

        {/* Notifications Drawer */}
        <Drawer
            anchor="right"
            open={openNotifications}
            onClose={() => setOpenNotifications(false)}
            sx={{
                "& .MuiDrawer-paper": {
                    zIndex: 1300,
                },
            }}
        >
            <Box sx={{ width: 300, padding: "1rem" }}>
                <Box
                    sx={{
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                        marginBottom: "1rem",
                    }}
                >
                    <Typography variant="h6">Notifications</Typography>
                    <IconButton onClick={() => setOpenNotifications(false)}>
                        <CloseIcon />
                    </IconButton>
                </Box>
                {notifications.length === 0 ? (
                    <Typography
                        variant="body2"
                        color="textSecondary"
                        sx={{ textAlign: "center" }}
                    >
                        No notifications
                    </Typography>
                ) : (
                    notifications.map((notification) => (
                        <Box
                            key={notification.id}
                            sx={{
                                padding: "0.5rem",
                                borderBottom: "1px solid #ccc",
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                            }}
                        >
                            <Typography>{notification.message}</Typography>
                            <IconButton
                                size="small"
                                onClick={() => handleDeleteNotification(notification.id)}
                            >
                                <CloseIcon fontSize="small" />
                            </IconButton>
                        </Box>
                    ))
                )}
            </Box>
        </Drawer>

        <Container className="module-container">
            <Button
            variant="contained"
            sx={{
                backgroundColor: "#0056b3",
                "&:hover": {
                backgroundColor: "#004494",
                },
                position: "absolute",
                top: "10px",
                left: "10px",
                color: "white", // Ensure text remains visible
            }}
            onClick={() => { void navigate("/landing"); }}
            >
            Back to Landing
            </Button>
            <IconButton 
            onClick={() => setOpenDialog(true)} 
            className="plus-button">
                <AddIcon />
            </IconButton>
                <>
                <Typography variant="h4" className="module-title">
                Select Your Modules
            </Typography>

            <TextField
                label="Search Modules"
                variant="outlined"
                fullWidth
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="module-search-bar"
            />

            {loading && <Typography variant="body1">Loading modules...</Typography>}

                {!loading && filteredModules.length === 0 && (
                    <Typography variant="body1" className="module-no-modules">
                        No modules found.
                    </Typography>
                )}

            <Grid2 container spacing={2} justifyContent="center">
                {filteredModules.map((module) => (
                    <Grid2 key={module.id}>
                        <Card
                            onClick={() => handleToggleModule(module.id)}
                            className={`module-card ${selectedModules.includes(module.id) ? "selected" : ""}`}
                        >
                                <CardActionArea>
                                    <CardContent>
                                        <Typography variant="h6" className="module-text">
                                            {module.name}
                                        </Typography>
                                    </CardContent>
                                </CardActionArea>
                            </Card>
                            </Grid2>
                    ))}
                </Grid2>

                <Button
                variant="contained"
                sx={{
                    backgroundColor: "#0056b3",
                    "&:hover": {
                    backgroundColor: "#004494",
                    },
                    marginTop: "20px",
                    width: "100%",
                    color: "white",
                }}
                onClick={() => void handleSave()}
                >
                Save Preferences
                </Button>
            {/* Dialog for Creating a New Module */}
            <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
                <DialogTitle>Create a New Module</DialogTitle>
                <DialogContent>
                    <TextField
                        label="Module ID"
                        fullWidth
                        value={moduleID}
                        onChange={(e) => setModuleID(e.target.value)}
                        margin="dense"
                    />
                    <TextField
                        label="Module Name"
                        fullWidth
                        value={moduleName}
                        onChange={(e) => setModuleName(e.target.value)}
                        margin="dense"
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenDialog(false)} color="error" variant="contained">
                        Cancel
                    </Button>
                    <Button onClick={handleCreateModuleClick} color="primary" variant="contained"
                            sx={{
                                backgroundColor: "#0056b3",
                                "&:hover": {
                                    backgroundColor: "#004494",
                                }}}>
                        Create
                    </Button>
                    </DialogActions>
                </Dialog>
                </>
            </Container>
            </Box>
    );
};

const styles = {
    container: {
        width: "100%",
        maxWidth: "1200px",
        backgroundColor: "#ffffff",
        borderRadius: "10px",
        boxShadow: "0px 4px 10px rgba(0, 0, 0, 0.1)",
        padding: "30px",
        position: "relative" as const,
        textAlign: "center" as const,
    },
    title: {
        marginBottom: "20px",
        color: "#333",
        fontWeight: "bold",
    },
    searchBar: {
        marginBottom: "20px",
    },
    card: {
        padding: "10px",
        borderRadius: "10px",
        cursor: "pointer",
        transition: "background-color 0.3s, transform 0.2s",
        boxShadow: "0px 2px 5px rgba(0, 0, 0, 0.2)",
        textAlign: "center" as const,
        "&:hover": {
            transform: "scale(1.03)",
        },
    },
    moduleText: {
        fontWeight: "bold",
        textAlign: "center" as const,
        color: "#333",
    },
    saveButton: {
        marginTop: "20px",
        width: "100%",
    },
    noModules: {
        color: "#555",
        fontStyle: "italic",
        marginBottom: "20px",
    },
    plusButton: {
        position: "absolute" as const,
        top: "10px",
        right: "10px",
        backgroundColor: "#1976D2",
        color: "white",
    },
};

export default ModuleSettings;