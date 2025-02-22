// src/StudyGroupsPage/StudyGroupsPage.tsx
import React, { useState, useEffect } from "react";
import { useAuth } from "../auth/useAuth";
import {
  Container,
  Card,
  CardContent,
  Typography,
  Button,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  Select,
  SelectChangeEvent,
  FormControl,
  InputLabel,
  Tooltip,
  Drawer,
  IconButton,
  Box
} from "@mui/material";
import NotificationsIcon from '@mui/icons-material/Notifications';
import CloseIcon from '@mui/icons-material/Close';
import { useNavigate } from "react-router-dom";

interface StudyGroupMember {
  userID: string;
  role: "admin" | "member" | "invitee" | "requester";
}

interface StudyGroupDetails {
  name: string;
  description: string;
  type: "public" | "closed" | "invite-only";
  moduleID: string;
}

interface StudyGroup {
  id: number;
  studyGroupDetails: StudyGroupDetails;
  members: StudyGroupMember[];
}

interface HardcodedStudyGroup {
  id: number;
  moduleID: string
  members: StudyGroupMember[];
}

interface Notification {
  id: number;
  message: string;
}

interface APIMember {
  id: string;
  name: string;
  role: "admin" | "member" | "invitee" | "requester";
}

interface APIStudyGroup {
  id: number;
  name: string;
  description: string;
  type: "public" | "closed" | "invite-only";
  members: APIMember[];
}

interface Module {
  id: string;
  name: string;
}

const initialGroups: HardcodedStudyGroup[] = [
  {
    id: 1,
    moduleID: "CSU44052",
    members: [],
  },
  {
    id: 2,
    moduleID: "CSU44052",
    members: [],
  },
  {
    id: 3,
    moduleID: "CSU44099",
    members: [],
  },
];

