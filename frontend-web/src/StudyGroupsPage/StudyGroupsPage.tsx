import React, { useState, useEffect } from "react";
import { Container, Card, CardContent, Typography, Button, Grid } from "@mui/material";

// Define a TypeScript interface for Study Groups
interface StudyGroup {
  id: number;
  name: string;
  members: number;
  maximumMembers: number;
}

const StudyGroupsPage: React.FC = () => {
  // Define state with TypeScript
  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);

  useEffect(() => {
    // Simulated API call (Replace with actual fetch request)
    setTimeout(() => {
      setStudyGroups([
        { id: 1, name: "Tech Nerds", members: 5, maximumMembers: 6 },
        { id: 2, name: "CS Wizards", members: 3, maximumMembers: 5 },
        { id: 3, name: "The Elites", members: 3, maximumMembers: 3 },
        { id: 4, name: "The Fun Group", members: 6, maximumMembers: 6 },
        { id: 5, name: "The Prefects", members: 8, maximumMembers: 10 },
        { id: 6, name: "Trinners for Winners", members: 7, maximumMembers: 8 },
      ]);
    }, 500);
  }, []);

  // Handle join group
  const handleJoinGroup = (id: number) => {
    setStudyGroups((prevGroups) =>
      prevGroups.map((group) =>
        group.id === id && group.members < group.maximumMembers
          ? { ...group, members: group.members + 1 }
          : group
      )
    );
  };

  const handleCreateGroup = () => {
    alert("Create a new study group!");
  };

  return (
    <Container maxWidth="md" style={{ marginTop: "20px" }}>
      <Typography variant="h4" gutterBottom>
        Available Study Groups
      </Typography>
      <Grid container spacing={2}>
        {studyGroups.map((group) => (
          <Grid item xs={12} sm={6} md={4} key={group.id}>
            <Card>
              <CardContent>
                <Typography variant="h6">{group.name}</Typography>
                <Typography color="textSecondary">
                  Members: {group.members} / {group.maximumMembers}
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
      <Button
        variant="contained"
        color="success"
        onClick={handleCreateGroup}
        style={{ marginTop: "20px", display: "block", width: "100%" }}
      >
        Create a Study Group
      </Button>
    </Container>
  );
};

export default StudyGroupsPage;
