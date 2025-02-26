import { useRef } from "react";
import { useNavigate } from "react-router-dom";
import {
  Card,
  CardContent,
  Typography,
  Button,
  Tooltip,
  Divider,
} from "@mui/material";

import { Module } from "../api/modules";
import { useAuth } from "../auth/useAuth";
import { requestToJoinStudyGroup, StudyGroup } from "../api/studyGroups";
import DialogStudyGroupMembers from "./DialogStudyGroupMembers";
import { DialogRef } from "../utils/useDialog";

interface StudyGroupCardProps {
  group: StudyGroup;
  module: Module;
  currentUserId: string;
  onUpdate?: () => void;
}

export default function StudyGroupCard({ group, module, currentUserId, onUpdate }: StudyGroupCardProps) {
  const navigate = useNavigate();
  const { token } = useAuth();

  const membersDialogRef = useRef<DialogRef>(null);

  const isMember = group.members.some((member) => member.id === currentUserId);
  const isAdmin = group.members.some((member) => member.id === currentUserId && member.role === "admin");
  const isFull = group.members.length >= group.maxMembers;

  const handleJoinGroup = (id: number) => {
    void (async () => {
      try {
        await requestToJoinStudyGroup(token, id);
        onUpdate?.();
      } catch (error) {
        console.error("Error joining study group:", error);
        alert(`Error joining study group.`);
      }
    })();
  };  

  return (
    <Card>
      <CardContent>
        <Typography variant="h6">{group.name}</Typography>
        <Divider sx={{ my: 1 }} />
        <Typography color="textPrimary" style={{ marginBottom: "5px" }}>
          {`${module.code} ${module.name}`}
        </Typography>
        <Typography variant="body2" color="textSecondary" style={{ marginBottom: "5px" }}>
          {group.description}
        </Typography>
        <Typography variant="body2" color="textSecondary" style={{ marginBottom: "10px" }}>
          {`Group Capacity: ${group.maxMembers} members`}
        </Typography>
        <Tooltip
          arrow
          title={
            group.type === "closed"
              ? "You cannot view the members of this group as it is a closed group."
              : group.members.length === 0
                ? "No members yet"
                : ""
          }
        >

          {/* Wrap button in a <span> so Tooltip works on a disabled button */}
          <span>
            <Button
              variant="contained"
              sx={{
                backgroundColor: "#0056b3",
                "&:hover": {
                  backgroundColor: "#004494",
                }
              }}
              onClick={() => membersDialogRef.current?.openDialog()}
              disabled={group.type !== "public"}
            >
              View Members
            </Button>
          </span>
        </Tooltip>
        <Divider sx={{ my: 1 }} />
        
        <Button
          variant="contained"
          sx={{
            backgroundColor: "#0056b3",
            "&:hover": {
              backgroundColor: "#004494",
            }
          }}
          fullWidth
          onClick={() => handleJoinGroup(group.id)}
          style={{ marginTop: "10px" }}
          disabled={
            isMember || isFull
          }
        >
          {isMember ? "Joined" : isFull ? "Full" : "Request to Join"}
        </Button>

        {!isFull && (
          <Button
            variant="contained"
            sx={{
              backgroundColor: "#0056b3",
              "&:hover": {
                backgroundColor: "#004494",
              }
            }}
            fullWidth
            onClick={() => void navigate(`/chat/${group.id}`)}
            disabled={!isMember}
            style={{ marginTop: "10px" }}
          >
            {isMember ? "Open Chat" : "Join to Chat!"}
          </Button>
        )}

        <DialogStudyGroupMembers
          studyGroup={group}
          isCurrentUserAdmin={isAdmin}
          ref={membersDialogRef}
          onUpdate={onUpdate}
        />

      </CardContent>
    </Card>
  );
}
