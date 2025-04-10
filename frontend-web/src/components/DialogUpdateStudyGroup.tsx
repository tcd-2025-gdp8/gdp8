import { useState, useEffect } from "react";
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    Button,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
} from "@mui/material";
import { updateStudyGroup, StudyGroup, StudyGroupCreationDetails } from "../api/studyGroups";
import { useAuth } from "../auth/useAuth";

interface DialogUpdateStudyGroupProps {
    open: boolean;
    onClose: () => void;
    groupData: StudyGroup;
    onUpdated: (updated: StudyGroup) => void;
}

export default function DialogUpdateStudyGroup({
    open,
    onClose,
    groupData,
    onUpdated,
}: DialogUpdateStudyGroupProps) {
    const { token } = useAuth();

    const [name, setName] = useState(groupData.name);
    const [description, setDescription] = useState(groupData.description);
    const [type, setType] = useState<"public" | "closed" | "invite-only">(groupData.type);
    const [maxMembers, setMaxMembers] = useState<number>(groupData.maxMembers);

    useEffect(() => {
        setName(groupData.name);
        setDescription(groupData.description);
        setType(groupData.type);
        setMaxMembers(groupData.maxMembers);
    }, [groupData]);

    const handleSave = async () => {
        try {
            const updates: StudyGroupCreationDetails = {
                name,
                description,
                type,
                maxMembers,
                moduleId: groupData.moduleId
            };
            const updatedGroup = await updateStudyGroup(token, groupData.id, updates);
            onUpdated(updatedGroup);
        } catch (error) {
            console.error("Error updating study group:", error);
        }
    };

    return (
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
            <DialogTitle>Update Study Group</DialogTitle>
            <DialogContent>
                <TextField
                    label="Group Name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    fullWidth
                    margin="normal"
                />
                <TextField
                    label="Description"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    fullWidth
                    margin="normal"
                    multiline
                    rows={3}
                />
                <FormControl fullWidth margin="normal">
                    <InputLabel>Type</InputLabel>
                    <Select
                        value={type}
                        label="Type"
                        onChange={(e) => setType(e.target.value as "public" | "closed" | "invite-only")}
                    >
                        <MenuItem value="public">Public</MenuItem>
                        <MenuItem value="closed">Closed</MenuItem>
                        <MenuItem value="invite-only">Invite Only</MenuItem>
                    </Select>
                </FormControl>
                <TextField
                    label="Max Members"
                    type="number"
                    value={maxMembers}
                    onChange={(e) => setMaxMembers(Number(e.target.value))}
                    fullWidth
                    margin="normal"
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button variant="contained" onClick={() => void handleSave()} color="primary">
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    );
}
