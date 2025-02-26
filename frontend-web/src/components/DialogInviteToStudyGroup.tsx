import { useRef } from "react";
import {
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from "@mui/material";

import { useDialog, DialogRef } from "../utils/useDialog";


interface DialogInviteToStudyGroupProps {
  ref: React.RefObject<DialogRef | null>;
  onClose?: () => void;
  onUpdate?: () => void;
}


export default function DialogInviteToStudyGroup({ ref, onClose, onUpdate }: DialogInviteToStudyGroupProps) {
  //const { token } = useAuth();
  const { isOpen, handleClose } = useDialog(ref, onClose);

  const emailRef = useRef<HTMLInputElement>(null);


  const handleInvite = () => {
    const email = emailRef.current?.value ?? ""
    if (email.trim() === "") {
      alert("Please enter a valid email.");
      return;
    }
    // TODO implement
    console.log(`Inviting ${email} to study group`);
    onUpdate?.();
    handleClose();
  };


  return (
    <Dialog open={isOpen} onClose={handleClose}>
      <DialogTitle>Invite a User</DialogTitle>
      <DialogContent>
        <TextField
          label="Invite Email"
          fullWidth
          inputRef={emailRef}
          margin="dense"
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="error" variant="contained">
          Cancel
        </Button>
        <Button onClick={handleInvite} variant="contained"
          sx={{
            backgroundColor: "#0056b3",
            "&:hover": {
              backgroundColor: "#004494",
            }
          }}>
          Invite
        </Button>
      </DialogActions>
    </Dialog>
  );
}
