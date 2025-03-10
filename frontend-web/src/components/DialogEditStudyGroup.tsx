// components/DialogEditStudyGroup.tsx

import { useState, useEffect } from "react";
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    MenuItem,
    Select,
    FormControl,
    InputLabel,
    Typography,
    Slider,
    Button,
} from "@mui/material";

export interface EditStudyGroupData {
    name: string;
    description: string;
    type: "public" | "closed" | "invite-only";
    moduleId: number;
    maxMembers: number;
}

interface Module {
    id: number;
    code: string;
    name: string;
}

interface DialogEditStudyGroupProps {
    open: boolean;
    onClose: () => void;
    /** The existing data to populate the fields (if you have it). */
    existingData: EditStudyGroupData;
    /** The list of modules to show in the dropdown. */
    modules: Module[];
    /** Callback when the user clicks "Save" (just logs or updates local state in this example). */
    onSave: (updatedData: EditStudyGroupData) => void;
}

export default function DialogEditStudyGroup({
    open,
    onClose,
    existingData,
    modules,
    onSave,
}: DialogEditStudyGroupProps) {
    // Local state to hold the edited values
    const [name, setName] = useState(existingData.name);
    const [description, setDescription] = useState(existingData.description);
    const [type, setType] = useState<"public" | "closed" | "invite-only">(existingData.type);
    const [moduleId, setModuleId] = useState<number>(existingData.moduleId);
    const [maxMembers, setMaxMembers] = useState<number>(existingData.maxMembers);

    // Whenever `existingData` changes, reset the local state
    useEffect(() => {
        setName(existingData.name);
        setDescription(existingData.description);
        setType(existingData.type);
        setModuleId(existingData.moduleId);
        setMaxMembers(existingData.maxMembers);
    }, [existingData]);

    const handleSave = () => {
        // Build an updated data object
        const updatedData: EditStudyGroupData = {
            name,
            description,
            type,
            moduleId,
            maxMembers,
        };
        // Call the onSave callback to handle it in the parent
        onSave(updatedData);
    };

    const handleInviteMembers = () => {
        // For now, just log or do nothing. You can open a separate dialog if needed.
        alert("Invite members clicked (frontend only).");
    };

    return (
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
            <DialogTitle>Edit Study Group</DialogTitle>
            <DialogContent>
                <TextField
                    label="Group Name"
                    fullWidth
                    margin="dense"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
                <TextField
                    label="Group Description"
                    fullWidth
                    margin="dense"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                />
                <FormControl fullWidth sx={{ mt: 2 }}>
                    <InputLabel>Group Type</InputLabel>
                    <Select
                        value={type}
                        label="Group Type"
                        onChange={(e) => setType(e.target.value as "public" | "closed" | "invite-only")}
                    >
                        <MenuItem value="public">Public</MenuItem>
                        <MenuItem value="closed">Closed</MenuItem>
                        <MenuItem value="invite-only">Invite Only</MenuItem>
                    </Select>
                </FormControl>

                <FormControl fullWidth sx={{ mt: 2 }}>
                    <InputLabel>Select Module</InputLabel>
                    <Select
                        value={moduleId}
                        label="Select Module"
                        onChange={(e) => setModuleId(Number(e.target.value))}
                    >
                        {modules.map((m) => (
                            <MenuItem key={m.id} value={m.id}>
                                {`${m.code} ${m.name}`}
                            </MenuItem>
                        ))}
                    </Select>
                </FormControl>

                <Typography gutterBottom sx={{ mt: 2 }}>
                    Max Members: {maxMembers}
                </Typography>
                <Slider
                    value={maxMembers}
                    onChange={(_, val) => setMaxMembers(val as number)}
                    min={2}
                    max={10}
                    step={1}
                    marks
                    valueLabelDisplay="auto"
                />
            </DialogContent>
            <DialogActions>
                <Button color="error" variant="contained" onClick={onClose}>
                    Cancel
                </Button>
                <Button
                    variant="contained"
                    onClick={handleInviteMembers}
                    sx={{
                        backgroundColor: "#0056b3",
                        "&:hover": {
                            backgroundColor: "#004494",
                        },
                    }}
                >
                    Invite Members
                </Button>
                <Button
                    variant="contained"
                    onClick={handleSave}
                    sx={{
                        backgroundColor: "#0056b3",
                        "&:hover": {
                            backgroundColor: "#004494",
                        },
                    }}
                >
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    );
}
