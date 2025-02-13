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
  Slider,
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

interface StudyGroup {
  id: number;
  name: string;
  members: number;
  maximumMembers: number;
  module: string;
  membersList?: string[];
}

interface Notification {
  id: number;
  message: string;
}

const modulesList = [
  "CSU44052: Computer Graphics",
  "CSU44061: Machine Learning",
  "CSU44051: Human Factors",
  "CSU44000: Internet Applications",
  "CSU44012: Topics in Functional Programming",
  "CSU44099: Final Year Project",
  "CSU44098: Group Design Project",
  "CSU44081: Entrepreneurship & High Tech Venture Creation",
];

const initialGroups: StudyGroup[] = [
  {
    id: 1,
    name: "Tech Nerds",
    members: 5,
    maximumMembers: 6,
    module: "CSU44052: Computer Graphics",
    membersList: ["Alice", "Bob", "Charlie", "Maria", "Catriona"],
  },
  {
    id: 2,
    name: "CS Wizards",
    members: 3,
    maximumMembers: 5,
    module: "CSU44051: Human Factors",
    membersList: ["David", "Eve", "Frank"],
  },
  {
    id: 3,
    name: "The Elites",
    members: 3,
    maximumMembers: 3,
    module: "CSU44052: Computer Graphics",
    membersList: ["Grace", "Hannah", "Ian"],
  },
  {
    id: 4,
    name: "The Fun Group",
    members: 6,
    maximumMembers: 6,
    module: "CSU44061: Machine Learning",
    membersList: ["Jack", "Kate", "Leo", "Blake", "Robert", "Marco"],
  },
  {
    id: 5,
    name: "The Prefects",
    members: 8,
    maximumMembers: 10,
    module: "CSU44051: Human Factors",
    membersList: ["Mike", "Nina", "Oscar", "Alessandro", "Alice", "David", "Grace", "Ava"],
  },
  {
    id: 6,
    name: "Trinners for Winners",
    members: 7,
    maximumMembers: 8,
    module: "CSU44099: Final Year Project",
    membersList: ["Paul", "Quinn", "Rachel", "Jade", "Robert", "Bob", "Hannah"],
  },
];

