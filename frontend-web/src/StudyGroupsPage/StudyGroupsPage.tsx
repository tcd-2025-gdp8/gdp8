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
  Box,
  List,
  ListItem,
  ListItemText,
  Badge,
  Toolbar,
  AppBar
} from "@mui/material";
import NotificationsIcon from '@mui/icons-material/Notifications';
import CloseIcon from '@mui/icons-material/Close';
import { useNavigate } from "react-router-dom";
import DeleteIcon from "@mui/icons-material/Delete";
import { removeMemberFromStudyGroup } from "../api/studyGroups";

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
  id: number;
  code: string;
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
  const { token, user } = useAuth();
  const navigate = useNavigate();
  const currentUserID = user?.uid;
  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [filteredGroups, setFilteredGroups] = useState<StudyGroup[]>([]);
  const [selectedModule, setSelectedModule] = useState<string>("");
  const [openDialog, setOpenDialog] = useState(false);
  const [groupName, setGroupName] = useState("");
  const [groupDescription, setGroupDescription] = useState("");
  const [groupType, setGroupType] = useState<"public" | "closed" | "invite-only">("public");
  const [selectedGroupModuleCode, setSelectedGroupModuleCode] = useState<string>("");
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
  const [openMemberDialog, setOpenMemberDialog] = useState(false);
  const [selectedGroupMembers, setSelectedGroupMembers] = useState<StudyGroupMember[]>([]);
  const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);
  const isCurrentUserAdmin = selectedGroupMembers.some(
    (member) => member.userID === currentUserID && member.role === "admin"
  );

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
    if (!token) {
      alert("You are not authorised. Please log in.");
      return;
    }

    void (async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/study-groups/${id}/request-to-join`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });

        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`Failed to join study group: ${errorText}`);
        }

        setStudyGroups((prevGroups) =>
          prevGroups.map((group) => {
            if (group.id === id) {
              return {
                ...group,
                members: [
                  ...group.members,
                  { userID: currentUserID ?? "", role: "member" as const },
                ],
              };
            }
            return group;
          })
        );

        setNotifications((prev) => [
          ...prev,
          { id: Date.now(), message: `You joined the study group successfully.` },
        ]);
      } catch (error) {
        console.error("Error joining study group:", error);
        alert(`Error joining study group: ${error instanceof Error ? error.message : String(error)}`);
      }
    })();
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
    setSelectedGroupModuleCode("");
  };

  const handleCreateGroup = async () => {
    if (groupName.trim() === "" || selectedGroupModuleCode === "") {
      alert("Please enter a valid group name and select a module.");
      return;
    }

    if (!token) {
      alert("You are not authorised. Please log in.");
      return;
    }

    const selectedModuleId = modulesList.find(module => module.code === selectedGroupModuleCode)?.id;

    const newGroupDetails = {
      name: groupName,
      description: groupDescription,
      type: groupType,
      moduleID: selectedModuleId
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
          moduleID: selectedGroupModuleCode,
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
          message: `New study group "${groupName}" has been created for ${modulesList.find((module) => module.code === selectedGroupModuleCode)?.name
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
    if (!token || !currentUserID) return;

    const fetchModules = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/user/${currentUserID}/modules`, {
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

        // Add "All" option at the beginning
        setModulesList([{ id: 0, code: "All", name: "All" }, ...data]);

      } catch (err) {
        console.error("Error fetching modules:", err);
      }
    };

    // Properly handle the async function call
    void fetchModules();
  }, [token, currentUserID]);

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

  const handleOpenMembersDialog = (groupId: number, members: StudyGroupMember[]) => {
    setSelectedGroupId(groupId);
    setSelectedGroupMembers(members);
    setOpenMemberDialog(true);
  };

  const handleCloseMembersDialog = () => {
    setOpenMemberDialog(false);
    setSelectedGroupMembers([]);
  };

  const handleRemoveMember = async (memberId: string) => {
    if (!token || selectedGroupId === null) return;

    try {
      // Call the API function that handles member removal.
      await removeMemberFromStudyGroup(token, selectedGroupId, memberId);
      setSelectedGroupMembers(prevMembers =>
        prevMembers.filter((member) => member.userID !== memberId)
      );
    } catch (err) {
      console.error("Error removing member:", err);
    }
  };




  return (
    <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh" width="100vw">
      {/* Back to Landing Button */}
      <Button
        variant="contained"
        sx={{
          backgroundColor: "#0056b3",
          "&:hover": {
            backgroundColor: "#004494",
          }
        }}
        onClick={() => { void navigate("/landing"); }}
        style={{ position: "absolute", top: "10px", left: "10px" }}
      >
        Back to Landing
      </Button>
      <AppBar
        position="fixed"
        sx={{
          backgroundColor: "#ffffff",
          color: "#000000",
          zIndex: 1201,
          width: "calc(100% - 240px)",
          marginLeft: "240px",
        }}
      >
        <Toolbar sx={{ flexDirection: "row", alignItems: "center" }}>
          <Box sx={{ flexGrow: 1, display: "flex", alignItems: "center" }}>
            <Typography variant="h6">
              WeStudy
            </Typography>
          </Box>
          {user && (
            <Typography variant="subtitle1" sx={{ mr: 2 }}>
              {user.email}
            </Typography>
          )}
          <IconButton color="inherit" onClick={() => setOpenNotifications(true)}>
            <Badge badgeContent={notifications.length} color="error">
              <NotificationsIcon />
            </Badge>
          </IconButton>
        </Toolbar>
      </AppBar>

      {/* Notifications Drawer */}
      <Drawer
        anchor="right"
        open={openNotifications}
        onClose={() => setOpenNotifications(false)}
        sx={{
          "& .MuiDrawer-paper": {
            zIndex: 1300,
          },
        }}
      >
        <Box sx={{ width: 300, padding: "1rem" }}>
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              marginBottom: "1rem",
            }}
          >
            <Typography variant="h6">Notifications</Typography>
            <IconButton onClick={() => setOpenNotifications(false)}>
              <CloseIcon />
            </IconButton>
          </Box>
          {notifications.length === 0 ? (
            <Typography
              variant="body2"
              color="textSecondary"
              sx={{ textAlign: "center" }}
            >
              No notifications
            </Typography>
          ) : (
            notifications.map((notification) => (
              <Box
                key={notification.id}
                sx={{
                  padding: "0.5rem",
                  borderBottom: "1px solid #ccc",
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                }}
              >
                <Typography>{notification.message}</Typography>
                <IconButton
                  size="small"
                  onClick={() => handleDeleteNotification(notification.id)}
                >
                  <CloseIcon fontSize="small" />
                </IconButton>
              </Box>
            ))
          )}
        </Box>
      </Drawer>

      <Container maxWidth="md" style={{ marginTop: "20px", position: "relative", textAlign: "center" }}>
        <Typography variant="h4" gutterBottom>
          Study Groups
        </Typography>

        <FormControl fullWidth style={{ marginBottom: "20px" }}>
          <InputLabel>Filter by Module</InputLabel>
          <Select
            value={selectedModule}
            onChange={(e: SelectChangeEvent<string>) => setSelectedModule(e.target.value)}
          >
            {modulesList.map((module) => (
              <MenuItem key={module.code} value={module.code}>
                {module.code === "All" ? module.name : `${module.code}: ${module.name}`}
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
                      arrow
                      title={
                        group.studyGroupDetails.type === "closed"
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
                          onClick={() => {
                            if (group.studyGroupDetails.type !== "closed") {
                              handleOpenMembersDialog(group.id, group.members);
                            }
                          }}
                          disabled={group.studyGroupDetails.type === "closed"}
                        >
                          View Members
                        </Button>
                      </span>
                    </Tooltip>

                    <Typography color="textSecondary">
                      Module: {modulesList.find((module) => module.code === group.studyGroupDetails.moduleID)?.name}
                    </Typography>
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
                        sx={{
                          backgroundColor: "#0056b3",
                          "&:hover": {
                            backgroundColor: "#004494",
                          }
                        }}
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
        <Dialog open={openMemberDialog} onClose={handleCloseMembersDialog} maxWidth="sm" fullWidth>
          <DialogContent>
            {selectedGroupMembers.length > 0 && (
              <List>
                {selectedGroupMembers.map((member) => (
                  <ListItem key={member.userID}>
                    <ListItemText primary={member.userID} secondary={member.role} />
                    {isCurrentUserAdmin && member.role !== "admin" && (
                      <IconButton edge="end" onClick={() => void handleRemoveMember(member.userID)}>
                        <DeleteIcon color="error" />
                      </IconButton>
                    )}
                  </ListItem>
                ))}
              </List>
            )}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseMembersDialog} color="primary">
              Close
            </Button>
          </DialogActions>
        </Dialog>

        <Button
          variant="contained"
          color="success"
          onClick={handleOpenDialog}
          sx={{
            backgroundColor: "#0056b3",
            "&:hover": {
              backgroundColor: "#004494",
            }
          }}
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
                value={selectedGroupModuleCode}
                onChange={(e: SelectChangeEvent<string>) => setSelectedGroupModuleCode(e.target.value)}

              >
                {modulesList.map((module) => (
                  <MenuItem key={module.code} value={module.code}>
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
            <Button onClick={handleOpenInviteDialog} variant="contained"
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
      </Container>
    </Box>
  );
};

export default StudyGroupsPage;