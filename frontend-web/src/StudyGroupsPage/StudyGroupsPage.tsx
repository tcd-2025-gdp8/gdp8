import React, { useState, useEffect } from "react";
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
  FormControl,
  InputLabel
} from "@mui/material";

// Define TypeScript interfaces
interface StudyGroup {
  id: number;
  name: string;
  members: number;
  maximumMembers: number;
  module: string;
}

const modulesList = ["CSU44052: Computer Graphics",
                     "CSU44061: Machine Learning",
                     "CSU44051: Human Factors",
                     "CSU44000: Internet Applications",
                     "CSU44012: Topics in Functional Programming",
                     "CSU44099: Final Year Project",
                     "CSU44098: Group Design Project",
                     "CSU44081: Entrepreneurship & High Tech Venture Creation"];

const initialGroups: StudyGroup[] = [
  { id: 1, name: "Tech Nerds", members: 5, maximumMembers: 6, module: "CSU44052: Computer Graphics" },
  { id: 2, name: "CS Wizards", members: 3, maximumMembers: 5, module: "CSU44051: Human Factors" },
  { id: 3, name: "The Elites", members: 3, maximumMembers: 3, module: "CSU44052: Computer Graphics" },
  { id: 4, name: "The Fun Group", members: 6, maximumMembers: 6, module: "CSU44061: Machine Learning" },
  { id: 5, name: "The Prefects", members: 8, maximumMembers: 10, module: "CSU44051: Human Factors" },
  { id: 6, name: "Trinners for Winners", members: 7, maximumMembers: 8, module: "CSU44099: Final Year Project" },
];

const StudyGroupsPage: React.FC = () => {
  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [filteredGroups, setFilteredGroups] = useState<StudyGroup[]>([]);
  const [selectedModule, setSelectedModule] = useState<string>("");
  const [openDialog, setOpenDialog] = useState(false);
  const [groupName, setGroupName] = useState("");
  const [maxMembers, setMaxMembers] = useState(5);
  const [selectedGroupModule, setSelectedGroupModule] = useState("");

  useEffect(() => {
    // Simulated API call (Replace with actual fetch request)
    setTimeout(() => {
      setStudyGroups(initialGroups);
      setFilteredGroups(initialGroups);
    }, 500);
  }, []);

  // Handle filtering study groups by module
  useEffect(() => {
    if (selectedModule === "All" || selectedModule === "") {
      setFilteredGroups(studyGroups);
    } else {
      setFilteredGroups(studyGroups.filter((group) => group.module === selectedModule));
    }
  }, [selectedModule, studyGroups]);

  // Handle joining a group
  const handleJoinGroup = (id: number) => {
    setStudyGroups((prevGroups) =>
      prevGroups.map((group) =>
        group.id === id && group.members < group.maximumMembers
          ? { ...group, members: group.members + 1 }
          : group
      )
    );
  };

  // Handle opening and closing the create group dialog
  const handleOpenDialog = () => setOpenDialog(true);
  const handleCloseDialog = () => {
    setOpenDialog(false);
    setGroupName(""); // Reset input fields
    setMaxMembers(5);
    setSelectedGroupModule("");
  };

  // Handle creating a new group
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
    };

    setStudyGroups([...studyGroups, newGroup]);
    handleCloseDialog();
  };

  return (
    <Container maxWidth="md" style={{ marginTop: "20px" }}>
      <Typography variant="h4" gutterBottom>
        Study Groups
      </Typography>

      {/* Module Filter Dropdown */}
      <FormControl fullWidth style={{ marginBottom: "20px" }}>
        <InputLabel>Filter by Module</InputLabel>
        <Select
          value={selectedModule}
          onChange={(e) => setSelectedModule(e.target.value)}
        >
          <MenuItem value="All">All</MenuItem>
          {modulesList.map((module) => (
            <MenuItem key={module} value={module}>
              {module}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {/* Study Groups List */}
      <Grid container spacing={2}>
        {filteredGroups.map((group) => (
          <Grid item xs={12} sm={6} md={4} key={group.id} style={{ minWidth: "280px" }}>
            <Card>
              <CardContent>
                <Typography variant="h6">{group.name}</Typography>
                <Typography color="textSecondary">
                  Members: {group.members} / {group.maximumMembers}
                </Typography>
                <Typography color="textSecondary">
                  Module: {group.module}
                </Typography>
                <Button
                  variant="contained"
                  color="primary"
                  fullWidth
                  onClick={() => handleJoinGroup(group.id)}
                  style={{ marginTop: "10px" }}
                  disabled={group.members >= group.maximumMembers}
                >
                  {group.members >= group.maximumMembers ? "Full" : "Request to Join"}
                </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Create Group Button */}
      <Button
        variant="contained"
        color="success"
        onClick={handleOpenDialog}
        style={{ marginTop: "20px", display: "block", width: "100%" }}
      >
        Create a Study Group
      </Button>

      {/* Dialog for Creating a New Study Group */}
      <Dialog open={openDialog} onClose={handleCloseDialog}>
        <DialogTitle>Create a Study Group</DialogTitle>
        <DialogContent>
          <TextField
            label="Group Name"
            fullWidth
            value={groupName}
            onChange={(e) => setGroupName(e.target.value)}
            margin="dense"
          />
          <FormControl fullWidth style={{ marginTop: "10px" }}>
            <InputLabel>Select Module</InputLabel>
            <Select
              value={selectedGroupModule}
              onChange={(e) => setSelectedGroupModule(e.target.value)}
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
          <Button onClick={handleCreateGroup} color="success" variant="contained">
            Create
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default StudyGroupsPage;
