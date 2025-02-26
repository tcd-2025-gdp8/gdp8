import {
  Button,
  Dialog,
  DialogContent,
  DialogActions,
  IconButton,
  List,
  ListItem,
  ListItemText,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";

import { removeMemberFromStudyGroup, StudyGroup } from "../api/studyGroups";
import { useAuth } from "../auth/useAuth";
import { DialogRef, useDialog } from "../utils/useDialog";


interface DialogStudyGroupMembersProps {
  studyGroup: StudyGroup;
  isCurrentUserAdmin: boolean;
  ref: React.RefObject<DialogRef | null>;
  onClose?: () => void;
  onUpdate?: () => void;
}


export default function DialogStudyGroupMembers({ studyGroup, isCurrentUserAdmin, ref, onClose, onUpdate }: DialogStudyGroupMembersProps) {
  const { token } = useAuth();
  const { isOpen, handleClose } = useDialog(ref, onClose);

  const handleRemoveMember = async (memberId: string) => {
    try {
      await removeMemberFromStudyGroup(token, studyGroup.id, memberId);
      onUpdate?.();
    } catch (err) {
      console.error("Error removing member:", err);
      alert(`Error removing member.`);
    }
  };

  return (
    <Dialog open={isOpen} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogContent>
          {studyGroup.members.length > 0 && (
            <List>
              {studyGroup.members.map((member) => (
                <ListItem key={member.id}>
                  <ListItemText primary={member.name} secondary={member.role} />
                  {isCurrentUserAdmin && member.role !== "admin" && (
                    <IconButton edge="end" onClick={() => void handleRemoveMember(member.id)}>
                      <DeleteIcon color="error" />
                    </IconButton>
                  )}
                </ListItem>
              ))}
            </List>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}
          color="primary"
          variant="contained"
          sx={{
            backgroundColor: "#0056b3",
            "&:hover": {
              backgroundColor: "#004494",
            }
          }}>
            Close
          </Button>
        </DialogActions>
      </Dialog>
  );
}