const StudyGroupsPage: React.FC = () => {
  const { token } = useAuth();
  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [filteredGroups, setFilteredGroups] = useState<StudyGroup[]>([]);
  const [selectedModule, setSelectedModule] = useState<string>("");
  const [openDialog, setOpenDialog] = useState(false);
  const [groupName, setGroupName] = useState("");
  const [maxMembers, setMaxMembers] = useState(5);
  const [selectedGroupModule, setSelectedGroupModule] = useState("");
  const [openInviteDialog, setOpenInviteDialog] = useState(false); // State to manage the invite dialog
  const [inviteName, setInviteName] = useState("");
  const [inviteEmail, setInviteEmail] = useState("");
  const [notifications, setNotifications] = useState<Notification[]>([{
    id: 1,
    message: 'Your request to join "The Prefects" has been accepted.'
  }, {
    id: 2,
    message: 'New study group "CS Wizards" has been created for CSU44051: Human Factors.'
  }, {
    id: 3,
    message: 'New study group "The Elites" has been created for CSU44052: Computer Graphics.'
  }]);
  const [openNotifications, setOpenNotifications] = useState(false);

  // When the access token is available, log it and simulate fetching data.
  useEffect(() => {
    if (token) {
      console.log("Access token available:", token);
      // Simulate an API call with a slight delay
      setTimeout(() => {
        setStudyGroups(initialGroups);
        setFilteredGroups(initialGroups);
      }, 500);
    }
  }, [token]);

  useEffect(() => {
    if (selectedModule === "All" || selectedModule === "") {
      setFilteredGroups(studyGroups);
    } else {
      setFilteredGroups(
        studyGroups.filter((group) => group.module === selectedModule)
      );
    }
  }, [selectedModule, studyGroups]);

  const handleJoinGroup = (id: number) => {
    let joinedGroupName = "";
  
    const updatedGroups = studyGroups.map((group) => {
      if (group.id === id && group.members < group.maximumMembers) {
        if (!group.membersList?.includes("Alessandro")) {
          joinedGroupName = group.name; // Capture the name of the group joined
          return {
            ...group,
            members: group.members + 1,
            membersList: [...(group.membersList ?? []), "Alessandro"],
          };
        }
      }
      return group;
    });
  
    if (joinedGroupName) {
      setStudyGroups(updatedGroups);
      setNotifications((prev) => [
        ...prev,
        { id: Date.now(), message: `You joined '${joinedGroupName}'.` },
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
    setMaxMembers(5);
    setSelectedGroupModule("");
  };

  const handleCreateGroup = () => {
    if (groupName.trim() === "" || selectedGroupModule === "") {
      alert("Please enter a valid group name and select a module.");
      return;
    }

    const newGroup: StudyGroup = {
      id: studyGroups.length + 1,
      name: groupName,
      members: 1,
      maximumMembers: maxMembers,
      module: selectedGroupModule,
      membersList: ["Alessandro"],
    };

    setStudyGroups([...studyGroups, newGroup]);
    setNotifications((prev) => [...prev, { id: Date.now(), message: `Group '${groupName}' created.` }]);
    handleCloseDialog();
  };

  const handleOpenInviteDialog = () => setOpenInviteDialog(true); // Open invite dialog
  const handleCloseInviteDialog = () => setOpenInviteDialog(false); // Close invite dialog

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
        <Select value={selectedModule} onChange={(e: SelectChangeEvent<string>) => setSelectedModule(e.target.value)} >
          <MenuItem value="All">All</MenuItem>
          {modulesList.map((module) => (
            <MenuItem key={module} value={module}>
              {module}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      <Grid container spacing={2}>
        {filteredGroups.map((group) => (
          <Grid
            item
            xs={12}
            sm={6}
            md={4}
            key={group.id}
            style={{ minWidth: "280px" }}
          >
            <Card>
              <CardContent>
                <Typography variant="h6">{group.name}</Typography>
                <Tooltip
                  title={
                    group.membersList
                      ? group.membersList.join(", ")
                      : "No members yet"
                  }
                  arrow
                >
                  <Typography color="textSecondary" style={{ cursor: "pointer" }}>
                    Members: {group.members} / {group.maximumMembers}
                  </Typography>
                </Tooltip>
                <Typography color="textSecondary">
                  Module: {group.module}
                </Typography>
                <Button
                  variant="contained"
                  color="primary"
                  fullWidth
                  onClick={() => handleJoinGroup(group.id)}
                  style={{ marginTop: "10px" }}
                  disabled={
                    group.members >= group.maximumMembers ||
                    group.membersList?.includes("Alessandro")
                  }
                >
                  {group.membersList?.includes("Alessandro")
                    ? "Joined"
                    : group.members >= group.maximumMembers
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
            onChange={(e: React.ChangeEvent<{ value: unknown }>) => setSelectedGroupModule(e.target.value as string)}
            margin="dense"
          />
          <FormControl fullWidth style={{ marginTop: "10px" }}>
            <InputLabel>Select Module</InputLabel>
            <Select
              value={selectedGroupModule}
              onChange={(e:SelectChangeEvent<string>) => setSelectedGroupModule(e.target.value)}
            >
              {modulesList.map((module) => (
                <MenuItem key={module} value={module}>
                  {module}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
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
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog} color="error" variant="contained">
            Cancel
          </Button>
          <Button onClick={handleOpenInviteDialog} color="primary" variant="contained">
            Invite Members
          </Button>
          <Button onClick={handleCreateGroup} color="success" variant="contained">
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
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setGroupName(e.target.value)}
            margin="dense"
          />
          <TextField
            label="Invite Email"
            fullWidth
            value={inviteEmail}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setInviteName(e.target.value)}
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