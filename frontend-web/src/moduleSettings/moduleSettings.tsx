import React, { useState, useEffect, useCallback } from "react";
import { 
    Container, 
    Typography, 
    TextField, 
    Grid2, 
    Card, 
    CardContent, 
    CardActionArea, 
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    IconButton, 
} from "@mui/material";
import { useAuth } from "../auth/useAuth";
import AddIcon from "@mui/icons-material/Add";

interface Module {
    id: string;
    name: string;
}

const ModuleSettings: React.FC = () => {
    const { token } = useAuth();
    const [searchQuery, setSearchQuery] = useState("");
    const [selectedModules, setSelectedModules] = useState<string[]>([]);
    const [modulesList, setModulesList] = useState<Module[]>([]);
    const [loading, setLoading] = useState(true);

    const [openDialog, setOpenDialog] = useState(false);
    const [moduleID, setModuleID] = useState("");
    const [moduleName, setModuleName] = useState("");

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

    useEffect(() => {
        void fetchModules();
    }, [token, fetchModules]);


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
            const response = await fetch("http://localhost:8080/api/save-modules", {
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

        <Container style={styles.container}>
            <IconButton 
            onClick={() => setOpenDialog(true)} 
            style={styles.plusButton}>
                <AddIcon />
            </IconButton>
                <Typography variant="h4" style={styles.title}>
                Select Your Modules
            </Typography>

                <TextField
                label="Search Modules"
                variant="outlined"
                fullWidth
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                style={styles.searchBar}
            />

                {loading && <Typography variant="body1">Loading modules...</Typography>}

                {!loading && filteredModules.length === 0 && (
                    <Typography variant="body1" style={styles.noModules}>
                        No modules found.
                    </Typography>
                )}

                <Grid2 container spacing={2} justifyContent="center">
                    {filteredModules.map((module) => (
                        <Grid2 key={module.id}>
                            <Card 
                            onClick={() => handleToggleModule(module.id)}
                            style={{
                                ...styles.card,
                                backgroundColor: selectedModules.includes(module.id) ? "#d4edda" : "#ffffff",
                            }}
                        >
                                <CardActionArea>
                                    <CardContent>
                                        <Typography variant="h6" style={styles.moduleText}>
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
                color="primary"
                onClick={() => void handleSave()}
                style={styles.saveButton}
            >
                Save Preferences
            </Button>
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
                        <Button onClick={() => setOpenDialog(false)} color="error">
                        Cancel
                    </Button>
                        <Button onClick={() => void handleCreateModule()} color="primary">
                        Create
                    </Button>
                    </DialogActions>
                </Dialog>
            </Container>
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
