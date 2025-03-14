import { useRef } from "react";
import { useNavigate, Link } from "react-router-dom";
import {
  Card,
  CardContent,
  Typography,
  Button,
  Tooltip,
  Divider,
  Box,
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
    <Card sx={{
      height: '100%',
      display: 'flex',
      flexDirection: 'column',
      transition: 'transform 0.2s, box-shadow 0.2s',
      '&:hover': {
        transform: 'translateY(-4px)',
        boxShadow: 4,
      }
    }}>
      <CardContent sx={{
        flexGrow: 1,
        display: 'flex',
        flexDirection: 'column',
        gap: 2
      }}>
        {/* Study group name wrapped in a Link */}
        <Typography variant="h6" sx={{ fontWeight: 600 }}>
          {isMember ? (
            <Link
              to={`/study-groups/${group.id}`}
              style={{ textDecoration: "none", color: "inherit" }}
            >
              {group.name}
            </Link>
          ) : (
            <span>{group.name}</span>
          )}
        </Typography>

        <Divider />

        <div>
          <Typography variant="subtitle1" sx={{ fontWeight: 500, mb: 1 }}>
            {`${module.code} ${module.name}`}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {group.description}
          </Typography>
          <Typography
            variant="body2"
            sx={{
              mt: 1,
              color: isFull ? 'error.main' : 'success.main',
              fontWeight: 500
            }}
          >
            {`${group.members.length}/${group.maxMembers} members`}
          </Typography>
        </div>

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
          <span>
            <Button
              variant="outlined"
              onClick={() => membersDialogRef.current?.openDialog()}
              disabled={group.type !== "public"}
              fullWidth
            >
              View Members
            </Button>
          </span>
        </Tooltip>

        <Box sx={{ mt: 'auto' }}>
          <Button
            variant="contained"
            fullWidth
            onClick={() => handleJoinGroup(group.id)}
            disabled={isMember || isFull}
            sx={{ mb: 1 }}
          >
            {isMember ? "Joined" : isFull ? "Full" : "Request to Join"}
          </Button>

          {(isMember || !isFull) && (
            <Button
              variant={isMember ? "contained" : "outlined"}
              color={isMember ? "primary" : "inherit"}
              fullWidth
              onClick={() => void navigate(`/chat/${group.id}`)}
              disabled={!isMember}
            >
              {isMember ? "Open Chat" : "Join to Chat!"}
            </Button>
          )}
        </Box>

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
