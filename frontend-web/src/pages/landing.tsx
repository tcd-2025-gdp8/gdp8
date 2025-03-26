import { useState, useEffect, useCallback } from "react";
import {
  Container,
  Typography,
  Grid,
  CircularProgress,
  Alert,
  Paper,
  List,
  ListItem,
  ListItemText,
  Box,
} from "@mui/material";
import { fetchUserStudyGroups, StudyGroup } from "../api/studyGroups";
import { getUserModules } from "../api/users";
import { useAuth } from "../auth/useAuth";
import StudyGroupCard from "../components/StudyGroupCard";
import { Module } from "../api/modules";
import { fetchAllUserStudySessions, Session } from "../api/studySessions";

export default function LandingPage() {
  const { token, user } = useAuth();
  const currentUserId = user?.uid;

  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [modulesList, setModulesList] = useState<Module[]>([]);
  const [studySessions, setStudySessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    if (!token || !currentUserId) return;

    try {
      const [userStudyGroups, userModules] = await Promise.all([
        fetchUserStudyGroups(token, currentUserId),
        getUserModules(token, currentUserId),
      ]);

      const userStudySessions = await fetchAllUserStudySessions(token, userStudyGroups);
      setStudySessions(userStudySessions);
      setStudyGroups(userStudyGroups);
      setModulesList(userModules || []);
    } catch (err) {
      console.error("Error fetching study groups or modules:", err);
      setError("Failed to load study groups.");
    } finally {
      setLoading(false);
    }
  }, [token, currentUserId]);

  useEffect(() => {
    void fetchData();
  }, [fetchData]);

  return (
    <Container maxWidth="md" style={{ marginTop: "20px", textAlign: "center" }}>
      <Typography variant="h4" style={{ marginBottom: "2rem" }}>
        Welcome to Your Study Hub
      </Typography>

      {/* Study Groups Section */}
      <Typography variant="h5" style={{ marginBottom: "1.5rem" }}>
        Your Study Groups
      </Typography>
      {loading && <CircularProgress />}
      {error && <Alert severity="error">{error}</Alert>}
      {!loading && studyGroups.length === 0 && (
        <Typography>No study groups found. Join or create one!</Typography>
      )}

      <Grid container spacing={2} style={{ marginBottom: "2rem" }}>
        {studyGroups.map((group) => {
          const module = modulesList.find((mod) => mod.id === group.moduleId);
          if (!module) return null;
          return (
            <Grid item xs={12} sm={6} md={4} key={group.id}>
              <StudyGroupCard
                group={group}
                module={module}
                currentUserId={currentUserId!}
                onUpdate={() => void fetchData()}
              />
            </Grid>
          );
        })}
      </Grid>

      {/* Study Sessions Section */}
      <Typography variant="h5" style={{ marginBottom: "1.5rem", marginTop: "3rem" }}>
        Upcoming Study Sessions
      </Typography>
      {loading && <CircularProgress />}
      {error && <Alert severity="error">{error}</Alert>}
      {!loading && studySessions.length === 0 && (
        <Typography>No upcoming study sessions found.</Typography>
      )}

      <Paper
        sx={{
          p: 2,
          bgcolor: "#3b5998",
          borderRadius: "8px",
          maxHeight: "220px",
          overflowY: "auto",
          width: "100%",
          maxWidth: "1200px",
          marginTop: "20px",
        }}
      >
        <List>
          {studySessions.map((session) => (
            <ListItem
              key={session.id}
              sx={{
                backgroundColor: "white",
                borderRadius: "4px",
                mb: 1,
                padding: 1.5,
                "&:hover": { backgroundColor: "#f5f5f5" },
              }}
            >
              <ListItemText
                primary={`${session.name}, ${
                  session.date ? session.date.toLocaleDateString() : ""
                } ${session.earliestTime} - ${session.latestTime}`}
              />
            </ListItem>
          ))}
        </List>
      </Paper>
    </Container>
  );
}