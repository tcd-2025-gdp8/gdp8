import { useRef, useState } from "react";
import {
  Button,
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
  Slider
} from "@mui/material";

import { Module } from "../api/modules";
import { createStudyGroup, StudyGroupCreationDetails } from "../api/studyGroups";
import { useAuth } from "../auth/useAuth";
import DialogInviteToStudyGroup from "./DialogInviteToStudyGroup";
import { useDialog, DialogRef } from "../utils/useDialog";


interface DialogCreateStudyGroupProps {
  modules: Module[];
  ref: React.RefObject<DialogRef | null>;
  onClose?: () => void;
  onUpdate?: () => void;
}


export default function DialogCreateStudyGroup({ modules, ref, onClose, onUpdate }: DialogCreateStudyGroupProps) {
  const { token } = useAuth();
  const { isOpen, handleClose } = useDialog(ref, onClose);

  const inviteMembersDialogRef = useRef<DialogRef>(null);

  const groupNameRef = useRef<HTMLInputElement>(null);
  const groupDescriptionRef = useRef<HTMLInputElement>(null);
  const groupTypeRef = useRef<HTMLSelectElement>(null);
  const moduleIdRef = useRef<HTMLSelectElement>(null);

  const [maxMembers, setMaxMembers] = useState(5);


  const handleCreateGroup = async () => {
    const groupDetails: StudyGroupCreationDetails = {
      name: groupNameRef.current?.value ?? "",
      description: groupDescriptionRef.current?.value ?? "",
      type: groupTypeRef.current?.value as "public" | "closed" | "invite-only" ?? "public",
      moduleId: Number(moduleIdRef.current?.value) || -1,
      maxMembers: maxMembers
    };

    if (groupDetails.name.trim() === "" || groupDetails.description.trim() === "") {
      alert("Please enter a valid group name and description.");
      return;
    }

    try {
      await createStudyGroup(token, groupDetails);
      onUpdate?.();
      handleClose();
    } catch (error) {
      console.error("Error creating study group:", error);
      alert("Error creating study group.");
    }
  };


  return (
    <Dialog open={isOpen} onClose={handleClose}>
      <DialogTitle>Create a Study Group</DialogTitle>
      <DialogContent>
        <TextField
          label="Group Name"
          fullWidth
          inputRef={groupNameRef}
          margin="dense"
        />
        <TextField
          label="Group Description"
          fullWidth
          inputRef={groupDescriptionRef}
          margin="dense"
        />
        <FormControl fullWidth style={{ marginTop: "10px" }}>
          <InputLabel>Group Type</InputLabel>
          <Select
            defaultValue="public"
            inputRef={groupTypeRef}
          >
            <MenuItem value="public">Public</MenuItem>
            <MenuItem value="closed">Closed</MenuItem>
            <MenuItem value="invite-only">Invite Only</MenuItem>
          </Select>
        </FormControl>
        <FormControl fullWidth style={{ marginTop: "10px" }}>
          <InputLabel>Select Module</InputLabel>
          <Select
            defaultValue=""
            inputRef={moduleIdRef}
          >
            {modules.map((module) => (
              <MenuItem key={module.id} value={module.id}>
                {`${module.code} ${module.name}`}
              </MenuItem>
            ))}
          </Select>
          <Typography gutterBottom style={{ marginTop: "10px" }}>
            Max Members: {maxMembers}
          </Typography>
          <Slider
            value={maxMembers}
            onChange={(_, value) => setMaxMembers(value as number)}
            min={2}
            max={10}
            step={1}
            marks
            valueLabelDisplay="auto"
          />
        </FormControl>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="error" variant="contained">
          Cancel
        </Button>
        <Button onClick={() => { inviteMembersDialogRef.current?.openDialog(); }} variant="contained"
          sx={{
            backgroundColor: "#0056b3",
            "&:hover": {
              backgroundColor: "#004494",
            }
          }}>
          Invite Members
        </Button>
        <Button onClick={() => { void handleCreateGroup(); }} variant="contained"
          sx={{
            backgroundColor: "#0056b3",
            "&:hover": {
              backgroundColor: "#004494",
            }
          }}>
          Create
        </Button>
      </DialogActions>
      <DialogInviteToStudyGroup ref={inviteMembersDialogRef} onUpdate={onUpdate} />
    </Dialog>
  );
}