const StudyGroupsPage: React.FC = () => {
  const { token } = useAuth();
  const navigate = useNavigate();
  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [filteredGroups, setFilteredGroups] = useState<StudyGroup[]>([]);
  const [selectedModule, setSelectedModule] = useState<string>("");
  const [openDialog, setOpenDialog] = useState(false);
  const [groupName, setGroupName] = useState("");
  const [groupDescription, setGroupDescription] = useState("");
  const [groupType, setGroupType] = useState<"public" | "closed" | "invite-only">("public");
  const [selectedGroupModule, setSelectedGroupModule] = useState<string>("");
  const [openInviteDialog, setOpenInviteDialog] = useState(false);
  const [inviteName, setInviteName] = useState("");
  const [inviteEmail, setInviteEmail] = useState("");
  const [notifications, setNotifications] = useState<Notification[]>([
    { id: 1, message: 'Your request to join "The Prefects" has been accepted.' },
    { id: 2, message: 'New study group "CS Wizards" has been created for CSU44051: Human Factors.' },
    { id: 3, message: 'New study group "The Elites" has been created for CSU44052: Computer Graphics.' }
  ]);
  const [openNotifications, setOpenNotifications] = useState(false);
  const [modulesList, setModulesList] = useState<Module[]>([]);

  useEffect(() => {
    const fetchStudyGroups = async (): Promise<void> => {
      if (!token) return;

      try {
        const response = await fetch("http://localhost:8080/api/study-groups", {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });

        if (!response.ok) {
          throw new Error(`Failed to fetch study groups: ${response.statusText}`);
        }

        const data: APIStudyGroup[] = (await response.json()) as APIStudyGroup[];

        const formattedGroups: StudyGroup[] = data.map((group: APIStudyGroup) => {
          const matchedGroup = initialGroups.find((g) => g.id === group.id);

          return {
            id: group.id,
            studyGroupDetails: {
              name: group.name,
              description: group.description,
              type: group.type,
              moduleID: matchedGroup?.moduleID ?? "",
            },
            members: group.members ? group.members.map((member) => ({
              userID: member.id,
              role: member.role,
            })) : [],
          };
        });
        setStudyGroups(formattedGroups);
      } catch (err) {
        console.error("Error fetching study groups:", err);
      }
    };

    void fetchStudyGroups();
  }, [token]);

  useEffect(() => {
    if (selectedModule === "" || selectedModule === "All") {
      setFilteredGroups(studyGroups);
    } else {
      setFilteredGroups(
        studyGroups.filter((group) => group.studyGroupDetails.moduleID === selectedModule)
      );
    }
  }, [selectedModule, studyGroups]);

  const handleJoinGroup = (id: number) => {
    let joinedGroupName = "";

    const updatedGroups = studyGroups.map((group) => {
      if (group.id === id) {
        if (!group.members.some((member) => member.userID === "Alessandro")) {
          joinedGroupName = group.studyGroupDetails.name;
          return {
            ...group,
            members: [
              ...group.members,
              { userID: "Alessandro", role: "member" as const },
            ],
          };
        }
      }
      return group;
    });

    if (joinedGroupName) {
      setStudyGroups(updatedGroups);
      setNotifications((prev) => [
        ...prev,
        { id: Date.now(), message: `You joined Study Group: '${joinedGroupName}'.` },
      ]);
    }
  };


  const handleDeleteNotification = (notificationId: number) => {
    setNotifications((prevNotifications) =>
      prevNotifications.filter((notification) => notification.id !== notificationId)
    );
  };

  const handleOpenChat = (groupId: number) => {
    const group = studyGroups.find((g) => g.id === groupId);
    const isMember = group?.members.some(
      (member) => member.userID === "Alessandro" && ["member", "admin"].includes(member.role)
    );
    if (isMember) {
      void navigate(`/chat/${groupId}`);
    }
  };

  const handleOpenDialog = () => setOpenDialog(true);
  const handleCloseDialog = () => {
    setOpenDialog(false);
    setGroupName("");
    setGroupDescription("");
    setGroupType("public");
    setSelectedGroupModule("");
  };

  const handleCreateGroup = async () => {
    if (groupName.trim() === "" || selectedGroupModule === "") {
      alert("Please enter a valid group name and select a module.");
      return;
    }
  
    if (!token) {
      alert("You are not authorised. Please log in.");
      return;
    }
  
    const newGroupDetails = {
      name: groupName,
      description: groupDescription,
      type: groupType,
    };
  
    try {
      const response = await fetch("http://localhost:8080/api/study-groups", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newGroupDetails),
      });
  
      if (!response.ok) {
        throw new Error(`Failed to create group: ${response.statusText}`);
      }
  
      const createdGroup = (await response.json()) as APIStudyGroup;
  
      const newGroup: StudyGroup = {
        id: createdGroup.id,
        studyGroupDetails: {
          name: createdGroup.name,
          description: createdGroup.description,
          type: createdGroup.type,
          moduleID: selectedGroupModule,
        },
        members: createdGroup.members.map((member) => ({
          userID: member.id,
          role: member.role,
        })),
      };
  
      setStudyGroups([...studyGroups, newGroup]);
  
      setNotifications((prev) => [
        ...prev,
        {
          id: Date.now(),
          message: `New study group "${groupName}" has been created for ${
            modulesList.find((module) => module.id === selectedGroupModule)?.name
          }.`,
        },
      ]);
  
      handleCloseDialog();
    } catch (error) {
      console.error("Error creating study group:", error);
      alert(
        `Error creating study group: ${error instanceof Error ? error.message : String(error)}`
      );
    }
  };

  useEffect(() => {
    const fetchModules = async () => {
      if (!token) return;
  
      try {
        const response = await fetch("http://localhost:8080/api/modules", {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });
  
        if (!response.ok) {
          throw new Error(`Failed to fetch modules: ${response.statusText}`);
        }
  
        const data = (await response.json()) as Module[];
  
        // Transform module names
        const formattedModules = data.map((module) => ({
          id: module.id,
          name: `${module.id}: ${module.name}`,
        }));
  
        // Add "All" option at the beginning
        setModulesList([{ id: "All", name: "All" }, ...formattedModules]);
      } catch (err) {
        console.error("Error fetching modules:", err);
      }
    };
  
    void fetchModules();
  }, [token]);  

  const handleOpenInviteDialog = () => setOpenInviteDialog(true);
  const handleCloseInviteDialog = () => setOpenInviteDialog(false);

  const handleInvite = () => {
    if (inviteName.trim() === "" || inviteEmail.trim() === "") {
      alert("Please enter a valid name and email.");
      return;
    }
    alert(`Invite sent to ${inviteName} at ${inviteEmail}`);
    setInviteName("");
    setInviteEmail("");
    handleCloseInviteDialog();
  };

  return (
    <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh" width="100vw">
      <Container maxWidth="md" style={{ marginTop: "20px", position: "relative", textAlign: "center" }}>
        {/* Back to Landing Button */}
        <Button
          variant="contained"
          color="primary"
          onClick={() => { void navigate("/landing"); }}
          style={{ position: "absolute", top: "10px", left: "10px" }}
        >
          Back to Landing
        </Button>

        <Typography variant="h4" gutterBottom>
          Study Groups
          <IconButton onClick={() => setOpenNotifications(true)}>
            <NotificationsIcon />
          </IconButton>
        </Typography>

        <Drawer anchor="right" open={openNotifications} onClose={() => setOpenNotifications(false)}>
          <div style={{ width: 300, padding: 16 }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <Typography variant="h6">Notifications</Typography>
              <IconButton onClick={() => setOpenNotifications(false)}>
                <CloseIcon />
              </IconButton>
            </div>
            {notifications.length === 0 ? (
              <Typography variant="body2" color="textSecondary" style={{ textAlign: "center", marginTop: "20px" }}>
                No notifications
              </Typography>
            ) : (
              notifications.map((notification) => (
                <Card key={notification.id} style={{ marginBottom: "10px" }}>
                  <CardContent>
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                      <Typography>{notification.message}</Typography>
                      <IconButton onClick={() => handleDeleteNotification(notification.id)} size="small">
                        <CloseIcon />
                      </IconButton>
                    </div>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
        </Drawer>

        <FormControl fullWidth style={{ marginBottom: "20px" }}>
          <InputLabel>Filter by Module</InputLabel>
          <Select
            value={selectedModule}
            onChange={(e: SelectChangeEvent<string>) => setSelectedModule(e.target.value)}
          >
            {modulesList.map((module) => (
              <MenuItem key={module.id} value={module.id}>
                {module.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <Grid container spacing={2}>
          {filteredGroups.map((group) => {
            const isMember = group.members.some((member) => member.userID === "Alessandro");
            const isFull = group.members.length >= 10;


            return ( // ✅ Added return statement
              <Grid item xs={12} sm={6} md={4} key={group.id} style={{ minWidth: "280px" }}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">{group.studyGroupDetails.name}</Typography>
                    <Typography variant="body2" color="textSecondary" style={{ marginBottom: "10px" }}>
                      {group.studyGroupDetails.description}
                    </Typography>
                    <Tooltip
                      title={
                        group.members.length > 0
                          ? group.members.map((member) => member.userID).join(", ")
                          : (group.studyGroupDetails.type === "closed" ? "You cannot view the members of this group as it is a closed group."
                            : "No members yet")
                      }
                      arrow
                    >
                      <Typography color="textSecondary" style={{ cursor: "pointer" }}>
                        Members: {group.studyGroupDetails.type === "closed" ? "confidential" : group.members.length}
                      </Typography>
                    </Tooltip>
                    <Typography color="textSecondary">
                      Module: {modulesList.find((module) => module.id === group.studyGroupDetails.moduleID)?.name}
                    </Typography>
                    <Button
                      variant="contained"
                      color="primary"
                      fullWidth
                      onClick={() => handleJoinGroup(group.id)}
                      style={{ marginTop: "10px" }}
                      disabled={
                        group.members.length >= 10 ||
                        group.members.some((member) => member.userID === "Alessandro")
                      }
                    >
                      {group.members.some((member) => member.userID === "Alessandro")
                        ? "Joined"
                        : group.members.length >= 10
                          ? "Full"
                          : "Request to Join"}
                    </Button>
                    {!isFull && (
                      <Button
                        variant="contained"
                        color="primary"
                        fullWidth
                        onClick={() => handleOpenChat(group.id)}
                        disabled={!isMember}
                        style={{ marginTop: "10px" }}
                      >
                        {isMember ? "Open Chat" : "Join to Chat!"}
                      </Button>
                    )}

                  </CardContent>
                </Card>
              </Grid>
            );
          })}
        </Grid>


        <Button
          variant="contained"
          color="success"
          onClick={handleOpenDialog}
          style={{ marginTop: "20px", display: "block", width: "100%" }}
        >
          Create a Study Group
        </Button>

        <Dialog open={openDialog} onClose={handleCloseDialog}>
          <DialogTitle>Create a Study Group</DialogTitle>
          <DialogContent>
            <TextField
              label="Group Name"
              fullWidth
              value={groupName}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setGroupName(e.target.value)}
              margin="dense"
            />
            <TextField
              label="Group Description"
              fullWidth
              value={groupDescription}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setGroupDescription(e.target.value)}
              margin="dense"
            />
            <FormControl fullWidth style={{ marginTop: "10px" }}>
              <InputLabel>Group Type</InputLabel>
              <Select
                value={groupType}
                onChange={(e: SelectChangeEvent<string>) => setGroupType(e.target.value as "public" | "closed" | "invite-only")}
              >
                <MenuItem value="public">Public</MenuItem>
                <MenuItem value="closed">Closed</MenuItem>
                <MenuItem value="invite-only">Invite Only</MenuItem>
              </Select>
            </FormControl>
            <FormControl fullWidth style={{ marginTop: "10px" }}>
              <InputLabel>Select Module</InputLabel>
              <Select
                value={selectedGroupModule}
                onChange={(e: SelectChangeEvent<string>) => setSelectedGroupModule(e.target.value)}

              >
                {modulesList.map((module) => (
                  <MenuItem key={module.id} value={module.id}>
                    {module.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseDialog} color="error" variant="contained">
              Cancel
            </Button>
            <Button onClick={handleOpenInviteDialog} color="primary" variant="contained">
              Invite Members
            </Button>
            <Button onClick={() => { void handleCreateGroup(); }} color="success" variant="contained">
              Create
            </Button>
          </DialogActions>
        </Dialog>

        <Dialog open={openInviteDialog} onClose={handleCloseInviteDialog}>
          <DialogTitle>Invite a User</DialogTitle>
          <DialogContent>
            <TextField
              label="Invite Name"
              fullWidth
              value={inviteName}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setInviteName(e.target.value)}
              margin="dense"
            />
            <TextField
              label="Invite Email"
              fullWidth
              value={inviteEmail}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setInviteEmail(e.target.value)}
              margin="dense"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseInviteDialog} color="error" variant="contained">
              Cancel
            </Button>
            <Button onClick={handleInvite} color="primary" variant="contained">
              Invite
            </Button>
          </DialogActions>
        </Dialog>
      </Container>
    </Box>
  );
};

export default StudyGroupsPage;