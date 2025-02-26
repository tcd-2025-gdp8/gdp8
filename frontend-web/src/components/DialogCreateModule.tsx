import { useRef } from "react";
import {
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from "@mui/material";

import { createModule } from "../api/modules";
import { useAuth } from "../auth/useAuth";
import { useDialog, DialogRef } from "../utils/useDialog";


interface DialogCreateModuleProps {
  ref: React.RefObject<DialogRef | null>;
  onClose?: () => void;
  onUpdate?: () => void;
}


export default function DialogCreateModule({ ref, onClose, onUpdate }: DialogCreateModuleProps) {
  const { token } = useAuth();
  const { isOpen, handleClose } = useDialog(ref, onClose);

  const moduleCodeRef = useRef<HTMLInputElement>(null);
  const moduleNameRef = useRef<HTMLInputElement>(null);


  const handleCreateModule = async () => {
    const moduleCode = moduleCodeRef.current?.value ?? "";
    const moduleName = moduleNameRef.current?.value ?? "";
    if (!moduleCode.trim() || !moduleName.trim()) {
      alert("Both Module Code and Module Name are required.");
      return;
    }
    try {
      await createModule(token, {
          code: moduleCode,
          name: moduleName,
      });
      onUpdate?.();
      handleClose();
    } catch (error) {
      console.error("Error:", error);
      alert("Failed to create module.");
    }
  };


  return (
    <Dialog open={isOpen} onClose={handleClose}>
      <DialogTitle>Create a New Module</DialogTitle>
      <DialogContent>
        <TextField
          label="Module Code"
          fullWidth
          inputRef={moduleCodeRef}
          margin="dense"
        />
        <TextField
          label="Module Name"
          fullWidth
          inputRef={moduleNameRef}
          margin="dense"
        />
      </DialogContent>
      <DialogActions>
        <Button 
          onClick={handleClose}
          color="error" 
          variant="contained"
        >
            Cancel
        </Button>
        <Button 
          onClick={() => void handleCreateModule()} 
          color="primary" 
          variant="contained"
              sx={{
                  backgroundColor: "#0056b3",
                  "&:hover": {
                      backgroundColor: "#004494",
                  }
              }}
        >
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
}
