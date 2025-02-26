import { useState, useEffect, useCallback, useRef } from "react";
import { Navigate } from "react-router-dom";
import {
    Container,
    Typography,
    TextField,
    Button,
    IconButton,
    Card,
    CardContent,
    CardActionArea,
    Grid2,
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";

import { useAuth } from "../auth/useAuth";
import { fetchAllModules, Module } from "../api/modules";
import { getUserModules, updateUserModules } from "../api/users";
import DialogCreateModule from "../components/DialogCreateModule";
import { DialogRef } from "../utils/useDialog";

import "./moduleSettings.css";


export default function ModuleSettingsPage() {
    const { token, user } = useAuth();
    const userID = user?.uid;

    const moduleDialogRef = useRef<DialogRef>(null);

    const [searchQuery, setSearchQuery] = useState("");
    const [selectedModules, setSelectedModules] = useState<number[]>([]);
    const [modulesList, setModulesList] = useState<Module[]>([]);
    
    const [loading, setLoading] = useState(true);

    const fetchData = useCallback(async () => {
        const fetchModules = async () => {
            const modules = await fetchAllModules(token);
            setModulesList(modules);
        };
    
        const fetchUserModules = async () => {
            if (!userID) return;
            const userModules = await getUserModules(token, userID);
            const ids = userModules.map((mod) => mod.id);
            setSelectedModules(ids);
        };

        try {
            await Promise.all([
                fetchModules(), 
                fetchUserModules()
            ]);
        } catch (error) {
            console.error("Error fetching data:", error);
        } finally {
            setLoading(false);
        }
    }, [token, userID]);

    useEffect(() => {
        void fetchData();
    }, [fetchData]);

    const handleToggleModule = (moduleId: number) => {
        setSelectedModules((prevSelected) =>
            prevSelected.includes(moduleId)
                ? prevSelected.filter((id) => id !== moduleId)
                : [...prevSelected, moduleId]
        );
    };

    const handleSave = async () => {
        try {
            if (!userID) return;
            await updateUserModules(token, userID, { selectedModules });
        } catch (error) {
            console.error("Error:", error);
            alert("Failed to save modules.");
        }
    };

    if (!token || !user) {
        return <Navigate to="/login" />
    }

    const filteredModules = modulesList.filter((module) =>
        module.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <Container className="module-container">
            <IconButton 
                onClick={() => moduleDialogRef.current?.openDialog()} 
                className="plus-button"
            >
                <AddIcon />
            </IconButton>

            <Typography variant="h4" className="module-title">
                Select Your Modules
            </Typography>

            <div className="search-container">
                <TextField
                    label="Search Modules"
                    variant="outlined"
                    fullWidth
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="module-search-bar"
                />
            </div>

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
                                    <Typography variant="h6" className="module-text" style={{ fontSize: "1rem" }}>
                                        {`${module.code} ${module.name}`}
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

            <DialogCreateModule ref={moduleDialogRef} onUpdate={() => { void fetchData(); }} />
        </Container>
    );
};
