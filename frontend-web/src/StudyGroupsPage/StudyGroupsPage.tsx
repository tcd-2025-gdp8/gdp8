// src/StudyGroupsPage/StudyGroupsPage.tsx
import React, { useState, useEffect } from "react";
import { useAuth } from "../firebase/useAuth";
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
  IconButton
} from "@mui/material";
import NotificationsIcon from '@mui/icons-material/Notifications';
import CloseIcon from '@mui/icons-material/Close';

interface StudyGroupMember {
  userID: string;
  role: "admin" | "member" | "invitee" | "requester";
}

interface StudyGroupDetails {
  name: string;
  description: string;
  type: "public" | "closed" | "invite-only";
  moduleID: number;
}

interface StudyGroup {
  id: number;
  studyGroupDetails: StudyGroupDetails;
  members: StudyGroupMember[];
}

interface HardcodedStudyGroup {
  id: number;
  moduleID: number
  members: StudyGroupMember[];
}

interface Notification {
  id: number;
  message: string;
}

interface APIStudyGroup {
  id: number;
  name: string;
  description: string;
  type: "public" | "closed" | "invite-only";
}

const modulesList = [
  { id: 0, name: "All" },
  { id: 1, name: "CSU44052: Computer Graphics" },
  { id: 2, name: "CSU44061: Machine Learning" },
  { id: 3, name: "CSU44051: Human Factors" },
  { id: 4, name: "CSU44000: Internet Applications" },
  { id: 5, name: "CSU44012: Topics in Functional Programming" },
  { id: 6, name: "CSU44099: Final Year Project" },
  { id: 7, name: "CSU44098: Group Design Project" },
  { id: 8, name: "CSU44081: Entrepreneurship & High Tech Venture Creation" },
];

const initialGroups: HardcodedStudyGroup[] = [
  {
    id: 1,
    moduleID: 1,
    members: [
      { userID: "Alice", role: "member" },
      { userID: "Bob", role: "member" },
      { userID: "Charlie", role: "member" },
      { userID: "Maria", role: "member" },
      { userID: "Catriona", role: "member" },
    ],
  },
  {
    id: 2,
    moduleID: 3,
    members: [
      { userID: "David", role: "member" },
      { userID: "Eve", role: "member" },
      { userID: "Frank", role: "member" },
    ],
  },
  {
    id: 3,
    moduleID: 1,
    members: [
      { userID: "Grace", role: "member" },
      { userID: "Hannah", role: "member" },
      { userID: "Ian", role: "member" },
    ],
  },
  {
    id: 4,
    moduleID: 2,
    members: [
      { userID: "Jack", role: "member" },
      { userID: "Kate", role: "member" },
      { userID: "Leo", role: "member" },
      { userID: "Blake", role: "member" },
      { userID: "Robert", role: "member" },
      { userID: "Marco", role: "member" },
    ],
  },
  {
    id: 5,
    moduleID: 3,
    members: [
      { userID: "Mike", role: "member" },
      { userID: "Nina", role: "member" },
      { userID: "Oscar", role: "member" },
      { userID: "Alessandro", role: "member" },
      { userID: "Alice", role: "member" },
      { userID: "David", role: "member" },
      { userID: "Grace", role: "member" },
      { userID: "Ava", role: "member" },
    ],
  },
  {
    id: 6,
    moduleID: 6,
    members: [
      { userID: "Paul", role: "member" },
      { userID: "Quinn", role: "member" },
      { userID: "Rachel", role: "member" },
      { userID: "Jade", role: "member" },
      { userID: "Robert", role: "member" },
      { userID: "Bob", role: "member" },
      { userID: "Hannah", role: "member" },
      { userID: "Bianca", role: "member" },
      { userID: "Oscar", role: "member" },
      { userID: "Ava", role: "member" }
    ],
  },
];

const StudyGroupsPage: React.FC = () => {
  const { token } = useAuth();
  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [filteredGroups, setFilteredGroups] = useState<StudyGroup[]>([]);
  const [selectedModule, setSelectedModule] = useState<number | "">("");
  const [openDialog, setOpenDialog] = useState(false);
  const [groupName, setGroupName] = useState("");
  const [groupDescription, setGroupDescription] = useState("");
  const [groupType, setGroupType] = useState<"public" | "closed" | "invite-only">("public");
  const [selectedGroupModule, setSelectedGroupModule] = useState<number | "">("");
  const [openInviteDialog, setOpenInviteDialog] = useState(false);
  const [inviteName, setInviteName] = useState("");
  const [inviteEmail, setInviteEmail] = useState("");
  const [notifications, setNotifications] = useState<Notification[]>([
    { id: 1, message: 'Your request to join "The Prefects" has been accepted.' },
    { id: 2, message: 'New study group "CS Wizards" has been created for CSU44051: Human Factors.' },
    { id: 3, message: 'New study group "The Elites" has been created for CSU44052: Computer Graphics.' }
  ]);
  const [openNotifications, setOpenNotifications] = useState(false);

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
              moduleID: matchedGroup?.moduleID ?? 0
            },
            members: matchedGroup?.members ?? [],
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
    if (selectedModule === "" || selectedModule === 0) {
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
      alert("You are not authorized. Please log in.");
      return;
    }
  
    const newGroupDetails = {
      name: groupName,
      description: groupDescription,
      type: groupType,
      userID: "Alessandro",
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
        members: [{ userID: "Alessandro", role: "admin" }],
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
      alert(`Error creating study group: ${error instanceof Error ? error.message : String(error)}`);
    }
  };
  

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
    <Container maxWidth="md" style={{ marginTop: "20px" }}>
      <Typography variant="h4" gutterBottom>
        Study Groups
        <IconButton onClick={() => setOpenNotifications(true)}>
          <NotificationsIcon />
        </IconButton>
      </Typography>

      <Drawer anchor="right" open={openNotifications} onClose={() => setOpenNotifications(false)}>
        <div style={{ width: 300, padding: 16 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">Notifications</Typography>
            <IconButton onClick={() => setOpenNotifications(false)}>
              <CloseIcon />
            </IconButton>
          </div>
          {notifications.length === 0 ? (
            <Typography variant="body2" color="textSecondary" style={{ textAlign: 'center', marginTop: '20px' }}>
              No notifications
            </Typography>
          ) : (
            notifications.map((notification) => (
              <Card key={notification.id} style={{ marginBottom: '10px' }}>
                <CardContent>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
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
          onChange={(e: SelectChangeEvent<number | "">) => setSelectedModule(e.target.value as number | "")}
        >
          {modulesList.map((module) => (
            <MenuItem key={module.id} value={module.id}>
              {module.name}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      <Grid container spacing={2}>
        {filteredGroups.map((group) => (
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
                      : "No members yet"
                  }
                  arrow
                >
                  <Typography color="textSecondary" style={{ cursor: "pointer" }}>
                    Members: {group.members.length}
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
              </CardContent>
            </Card>
          </Grid>
        ))}
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
              onChange={(e: SelectChangeEvent<string | number>) => setSelectedGroupModule(e.target.value as number | "")}
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
  );
};

export default StudyGroupsPage;